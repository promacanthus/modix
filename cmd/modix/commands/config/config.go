package config

import (
	"fmt"

	"github.com/promacanthus/modix/internal/config"
	"github.com/spf13/cobra"
)

// initCmd initializes configuration
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize default configuration",
	Long: `Initialize default configuration with common LLM providers.

Creates a configuration file with pre-configured vendors including:
- Anthropic (Claude)
- DeepSeek
- Alibaba (Qwen)
- ByteDance (Doubao)
- Moonshot AI
- Kuaishou
- MiniMax
- ZHIPU AI
- Xiaomi

Examples:
  modix config init`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initConfig()
	},
}

// pathCmd shows configuration path
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `Display the path to the modix configuration file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showPath()
	},
}

// showCmd shows current configuration
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long: `Display the current configuration in a readable format.

Shows all vendors, models, current selections, and agent configurations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return showConfig()
	},
}

// resetCmd resets configuration
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	Long: `Reset configuration to default values.

⚠️  Warning: This will overwrite your current configuration.

Examples:
  modix config reset --force`,
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")
		return resetConfig(force)
	},
}

// checkCmd validates configuration
var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Validate configuration",
	Long: `Validate the current configuration for errors and completeness.

Checks for:
- Valid vendor configurations
- Model availability
- API key and endpoint settings
- Current model consistency

Examples:
  modix config check`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return checkConfig()
	},
}

func init() {
	resetCmd.Flags().BoolP("force", "f", false, "Force reset without confirmation")
}

// Implementation functions
func initConfig() error {
	configPath := config.GetConfigPath()

	// Check if config already exists
	if config.ConfigExists() {
		fmt.Printf("Configuration already exists at: %s\n", configPath)
		fmt.Println("Use 'modix config show' to view current configuration")
		fmt.Println("Use 'modix config reset' to reset to defaults")
		return nil
	}

	// Create default configuration
	modixConfig := config.SetupDefaultModels()

	// Save configuration
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to create default configuration: %w", err)
	}

	fmt.Printf("✓ Initialized configuration at: %s\n", configPath)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Add API keys to vendors: modix config show")
	fmt.Println("  2. Add models: modix llm model add <vendor> <model>")
	fmt.Println("  3. Switch to a model: modix llm model switch <model>")
	fmt.Println("  4. Configure agents: modix agent add claude-code")

	return nil
}

func showPath() error {
	path := config.GetConfigFilePath()
	fmt.Printf("Configuration file: %s\n", path)
	return nil
}

func showConfig() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println("=== Modix Configuration ===")
	fmt.Println()

	// Current selections
	fmt.Println("Current Selections:")
	fmt.Printf("  Model: %s\n", modixConfig.CurrentModel)
	fmt.Printf("  Vendor: %s\n", modixConfig.CurrentVendor)
	fmt.Printf("  Agent: %s\n", modixConfig.CurrentAgent)
	fmt.Println()

	// Vendors and models
	fmt.Println("Vendors & Models:")
	for vendorID, vendorConfig := range modixConfig.Vendors {
		fmt.Printf("  %s (%s):\n", vendorID, vendorConfig.Company)
		if vendorConfig.APIEndpoint != "" {
			fmt.Printf("    Endpoint: %s\n", vendorConfig.APIEndpoint)
		}
		if vendorConfig.APIKey != "" {
			keyDisplay := vendorConfig.APIKey
			if len(keyDisplay) > 10 {
				keyDisplay = keyDisplay[:10] + "..."
			}
			fmt.Printf("    API Key: %s\n", keyDisplay)
		}
		if len(vendorConfig.Models) > 0 {
			fmt.Printf("    Models: %v\n", vendorConfig.Models)
		}
		fmt.Println()
	}

	// Agents
	if len(modixConfig.Agents) > 0 {
		fmt.Println("Agents:")
		for agentID, agentConfig := range modixConfig.Agents {
			status := "Enabled"
			if !agentConfig.Enabled {
				status = "Disabled"
			}
			fmt.Printf("  %s: %s (%s) - %s\n", agentID, agentConfig.Name, agentConfig.Provider, status)
		}
	} else {
		fmt.Println("Agents: None configured")
	}

	return nil
}

func resetConfig(force bool) error {
	if !force {
		return fmt.Errorf("use --force to confirm reset")
	}

	// Create default configuration
	modixConfig := config.SetupDefaultModels()

	// Save configuration
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Println("✓ Configuration reset to defaults")
	return nil
}

func checkConfig() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println("=== Configuration Check ===")
	fmt.Println()

	// Check vendors
	vendorCount := len(modixConfig.Vendors)
	fmt.Printf("Vendors: %d\n", vendorCount)
	if vendorCount == 0 {
		fmt.Println("  ✗ No vendors configured")
	} else {
		fmt.Println("  ✓ Vendors configured")
	}

	// Check models
	modelCount := 0
	configuredVendors := 0
	for _, vendorConfig := range modixConfig.Vendors {
		modelCount += len(vendorConfig.Models)
		if vendorConfig.APIEndpoint != "" && vendorConfig.APIKey != "" {
			configuredVendors++
		}
	}
	fmt.Printf("Models: %d\n", modelCount)
	if modelCount == 0 {
		fmt.Println("  ✗ No models configured")
	} else {
		fmt.Println("  ✓ Models configured")
	}

	// Check configured vendors
	fmt.Printf("Configured vendors: %d\n", configuredVendors)
	if configuredVendors == 0 {
		fmt.Println("  ⚠ No vendors with complete API configuration")
	} else {
		fmt.Println("  ✓ Some vendors have complete API configuration")
	}

	// Check current model
	if modixConfig.CurrentModel != "" && modixConfig.CurrentVendor != "" {
		fmt.Printf("Current model: %s@%s\n", modixConfig.CurrentModel, modixConfig.CurrentVendor)
		fmt.Println("  ✓ Current model set")
	} else {
		fmt.Println("Current model: Not set")
		fmt.Println("  ⚠ No current model selected")
	}

	// Check agents
	if len(modixConfig.Agents) > 0 {
		fmt.Printf("Agents: %d configured\n", len(modixConfig.Agents))
		fmt.Println("  ✓ Agents configured")
	} else {
		fmt.Println("Agents: None")
		fmt.Println("  ℹ No agents configured (use 'modix agent add' to add)")
	}

	// Overall status
	fmt.Println()
	if configuredVendors > 0 && modelCount > 0 && modixConfig.CurrentModel != "" {
		fmt.Println("✓ Configuration is ready to use")
	} else {
		fmt.Println("⚠ Configuration needs setup")
		fmt.Println("\nRecommended actions:")
		if configuredVendors == 0 {
			fmt.Println("  - Add API keys to vendors: modix config show")
		}
		if modelCount == 0 {
			fmt.Println("  - Add models: modix llm model add <vendor> <model>")
		}
		if modixConfig.CurrentModel == "" {
			fmt.Println("  - Switch to a model: modix llm model switch <model>")
		}
	}

	return nil
}
