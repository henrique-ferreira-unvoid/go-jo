package main

import (
	"fmt"
	"os"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-integration-installer/cli"
)

func main() {
	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
