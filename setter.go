package gojo

import "github.com/iancoleman/orderedmap"

type setter interface {
	set(*orderedmap.OrderedMap) error
}

type keyValueSetter struct {
	key   string
	value interface{}
}

func (s *keyValueSetter) set(om *orderedmap.OrderedMap) error {
	om.Set(s.key, s.value)
	return nil
}

type arraySetter struct {
	key   string
	value interface{}
}

func (s *arraySetter) set(om *orderedmap.OrderedMap) error {
	v, ok := om.Get(s.key)
	if !ok {
		v = []interface{}{}
	}
	arr, ok := v.([]interface{})
	if !ok {
		return errArray{s.key, v}
	}
	arr = append(arr, s.value)
	om.Set(s.key, arr)
	return nil
}

type objectSetter struct {
	key   string
	inner string
	value interface{}
}

func (s *objectSetter) set(om *orderedmap.OrderedMap) error {
	v, ok := om.Get(s.key)
	if !ok {
		v = map[string]interface{}{}
	}
	obj, ok := v.(map[string]interface{})
	if !ok {
		return errObject{s.key, v}
	}
	obj[s.inner] = s.value
	om.Set(s.key, obj)
	return nil
}
