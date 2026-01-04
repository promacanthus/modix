package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration",
	Long: `Initialize default configuration.

Creates a default configuration file with common LLM providers and their settings.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit()
	},
}

func runInit() error {
	configPath := config.GetConfigPath()

	// Check if config already exists
	if config.ConfigExists() {
		fmt.Printf("Configuration already exists at: %s\n", configPath)
		fmt.Println("Use 'modix list' to see available models")
		return nil
	}

	// Create default configuration
	modixConfig := config.SetupDefaultModels()

	// Save configuration
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to create default configuration: %w", err)
	}

	fmt.Printf("Initialized default configuration at: %s\n", configPath)
	fmt.Println("Edit the configuration file to add your API keys and enable models")
	fmt.Println("Use 'modix list' to see available models")

	return nil
}