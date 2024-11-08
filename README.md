 
## EnglishProfile.org scraper

Small scraper to collect word data from https://www.englishprofile.org/wordlists/evp

Inspired by https://github.com/Granitosaurus/englishprofile-scraper

### Data:

englishprofile.json for word previews

worddata.json for full dataset

(last scraped 2024-11-08)

### Build:

```
$ go version
go version go1.23.3 linux/amd64

$ go mod tidy

$ go build
```

### Run:

```
$ ./englishprofile --help
Usage:
  englishprofile [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  discover    Discover word previews from word pagination
  help        Help about any command
  worddata    Collect word data from discovered word previews

Flags:
  -h, --help   help for englishprofile

Use "englishprofile [command] --help" for more information about a command.

$ ./englishprofile discover
$ ./englishprofile worddata
```

### Check and View:

```
$ cat ~/.jq
def profile_object:
    to_entries | def parse_entry: {"key": .key, "value": .value | type}; map(parse_entry)
        | sort_by(.key) | from_entries;

def profile_array_objects:
    map(profile_object) | map(to_entries) | reduce .[] as $item ([]; . + $item) | sort_by(.key) | from_entries;

$ cat englishprofile.json | jq "profile_array_objects"
{
  "baseword": "string",
  "guideword": "string",
  "level": "string",
  "partofspeech": "string",
  "topic": "string",
  "url": "string"
}

$ cat worddata.json | jq "profile_array_objects"
{
  "baseword": "string",
  "guideword": "string",
  "level": "string",
  "partofspeech": "string",
  "pronunciation": "string",
  "senses": "array",
  "topic": "string",
  "url": "string",
  "word_type": "string"
}

$ jq .[].baseword englishprofile.json | wc -l
   15696

$ jq .[].baseword worddata.json | wc -l
   15696

$ jq '.[] | select(.level == "A1")' worddata.json | jq -r '.baseword' | wc -l
     784

$ jq '.[] | select(.level == "A2")' worddata.json | jq -r '.baseword' | wc -l
    1594

$ jq '.[] | select(.level == "B1")' worddata.json | jq -r '.baseword' | wc -l
    2937

$ jq '.[] | select(.level == "B2")' worddata.json | jq -r '.baseword' | wc -l
    4164

$ jq '.[] | select(.level == "C1")' worddata.json | jq -r '.baseword' | wc -l
    2410

$ jq '.[] | select(.level == "C2")' worddata.json | jq -r '.baseword' | wc -l
    3807

$ echo 784+1594+2937+4164+2410+3807 | bc -l
15696
```
