package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-integration-installer/api"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-integration-installer/config"
)

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

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to load configuration: %v\033[0m\n", err)
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Parse command line arguments
	licensePath, err := parseLicenseFlag()
	if err != nil {
		fmt.Printf("\033[31m❌ %v\033[0m\n", err)
		return err
	}

	// Read license key from file
	licenseKey, err := readLicenseKey(licensePath)
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to read license key: %v\033[0m\n", err)
		return fmt.Errorf("failed to read license key: %w", err)
	}

	// Initialize API client
	client := api.NewClient(cfg.APIURL, licenseKey)

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
	selectedVersion, err := interactiveSelection(versions, "\033[32mSelect version\033[0m", &versions[0])
	if err != nil {
		fmt.Printf("\033[31m❌ Version selection failed: %v\033[0m\n", err)
		return err
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

	outputPath := fmt.Sprintf("go-jo-%s.zip", selectedIntegration)

	err = client.DownloadPackage(selectedVersion, selectedIntegration, outputPath)
	if err != nil {
		fmt.Printf("\033[31m❌ Failed to download package: %v\033[0m\n", err)
		return fmt.Errorf("failed to download package: %w", err)
	}

	fmt.Printf("\n\033[32m✅ Package downloaded successfully: %s\033[0m\n", outputPath)
	fmt.Printf("\033[36m📁 Location: %s\033[0m\n", filepath.Join(".", outputPath))

	return nil
}

// parseLicenseFlag parses the --license flag from command line arguments
func parseLicenseFlag() (string, error) {
	if len(os.Args) < 1 {
		return "", fmt.Errorf("usage: %s --license=<path-to-license-file>", os.Args[0])
	}

	for _, arg := range os.Args {
		// Check for --license=<path> format
		if strings.HasPrefix(arg, "--license=") {
			licensePath := strings.TrimPrefix(arg, "--license=")
			if licensePath == "" {
				return "", fmt.Errorf("--license flag requires a file path")
			}
			return licensePath, nil
		}
	}

	// Check for --license <path> format (separate arguments)
	for i, arg := range os.Args {
		if arg == "--license" {
			if i+1 >= len(os.Args) {
				return "", fmt.Errorf("--license flag requires a file path")
			}
			return os.Args[i+1], nil
		}
	}

	return "", fmt.Errorf("--license flag is required")
}

// readLicenseKey reads the license key from the specified file
func readLicenseKey(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read license file: %w", err)
	}

	licenseKey := strings.TrimSpace(string(content))
	if licenseKey == "" {
		return "", fmt.Errorf("license file is empty")
	}

	return licenseKey, nil
}

// interactiveSelection provides a robust interactive selection using bubbletea
func interactiveSelection(options []string, prompt string, latest *string) (string, error) {
	if len(options) == 0 {
		return "", fmt.Errorf("no options available")
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
