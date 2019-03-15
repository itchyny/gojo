package main

import (
	"os"

	"github.com/itchyny/gojo/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		os.Exit(1)
	}
}
