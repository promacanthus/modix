package config

import (
	"testing"
	"time"
)

func TestNewModixConfig(t *testing.T) {
	config := NewModixConfig()

	if config.CurrentVendor != "anthropic" {
		t.Errorf("Expected CurrentVendor to be 'anthropic', got '%s'", config.CurrentVendor)
	}

	if config.CurrentModel != "Claude" {
		t.Errorf("Expected CurrentModel to be 'Claude', got '%s'", config.CurrentModel)
	}

	if config.Vendors == nil {
		t.Error("Expected Vendors map to be initialized")
	}

	if config.Agents == nil {
		t.Error("Expected Agents map to be initialized")
	}

	if config.ConfigVersion != "1.0.0" {
		t.Errorf("Expected ConfigVersion to be '1.0.0', got '%s'", config.ConfigVersion)
	}
}

func TestModixConfig_AddVendor(t *testing.T) {
	config := NewModixConfig()

	vendorConfig := VendorConfig{
		Company:     "TestCorp",
		APIEndpoint: "https://api.test.com",
		APIKey:      "test-key",
		Models:      []string{"test-model"},
	}

	config.AddVendor("test", vendorConfig)

	if _, exists := config.Vendors["test"]; !exists {
		t.Error("Vendor not added correctly")
	}

	if config.Vendors["test"].Company != "TestCorp" {
		t.Error("Vendor company not set correctly")
	}
}

func TestModixConfig_AddModelToVendor(t *testing.T) {
	config := NewModixConfig()

	// Add a vendor first
	vendorConfig := VendorConfig{
		Company:     "TestCorp",
		APIEndpoint: "https://api.test.com",
		APIKey:      "test-key",
		Models:      []string{"model1"},
	}
	config.AddVendor("test", vendorConfig)

	// Add another model
	err := config.AddModelToVendor("test", "model2")
	if err != nil {
		t.Errorf("AddModelToVendor failed: %v", err)
	}

	if len(config.Vendors["test"].Models) != 2 {
		t.Errorf("Expected 2 models, got %d", len(config.Vendors["test"].Models))
	}

	// Try to add duplicate
	err = config.AddModelToVendor("test", "model1")
	if err != nil {
		t.Errorf("Adding duplicate should not error: %v", err)
	}

	// Try to add to non-existent vendor
	err = config.AddModelToVendor("nonexistent", "model")
	if err == nil {
		t.Error("Expected error when adding to non-existent vendor")
	}
}

func TestModixConfig_GetModels(t *testing.T) {
	config := NewModixConfig()

	config.AddVendor("vendor1", VendorConfig{Models: []string{"model1", "model2"}})
	config.AddVendor("vendor2", VendorConfig{Models: []string{"model3"}})

	models := config.GetModels()

	if len(models) != 3 {
		t.Errorf("Expected 3 models, got %d", len(models))
	}
}

func TestAgentConfig(t *testing.T) {
	config := NewModixConfig()

	agentConfig := AgentConfig{
		Name:        "Claude Code",
		Provider:    "Anthropic",
		ConfigPath:  "~/.claude/settings.json",
		Enabled:     true,
		Description: "Anthropic's Claude Code assistant",
	}

	if config.Agents == nil {
		t.Error("Agents map should be initialized")
	}

	config.Agents["claude-code"] = agentConfig

	if config.Agents["claude-code"].Name != "Claude Code" {
		t.Error("Agent config not set correctly")
	}

	if config.CurrentAgent != "" {
		t.Error("CurrentAgent should be empty initially")
	}

	config.CurrentAgent = "claude-code"
	if config.CurrentAgent != "claude-code" {
		t.Error("CurrentAgent not set correctly")
	}
}

func TestModixConfig_Validate(t *testing.T) {
	// Test empty config
	config := &ModixConfig{}
	err := config.Validate()
	if err == nil {
		t.Error("Empty config should not validate")
	}

	// Test valid config
	validConfig := NewModixConfig()
	validConfig.AddVendor("test", VendorConfig{
		Company:     "Test",
		APIEndpoint: "https://api.test.com",
		APIKey:      "key",
		Models:      []string{"model"},
	})
	validConfig.CurrentVendor = "test"
	validConfig.CurrentModel = "model"
	validConfig.DefaultVendor = "test"
	validConfig.DefaultModel = "model"

	err = validConfig.Validate()
	if err != nil {
		t.Errorf("Valid config should validate: %v", err)
	}
}

func TestTimestamps(t *testing.T) {
	config := NewModixConfig()

	// Before save, CreatedAt should be zero
	if !config.CreatedAt.IsZero() {
		t.Error("CreatedAt should be zero before save")
	}

	// Simulate save by setting timestamps
	now := time.Now()
	config.CreatedAt = now
	config.UpdatedAt = now

	if config.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero after setting")
	}

	if config.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero after setting")
	}
}
