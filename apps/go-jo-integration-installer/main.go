package main

import (
	"fmt"
	"os"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-integration-installer/cli"
)

func main() {
	// Check if running as root/sudo
	if os.Geteuid() == 0 {
		fmt.Fprintf(os.Stderr, "Error: This application should not be run as root or with sudo\n")
		os.Exit(1)
	}

	if err := cli.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
