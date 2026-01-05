package agent

import (
	"github.com/spf13/cobra"
)

// AgentCmd represents the agent command family
var AgentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage coding agents",
	Long: `Manage coding agents and their configurations.

This command manages AI coding assistants like Claude Code, Gemini CLI, and Codex.
Each agent can be configured with different LLM backends and settings.

Available subcommands:
  list      List all configured agents
  check     Check agent configuration and status
  switch    Switch to a different agent

Supported agents:
  - claude-code: Anthropic's Claude Code
  - gemini-cli: Google's Gemini CLI (future)
  - codex: OpenAI's Codex (future)

Examples:
  modix agent add claude-code
  modix agent check claude-code
  modix agent list`,
}

func init() {
	AgentCmd.AddCommand(addCmd)
	AgentCmd.AddCommand(removeCmd)
	AgentCmd.AddCommand(listCmd)
	AgentCmd.AddCommand(configCmd)
	AgentCmd.AddCommand(checkCmd)
	AgentCmd.AddCommand(switchCmd)
}
