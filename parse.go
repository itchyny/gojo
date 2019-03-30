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
	return buildSetter(key, value, false), nil
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

func buildSetter(key string, value interface{}, inner bool) setter {
	i, j := strings.IndexRune(key, '['), strings.IndexRune(key, ']')
	if i < 0 || j < 0 || j < i {
		if inner {
			return nil
		}
		return &objectSetter{key, value}
	}
	if i > 0 {
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
