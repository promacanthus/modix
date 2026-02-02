package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/promacanthus/modix/internal/tui/models"
)

// Run starts the TUI application
func Run() error {
	// Check if we're in a TTY environment
	if !isTTY() {
		fmt.Println("Error: TUI requires a TTY (terminal) environment.")
		fmt.Println("Please run this command in a terminal, not in a script or IDE.")
		fmt.Println("\nAlternatively, use the CLI commands directly:")
		fmt.Println("  mx project init <name>")
		fmt.Println("  mx project list")
		fmt.Println("  mx shell register <name>")
		fmt.Println("  mx agent define <role>")
		fmt.Println("  mx --help")
		os.Exit(1)
	}

	p := tea.NewProgram(models.NewMainModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running TUI: %v\n", err)
		os.Exit(1)
	}
	return nil
}

// RunWithArgs starts the TUI application with command-line arguments
func RunWithArgs(args []string) error {
	// Parse arguments if needed
	// For now, just run the TUI
	return Run()
}

// isTTY checks if the current environment is a TTY (terminal)
func isTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
