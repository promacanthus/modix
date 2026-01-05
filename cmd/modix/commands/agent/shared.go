package agent

// cSpell:ignore modix promacanthus Claude Code Anthropic Gemini CLI OpenAI Codex gemini claude codex streamlake KAT-Coder

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/promacanthus/modix/internal/config"
)

// Supported agents
var supportedAgents = map[string]AgentInfo{
	"claude-code": {
		Name:        "Claude Code",
		Provider:    "Anthropic",
		ConfigPath:  "~/.claude/settings.json",
		Description: "Anthropic's Claude Code assistant",
	},
	"gemini-cli": {
		Name:        "Gemini CLI",
		Provider:    "Google",
		ConfigPath:  "~/.config/gemini-cli/settings.json",
		Description: "Google's Gemini CLI (coming soon)",
	},
	"codex": {
		Name:        "Codex",
		Provider:    "OpenAI",
		ConfigPath:  "~/.codex/config.json",
		Description: "OpenAI's Codex (coming soon)",
	},
}

// AgentInfo represents information about a coding agent
type AgentInfo struct {
	Name        string
	Provider    string
	ConfigPath  string
	Description string
}

// Implementation functions
func addAgent(agentName string) error {
	agentInfo, supported := supportedAgents[agentName]
	if !supported {
		names := make([]string, 0, len(supportedAgents))
		for name := range supportedAgents {
			names = append(names, name)
		}
		return fmt.Errorf("unsupported agent '%s'. Supported agents: %v", agentName, names)
	}

	// Load modix config
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check if agent already exists in config
	if modixConfig.Agents == nil {
		modixConfig.Agents = make(map[string]config.AgentConfig)
	}

	if _, exists := modixConfig.Agents[agentName]; exists {
		return fmt.Errorf("agent '%s' already configured", agentName)
	}

	// Add agent configuration
	modixConfig.Agents[agentName] = config.AgentConfig{
		Name:        agentInfo.Name,
		Provider:    agentInfo.Provider,
		ConfigPath:  agentInfo.ConfigPath,
		Enabled:     true,
		Description: agentInfo.Description,
	}

	// Save configuration
	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Added agent '%s' (%s)\n", agentName, agentInfo.Name)
	fmt.Printf("Configure it with: modix agent config %s\n", agentName)
	return nil
}

func removeAgent(agentName string) error {
	modixConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	if modixConfig.Agents == nil {
		return fmt.Errorf("no agents configured")
	}

	if _, exists := modixConfig.Agents[agentName]; !exists {
		return fmt.Errorf("agent '%s' not found", agentName)
	}

	delete(modixConfig.Agents, agentName)

	if err := config.SaveConfig(modixConfig); err != nil {
		return fmt.Errorf("failed to save configuration: %w", err)
	}

	fmt.Printf("Removed agent: %s\n", agentName)
	return nil
}

func configAgent(agentName string) error {
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

	// Get current LLM model
	if modixConfig.CurrentModel == "" {
		return fmt.Errorf("no LLM model selected. Use 'modix llm model switch <model>' first")
	}

	// Configure based on agent type
	switch agentName {
	case "claude-code":
		return configureClaudeCode(modixConfig, agentConfig)
	case "gemini-cli":
		return fmt.Errorf("Gemini CLI configuration not yet implemented")
	case "codex":
		return fmt.Errorf("Codex configuration not yet implemented")
	default:
		return fmt.Errorf("configuration for '%s' not implemented", agentName)
	}
}

// Helper functions
func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(homeDir, strings.TrimPrefix(path, "~/"))
	}
	return path
}

func configureClaudeCode(modixConfig *config.ModixConfig, agentConfig config.AgentConfig) error {
	// This uses the existing UpdateClaudeEnvConfig function
	if modixConfig.CurrentModel == "" || modixConfig.CurrentVendor == "" {
		return fmt.Errorf("no current model/vendor set")
	}

	err := config.UpdateClaudeEnvConfig(modixConfig.CurrentModel, modixConfig.CurrentVendor)
	if err != nil {
		return fmt.Errorf("failed to configure Claude Code: %w", err)
	}

	fmt.Printf("Configured Claude Code to use: %s@%s\n", modixConfig.CurrentModel, modixConfig.CurrentVendor)
	fmt.Printf("Config file: %s\n", agentConfig.ConfigPath)
	return nil
}