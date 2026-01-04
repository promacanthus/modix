package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/config"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [model-name]",
	Short: "Switch to a different model",
	Long: `Switch to a different model.

Examples:
  modix switch "claude-official"
  modix switch "deepseek-reasoner"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return runSwitch(modelName)
	},
}

func runSwitch(modelName string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Find the model in all vendors
	var foundVendor string
	var foundModel string

	for vendor, vendorConfig := range modixConfig.Vendors {
		for _, model := range vendorConfig.Models {
			if model == modelName {
				foundVendor = vendor
				foundModel = model
				break
			}
		}
		if foundVendor != "" {
			break
		}
	}

	if foundVendor == "" {
		return fmt.Errorf("model '%s' not found in any vendor", modelName)
	}

	// Switch to the model
	err = modixConfig.SetCurrentVendorAndModel(foundVendor, foundModel)
	if err != nil {
		return fmt.Errorf("failed to switch to model: %w", err)
	}

	// Update Claude configuration if needed
	if err := config.UpdateClaudeEnvConfig(foundModel, foundVendor); err != nil {
		return fmt.Errorf("failed to update Claude configuration: %w", err)
	}

	// Save the configuration
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Switched to model: %s\n", modelName)
	return nil
}