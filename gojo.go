package gojo

import (
	"encoding/json"
	"os"
)

// Gojo ...
type Gojo struct {
	args []string
}

// New Gojo
func New(args []string) *Gojo {
	return &Gojo{args: args}
}

// Run ...
func (g *Gojo) Run() error {
	ms := make(map[string]string, len(g.args))
	for _, arg := range g.args {
		m, err := parse(arg)
		if err != nil {
			return err
		}
		for k, v := range m {
			ms[k] = v
		}
	}
	return json.NewEncoder(os.Stdout).Encode(ms)
}
