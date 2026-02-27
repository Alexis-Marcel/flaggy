package main

import (
	"os"

	"github.com/alexis/flaggy/cmd/flaggy/cli"
)

var version = "dev"

func main() {
	cli.Version = version
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
