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

## Bug Tracker
Report bug at [Issues・itchyny/gojo - GitHub](https://github.com/itchyny/gojo/issues).

## Author
itchyny (https://github.com/itchyny)

## License
This software is released under the MIT License, see LICENSE.
