package agent

// cSpell:ignore modix promacanthus Claude Code Anthropic Gemini CLI OpenAI Codex gemini claude codex streamlake KAT-Coder

import (
	"fmt"

	"github.com/promacanthus/modix/internal/config"
	"github.com/spf13/cobra"
)

// switchCmd switches to a different agent
var switchCmd = &cobra.Command{
	Use:   "switch [agent-name]",
	Short: "Switch to a different agent",
	Long: `Switch to a different coding agent.

Examples:
  modix agent switch claude-code`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		agentName := args[0]
		return switchAgent(agentName)
	},
}

// switchAgent switches to a different agent
func switchAgent(agentName string) error {
	_, supported := supportedAgents[agentName]
	if !supported {
		return fmt.Errorf("unsupported agent '%s'", agentName)
	}

	// Load modix config
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if agent is configured
	if modixConfig.Agents == nil {
		return fmt.Errorf("agent '%s' not configured. Add it first with 'modix agent add %s'", agentName, agentName)
	}

	agentConfig, exists := modixConfig.Agents[agentName]
	if !exists {
		return fmt.Errorf("agent '%s' not configured. Add it first with 'modix agent add %s'", agentName, agentName)
	}

	// Update current agent
	modixConfig.CurrentAgent = agentName
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Switched to agent: %s (%s)\n", agentName, agentConfig.Name)
	return nil
}