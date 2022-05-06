package gojo

import "github.com/iancoleman/orderedmap"

type setter interface {
	set([]string, any) error
}

type arraySetter struct {
	value any
}

func (s *arraySetter) set(keys []string, v any) error {
	arr, ok := v.(*[]any)
	if !ok {
		return errArray{keys, v}
	}
	if st, ok := s.value.(setter); ok {
		if ast, ok := st.(*arraySetter); ok {
			ar := []any{}
			if err := ast.set(keys, &ar); err != nil {
				return err
			}
			*arr = append(*arr, ar)
		} else {
			om := orderedmap.New()
			if err := st.set(keys, om); err != nil {
				return err
			}
			*arr = append(*arr, om)
		}
	} else {
		*arr = append(*arr, s.value)
	}
	return nil
}

type objectSetter struct {
	key   string
	value any
}

func (s *objectSetter) set(keys []string, v any) error {
	if !isMap(v) {
		return errObject{keys, v}
	}
	if st, ok := s.value.(setter); ok {
		val, ok := getKey(v, s.key)
		switch v := val.(type) {
		case orderedmap.OrderedMap:
			val = &v
		case []any:
			val = &v
		}
		if !ok {
			if _, ok := st.(*arraySetter); ok {
				val = &[]any{}
			} else {
				val = orderedmap.New()
			}
		}
		if err := st.set(append(keys, s.key), val); err != nil {
			return err
		}
		setKey(v, s.key, val)
	} else {
		setKey(v, s.key, s.value)
	}
	return nil
}

func isMap(t any) bool {
	switch t.(type) {
	case *orderedmap.OrderedMap:
		return true
	case map[string]any:
		return true
	default:
		return false
	}
}

func getKey(t any, key string) (any, bool) {
	switch t := t.(type) {
	case *orderedmap.OrderedMap:
		return t.Get(key)
	case map[string]any:
		v, ok := t[key]
		return v, ok
	default:
		panic(t)
	}
}

func setKey(t any, key string, value any) {
	switch t := t.(type) {
	case *orderedmap.OrderedMap:
		t.Set(key, value)
	case map[string]any:
		t[key] = value
	default:
		panic(t)
	}
}
