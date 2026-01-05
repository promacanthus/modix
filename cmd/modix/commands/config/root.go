package config

import (
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config command family
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage modix configuration",
	Long: `Manage modix configuration files and settings.

This command provides utilities for managing the modix configuration
including initialization, path discovery, and reset operations.

Available subcommands:
  init      Initialize default configuration
  path      Show configuration file path
  show      Display current configuration
  reset     Reset configuration to defaults
  check     Validate configuration

Examples:
  modix config init
  modix config path
  modix config show
  modix config reset`,
}

func init() {
	ConfigCmd.AddCommand(initCmd)
	ConfigCmd.AddCommand(pathCmd)
	ConfigCmd.AddCommand(showCmd)
	ConfigCmd.AddCommand(resetCmd)
	ConfigCmd.AddCommand(checkCmd)
}
