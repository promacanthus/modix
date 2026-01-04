package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/config"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [model-name]",
	Short: "Add a new model configuration",
	Long: `Add a new model configuration to a vendor.

Examples:
  modix add "My-Model" -c "MyCorp" -v "my-vendor" -u "https://api.mycorp.com" -k "my-api-key"
  modix add "deepseek-reasoner" -c "DeepSeek" -v "deepseek" -u "https://api.deepseek.com/v1" -k "sk-xxx"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		company, _ := cmd.Flags().GetString("company")
		vendor, _ := cmd.Flags().GetString("vendor")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		apiKey, _ := cmd.Flags().GetString("api-key")

		return runAdd(modelName, company, vendor, endpoint, apiKey)
	},
}

func runAdd(modelName, company, vendor, endpoint, apiKey string) error {
	configPath := config.GetConfigPath()

	// Load existing configuration
	modixConfig, err := config.LoadConfigFromPath(configPath)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if this model name already exists in any vendor
	allModels := modixConfig.GetModels()
	for _, existingModel := range allModels {
		if existingModel == modelName {
			return fmt.Errorf("model '%s' already exists. Please use a different name or remove the existing one first.", modelName)
		}
	}

	// First try to add the model to an existing vendor
	err = modixConfig.AddModelToVendor(vendor, modelName)
	if err == nil {
		// Update the API endpoint and key if they changed
		if vendorConfig, exists := modixConfig.GetVendor(vendor); exists {
			vendorConfig.Company = company
			vendorConfig.APIEndpoint = endpoint
			vendorConfig.APIKey = apiKey
			modixConfig.Vendors[vendor] = *vendorConfig
		}
		if err := config.SaveConfig(modixConfig); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}
		fmt.Printf("Added model '%s' to existing vendor '%s'\n", modelName, vendor)
	} else {
		// If vendor doesn't exist, create a new vendor config
		modelConfig := config.ModelConfig{
			Company:     company,
			APIEndpoint: endpoint,
			APIKey:      apiKey,
			Models:      []string{modelName},
		}

		modixConfig.AddVendor(vendor, modelConfig)
		if err := config.SaveConfig(modixConfig); err != nil {
			return fmt.Errorf("failed to save configuration: %w", err)
		}

		fmt.Printf("Created new vendor '%s' with model: %s\n", vendor, modelName)
	}

	fmt.Printf("Switch to it with: modix switch %s\n", modelName)
	return nil
}

func init() {
	addCmd.Flags().StringP("company", "c", "", "Company that develops the model")
	addCmd.Flags().StringP("vendor", "v", "", "API provider/vendor identifier")
	addCmd.Flags().StringP("endpoint", "u", "", "API endpoint URL")
	addCmd.Flags().StringP("api-key", "k", "", "API key")

	addCmd.MarkFlagRequired("company")
	addCmd.MarkFlagRequired("vendor")
	addCmd.MarkFlagRequired("endpoint")
	addCmd.MarkFlagRequired("api-key")
}