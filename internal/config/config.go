package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// ModixConfig represents the main configuration structure
type ModixConfig struct {
	CurrentVendor string                  `json:"current_vendor,omitempty"`
	CurrentModel  string                  `json:"current_model,omitempty"`
	DefaultVendor string                  `json:"default_vendor,omitempty"`
	DefaultModel  string                  `json:"default_model,omitempty"`
	CurrentAgent  string                  `json:"current_agent,omitempty"`
	Vendors       map[string]VendorConfig `json:"vendors,omitempty"`
	Agents        map[string]AgentConfig  `json:"agents,omitempty"`
	ConfigVersion string                  `json:"config_version,omitempty"`
	CreatedAt     time.Time               `json:"created_at,omitempty"`
	UpdatedAt     time.Time               `json:"updated_at,omitempty"`
}

// VendorConfig represents the configuration for a single model
type VendorConfig struct {
	Company     string   `json:"company,omitempty"`
	APIEndpoint string   `json:"api_endpoint,omitempty"`
	APIKey      string   `json:"api_key,omitempty"`
	Models      []string `json:"models,omitempty"`
}

// AgentConfig represents the configuration for a coding agent
type AgentConfig struct {
	Name        string `json:"name,omitempty"`
	Provider    string `json:"provider,omitempty"`
	ConfigPath  string `json:"config_path,omitempty"`
	Enabled     bool   `json:"enabled,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewModixConfig creates a new default Modix configuration
func NewModixConfig() *ModixConfig {
	return &ModixConfig{
		CurrentVendor: "anthropic",
		CurrentModel:  "Claude",
		DefaultVendor: "anthropic",
		DefaultModel:  "Claude",
		Vendors:       make(map[string]VendorConfig),
		Agents:        make(map[string]AgentConfig),
		ConfigVersion: "1.0.0",
	}
}

// AddVendor adds a new vendor configuration
func (c *ModixConfig) AddVendor(vendor string, config VendorConfig) {
	c.Vendors[vendor] = config
}

// AddModelToVendor adds a model to an existing vendor
func (c *ModixConfig) AddModelToVendor(vendor, modelName string) error {
	if config, exists := c.Vendors[vendor]; exists {
		// Check if model already exists
		for _, existingModel := range config.Models {
			if existingModel == modelName {
				return nil // Model already exists
			}
		}
		config.Models = append(config.Models, modelName)
		c.Vendors[vendor] = config
		return nil
	}
	return fmt.Errorf("vendor '%s' not found", vendor)
}

// RemoveVendor removes a vendor configuration
func (c *ModixConfig) RemoveVendor(vendor string) {
	delete(c.Vendors, vendor)
}

// GetModel returns a model configuration by vendor and model name
func (c *ModixConfig) GetModel(vendor, modelName string) (*VendorConfig, bool) {
	if config, exists := c.Vendors[vendor]; exists {
		for _, model := range config.Models {
			if model == modelName {
				return &config, true
			}
		}
	}
	return nil, false
}

// GetVendor returns a vendor configuration
func (c *ModixConfig) GetVendor(vendor string) (*VendorConfig, bool) {
	config, exists := c.Vendors[vendor]
	return &config, exists
}

// SetCurrentVendorAndModel sets the current active vendor and model
func (c *ModixConfig) SetCurrentVendorAndModel(vendor, modelName string) error {
	if config, exists := c.Vendors[vendor]; exists {
		for _, model := range config.Models {
			if model == modelName {
				c.CurrentVendor = vendor
				c.CurrentModel = modelName
				return nil
			}
		}
		return fmt.Errorf("model '%s' not found for vendor '%s'", modelName, vendor)
	}
	return fmt.Errorf("vendor '%s' not found in configuration", vendor)
}

// GetCurrentModel returns the current active vendor and model configuration
func (c *ModixConfig) GetCurrentModel() (*string, *VendorConfig, bool) {
	if config, exists := c.Vendors[c.CurrentVendor]; exists {
		for _, model := range config.Models {
			if model == c.CurrentModel {
				return &c.CurrentModel, &config, true
			}
		}
	}
	return nil, nil, false
}

// GetVendors returns all available vendors
func (c *ModixConfig) GetVendors() map[string]VendorConfig {
	return c.Vendors
}

// GetModels returns all available models across all vendors
func (c *ModixConfig) GetModels() []string {
	var models []string
	for _, vendorConfig := range c.Vendors {
		models = append(models, vendorConfig.Models...)
	}
	return models
}

// GetModelsByVendor returns all available models for a specific vendor
func (c *ModixConfig) GetModelsByVendor(vendor string) []string {
	if config, exists := c.Vendors[vendor]; exists {
		return config.Models
	}
	return []string{}
}

// Validate validates the configuration
func (c *ModixConfig) Validate() error {
	if len(c.Vendors) == 0 {
		return fmt.Errorf("no vendors configured")
	}

	if _, exists := c.Vendors[c.CurrentVendor]; !exists {
		return fmt.Errorf("current vendor '%s' not found in configuration", c.CurrentVendor)
	}

	if config, exists := c.Vendors[c.CurrentVendor]; exists {
		found := false
		for _, model := range config.Models {
			if model == c.CurrentModel {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("current model '%s' not available for vendor '%s'", c.CurrentModel, c.CurrentVendor)
		}
	}

	if _, exists := c.Vendors[c.DefaultVendor]; !exists {
		return fmt.Errorf("default vendor '%s' not found in configuration", c.DefaultVendor)
	}

	if config, exists := c.Vendors[c.DefaultVendor]; exists {
		found := false
		for _, model := range config.Models {
			if model == c.DefaultModel {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("default model '%s' not available for vendor '%s'", c.DefaultModel, c.DefaultVendor)
		}
	}

	return nil
}

// GetConfigPath returns the default configuration file path
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory
		homeDir = "."
	}

	// On Windows, use %APPDATA%\modix\settings.json
	// On Unix-like systems, use ~/.modix/settings.json
	if homeDir == "." {
		return filepath.Join(".", "modix", "settings.json")
	}

	return filepath.Join(homeDir, ".modix", "settings.json")
}

// LoadConfig loads configuration from the default path
func LoadConfig() (*ModixConfig, error) {
	configPath := GetConfigPath()

	// If config file doesn't exist, create default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := NewModixConfig()
		if err := SaveConfig(config); err != nil {
			return nil, fmt.Errorf("failed to create default config: %w", err)
		}
		return config, nil
	}

	// Read the file directly and use custom unmarshaling
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config ModixConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration JSON: %w", err)
	}

	// Validate the loaded configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// LoadConfigFromPath loads configuration from a custom path
func LoadConfigFromPath(path string) (*ModixConfig, error) {
	// Read the file directly and use custom unmarshaling
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %w", err)
	}

	var config ModixConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse configuration JSON: %w", err)
	}

	// Validate the loaded configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// SaveConfig saves configuration to the default path
func SaveConfig(config *ModixConfig) error {
	configPath := GetConfigPath()

	// Update timestamps
	now := time.Now()
	config.UpdatedAt = now
	if config.CreatedAt.IsZero() {
		config.CreatedAt = now
	}

	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return fmt.Errorf("failed to create configuration directory: %w", err)
	}

	// Marshal to JSON with indentation
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize configuration to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, configJSON, 0600); err != nil {
		return fmt.Errorf("failed to write configuration file: %w", err)
	}

	return nil
}

// ConfigExists checks if configuration file exists
func ConfigExists() bool {
	configPath := GetConfigPath()
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

// GetConfigFilePath returns the path to the configuration file
func GetConfigFilePath() string {
	return GetConfigPath()
}

// ResetConfig resets configuration to default values
func ResetConfig() error {
	config := NewModixConfig()
	return SaveConfig(config)
}
