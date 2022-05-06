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

func parseValue(s string) (any, error) {
	if s == "" {
		return "", nil
	}
	switch s[0] {
	case '@':
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
	case ':':
		cnt, err := ioutil.ReadFile(s[1:])
		if err != nil {
			return nil, err
		}
		return decodeJSON(cnt)
	case '%':
		cnt, err := ioutil.ReadFile(s[1:])
		if err != nil {
			return nil, err
		}
		return base64.StdEncoding.EncodeToString(cnt), nil
	default:
		if v, err := decodeJSON([]byte(s)); err == nil {
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
		return s, nil
	}
}

func decodeJSON(bs []byte) (any, error) {
	var m orderedmap.OrderedMap
	if err := json.Unmarshal(bs, &m); err == nil {
		return m, nil
	}
	var vs []orderedmap.OrderedMap
	if err := json.Unmarshal(bs, &vs); err == nil {
		if vs == nil {
			return nil, nil
		}
		ws := make([]any, len(vs))
		for i, v := range vs {
			ws[i] = v
		}
		return ws, nil
	}
	var v any
	if err := json.Unmarshal(bs, &v); err != nil {
		return nil, err
	}
	return v, nil
}

func buildSetter(key string, value any, inner bool) setter {
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
