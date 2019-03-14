package gojo

import (
	"encoding/json"
	"strings"

	"github.com/iancoleman/orderedmap"
)

func parseKeyValue(s string) (map[string]interface{}, error) {
	i := strings.IndexRune(s, '=')
	if i < 0 {
		return nil, errParse(s)
	}
	key, value := s[:i], parseValue(s[i+1:])
	return map[string]interface{}{key: value}, nil
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
	return s
}
