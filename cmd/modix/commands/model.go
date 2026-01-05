package commands

import (
	"fmt"

	"github.com/promacanthus/modix/internal/config"
	"github.com/spf13/cobra"
)

// modelCmd represents the model command family
var modelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage and switch LLM models",
	Long: `Manage and switch between LLM models.

This command provides functionality to view all available models
across all vendors and switch between them.

Available subcommands:
  list      List all models across all vendors
  switch    Switch to a specific model
  status    Show current model status

Examples:
  # List all available models
  modix model list

  # Switch to a model
  modix model switch deepseek-reasoner

  # Show current model status
  modix model status`,
}

// modelListCmd lists all models
var modelListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all models",
	Long:  `List all configured models grouped by vendor.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listAllModels()
	},
}

// modelSwitchCmd switches to a model
var modelSwitchCmd = &cobra.Command{
	Use:   "switch [model-name]",
	Short: "Switch to a model",
	Long: `Switch to a different model.

Examples:
  modix model switch deepseek-reasoner

This will update Claude Code configuration to use the selected model.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return switchModel(modelName)
	},
}

// modelStatusCmd shows current model status
var modelStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current model status",
	Long:  `Show the currently active model and its configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showModelStatus()
	},
}

func init() {
	modelCmd.AddCommand(modelListCmd)
	modelCmd.AddCommand(modelSwitchCmd)
	modelCmd.AddCommand(modelStatusCmd)
}

// Implementation functions
func listAllModels() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(modixConfig.Vendors) == 0 {
		fmt.Println("No vendors configured")
		return nil
	}

	fmt.Printf("%-30s %-15s %-25s %-10s\n", "MODEL", "VENDOR", "COMPANY", "STATUS")
	fmt.Println("-------------------------------------------------------------------")

	for vendorID, vendorConfig := range modixConfig.Vendors {
		for _, model := range vendorConfig.Models {
			status := "Ready"
			if vendorConfig.APIKey == "" {
				status = "No API Key"
			} else if vendorConfig.APIEndpoint == "" {
				status = "No Endpoint"
			}

			// Highlight current model
			modelDisplay := model
			if modixConfig.CurrentVendor == vendorID && modixConfig.CurrentModel == model {
				modelDisplay = model + " (current)"
			}

			fmt.Printf("%-30s %-15s %-25s %-10s\n", modelDisplay, vendorID, vendorConfig.Company, status)
		}
	}

	return nil
}

func switchModel(modelName string) error {
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

	// Update Claude configuration
	if err := config.UpdateClaudeEnvConfig(foundModel, foundVendor); err != nil {
		return fmt.Errorf("failed to update Claude configuration: %w", err)
	}

	// Save the configuration
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Switched to model: %s@%s\n", foundModel, foundVendor)
	return nil
}

func showModelStatus() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if currentModel, modelConfig, exists := modixConfig.GetCurrentModel(); exists {
		fmt.Println("Current model status:")
		fmt.Println()
		fmt.Printf("Model: %s\n", *currentModel)
		fmt.Printf("Vendor: %s\n", modixConfig.CurrentVendor)
		fmt.Printf("Company: %s\n", modelConfig.Company)
		fmt.Printf("API Endpoint: %s\n", modelConfig.APIEndpoint)
		fmt.Printf("API Key: %s\n", maskAPIKey(modelConfig.APIKey))
	} else {
		fmt.Println("No current model configured")
	}

	return nil
}

func maskAPIKey(key string) string {
	if key == "" {
		return "[Not set]"
	}
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
