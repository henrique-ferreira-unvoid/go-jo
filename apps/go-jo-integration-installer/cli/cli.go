package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"apps/go-jo-integration-installer/api"
	"apps/go-jo-integration-installer/config"
	"apps/go-jo-integration-installer/docker"
	"apps/go-jo-integration-installer/utils"

	tea "github.com/charmbracelet/bubbletea"
)

const MAX_OPTIONS = 15

// selectionModel represents the state of the selection interface
type selectionModel struct {
	items    []string
	cursor   int
	selected string
	done     bool
	title    string
	latest   *string
}

// initialSelectionModel creates a new selection model
func initialSelectionModel(items []string, title string, latest *string) selectionModel {
	return selectionModel{
		items:  items,
		cursor: 0,
		title:  title,
		latest: latest,
	}
}

// Init initializes the model
func (m selectionModel) Init() tea.Cmd {
	return nil
}

// Update handles user input
func (m selectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected = m.items[m.cursor]
			m.done = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the selection interface
func (m selectionModel) View() string {
	if m.done {
		return ""
	}

	var s = fmt.Sprintf("\n%s\n\n", m.title)

	for i, choice := range m.items {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = "\033[36m>\033[0m" // cursor with cyan color!
		}

		// Is this choice selected?
		checked := " " // not selected
		if m.cursor == i {
			checked = "\033[32mx\033[0m" // selected with green color!
		}

		latest := " " // not latest
		if m.latest != nil && *m.latest == choice {
			latest = "\033[33m(latest)\033[0m" // latest with yellow color!
		}

		// Render the row with colors
		choiceColor := ""
		resetColor := ""
		if m.cursor == i {
			choiceColor = "\033[1;37m" // bold white for selected item
			resetColor = "\033[0m"
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s%s%s %s\n", cursor, checked, choiceColor, choice, resetColor, latest)
	}

	s += "\n(press q to quit)\n"

	return s
}

// Run executes the main CLI logic
func Run() error {
	// Display ASCII art banner
	fmt.Printf("\033[36m")
	fmt.Println(`

  ██████╗  ██████╗            ██╗ ██████╗ 
 ██╔════╝ ██╔═══██╗           ██║██╔═══██╗
 ██║  ███╗██║   ██║█████╗     ██║██║   ██║
 ██║   ██║██║   ██║╚════╝██   ██║██║   ██║
 ╚██████╔╝╚██████╔╝      ╚█████╔╝╚██████╔╝
 ╚═════╝  ╚═════╝        ╚════╝  ╚═════`)

	// Check if Docker is running
	fmt.Printf("\033[36m🔍 Checking Docker status...\033[0m\n")
	if err := docker.CheckDocker(); err != nil {
		fmt.Printf("\033[31m❌ Docker is not running: %v\033[0m\n", err)
		return fmt.Errorf("Docker is required but not running: %w", err)
	}
	fmt.Printf("\033[32m✅ Docker is running\033[0m\n")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to load configuration: %v\033[0m\n", err)
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Parse command line arguments
	licensePath, err := utils.ParseLicenseFlag()
	if err != nil {
		fmt.Printf("\033[31m❌ %v\033[0m\n", err)
		return err
	}

	requestedVersion, err := utils.ParseVersionFlag()
	if err != nil {
		fmt.Printf("\033[31m❌ %v\033[0m\n", err)
		return err
	}

	apiURL, err := utils.ParseAPIURLFlag()
	if err != nil {
		fmt.Printf("\033[31m❌ %v\033[0m\n", err)
		return err
	}

	if apiURL != "" {
		cfg.APIURL = apiURL
	}

	// Read license key from file
	licenseKey, err := utils.ReadLicenseKey(licensePath)
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to read license key: %v\033[0m\n", err)
		return fmt.Errorf("failed to read license key: %w", err)
	}

	// Initialize API client
	client := api.NewClient(cfg.APIURL, licenseKey)

	selectedVersion := requestedVersion
	if selectedVersion == "" {
		// Get available versions
		fmt.Printf("\033[36m🔍 Fetching available versions...\033[0m\n")
		versions, err := client.GetVersions()
		if err != nil {
			fmt.Printf("\033[31m❌ Failed to fetch versions: %v\033[0m\n", err)
			return fmt.Errorf("failed to fetch versions: %w", err)
		}

		if len(versions) == 0 {
			fmt.Printf("\033[31m❌ No versions available\033[0m\n")
			return fmt.Errorf("no versions available")
		}

		// Display versions and get user selection
		fmt.Printf("\033[33m📦 Available versions:\033[0m\n")
		selectedVersion, err = interactiveSelection(versions, "\033[32mSelect version\033[0m", &versions[0])
		if err != nil {
			fmt.Printf("\033[31m❌ Version selection failed: %v\033[0m\n", err)
			return err
		}
	} else {
		fmt.Printf("\033[36m🎯 Using provided version tag: %s\033[0m\n", selectedVersion)
	}

	fmt.Printf("\033[32m✅ Selected version: %s\033[0m\n", selectedVersion)

	// Get available integrations
	fmt.Printf("\033[36m🔍 Fetching available integrations...\033[0m\n")
	integrations, err := client.GetIntegrations()
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to fetch integrations: %v\033[0m\n", err)
		return fmt.Errorf("failed to fetch integrations: %w", err)
	}

	if len(integrations) == 0 {
		fmt.Printf("\033[31m❌ No integrations available\033[0m\n")
		return fmt.Errorf("no integrations available")
	}

	// Display integrations and get user selection
	fmt.Printf("\033[33m🔌 Available integrations:\033[0m\n")
	selectedIntegration, err := interactiveSelection(integrations, "\033[32mSelect integration\033[0m", nil)
	if err != nil {
		fmt.Printf("\033[31m❌ Integration selection failed: %v\033[0m\n", err)
		return err
	}

	fmt.Printf("\033[32m✅ Selected integration: %s\033[0m\n", selectedIntegration)

	// Download the package
	fmt.Printf("\033[35m⬇️  Downloading package for version %s with integration %s...\033[0m\n",
		selectedVersion, selectedIntegration)

	safeSelectedIntegration := strings.ReplaceAll(selectedIntegration, "/", "@")

	outputPath := fmt.Sprintf("go-jo-%s.zip", safeSelectedIntegration)

	err = client.DownloadPackage(selectedVersion, safeSelectedIntegration, outputPath)
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to download package: %v\033[0m\n", err)
		return fmt.Errorf("failed to download package: %w", err)
	}

	fmt.Printf("\n\033[32m✅ Package downloaded successfully: %s\033[0m\n", outputPath)
	fmt.Printf("\033[36m📁 Location: %s\033[0m\n", filepath.Join(".", outputPath))

	// Extract and deploy the package
	fmt.Printf("\033[35m🚀 Extracting and deploying package...\033[0m\n")
	if err := extractAndDeploy(outputPath); err != nil {
		fmt.Printf("\033[31m❌ Failed to deploy package: %v\033[0m\n", err)
		return fmt.Errorf("failed to deploy package: %w", err)
	}

	fmt.Printf("\n\033[32m🎉 Deployment completed successfully!\033[0m\n")
	fmt.Printf("\033[36m📋 The go-jo application is now running in Docker containers.\033[0m\n")

	return nil
}

// extractAndDeploy extracts the zip file and runs the deployment
func extractAndDeploy(zipPath string) error {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "go-jo-deploy-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory on exit
	defer os.Remove(zipPath)    // Clean up zip file on exit

	fmt.Printf("\033[36m📂 Extracting to: %s\033[0m\n", tempDir)

	// Extract the zip file
	if err := utils.ExtractZip(zipPath, tempDir); err != nil {
		return fmt.Errorf("failed to extract zip file: %w", err)
	}

	// Find the extracted directory (should contain docker-compose.yml and Makefile)
	deployDir, err := utils.FindDeployDirectory(tempDir)
	if err != nil {
		return fmt.Errorf("failed to find deployment directory: %w", err)
	}

	fmt.Printf("\033[36m🔧 Deploying from: %s\033[0m\n", deployDir)

	// Run make build and make start
	if err := docker.RunMakeCommands(deployDir); err != nil {
		return fmt.Errorf("failed to run make commands: %w", err)
	}

	return nil
}

// interactiveSelection provides a robust interactive selection using bubbletea
func interactiveSelection(options []string, prompt string, latest *string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options available")
	}

	if len(options) > MAX_OPTIONS {
		fmt.Printf("\033[33m🔍 Showing first %d options...\033[0m\n", MAX_OPTIONS)
		options = options[:MAX_OPTIONS]
	}

	m := initialSelectionModel(options, prompt, latest)
	p := tea.NewProgram(m)

	// Run the program
	finalModel, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("selection failed: %w", err)
	}

	// Get the final model
	final, ok := finalModel.(selectionModel)
	if !ok {
		return "", fmt.Errorf("unexpected model type")
	}

	if final.selected == "" {
		return "", fmt.Errorf("selection cancelled")
	}

	return final.selected, nil
}
