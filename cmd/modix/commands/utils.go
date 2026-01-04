package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/promacanthus/modix/internal/config"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `Show the path to the configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runPath()
	},
}

func runPath() error {
	path := config.GetConfigFilePath()
	fmt.Printf("Configuration file path: %s\n", path)
	return nil
}

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove [model-name]",
	Short: "Remove a model configuration",
	Long: `Remove a model configuration.

Examples:
  modix remove "my-model"
  modix remove "deepseek-reasoner"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName := args[0]
		return runRemove(modelName)
	},
}

func runRemove(modelName string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Find which vendor this model belongs to
	modelInfos := modixConfig.GetAllModelInfos()
	var targetVendor string

	for _, modelInfo := range modelInfos {
		for _, model := range modelInfo.Models {
			if model == modelName {
				targetVendor = modelInfo.Vendor
				break
			}
		}
		if targetVendor != "" {
			break
		}
	}

	if targetVendor == "" {
		return fmt.Errorf("model '%s' not found", modelName)
	}

	// Remove the model from the vendor
	if vendorConfig, exists := modixConfig.GetVendor(targetVendor); exists {
		newModels := []string{}
		for _, model := range vendorConfig.Models {
			if model != modelName {
				newModels = append(newModels, model)
			}
		}
		vendorConfig.Models = newModels
		modixConfig.Vendors[targetVendor] = *vendorConfig

		// If the vendor has no more models, remove the vendor entirely
		if len(vendorConfig.Models) == 0 {
			modixConfig.RemoveVendor(targetVendor)
			fmt.Printf("Vendor '%s' had no remaining models and was removed\n", targetVendor)
		}
	}

	// If we removed the current model, switch to default
	if modixConfig.CurrentModel == modelName {
		modixConfig.CurrentModel = modixConfig.DefaultModel
		modixConfig.CurrentVendor = modixConfig.DefaultVendor
		fmt.Printf("Removed current model. Switched to default: %s@%s\n", modixConfig.DefaultModel, modixConfig.DefaultVendor)
	}

	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}
	fmt.Printf("Removed model: %s\n", modelName)

	return nil
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show [vendor]",
	Short: "Show details for a specific vendor",
	Long: `Show details for a specific vendor including API key.

Examples:
  modix show "anthropic"
  modix show "deepseek"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendor := args[0]
		return runShow(vendor)
	},
}

func runShow(vendor string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if vendorConfig, exists := modixConfig.GetVendor(vendor); exists {
		fmt.Println("Vendor Details:")
		fmt.Println()
		fmt.Printf("Vendor ID: %s\n", vendor)
		fmt.Printf("Company: %s\n", vendorConfig.Company)
		fmt.Printf("API Endpoint: %s\n", vendorConfig.APIEndpoint)
		fmt.Printf("API Key: %s\n", vendorConfig.APIKey)
		fmt.Println()
		fmt.Println("Models:")
		for _, model := range vendorConfig.Models {
			currentMarker := ""
			if modixConfig.CurrentVendor == vendor && modixConfig.CurrentModel == model {
				currentMarker = " (current)"
			}
			fmt.Printf("  - %s%s\n", model, currentMarker)
		}
	} else {
		return fmt.Errorf("vendor '%s' not found", vendor)
	}

	return nil
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update [vendor]",
	Short: "Update an existing vendor configuration",
	Long: `Update an existing vendor configuration (model, company, api_endpoint, or api_key).

Examples:
  modix update "deepseek" --add-model "new-model"
  modix update "deepseek" --company "NewCorp"
  modix update "deepseek" --endpoint "https://new-api.com"
  modix update "deepseek" --api-key "new-api-key"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendor := args[0]
		addModel, _ := cmd.Flags().GetString("add-model")
		company, _ := cmd.Flags().GetString("company")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		apiKey, _ := cmd.Flags().GetString("api-key")

		return runUpdate(vendor, addModel, company, endpoint, apiKey)
	},
}

func runUpdate(vendor, addModel, company, endpoint, apiKey string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if vendor exists
	if _, exists := modixConfig.GetVendor(vendor); !exists {
		return fmt.Errorf("vendor '%s' not found. Use 'modix add' to create a new vendor first.", vendor)
	}

	var updates []string

	// Update company if provided
	if company != "" {
		if vendorConfig, exists := modixConfig.GetVendor(vendor); exists {
			vendorConfig.Company = company
			modixConfig.Vendors[vendor] = *vendorConfig
			updates = append(updates, fmt.Sprintf("Company: %s", company))
		}
	}

	// Update endpoint if provided
	if endpoint != "" {
		if vendorConfig, exists := modixConfig.GetVendor(vendor); exists {
			vendorConfig.APIEndpoint = endpoint
			modixConfig.Vendors[vendor] = *vendorConfig
			updates = append(updates, fmt.Sprintf("API Endpoint: %s", endpoint))
		}
	}

	// Update API key if provided
	if apiKey != "" {
		if vendorConfig, exists := modixConfig.GetVendor(vendor); exists {
			vendorConfig.APIKey = apiKey
			modixConfig.Vendors[vendor] = *vendorConfig
			updates = append(updates, "API Key: [updated]")
		}
	}

	// Add model if provided
	if addModel != "" {
		if err := modixConfig.AddModelToVendor(vendor, addModel); err == nil {
			updates = append(updates, fmt.Sprintf("Added model: %s", addModel))
		} else {
			return fmt.Errorf("model '%s' already exists in vendor '%s'", addModel, vendor)
		}
	}

	// Check if any updates were made
	if len(updates) == 0 {
		fmt.Println("No updates were specified. Use --help to see available options.")
		return nil
	}

	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Updated vendor '%s' with:\n", vendor)
	for _, update := range updates {
		fmt.Printf("  - %s\n", update)
	}

	return nil
}

func init() {
	pathCmd.Flags().BoolP("verbose", "v", false, "Show verbose output")
	RootCmd.AddCommand(pathCmd)

	removeCmd.Flags().BoolP("force", "f", false, "Force remove without confirmation")
	RootCmd.AddCommand(removeCmd)

	showCmd.Flags().BoolP("include-key", "k", false, "Include API key in output")
	RootCmd.AddCommand(showCmd)

	updateCmd.Flags().StringP("add-model", "m", "", "Add a model to the vendor")
	updateCmd.Flags().StringP("company", "c", "", "Update company name")
	updateCmd.Flags().StringP("endpoint", "u", "", "Update API endpoint URL")
	updateCmd.Flags().StringP("api-key", "k", "", "Update API key")
	RootCmd.AddCommand(updateCmd)
}