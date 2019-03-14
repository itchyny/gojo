# gojo [![Travis Build Status](https://travis-ci.org/itchyny/gojo.svg?branch=master)](https://travis-ci.org/itchyny/gojo)
Yet another Go implementation of [jo](https://github.com/jpmens/jo).

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
```

## Installation
### Homebrew
```sh
 $ brew install itchyny/gojo/gojo
```

### Build from source
```bash
 $ go get -u github.com/itchyny/gojo/cmd/gojo
```

## Difference to jo
- Implemented in Go and Go-gettable
- Drops support of `k@v` syntax (use `k=true` or `k=false`)
- Does not print duplicate keys (although duplicate key in JSON is valid but it's not that useful and overwritten by the latter)
- Reading file contents is not implemented yet

## Bug Tracker
Report bug at [Issues・itchyny/gojo - GitHub](https://github.com/itchyny/gojo/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
