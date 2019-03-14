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

// Array ...
func Array() Option {
	return func(g *Gojo) {
		g.array = true
	}
}

// Pretty ...
func Pretty() Option {
	return func(g *Gojo) {
		g.pretty = true
	}
}

// OutStream ...
func OutStream(outStream io.Writer) Option {
	return func(g *Gojo) {
		g.outStream = outStream
	}
}
