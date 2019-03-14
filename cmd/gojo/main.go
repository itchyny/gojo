package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/itchyny/gojo"
)

const name = "gojo"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
	}
}

func run() error {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	var pretty bool
	fs.BoolVar(&pretty, "p", false, "pretty print")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	opts := []gojo.Option{
		gojo.Args(fs.Args()),
		gojo.OutStream(os.Stdout),
	}
	if pretty {
		opts = append(opts, gojo.Pretty())
	}
	return gojo.New(opts...).Run()
}
