package main

import (
	"fmt"
)

var (
	// Version will be set by GoReleaser at build time
	Version = "dev"
	// GitCommit will be set by GoReleaser at build time
	GitCommit = "unknown"
	// BuildDate will be set by GoReleaser at build time
	BuildDate = "unknown"
)

func main() {
	fmt.Printf("go-jo version %s\n", Version)
	fmt.Printf("Git commit: %s\n", GitCommit)
	fmt.Printf("Build date: %s\n", BuildDate)
}
