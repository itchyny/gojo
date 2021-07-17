package gojo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type errParse string

func (err errParse) Error() string {
	return fmt.Sprintf("failed to parse: %q", string(err))
}

type errArray struct {
	keys  []string
	value interface{}
}

func (err errArray) Error() string {
	bs, _ := json.Marshal(err.value)
	return fmt.Sprintf("expected an array: .%s: %s", strings.Join(err.keys, "."), string(bs))
}

type errObject struct {
	keys  []string
	value interface{}
}

func (err errObject) Error() string {
	bs, _ := json.Marshal(err.value)
	return fmt.Sprintf("expected an object: .%s: %s", strings.Join(err.keys, "."), string(bs))
}
