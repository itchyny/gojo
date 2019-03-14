package gojo

import "strings"

func parse(s string) (map[string]string, error) {
	i := strings.IndexRune(s, '=')
	if i < 0 {
		return nil, errParse(s)
	}
	key, value := s[:i], s[i+1:]
	return map[string]string{key: value}, nil
}
