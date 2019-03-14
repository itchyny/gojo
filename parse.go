package gojo

import (
	"encoding/json"
	"strings"
)

func parse(s string) (map[string]interface{}, error) {
	i := strings.IndexRune(s, '=')
	if i < 0 {
		return nil, errParse(s)
	}
	key, value := s[:i], parseValue(s[i+1:])
	return map[string]interface{}{key: value}, nil
}

func parseValue(s string) interface{} {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err == nil {
		return v
	}
	return s
}
