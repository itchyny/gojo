package main

import (
	"fmt"
	"os"

	"github.com/itchyny/gojo"
)

const name = "gojo"

func main() {
	if err := gojo.New(
		gojo.Args(os.Args[1:]),
		gojo.OutStream(os.Stdout),
	).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
	}
}
