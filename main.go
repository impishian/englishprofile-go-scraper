package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "sync"
    "io/ioutil"

    "github.com/PuerkitoBio/goquery"
    "github.com/spf13/cobra"
)

type WordPreview struct {
    Baseword     string `json:"baseword"`
    Guideword    string `json:"guideword"`
    Level        string `json:"level"`
    PartOfSpeech string `json:"partofspeech"`
    Topic        string `json:"topic"`
    URL          string `json:"url"`
}

type WordSense struct {
    Definition        string `json:"definition"`
    Label             string `json:"label"`
    DictExample       string `json:"dict_example"`
    LearnerExample    string `json:"learner_example"`
    LearnerExampleCite string `json:"learner_example_cite"`
}

type WordData struct {
    WordPreview
    Pronunciation string      `json:"pronunciation"`
    WordType      string      `json:"word_type"`
    Senses        []WordSense `json:"senses"`
}

func discoverWords() ([]WordPreview, error) {
    body := "filter_search=&filter_custom_Topic=&filter_custom_Parts=&filter_custom_Category=&filter_custom_Grammar=&filter_custom_Usage=&filter_custom_Prefix=&filter_custom_Suffix=&limit=0&directionTable=asc&sortTable=base&task=&boxchecked=0&filter_order=pos_rank&filter_order_Dir=asc&ce91224c5693e21d15ac97cc105e6520=1"
    resp, err := http.Post("https://www.englishprofile.org/wordlists/evp", "application/x-www-form-urlencoded", strings.NewReader(body))
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return nil, err
    }

    var data []WordPreview
    doc.Find("#reportList>tbody>tr").Each(func(i int, s *goquery.Selection) {
        data = append(data, WordPreview{
            Baseword:     s.Find("td:nth-child(1)").Text(),
            Guideword:    s.Find("td:nth-child(2)").Text(),
            Level:        s.Find("td:nth-child(3)>span").Text(),
            PartOfSpeech: s.Find("td:nth-child(4)").Text(),
            Topic:        s.Find("td:nth-child(5)").Text(),
            URL:          s.Find("td:nth-child(6)>a").AttrOr("href", ""),
        })
    })

    return data, nil
}

func scrapeWordPage(word WordPreview, client *http.Client) (WordData, error) {
    url := "https://www.englishprofile.org/" + word.URL
    resp, err := client.Get(url)
    if err != nil {
        return WordData{}, err
    }
    defer resp.Body.Close()

    doc, err := goquery.NewDocumentFromReader(resp.Body)
    if err != nil {
        return WordData{}, err
    }

    wordData := WordData{
        WordPreview:  word,
        Pronunciation: doc.Find(".written").Text(),
        WordType:      doc.Find(".pos").Text(),
    }

    doc.Find(".info.sense").Each(func(i int, s *goquery.Selection) {
        wordData.Senses = append(wordData.Senses, WordSense{
            Definition:        s.Find("span.definition").Text(),
            Label:             s.Find(".label").Text(),
            DictExample:       s.Find(".example p.blockquote").Text(),
            LearnerExample:    s.Find(".learnerexamp").Text(),
            LearnerExampleCite: s.Find(".learnerexamp span").Text(),
        })
    })

    return wordData, nil
}

func scrapeWords(words []WordPreview, filename string, conLimit int, batchSize int) error {
    client := &http.Client{}
    var wg sync.WaitGroup
    var mu sync.Mutex

    f, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    f.WriteString("[\n")

    var first = true

    for i := 0; i < len(words); i += batchSize {
        end := i + batchSize
        if end > len(words) {
            end = len(words)
        }

        batch := words[i:end]
        wg.Add(len(batch))

        for _, word := range batch {
            go func(word WordPreview) {
                defer wg.Done()
                fullWord, err := scrapeWordPage(word, client)
                if err != nil {
                    fmt.Println("Error scraping word:", err)
                    return
                }

                mu.Lock()
                defer mu.Unlock()
                fullWordJSON, _ := json.Marshal(fullWord)
                if !first {
                    f.WriteString(",\n")
                }
                f.WriteString(string(fullWordJSON))
                first = false
            }(word)
        }

        wg.Wait()
    }

    f.WriteString("\n]")
    return nil
}

func main() {
    var rootCmd = &cobra.Command{Use: "englishprofile"}

    var discoverCmd = &cobra.Command{
        Use:   "discover",
        Short: "Discover word previews from word pagination",
        Run: func(cmd *cobra.Command, args []string) {
            words, err := discoverWords()
            if err != nil {
                fmt.Println("Error discovering words:", err)
                return
            }
            data, err := json.MarshalIndent(words, "", "  ")
            if err != nil {
                fmt.Println("Error marshalling data:", err)
                return
            }
            ioutil.WriteFile("englishprofile.json", data, 0644)
        },
    }

    var worddataCmd = &cobra.Command{
        Use:   "worddata",
        Short: "Collect word data from discovered word previews",
        Run: func(cmd *cobra.Command, args []string) {
            speed, _ := cmd.Flags().GetInt("speed")
            wordsFile, err := os.ReadFile("englishprofile.json")
            if err != nil {
                fmt.Println("Error reading file:", err)
                return
            }
            var words []WordPreview
            if err := json.Unmarshal(wordsFile, &words); err != nil {
                fmt.Println("Error unmarshalling JSON:", err)
                return
            }
            err = scrapeWords(words, "worddata.json", speed, 12)
            if err != nil {
                fmt.Println("Error scraping words:", err)
            }
        },
    }
    worddataCmd.Flags().Int("speed", 4, "connection speed")

    rootCmd.AddCommand(discoverCmd, worddataCmd)
    rootCmd.Execute()
}
