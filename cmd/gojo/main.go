package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/itchyny/gojo"
)

const name = "gojo"

const version = "0.0.3"

var revision = "HEAD"

func main() {
	if err := run(); err != nil && err != flag.ErrHelp {
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.Usage = func() {
		fs.SetOutput(os.Stdout)
		fmt.Printf(`gojo - Yet another Go implementation of jo

Version: %s (rev: %s/%s)

Synopsis:
    %% gojo key=value ...

Options:
`, version, revision, runtime.Version())
		fs.PrintDefaults()
	}
	var showVersion bool
	fs.BoolVar(&showVersion, "v", false, "print version")
	var array bool
	fs.BoolVar(&array, "a", false, "creates an array")
	var pretty bool
	fs.BoolVar(&pretty, "p", false, "pretty print")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return nil
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
		gojo.Output(os.Stdout),
	}
	if array {
		opts = append(opts, gojo.Array())
	}
	if pretty {
		opts = append(opts, gojo.Pretty())
	}
	if err := gojo.New(opts...).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
		return err
	}
	return nil
}
