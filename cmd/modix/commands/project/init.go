package project

import (
	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/project"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new modix project",
	Long: `Initialize a new modix project by creating the .modix/ directory and configuration files.

This command creates the following files:
  - .modix/shells.json
  - .modix/brains.json
  - .modix/agents.json
  - .modix/runtimes.json
  - .modix/projects.json
  - .modix/state.json
  - .modix/version.json

If the project is already initialized, this command will fail.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		return project.Init(format)
	},
}

func init() {
	ProjectCmd.AddCommand(initCmd)
}
