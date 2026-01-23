package project

import (
	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/project"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check dependencies for modix project",
	Long: `Check if required command-line tools are installed and available.

This command checks for the following tools:
  - git
  - claude-code
  - codex-cli
  - gemini-cli

It verifies that each tool is available in the system PATH.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		return project.Check(format)
	},
}

func init() {
	ProjectCmd.AddCommand(checkCmd)
}
