package main

import (
	"os"

	"github.com/promacanthus/modix/cmd/modix/commands"
	_ "github.com/promacanthus/modix/cmd/modix/commands/project"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
