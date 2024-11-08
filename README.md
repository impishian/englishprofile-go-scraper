 
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

$ cat englishprofile.json
...

$ ./englishprofile worddata

$ cat worddata.json
...
```
