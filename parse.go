package gojo

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func parseKeyValue(s string) (setter, error) {
	i := strings.IndexRune(s, '=')
	if i < 0 {
		return nil, errParse(s)
	}
	key, value := s[:i], parseValue(s[i+1:])
	i, j := strings.IndexRune(key, '['), strings.IndexRune(key, ']')
	if i < 0 || j < 0 || j < i || j < len(key)-1 {
		return &keyValueSetter{key, value}, nil
	}
	if i+1 == j && j == len(key)-1 {
		return &arraySetter{key[:i], value}, nil
	}
	return &objectSetter{key[:i], key[i+1 : j], value}, nil
}

func parseValue(s string) interface{} {
	m := orderedmap.New()
	if err := json.Unmarshal([]byte(s), m); err == nil {
		return m
	}
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err == nil {
		return v
	}
	if i, err := strconv.ParseInt(s, 10, 64); err == nil {
		return i
	}
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		if i, err := strconv.ParseInt(s[2:], 16, 64); err == nil {
			return i
		}
	}
	return s
}
