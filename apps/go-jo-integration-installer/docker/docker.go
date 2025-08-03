package docker

import (
	"fmt"
	"os"
	"os/exec"
)

// CheckDocker verifies that Docker is running
func CheckDocker() error {
	cmd := exec.Command("docker", "version")
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}

// RunMakeCommands runs make build and make start with live output
func RunMakeCommands(deployDir string) error {
	// Change to the deployment directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	if err := os.Chdir(deployDir); err != nil {
		return fmt.Errorf("failed to change to deployment directory: %w", err)
	}
	defer os.Chdir(originalDir) // Restore original directory

	// Run make build
	fmt.Printf("\033[33m🔨 Running 'make build'...\033[0m\n")
	if err := runCommandWithOutput("make", "build"); err != nil {
		return fmt.Errorf("make build failed: %w", err)
	}

	// Run make start
	fmt.Printf("\033[33m🚀 Running 'make start'...\033[0m\n")
	if err := runCommandWithOutput("make", "start"); err != nil {
		return fmt.Errorf("make start failed: %w", err)
	}

	return nil
}

// runCommandWithOutput runs a command and displays its output in real-time
func runCommandWithOutput(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
