package main

import (
	"os"

	"github.com/kish1n/usdt_listening/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
