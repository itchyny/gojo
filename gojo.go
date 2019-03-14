package gojo

import (
	"encoding/json"
	"io"
)

// Gojo ...
type Gojo struct {
	args      []string
	array     bool
	pretty    bool
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
	if g.array {
		return g.runArr()
	}
	return g.runObj()
}

func (g *Gojo) runObj() error {
	ms := make(map[string]interface{}, len(g.args))
	for _, arg := range g.args {
		m, err := parse(arg)
		if err != nil {
			return err
		}
		for k, v := range m {
			ms[k] = v
		}
	}
	enc := json.NewEncoder(g.outStream)
	if g.pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(ms)
}

func (g *Gojo) runArr() error {
	as := make([]interface{}, 0, len(g.args))
	for _, arg := range g.args {
		as = append(as, parseValue(arg))
	}
	enc := json.NewEncoder(g.outStream)
	if g.pretty {
		enc.SetIndent("", "  ")
	}
	return enc.Encode(as)
}
