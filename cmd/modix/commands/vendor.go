package commands

import (
	"fmt"
	"strings"

	"github.com/promacanthus/modix/internal/config"
	"github.com/spf13/cobra"
)

// vendorCmd represents the vendor command family
var vendorCmd = &cobra.Command{
	Use:   "vendor",
	Short: "Manage LLM vendors",
	Long: `Manage LLM vendor configurations.

Vendors are companies or providers that offer LLM models. Each vendor
can have multiple models and shared API configuration.

Available subcommands:
  add       Add a new vendor
  remove    Remove a vendor
  update    Update vendor configuration
  list      List all vendors
  show      Show vendor details
  model     Manage models for a vendor

Examples:
  # Add a vendor
  modix vendor add deepseek --company "DeepSeek" --endpoint "https://api.deepseek.com/v1" --api-key "sk-xxx"

  # List vendors
  modix vendor list

  # Show vendor details
  modix vendor show deepseek

  # Update vendor
  modix vendor update deepseek --company "DeepSeek AI"

  # Add model to vendor
  modix vendor model add deepseek deepseek-reasoner

  # Remove model from vendor
  modix vendor model remove deepseek deepseek-reasoner`,
}

// vendorAddCmd adds a new vendor
var vendorAddCmd = &cobra.Command{
	Use:   "add [vendor-id]",
	Short: "Add a new vendor",
	Long: `Add a new vendor configuration.

Examples:
  modix vendor add deepseek \
    --company "DeepSeek" \
    --endpoint "https://api.deepseek.com/v1" \
    --api-key "sk-xxx"

  modix vendor add custom \
    --company "MyCorp" \
    --endpoint "https://api.mycorp.com/v1" \
    --api-key "my-key"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendorID := args[0]
		company, _ := cmd.Flags().GetString("company")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		apiKey, _ := cmd.Flags().GetString("api-key")

		return addVendor(vendorID, company, endpoint, apiKey)
	},
}

// vendorRemoveCmd removes a vendor
var vendorRemoveCmd = &cobra.Command{
	Use:   "remove [vendor-id]",
	Short: "Remove a vendor",
	Long: `Remove a vendor configuration.

This will remove all models associated with the vendor.

Examples:
  modix vendor remove deepseek`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendorID := args[0]
		return removeVendor(vendorID)
	},
}

// vendorUpdateCmd updates a vendor
var vendorUpdateCmd = &cobra.Command{
	Use:   "update [vendor-id]",
	Short: "Update a vendor",
	Long: `Update vendor configuration.

Examples:
  modix vendor update deepseek --company "DeepSeek AI"
  modix vendor update deepseek --endpoint "https://new-api.com/v1"
  modix vendor update deepseek --api-key "new-key"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendorID := args[0]
		company, _ := cmd.Flags().GetString("company")
		endpoint, _ := cmd.Flags().GetString("endpoint")
		apiKey, _ := cmd.Flags().GetString("api-key")

		return updateVendor(vendorID, company, endpoint, apiKey)
	},
}

// vendorListCmd lists all vendors
var vendorListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all vendors",
	Long:  `List all configured vendors with their details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listVendors()
	},
}

// vendorShowCmd shows vendor details
var vendorShowCmd = &cobra.Command{
	Use:   "show [vendor-id]",
	Short: "Show vendor details",
	Long: `Show detailed vendor configuration including API key.

Examples:
  modix vendor show deepseek`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendorID := args[0]
		return showVendor(vendorID)
	},
}

// vendorModelCmd represents the model subcommand under vendor
var vendorModelCmd = &cobra.Command{
	Use:   "model",
	Short: "Manage models for a vendor",
	Long: `Manage models for a specific vendor.

This command allows you to add or remove models from a vendor.
Each vendor can have multiple models.

Examples:
  # Add a model to a vendor
  modix vendor model add deepseek deepseek-reasoner

  # Remove a model from a vendor
  modix vendor model remove deepseek deepseek-reasoner`,
}

// vendorModelAddCmd adds a model to a vendor
var vendorModelAddCmd = &cobra.Command{
	Use:   "add [vendor-id] [model-name]",
	Short: "Add a model to a vendor",
	Long: `Add a model to an existing vendor.

Examples:
  modix vendor model add deepseek deepseek-reasoner
  modix vendor model add custom my-model

If the vendor doesn't exist, you'll be prompted to create it first.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendorID := args[0]
		modelName := args[1]
		return addModelToVendor(vendorID, modelName)
	},
}

// vendorModelRemoveCmd removes a model from a vendor
var vendorModelRemoveCmd = &cobra.Command{
	Use:   "remove [vendor-id] [model-name]",
	Short: "Remove a model from a vendor",
	Long: `Remove a model from a vendor.

Examples:
  modix vendor model remove deepseek deepseek-reasoner

If the vendor has no remaining models, it will be kept but empty.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		vendorID := args[0]
		modelName := args[1]
		return removeModelFromVendor(vendorID, modelName)
	},
}

func init() {
	// Add flags for add command
	vendorAddCmd.Flags().StringP("company", "c", "", "Company name")
	vendorAddCmd.Flags().StringP("endpoint", "u", "", "API endpoint URL")
	vendorAddCmd.Flags().StringP("api-key", "k", "", "API key")
	vendorAddCmd.MarkFlagRequired("company")
	vendorAddCmd.MarkFlagRequired("endpoint")
	vendorAddCmd.MarkFlagRequired("api-key")

	// Add flags for update command
	vendorUpdateCmd.Flags().StringP("company", "c", "", "Company name")
	vendorUpdateCmd.Flags().StringP("endpoint", "u", "", "API endpoint URL")
	vendorUpdateCmd.Flags().StringP("api-key", "k", "", "API key")

	// Add subcommands to vendor command
	vendorCmd.AddCommand(vendorAddCmd)
	vendorCmd.AddCommand(vendorRemoveCmd)
	vendorCmd.AddCommand(vendorUpdateCmd)
	vendorCmd.AddCommand(vendorListCmd)
	vendorCmd.AddCommand(vendorShowCmd)

	// Add model subcommands
	vendorModelCmd.AddCommand(vendorModelAddCmd)
	vendorModelCmd.AddCommand(vendorModelRemoveCmd)
	vendorCmd.AddCommand(vendorModelCmd)
}

// Implementation functions
func addVendor(vendorID, company, endpoint, apiKey string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if vendor already exists
	if _, exists := modixConfig.GetVendor(vendorID); exists {
		return fmt.Errorf("vendor '%s' already exists", vendorID)
	}

	vendorConfig := config.VendorConfig{
		Company:     company,
		APIEndpoint: endpoint,
		APIKey:      apiKey,
		Models:      []string{},
	}

	modixConfig.AddVendor(vendorID, vendorConfig)
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Added vendor '%s' (%s)\n", vendorID, company)
	fmt.Printf("Next step: Add models with 'modix vendor model add %s <model-name>'\n", vendorID)
	return nil
}

func removeVendor(vendorID string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if _, exists := modixConfig.GetVendor(vendorID); !exists {
		return fmt.Errorf("vendor '%s' not found", vendorID)
	}

	// Check if this is the current vendor
	if modixConfig.CurrentVendor == vendorID {
		fmt.Printf("Warning: '%s' is the current vendor. Consider switching first.\n", vendorID)
	}

	// Get model count for confirmation
	vendorConfig, _ := modixConfig.GetVendor(vendorID)
	modelCount := len(vendorConfig.Models)

	modixConfig.RemoveVendor(vendorID)
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Removed vendor '%s' (%d models)\n", vendorID, modelCount)
	return nil
}

func updateVendor(vendorID, company, endpoint, apiKey string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	vendorConfig, exists := modixConfig.GetVendor(vendorID)
	if !exists {
		return fmt.Errorf("vendor '%s' not found", vendorID)
	}

	var updates []string

	if company != "" {
		vendorConfig.Company = company
		updates = append(updates, fmt.Sprintf("Company: %s", company))
	}
	if endpoint != "" {
		vendorConfig.APIEndpoint = endpoint
		updates = append(updates, fmt.Sprintf("Endpoint: %s", endpoint))
	}
	if apiKey != "" {
		vendorConfig.APIKey = apiKey
		updates = append(updates, "API Key: [updated]")
	}

	if len(updates) == 0 {
		return fmt.Errorf("no updates specified")
	}

	modixConfig.Vendors[vendorID] = *vendorConfig

	// Update Claude config if this is the current vendor
	if modixConfig.CurrentVendor == vendorID && len(vendorConfig.Models) > 0 {
		currentModel := modixConfig.CurrentModel
		if err := config.UpdateClaudeEnvConfig(currentModel, vendorID); err != nil {
			return fmt.Errorf("failed to update Claude configuration: %w", err)
		}
	}

	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Updated vendor '%s':\n", vendorID)
	for _, update := range updates {
		fmt.Printf("  - %s\n", update)
	}
	return nil
}

func listVendors() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(modixConfig.Vendors) == 0 {
		fmt.Println("No vendors configured")
		return nil
	}

	fmt.Printf("%-15s %-20s %-40s %-10s\n", "VENDOR", "COMPANY", "ENDPOINT", "MODELS")
	fmt.Println(strings.Repeat("-", 90))

	for vendorID, vendorConfig := range modixConfig.Vendors {
		endpoint := vendorConfig.APIEndpoint
		if endpoint == "" {
			endpoint = "[Not set]"
		}
		modelCount := len(vendorConfig.Models)
		fmt.Printf("%-15s %-20s %-40s %-10d\n", vendorID, vendorConfig.Company, endpoint, modelCount)
	}

	return nil
}

func showVendor(vendorID string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	vendorConfig, exists := modixConfig.GetVendor(vendorID)
	if !exists {
		return fmt.Errorf("vendor '%s' not found", vendorID)
	}

	fmt.Printf("Vendor: %s\n", vendorID)
	fmt.Printf("Company: %s\n", vendorConfig.Company)
	fmt.Printf("API Endpoint: %s\n", vendorConfig.APIEndpoint)
	fmt.Printf("API Key: %s\n", vendorConfig.APIKey)
	fmt.Println("\nModels:")
	for _, model := range vendorConfig.Models {
		marker := ""
		if modixConfig.CurrentVendor == vendorID && modixConfig.CurrentModel == model {
			marker = " (current)"
		}
		fmt.Printf("  - %s%s\n", model, marker)
	}

	return nil
}

func addModelToVendor(vendorID, modelName string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if model already exists in any vendor
	allModels := modixConfig.GetModels()
	for _, existingModel := range allModels {
		if existingModel == modelName {
			return fmt.Errorf("model '%s' already exists", modelName)
		}
	}

	// Try to add to existing vendor
	err = modixConfig.AddModelToVendor(vendorID, modelName)
	if err != nil {
		return fmt.Errorf("vendor '%s' not found. Use 'modix vendor add %s' to create it first", vendorID, vendorID)
	}

	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Added model '%s' to vendor '%s'\n", modelName, vendorID)
	fmt.Printf("Switch to it with: modix model switch %s\n", modelName)
	return nil
}

func removeModelFromVendor(vendorID, modelName string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Verify vendor exists
	vendorConfig, exists := modixConfig.GetVendor(vendorID)
	if !exists {
		return fmt.Errorf("vendor '%s' not found", vendorID)
	}

	// Check if model exists in this vendor
	found := false
	for _, model := range vendorConfig.Models {
		if model == modelName {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("model '%s' not found in vendor '%s'", modelName, vendorID)
	}

	// Remove the model from the vendor
	newModels := []string{}
	for _, model := range vendorConfig.Models {
		if model != modelName {
			newModels = append(newModels, model)
		}
	}
	vendorConfig.Models = newModels
	modixConfig.Vendors[vendorID] = *vendorConfig

	// If we removed the current model, switch to default
	if modixConfig.CurrentModel == modelName && modixConfig.CurrentVendor == vendorID {
		modixConfig.CurrentModel = modixConfig.DefaultModel
		modixConfig.CurrentVendor = modixConfig.DefaultVendor
		fmt.Printf("Removed current model. Switched to default: %s@%s\n", modixConfig.DefaultModel, modixConfig.DefaultVendor)
	}

	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Removed model '%s' from vendor '%s'\n", modelName, vendorID)
	return nil
}
