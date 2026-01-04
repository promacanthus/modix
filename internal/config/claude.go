package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ClaudeConfig represents the Claude Code configuration
type ClaudeConfig map[string]interface{}

// LoadClaudeConfig loads the Claude Code configuration from ~/.claude/settings.json
func LoadClaudeConfig() (ClaudeConfig, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("could not determine home directory: %w", err)
	}

	claudeConfigPath := filepath.Join(homeDir, ".claude", "settings.json")

	// If file doesn't exist, return empty config
	if _, err := os.Stat(claudeConfigPath); os.IsNotExist(err) {
		return ClaudeConfig{}, nil
	}

	// Read the file
	content, err := os.ReadFile(claudeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Claude configuration file: %w", err)
	}

	// If file is empty, return empty config
	if len(content) == 0 {
		return ClaudeConfig{}, nil
	}

	// Parse JSON
	var config ClaudeConfig
	if err := json.Unmarshal(content, &config); err != nil {
		return nil, fmt.Errorf("failed to parse Claude configuration JSON: %w", err)
	}

	return config, nil
}

// SaveClaudeConfig saves the Claude Code configuration to ~/.claude/settings.json
func SaveClaudeConfig(config ClaudeConfig) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not determine home directory: %w", err)
	}

	claudeConfigPath := filepath.Join(homeDir, ".claude", "settings.json")

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(claudeConfigPath), 0755); err != nil {
		return fmt.Errorf("failed to create Claude configuration directory: %w", err)
	}

	// Marshal to JSON with indentation
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize Claude configuration to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(claudeConfigPath, configJSON, 0600); err != nil {
		return fmt.Errorf("failed to write Claude configuration file: %w", err)
	}

	return nil
}

// SwitchToClaudeOfficial switches Claude configuration to use official backend (removes env field)
func SwitchToClaudeOfficial() error {
	config, err := LoadClaudeConfig()
	if err != nil {
		return err
	}

	// Remove the env field to switch back to official Claude
	delete(config, "env")

	return SaveClaudeConfig(config)
}

// SwitchToClaudeEnv switches Claude configuration to use specified environment
func SwitchToClaudeEnv(env string) error {
	config, err := LoadClaudeConfig()
	if err != nil {
		return err
	}

	// Set the env field
	config["env"] = env

	return SaveClaudeConfig(config)
}

// UpdateClaudeEnvConfig updates Claude configuration based on the selected model
func UpdateClaudeEnvConfig(modelName, vendor string) error {
	// Load the current Modix configuration to get model details
	config, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load Modix configuration: %w", err)
	}

	// Get the model configuration
	modelConfig, exists := config.GetModel(vendor, modelName)
	if !exists {
		return fmt.Errorf("model '%s' not found in configuration for vendor '%s'", modelName, vendor)
	}

	// Load existing Claude configuration
	claudeConfig, err := LoadClaudeConfig()
	if err != nil {
		return fmt.Errorf("failed to load Claude configuration: %w", err)
	}

	// If the config is empty, initialize with basic structure
	if len(claudeConfig) == 0 {
		claudeConfig["companyAnnouncements"] = []string{"Welcome to Claude code, the configuration managed by modix."}
	}

	// Check if the model is Claude (case-insensitive check for various Claude model names)
	isClaudeModel := containsIgnoreCase(modelName, "claude") || containsIgnoreCase(vendor, "anthropic")

	if isClaudeModel {
		// For Claude models, remove the env field to use official Claude API
		delete(claudeConfig, "env")
	} else {
		// For non-Claude models, set up the detailed env configuration
		envConfig := map[string]interface{}{
			"ANTHROPIC_BASE_URL":                       modelConfig.APIEndpoint,
			"ANTHROPIC_AUTH_TOKEN":                     modelConfig.APIKey,
			"API_TIMEOUT_MS":                           "3000000",
			"CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1,
			"ANTHROPIC_MODEL":                          modelName,
			"ANTHROPIC_SMALL_FAST_MODEL":               modelName,
			"ANTHROPIC_DEFAULT_SONNET_MODEL":           modelName,
			"ANTHROPIC_DEFAULT_OPUS_MODEL":             modelName,
			"ANTHROPIC_DEFAULT_HAIKU_MODEL":            modelName,
		}

		// Update the env field with the new configuration
		claudeConfig["env"] = envConfig
	}

	// Save the updated Claude configuration
	return SaveClaudeConfig(claudeConfig)
}

// Helper function to check if a string contains a substring (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	return containsSubstringIgnoreCase(s, substr)
}

func containsSubstringIgnoreCase(s, substr string) bool {
	s = toLower(s)
	substr = toLower(substr)
	return indexOf(s, substr) >= 0
}

func toLower(s string) string {
	var result []byte
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			result = append(result, c+32) // Convert to lowercase
		} else {
			result = append(result, c)
		}
	}
	return string(result)
}

func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

// GetClaudeConfigPath returns the path to the Claude configuration file
func GetClaudeConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "~/.claude/settings.json"
	}

	return filepath.Join(homeDir, ".claude", "settings.json")
}

// IsClaudeConfigured checks if Claude configuration exists
func IsClaudeConfigured() bool {
	path := GetClaudeConfigPath()
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetClaudeConfigSummary returns a summary of the Claude configuration
func GetClaudeConfigSummary() (map[string]interface{}, error) {
	config, err := LoadClaudeConfig()
	if err != nil {
		return nil, err
	}

	summary := make(map[string]interface{})
	for key, value := range config {
		summary[key] = value
	}
	return summary, nil
}
