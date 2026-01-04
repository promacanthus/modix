package commands

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/promacanthus/modix/internal/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration",
	Long: `Initialize default configuration with pre-configured models.

This command creates a default configuration file with common AI model vendors
pre-configured. You can then edit the configuration file to add your API keys
and enable the models you want to use.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	if config.ConfigExists() {
		fmt.Printf("Configuration already exists at: %s\n", config.GetConfigFilePath())
		fmt.Println("Use 'modix path' to show the configuration file path")
		return nil
	}

	modixConfig := config.SetupDefaultModels()
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	color.Cyan("Initialized default configuration at: %s\n", config.GetConfigFilePath())
	fmt.Println("Edit the configuration file to add your API keys and enable models")
	fmt.Println("Use 'modix list' to see available models")

	return nil
}

func init() {
	RootCmd.AddCommand(initCmd)
}