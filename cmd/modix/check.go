package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/promacanthus/modix/internal/config"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [tool]",
	Short: "Check configuration for different tools",
	Long: `Check configuration for different tools.

Examples:
  modix check claude-code
  modix check vscode`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		tool := args[0]
		return runCheck(tool)
	},
}

func runCheck(tool string) error {
	switch tool {
	case "claude-code":
		return checkClaudeCode()
	case "vscode":
		return checkVSCode()
	default:
		return fmt.Errorf("unknown tool: %s", tool)
	}
}

func checkClaudeCode() error {
	fmt.Println("Checking Claude Code configuration...")
	configPath := config.GetConfigFilePath()
	fmt.Printf("Config file path: %s\n", configPath)

	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	fmt.Println()
	fmt.Println("=== Claude Code Configuration ===")
	fmt.Println()

	// Check if current model is Claude
	if modixConfig.CurrentVendor == "anthropic" && modixConfig.CurrentModel == "Claude" {
		fmt.Println("✓ Current model is Claude (Anthropic)")
	} else {
		fmt.Printf("⚠ Current model is %s@%s (not Claude)\n", modixConfig.CurrentVendor, modixConfig.CurrentModel)
	}

	// Check if Anthropic vendor exists
	if vendorConfig, exists := modixConfig.GetVendor("anthropic"); exists {
		fmt.Printf("✓ Anthropic vendor configured: %s\n", vendorConfig.Company)
		fmt.Printf("  Models: %v\n", vendorConfig.Models)
		if vendorConfig.APIKey != "" {
			fmt.Printf("  API Key: ✓ Configured\n")
		} else {
			fmt.Printf("  API Key: ✗ Not configured\n")
		}
		if vendorConfig.APIEndpoint != "" {
			fmt.Printf("  API Endpoint: ✓ Configured\n")
		} else {
			fmt.Printf("  API Endpoint: ✗ Not configured\n")
		}
	} else {
		fmt.Println("⚠ Anthropic vendor not configured")
	}

	return nil
}

func checkVSCode() error {
	fmt.Println("Checking VS Code configuration...")
	configPath := config.GetConfigFilePath()
	fmt.Printf("Config file path: %s\n", configPath)
	return nil
}