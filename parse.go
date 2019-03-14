package gojo

import (
	"encoding/json"
	"math/big"
	"strconv"
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
	bi := new(big.Int)
	if err := json.Unmarshal([]byte(s), bi); err == nil {
		return bi
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
