package agent

// cSpell:ignore modix promacanthus Claude Code Anthropic Gemini CLI OpenAI Codex gemini claude codex streamlake KAT-Coder

import (
	"fmt"
	"os"
	"strings"

	"github.com/promacanthus/modix/internal/config"
	"github.com/spf13/cobra"
)

// addCmd adds a new agent configuration
var addCmd = &cobra.Command{
	Use:   "add [agent-name]",
	Short: "Add a new coding agent",
	Long: `Add a new coding agent configuration.

Supported agents:
  - claude-code: Anthropic's Claude Code
  - gemini-cli: Google's Gemini CLI (future)
  - codex: OpenAI's Codex (future)

Examples:
  modix agent add claude-code
  modix agent add gemini-cli`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		return addAgent(agentName)
	},
}

// removeCmd removes an agent configuration
var removeCmd = &cobra.Command{
	Use:   "remove [agent-name]",
	Short: "Remove an agent configuration",
	Long: `Remove an agent configuration.

This will remove the agent from modix tracking but won't delete
the agent's configuration files.

Examples:
  modix agent remove claude-code`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		return removeAgent(agentName)
	},
}

// configCmd configures an agent
var configCmd = &cobra.Command{
	Use:   "config [agent-name]",
	Short: "Configure an agent",
	Long: `Configure a coding agent.

This command helps set up the agent's configuration to work with
the currently selected LLM model.

Examples:
  modix agent config claude-code`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		return configAgent(agentName)
	},
}

// listCmd lists all configured agents
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	Long:  `List all configured coding agents and their status.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return listAgents()
	},
}

// listAgents lists all configured agents
func listAgents() error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if len(modixConfig.Agents) == 0 {
		fmt.Println("No agents configured")
		fmt.Println("\nSupported agents:")
		for name, info := range supportedAgents {
			fmt.Printf("  - %-15s %s\n", name, info.Description)
		}
		return nil
	}

	fmt.Printf("%-15s %-15s %-15s %-10s %-10s\n", "AGENT", "PROVIDER", "LLM MODEL", "STATUS", "CONFIG")
	fmt.Println(strings.Repeat("-", 80))

	for agentName, agentConfig := range modixConfig.Agents {
		status := "Enabled"
		if !agentConfig.Enabled {
			status = "Disabled"
		}

		// Get current LLM model
		currentModel := "None"
		if modixConfig.CurrentModel != "" {
			currentModel = modixConfig.CurrentModel
		}

		configExists := "No"
		expandedPath := expandPath(agentConfig.ConfigPath)
		if _, err := os.Stat(expandedPath); err == nil {
			configExists = "Yes"
		}

		fmt.Printf("%-15s %-15s %-15s %-10s %-10s\n",
			agentName, agentConfig.Provider, currentModel, status, configExists)
	}

	return nil
}