package project

import (
	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/project"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate modix project configuration",
	Long: `Validate the modix project configuration files.

This command checks:
  - .modix/ directory exists
  - All configuration files exist
  - JSON format is valid
  - Required fields are present

If validation fails, it will show detailed error messages.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		return project.Validate(format)
	},
}

func init() {
	ProjectCmd.AddCommand(validateCmd)
}
