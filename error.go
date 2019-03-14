package gojo

import "fmt"

type errParse string

func (err errParse) Error() string {
	return fmt.Sprintf("failed to parse: %s", string(err))
}
