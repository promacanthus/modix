package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/promacanthus/modix/internal/config"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [tool]",
	Short: "Check configuration for different tools",
	Long: `Check configuration for different tools.

Available tools:
  - claude-code: Check Claude Code configuration
  - modix: Check Modix configuration
  - codex: Check Codex configuration
  - gemini-cli: Check Gemini CLI configuration

Examples:
  modix check claude-code
  modix check modix`,
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
	case "modix":
		return checkModixConfig()
	case "codex":
		color.Yellow("Codex configuration check is not yet implemented")
		return nil
	case "gemini-cli":
		color.Yellow("Gemini CLI configuration check is not yet implemented")
		return nil
	default:
		return fmt.Errorf("unknown tool '%s'. Supported tools: claude-code, modix, codex, gemini-cli", tool)
	}
}

func checkClaudeCode() error {
	claudeConfigPath := config.GetClaudeConfigPath()

	color.Yellow("Checking Claude Code configuration...")
	color.Cyan("Config file path: %s\n", claudeConfigPath)

	if !config.IsClaudeConfigured() {
		color.Red("Configuration file not found: %s\n", claudeConfigPath)
		return nil
	}

	configSummary, err := config.GetClaudeConfigSummary()
	if err != nil {
		return fmt.Errorf("failed to load Claude configuration: %w", err)
	}

	color.Cyan("=== Claude Code Configuration ===")
	fmt.Println()

	// Pretty print JSON with syntax highlighting
	prettyJSON, err := json.MarshalIndent(configSummary, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format configuration: %w", err)
	}

	fmt.Println(string(prettyJSON))

	return nil
}

func checkModixConfig() error {
	modixConfigPath := config.GetConfigFilePath()

	color.Yellow("Checking Modix configuration...")
	color.Cyan("Config file path: %s\n", modixConfigPath)

	if !config.ConfigExists() {
		color.Red("Configuration file not found: %s\n", modixConfigPath)
		color.Blue("Run 'modix init' to create a default configuration")
		return nil
	}

	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load Modix configuration: %w", err)
	}

	color.Cyan("=== Modix Configuration ===")
	fmt.Println()

	// Pretty print JSON with syntax highlighting
	configMap := map[string]interface{}{
		"current_vendor": modixConfig.CurrentVendor,
		"current_model":  modixConfig.CurrentModel,
		"default_vendor": modixConfig.DefaultVendor,
		"default_model":  modixConfig.DefaultModel,
		"vendors":        modixConfig.Vendors,
		"config_version": modixConfig.ConfigVersion,
		"created_at":     modixConfig.CreatedAt,
		"updated_at":     modixConfig.UpdatedAt,
	}

	prettyJSON, err := json.MarshalIndent(configMap, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to format configuration: %w", err)
	}

	fmt.Println(string(prettyJSON))

	// Show configuration summary
	fmt.Println()
	color.Cyan("--- Configuration Summary ---")

	totalVendors, totalModels, configuredVendors, currentModelInfo := modixConfig.GetConfigStatus()
	color.Cyan("Total vendors: %s\n", totalVendors)
	color.Cyan("Total models: %s\n", totalModels)
	color.Cyan("Configured models: %s\n", configuredVendors)
	color.Cyan("Current selection: %s\n", currentModelInfo)

	// Check for common configuration issues
	fmt.Println()
	color.Cyan("--- Configuration Health Check ---")
	checkConfigHealth(modixConfig)

	return nil
}

func checkConfigHealth(config *config.ModixConfig) {
	var issues []string

	// Check if vendors section exists
	if len(config.Vendors) == 0 {
		issues = append(issues, "Missing vendors section")
	}

	// Check for vendors with empty endpoints or keys
	for vendorName, vendorConfig := range config.Vendors {
		// Skip health check for Anthropic since it's pre-configured in Claude-code
		if strings.ToLower(vendorName) == "anthropic" {
			continue
		}

		if vendorConfig.APIEndpoint == "" {
			issues = append(issues, fmt.Sprintf("Vendor '%s' has empty API endpoint", vendorName))
		}

		if vendorConfig.APIKey == "" {
			issues = append(issues, fmt.Sprintf("Vendor '%s' has empty API key", vendorName))
		}
	}

	// Check current model selection
	if currentModel, _, exists := config.GetCurrentModel(); exists {
		found := false
		for _, model := range config.GetVendorModels(config.CurrentVendor) {
			if model == *currentModel {
				found = true
				break
			}
		}
		if !found {
			issues = append(issues, fmt.Sprintf("Current model '%s' not found in vendor '%s' models", *currentModel, config.CurrentVendor))
		}
	} else {
		issues = append(issues, "Current vendor not found in configuration")
	}

	// Report issues
	if len(issues) == 0 {
		color.Green("Health Check: All checks passed! 🎉")
	} else {
		color.Red("Health Check: Found %d issue(s)", len(issues))
		for _, issue := range issues {
			color.Red("  - %s", issue)
		}
	}
}

func init() {
	// Command is already registered in root.go
}
