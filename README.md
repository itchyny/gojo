# gojo
[![CI Status](https://github.com/itchyny/gojo/workflows/CI/badge.svg)](https://github.com/itchyny/gojo/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/itchyny/gojo)](https://goreportcard.com/report/github.com/itchyny/gojo)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/itchyny/gojo/blob/main/LICENSE)
[![release](https://img.shields.io/github/release/itchyny/gojo/all.svg)](https://github.com/itchyny/gojo/releases)
[![pkg.go.dev](https://pkg.go.dev/badge/github.com/itchyny/gojo)](https://pkg.go.dev/github.com/itchyny/gojo)

### Yet another Go implementation of [jo](https://github.com/jpmens/jo)
This is an implementation of [jo command](https://github.com/jpmens/jo) written in Go language.

## Usage
```sh
 $ gojo foo=bar qux=quux
{"foo":"bar","qux":"quux"}
 $ gojo -p foo=bar qux=quux
{
  "foo": "bar",
  "qux": "quux"
}
 $ gojo -a foo bar baz
["foo","bar","baz"]
 $ seq 10 | gojo -a
[1,2,3,4,5,6,7,8,9,10]
 $ gojo -p foo=$(gojo bar=$(gojo baz=100))
{
  "foo": {
    "bar": {
      "baz": 100
    }
  }
}
 $ gojo -p foo[bar][baz][qux][quux]=128
{
  "foo": {
    "bar": {
      "baz": {
        "qux": {
          "quux": 128
        }
      }
    }
  }
}
 $ gojo -p res[foo][][id]=10 res[foo][][id]=20 res[cnt]=2
{
  "res": {
    "foo": [
      {
        "id": 10
      },
      {
        "id": 20
      }
    ],
    "cnt": 2
  }
}
 $ gojo foo=@testdata/file.txt    # read contents from file
{"foo":"a\nb\nc\nd\ne"}
 $ gojo foo=:testdata/file1.json  # json value of file
{"foo":{"x":1,"z":2,"y":3}}
 $ gojo foo=%testdata/file.txt    # base64 of file contents
{"foo":"YQpiCmMKZAplCg=="}
```

## Installation
### Homebrew
```sh
brew install itchyny/tap/gojo
```

### Build from source
```bash
go get github.com/itchyny/gojo/cmd/gojo
```

## Difference to jo
- Implemented in Go and Go-gettable, also useful for library usage.
- Implements nested paths (example: `foo[x][y][z]=1`, `foo[][][]=1`).
- Drops support of `k@v` syntax (use `k=true` or `k=false`) and type guessing flags.
- Does not print duplicate keys (although duplicate key in JSON is valid, it's not that useful and overwritten by the latter).

## Bug Tracker
Report bug at [Issues・itchyny/gojo - GitHub](https://github.com/itchyny/gojo/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
