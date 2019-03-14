package gojo

import "io"

// Option for Gojo
type Option func(*Gojo)

// Args ...
func Args(args []string) Option {
	return func(g *Gojo) {
		g.args = args
	}
}

// OutStream ...
func OutStream(outStream io.Writer) Option {
	return func(g *Gojo) {
		g.outStream = outStream
	}
}
