package gojo

import "github.com/iancoleman/orderedmap"

type setter interface {
	set(interface{}) error
}

type arraySetter struct {
	value interface{}
}

func (s *arraySetter) set(v interface{}) error {
	arr, ok := v.(*[]interface{})
	if !ok {
		return errArray{"", v}
	}
	if st, ok := s.value.(setter); ok {
		if ast, ok := st.(*arraySetter); ok {
			ar := []interface{}{}
			if err := ast.set(&ar); err != nil {
				return err
			}
			*arr = append(*arr, ar)
		} else {
			om := orderedmap.New()
			if err := st.set(om); err != nil {
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
	value interface{}
}

func (s *objectSetter) set(v interface{}) error {
	if !isMap(v) {
		return errObject{"", v}
	}
	if st, ok := s.value.(setter); ok {
		val, ok := getKey(v, s.key)
		if !ok {
			if _, ok := st.(*arraySetter); ok {
				val = &[]interface{}{}
			} else {
				val = orderedmap.New()
			}
		}
		if err := st.set(val); err != nil {
			return err
		}
		setKey(v, s.key, val)
	} else {
		setKey(v, s.key, s.value)
	}
	return nil
}

func isMap(t interface{}) bool {
	switch t.(type) {
	case *orderedmap.OrderedMap:
		return true
	case map[string]interface{}:
		return true
	default:
		return false
	}
}

func getKey(t interface{}, key string) (interface{}, bool) {
	switch t := t.(type) {
	case *orderedmap.OrderedMap:
		return t.Get(key)
	case map[string]interface{}:
		v, ok := t[key]
		return v, ok
	default:
		panic(t)
	}
}

func setKey(t interface{}, key string, value interface{}) {
	switch t := t.(type) {
	case *orderedmap.OrderedMap:
		t.Set(key, value)
	case map[string]interface{}:
		t[key] = value
	default:
		panic(t)
	}
}
