package gojo

import (
	"errors"
	"fmt"
)

type errParse string

func (err errParse) Error() string {
	return fmt.Sprintf("failed to parse: %q", string(err))
}

var errOutputNotSet = errors.New("output is not set")
