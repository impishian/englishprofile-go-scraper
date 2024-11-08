 
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

```

### By Level:

```
$ jq '.[] | select(.level == "A1")' worddata.json | jq -r '.baseword' | wc -l
     784
...

$ echo 784+1594+2937+4164+2410+3807 | bc -l
15696
```

| Level  | Count |
| ------------- | ------------- |
| A1  | 784  |
| A2  | 1594  |
| B1  | 2937  |
| B2  | 4164  |
| C1  | 2410  |
| C2  | 3807  |
| Total | 15696 |

```
$ jq '.[] | select(.level == "A1" or .level == "A2")' worddata.json | jq -r '.baseword' | wc -l
    2378
```
| Level  | Count |
| ------------- | ------------- |
| A1+A2  | 2378  |
| B1+B2  | 7101  |
| C1+C2  | 6217  |
| Total | 15696 |

### By Topic:

| Topic  | Count |
| ------------- | ------------- |
| animals | 94 |
| arts and media | 242 |
| body and health | 356 |
| clothes | 88 |
| communication | 1225 |
| crime | 84 |
| describing things | 748 |
| education | 92 |
| food and drink | 253 |
| homes and buildings | 170 |
| money | 153 |
| natural world | 261 |
| people: actions | 952 |
| people: appearance | 91 |
| people: personality | 1053 |
| politics | 72 |
| relationships | 170 |
| shopping | 205 |
| sports and games | 0 |
| technology | 161 |
| travel | 231 |
| work | 201 |
| Total | 6902 |

Topic is empty:

```
$ jq '.[] | select(.topic == "")' worddata.json | jq -r '.baseword' | wc -l
    8794

$ echo 6902+8794 | bc -l
15696
```
