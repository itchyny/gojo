package main

import (
	"bufio"
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
	var array bool
	fs.BoolVar(&array, "a", false, "creates an array")
	var pretty bool
	fs.BoolVar(&pretty, "p", false, "pretty print")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	args := fs.Args()
	if len(args) == 0 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			args = append(args, s.Text())
		}
		if err := s.Err(); err != nil {
			return err
		}
	}
	opts := []gojo.Option{
		gojo.Args(args),
		gojo.OutStream(os.Stdout),
	}
	if array {
		opts = append(opts, gojo.Array())
	}
	if pretty {
		opts = append(opts, gojo.Pretty())
	}
	return gojo.New(opts...).Run()
}
