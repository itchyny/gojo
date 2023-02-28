// gojo - Go implementation of jo
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/itchyny/gojo"
	"github.com/itchyny/json2yaml"
)

const name = "gojo"

const version = "0.3.1"

var revision = "HEAD"

func main() {
	os.Exit((&cli{
		inStream:  os.Stdin,
		outStream: os.Stdout,
		errStream: os.Stderr,
	}).run(os.Args[1:]))
}

const (
	exitCodeOK = iota
	exitCodeErr
)

type cli struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

func (cli *cli) run(args []string) (exitCode int) {
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
	var yaml bool
	fs.BoolVar(&yaml, "y", false, "YAML format output")
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
	var ret any
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
	var enc *json.Encoder
	if yaml {
		var buf bytes.Buffer
		enc = json.NewEncoder(&buf)
		defer func() {
			if exitCode == exitCodeOK {
				if err = json2yaml.Convert(cli.outStream, &buf); err != nil {
					fmt.Fprintf(cli.errStream, "%s: %s\n", name, err)
					exitCode = exitCodeErr
				}
			}
		}()
	} else {
		enc = json.NewEncoder(cli.outStream)
		if pretty {
			enc.SetIndent("", "  ")
		}
	}
	if err = enc.Encode(ret); err != nil {
		fmt.Fprintf(cli.errStream, "%s: %s\n", name, err)
		return exitCodeErr
	}
	return exitCodeOK
}
