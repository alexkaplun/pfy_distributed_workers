package main

import (
	"github.com/alexkaplun/pfy_distributed_workers/cli"
	"os"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
