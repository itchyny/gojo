package main

import (
	"fmt"
	"os"

	"github.com/itchyny/gojo"
)

const name = "gojo"

func main() {
	if err := gojo.New(os.Args[1:], os.Stdout).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
	}
}
