package gojo

import (
	"encoding/json"
	"io"
)

// Gojo ...
type Gojo struct {
	args      []string
	outStream io.Writer
}

// New Gojo
func New(opts ...Option) *Gojo {
	g := &Gojo{}
	for _, opt := range opts {
		g.Apply(opt)
	}
	return g
}

// Apply an Option
func (g *Gojo) Apply(opt Option) {
	opt(g)
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
	return json.NewEncoder(g.outStream).Encode(ms)
}
