package gojo

import "fmt"

type errParse string

func (err errParse) Error() string {
	return fmt.Sprintf("failed to parse: %q", string(err))
}

type errArray struct {
	key   string
	value interface{}
}

func (err errArray) Error() string {
	return fmt.Sprintf("expected an array: %q (%T)", err.key, err.value)
}

type errObject struct {
	key   string
	value interface{}
}

func (err errObject) Error() string {
	return fmt.Sprintf("expected an object: %q (%T)", err.key, err.value)
}
