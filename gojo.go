package gojo

import (
	"encoding/json"
	"io"
	"os"

	"github.com/iancoleman/orderedmap"
)

// Gojo represents the gojo printer
type Gojo struct {
	args   []string
	array  bool
	pretty bool
	output io.Writer
}

// New Gojo
func New(opts ...Option) *Gojo {
	g := &Gojo{output: os.Stdout}
	for _, opt := range opts {
		opt(g)
	}
	return g
}

// Run gojo
func (g *Gojo) Run() error {
	if g.array {
		return g.runArr()
	}
	return g.runObj()
}

const indent = "  "

func (g *Gojo) runObj() error {
	ms := orderedmap.New()
	for _, arg := range g.args {
		s, err := parseKeyValue(arg)
		if err != nil {
			return err
		}
		if err := s.set(nil, ms); err != nil {
			return err
		}
	}
	enc := json.NewEncoder(g.output)
	if g.pretty {
		enc.SetIndent("", indent)
	}
	return enc.Encode(ms)
}

func (g *Gojo) runArr() error {
	xs := make([]interface{}, len(g.args))
	for i, arg := range g.args {
		v, err := parseValue(arg)
		if err != nil {
			return err
		}
		xs[i] = v
	}
	enc := json.NewEncoder(g.output)
	if g.pretty {
		enc.SetIndent("", indent)
	}
	return enc.Encode(xs)
}
