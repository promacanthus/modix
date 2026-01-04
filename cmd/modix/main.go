package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/promacanthus/modix/cmd/modix/commands"
)

var (
	// Force colored output for all terminals
	bold   = color.New(color.FgCyan, color.Bold)
	blue   = color.New(color.FgBlue)
	green  = color.New(color.FgGreen, color.Bold)
	yellow = color.New(color.FgYellow, color.Bold)
	red    = color.New(color.FgRed, color.Bold)
)

func init() {
	// Force colored output for all terminals
	color.NoColor = false
}

func main() {
	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}