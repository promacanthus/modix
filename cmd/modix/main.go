package main

import (
	"os"

	"github.com/promacanthus/modix/cmd/modix/commands"
	_ "github.com/promacanthus/modix/cmd/modix/commands/project"
	"github.com/promacanthus/modix/internal/tui"
)

func main() {
	// Check if no arguments were provided
	if len(os.Args) == 1 {
		// Launch TUI by default
		if err := tui.Run(); err != nil {
			os.Exit(1)
		}
		return
	}

	// Otherwise, execute the CLI command
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
