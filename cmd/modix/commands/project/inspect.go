package project

import (
	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/project"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect modix project configuration",
	Long: `Inspect and display the modix project configuration.

This command shows:
  - Project information
  - Shell configuration
  - Brain configuration
  - Agent configuration
  - Runtime configuration
  - State information

Use --format json to get JSON output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		return project.Inspect(format)
	},
}

func init() {
	ProjectCmd.AddCommand(inspectCmd)
}
