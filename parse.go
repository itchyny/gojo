package gojo

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func parseKeyValue(s string) (setter, error) {
	i := strings.IndexByte(s, '=')
	if i < 0 {
		return nil, errParse(s)
	}
	key := s[:i]
	value, err := parseValue(s[i+1:])
	if err != nil {
		return nil, err
	}
	return buildSetter(key, value, false), nil
}

func parseValue(s string) (interface{}, error) {
	if s == "null" {
		return nil, nil
	}
	m := orderedmap.New()
	if err := json.Unmarshal([]byte(s), m); err == nil {
		return m, nil
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err == nil {
		return v, nil
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i, nil
	}
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		if i, err := strconv.ParseInt(s[2:], 16, 64); err == nil {
			return i, nil
		}
	}
	if strings.HasPrefix(s, "@") {
		cnt, err := ioutil.ReadFile(s[1:])
		if err != nil {
			return nil, err
		}
		for _, c := range []byte{'\n', '\r'} {
			if l := len(cnt); l > 0 && cnt[l-1] == c {
				cnt = cnt[:l-1]
			}
		}
		return string(cnt), nil
	}
	if strings.HasPrefix(s, "%") {
		cnt, err := ioutil.ReadFile(s[1:])
		if err != nil {
			return nil, err
		}
		return base64.StdEncoding.EncodeToString(cnt), nil
	}
	return s, nil
}

func buildSetter(key string, value interface{}, inner bool) setter {
	i, j := strings.IndexByte(key, '['), strings.IndexByte(key, ']')
	if i < 0 || j < 0 || j < i {
		if inner {
			return nil
		}
		return &objectSetter{key, value}
	}
	if i > 0 || !inner {
		s := buildSetter(key[i:], value, true)
		if s == nil {
			if inner {
				return nil
			}
			return &objectSetter{key, value}
		}
		return &objectSetter{key[:i], s}
	}
	if j < len(key)-1 {
		s := buildSetter(key[j+1:], value, true)
		if s == nil {
			if inner {
				return nil
			}
			return &objectSetter{key, value}
		}
		value = s
	}
	if i+1 == j {
		return &arraySetter{value}
	}
	return &objectSetter{key[i+1 : j], value}
}
