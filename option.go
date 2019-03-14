package gojo

import "io"

// Option for Gojo
type Option func(*Gojo)

// Args option
func Args(args []string) Option {
	return func(g *Gojo) {
		g.args = args
	}
}

// Array option
func Array() Option {
	return func(g *Gojo) {
		g.array = true
	}
}

// Pretty option
func Pretty() Option {
	return func(g *Gojo) {
		g.pretty = true
	}
}

// Output option
func Output(output io.Writer) Option {
	return func(g *Gojo) {
		g.output = output
	}
}
