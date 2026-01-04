package config

import "time"

// SetupDefaultModels initializes default models for common vendors
func SetupDefaultModels() *ModixConfig {
	config := NewModixConfig()

	// Anthropic
	anthropic := VendorConfig{
		Company: "Anthropic",
		Models:  []string{"Claude"},
	}
	config.AddVendor("anthropic", anthropic)

	// DeepSeek
	deepseek := VendorConfig{
		Company:     "DeepSeek",
		APIEndpoint: "https://api.deepseek.com/v1",
		Models:      []string{"deepseek-reasoner", "deepseek-chat"},
	}
	config.AddVendor("deepseek", deepseek)

	// Alibaba
	bailian := VendorConfig{
		Company:     "Alibaba",
		APIEndpoint: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Models:      []string{"qwen3-coder-plus", "qwen3-coder-flash"},
	}
	config.AddVendor("bailian", bailian)

	// ByteDance
	volcengine := VendorConfig{
		Company:     "ByteDance",
		APIEndpoint: "https://ark.cn-beijing.volces.com/api/coding",
		Models:      []string{"doubao-seed-code-preview-latest"},
	}
	config.AddVendor("volcengine", volcengine)

	// Moonshot AI
	moonshoot := VendorConfig{
		Company:     "Moonshot AI",
		APIEndpoint: "https://api.moonshot.cn/anthropic",
		Models:      []string{"kimi-k2-thinking-turbo"},
	}
	config.AddVendor("moonshot", moonshoot)

	// Kuaishou
	streamlake := VendorConfig{
		Company:     "Kuaishou",
		APIEndpoint: "https://wanqing.streamlakeapi.com/api/gateway/v1/endpoints/ep-xxx-xxx/claude-code-proxy",
		Models:      []string{"KAT-Coder"},
	}
	config.AddVendor("streamlake", streamlake)

	// MiniMax
	minimax := VendorConfig{
		Company:     "MiniMax",
		APIEndpoint: "https://api.minimaxi.com/anthropic",
		Models:      []string{"MiniMax-M2"},
	}
	config.AddVendor("minimax", minimax)

	// ZHIPU AI
	bigmodel := VendorConfig{
		Company:     "ZHIPU AI",
		APIEndpoint: "https://open.bigmodel.cn/api/anthropic",
		Models:      []string{"GLM-4.6"},
	}
	config.AddVendor("bigmodel", bigmodel)

	// XiaoMi
	xiaomi := VendorConfig{
		Company:     "Xiaomi",
		APIEndpoint: "https://api.xiaomimimo.com/anthropic",
		Models:      []string{"mimo-v2-flash"},
	}
	config.AddVendor("xiaomi", xiaomi)

	// Set current and default vendor/model to the first enabled model
	config.CurrentVendor = "anthropic"
	config.CurrentModel = "Claude"
	config.DefaultVendor = "anthropic"
	config.DefaultModel = "Claude"

	return config
}

// GetDefaultModels returns the list of default models that would be created
func GetDefaultModels() map[string]VendorConfig {
	defaultConfig := SetupDefaultModels()
	return defaultConfig.Vendors
}

// ModelInfo represents information about a model
type ModelInfo struct {
	Vendor      string
	Company     string
	Endpoint    string
	Models      []string
	HasAPIKey   bool
	HasEndpoint bool
}

// GetModelInfo returns detailed information about a model
func (c *ModixConfig) GetModelInfo(vendor, modelName string) (*ModelInfo, bool) {
	if vendorConfig, exists := c.GetVendor(vendor); exists {
		for _, model := range vendorConfig.Models {
			if model == modelName {
				return &ModelInfo{
					Vendor:      vendor,
					Company:     vendorConfig.Company,
					Endpoint:    vendorConfig.APIEndpoint,
					Models:      vendorConfig.Models,
					HasAPIKey:   vendorConfig.APIKey != "",
					HasEndpoint: vendorConfig.APIEndpoint != "",
				}, true
			}
		}
	}
	return nil, false
}

// GetAllModelInfos returns detailed information about all models
func (c *ModixConfig) GetAllModelInfos() []*ModelInfo {
	var modelInfos []*ModelInfo

	for vendor, vendorConfig := range c.Vendors {
		for _, modelName := range vendorConfig.Models {
			modelInfos = append(modelInfos, &ModelInfo{
				Vendor:      vendor,
				Company:     vendorConfig.Company,
				Endpoint:    vendorConfig.APIEndpoint,
				Models:      []string{modelName},
				HasAPIKey:   vendorConfig.APIKey != "",
				HasEndpoint: vendorConfig.APIEndpoint != "",
			})
		}
	}

	return modelInfos
}

// IsVendorConfigured checks if a vendor has proper configuration
func (c *ModixConfig) IsVendorConfigured(vendor string) bool {
	if vendorConfig, exists := c.GetVendor(vendor); exists {
		return vendorConfig.APIEndpoint != "" && vendorConfig.APIKey != ""
	}
	return false
}

// GetConfigStatus returns the current configuration status
func (c *ModixConfig) GetConfigStatus() (totalVendors, totalModels, configuredVendors int, currentModel string) {
	totalVendors = len(c.Vendors)
	currentModel = c.CurrentModel

	for _, vendorConfig := range c.Vendors {
		totalModels += len(vendorConfig.Models)
		if vendorConfig.APIEndpoint != "" && vendorConfig.APIKey != "" {
			configuredVendors++
		}
	}

	return
}

// GetVendorModels returns all models for a specific vendor
func (c *ModixConfig) GetVendorModels(vendor string) []string {
	if vendorConfig, exists := c.GetVendor(vendor); exists {
		return vendorConfig.Models
	}
	return []string{}
}

// GetCurrentVendorModels returns all models for the current vendor
func (c *ModixConfig) GetCurrentVendorModels() []string {
	return c.GetVendorModels(c.CurrentVendor)
}

// GetLastUpdated returns the last update timestamp
func (c *ModixConfig) GetLastUpdated() time.Time {
	return c.UpdatedAt
}

// GetCreatedAt returns the creation timestamp
func (c *ModixConfig) GetCreatedAt() time.Time {
	return c.CreatedAt
}
