package tui

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/promacanthus/modix/internal/tui/models"
)

// Run starts the TUI application
func Run() error {
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
