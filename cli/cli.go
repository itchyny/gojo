package cli

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"runtime"

	"github.com/itchyny/gojo"
)

const name = "gojo"

const version = "0.2.0"

var revision = "HEAD"

const (
	exitCodeOK = iota
	exitCodeErr
)

type cli struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

func (cli *cli) run(args []string) int {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(cli.errStream)
	fs.Usage = func() {
		fs.SetOutput(cli.outStream)
		fmt.Fprintf(cli.outStream, `%[1]s - Go implementation of jo

Version: %s (rev: %s/%s)

Synopsis:
  %% %[1]s key=value ...

Options:
`, name, version, revision, runtime.Version())
		fs.PrintDefaults()
	}
	var showVersion bool
	fs.BoolVar(&showVersion, "v", false, "print version")
	var array bool
	fs.BoolVar(&array, "a", false, "creates an array")
	var pretty bool
	fs.BoolVar(&pretty, "p", false, "pretty print")
	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return exitCodeOK
		}
		return exitCodeErr
	}
	if showVersion {
		fmt.Fprintf(cli.outStream, "%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return exitCodeOK
	}
	args = fs.Args()
	if len(args) == 0 {
		s := bufio.NewScanner(cli.inStream)
		for s.Scan() {
			args = append(args, s.Text())
		}
		if err := s.Err(); err != nil {
			return exitCodeErr
		}
	}
	var ret interface{}
	var err error
	if array {
		ret, err = gojo.Array(args)
	} else {
		ret, err = gojo.Map(args)
	}
	if err != nil {
		fmt.Fprintf(cli.errStream, "%s: %s\n", name, err)
		return exitCodeErr
	}
	enc := json.NewEncoder(cli.outStream)
	if pretty {
		enc.SetIndent("", "  ")
	}
	if err = enc.Encode(ret); err != nil {
		fmt.Fprintf(cli.errStream, "%s: %s\n", name, err)
		return exitCodeErr
	}
	return exitCodeOK
}
