package agent

// cSpell:ignore modix promacanthus Claude Code Anthropic Gemini CLI OpenAI Codex gemini claude codex streamlake KAT-Coder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/promacanthus/modix/internal/config"
	"github.com/spf13/cobra"
)

// checkCmd checks agent configuration
var checkCmd = &cobra.Command{
	Use:   "check [agent-name]",
	Short: "Check agent configuration",
	Long: `Check the configuration and status of a coding agent.

Examples:
  modix agent check claude-code`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		return checkAgent(agentName)
	},
}

// checkAgent checks the configuration and status of a coding agent
func checkAgent(agentName string) error {
	agentInfo, supported := supportedAgents[agentName]
	if !supported {
		return fmt.Errorf("unsupported agent '%s'", agentName)
	}

	fmt.Printf("Checking %s configuration...\n", agentInfo.Name)
	fmt.Println()

	// Load modix config
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if agent is configured in modix
	if modixConfig.Agents == nil {
		fmt.Printf("✗ Agent '%s' not configured in modix\n", agentName)
		return nil
	}

	agentConfig, exists := modixConfig.Agents[agentName]
	if !exists {
		fmt.Printf("✗ Agent '%s' not configured in modix\n", agentName)
		return nil
	}

	fmt.Printf("✓ Agent '%s' configured in modix\n", agentName)
	fmt.Printf("  Provider: %s\n", agentConfig.Provider)
	fmt.Printf("  Enabled: %v\n", agentConfig.Enabled)

	// Check if agent config file exists
	expandedPath := expandPath(agentConfig.ConfigPath)
	if _, err := os.Stat(expandedPath); err == nil {
		fmt.Printf("✓ Config file exists: %s\n", expandedPath)
	} else {
		fmt.Printf("✗ Config file not found: %s\n", expandedPath)
	}

	// Check current LLM model
	if modixConfig.CurrentModel != "" {
		fmt.Printf("✓ Current LLM model: %s@%s\n", modixConfig.CurrentModel, modixConfig.CurrentVendor)
	} else {
		fmt.Println("✗ No LLM model selected")
	}

	// Agent-specific checks
	switch agentName {
	case "claude-code":
		checkClaudeCodeSpecific()
	}

	return nil
}

// checkClaudeCodeSpecific performs Claude Code specific checks
func checkClaudeCodeSpecific() {
	// Check if Claude Code is installed
	homeDir, err := os.UserHomeDir()
	if err == nil {
		claudePath := filepath.Join(homeDir, ".claude")
		if _, err := os.Stat(claudePath); err == nil {
			fmt.Println("✓ Claude Code installation detected")
		} else {
			fmt.Println("⚠ Claude Code installation not found")
		}
	}

	// Check environment variables
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		fmt.Println("✓ ANTHROPIC_API_KEY environment variable set")
	} else {
		fmt.Println("⚠ ANTHROPIC_API_KEY environment variable not set")
	}
}