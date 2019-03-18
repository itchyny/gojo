package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"runtime"

	"github.com/itchyny/gojo"
)

const name = "gojo"

const version = "0.0.3"

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
	opts := []gojo.Option{
		gojo.Args(args),
		gojo.Output(cli.outStream),
	}
	if array {
		opts = append(opts, gojo.Array())
	}
	if pretty {
		opts = append(opts, gojo.Pretty())
	}
	if err := gojo.New(opts...).Run(); err != nil {
		fmt.Fprintf(cli.errStream, "%s: %s\n", name, err)
		return exitCodeErr
	}
	return exitCodeOK
}
