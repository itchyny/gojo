package gojo

import "github.com/iancoleman/orderedmap"

// Map builds a new ordered map from arguments.
func Map(args []string) (*orderedmap.OrderedMap, error) {
	ms := orderedmap.New()
	for _, arg := range args {
		s, err := parseKeyValue(arg)
		if err != nil {
			return nil, err
		}
		if err := s.set(nil, ms); err != nil {
			return nil, err
		}
	}
	return ms, nil
}

// Array builds a new slice from arguments.
func Array(args []string) ([]interface{}, error) {
	xs := make([]interface{}, len(args))
	for i, arg := range args {
		v, err := parseValue(arg)
		if err != nil {
			return nil, err
		}
		xs[i] = v
	}
	return xs, nil
}
