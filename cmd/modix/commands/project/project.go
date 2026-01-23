package project

import (
	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/cmd/modix/commands"
)

// ProjectCmd represents the project command
var ProjectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage modix projects",
	Long: `Manage modix projects, including initialization, validation, and inspection.

Available commands:
  init      Initialize a new modix project
  check     Check dependencies for modix project
  validate  Validate modix project configuration
  inspect   Inspect modix project configuration`,
}

func init() {
	commands.RootCmd.AddCommand(ProjectCmd)
}
