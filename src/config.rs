use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;

/// Model configuration structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ModelConfig {
    /// Company name
    pub company: String,
    /// API endpoint URL
    pub api_endpoint: String,
    /// API key for authentication
    pub api_key: String,
    /// List of model names available for this vendor
    pub models: Vec<String>,
}

/// Main configuration structure for Modix
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ModixConfig {
    /// Currently active vendor identifier
    pub current_vendor: String,
    /// Currently active model name
    pub current_model: String,
    /// Default vendor to use
    pub default_vendor: String,
    /// Default model name to use if current is invalid
    pub default_model: String,
    /// Map of vendor names to their configurations
    pub vendors: HashMap<String, ModelConfig>,
    /// Configuration version for compatibility checking
    pub config_version: String,
    /// Creation timestamp
    pub created_at: Option<String>,
    /// Last update timestamp
    pub updated_at: Option<String>,
}

impl Default for ModixConfig {
    fn default() -> Self {
        Self {
            current_vendor: "anthropic".to_string(),
            current_model: "Claude Code".to_string(),
            default_vendor: "anthropic".to_string(),
            default_model: "Claude Code".to_string(),
            vendors: HashMap::new(),
            config_version: "1.0.0".to_string(),
            created_at: None,
            updated_at: None,
        }
    }
}

impl ModixConfig {
    /// Create a new empty configuration
    pub fn new() -> Self {
        Self::default()
    }

    /// Add or update a model configuration for a vendor
    pub fn add_vendor(&mut self, vendor: String, config: ModelConfig) {
        self.vendors.insert(vendor, config);
    }

    /// Add a model to a vendor's model list
    pub fn add_model_to_vendor(&mut self, vendor: &str, model_name: String) -> Result<(), String> {
        if let Some(config) = self.vendors.get_mut(vendor) {
            if !config.models.contains(&model_name) {
                config.models.push(model_name);
            }
            Ok(())
        } else {
            Err(format!("Vendor '{}' not found", vendor))
        }
    }

    /// Remove a vendor configuration
    pub fn remove_vendor(&mut self, vendor: &str) -> Option<ModelConfig> {
        self.vendors.remove(vendor)
    }

    /// Get a model configuration by vendor and model name
    pub fn get_model(&self, vendor: &str, model_name: &str) -> Option<&ModelConfig> {
        self.vendors.get(vendor).filter(|config| config.models.contains(&model_name.to_string()))
    }

    /// Get a vendor configuration
    pub fn get_vendor(&self, vendor: &str) -> Option<&ModelConfig> {
        self.vendors.get(vendor)
    }

    /// Get a mutable reference to a vendor configuration
    pub fn get_vendor_mut(&mut self, vendor: &str) -> Option<&mut ModelConfig> {
        self.vendors.get_mut(vendor)
    }

    /// Set the current active vendor and model
    pub fn set_current_vendor_and_model(&mut self, vendor: &str, model_name: &str) -> Result<(), String> {
        if let Some(config) = self.vendors.get(vendor) {
            if config.models.contains(&model_name.to_string()) {
                self.current_vendor = vendor.to_string();
                self.current_model = model_name.to_string();
                Ok(())
            } else {
                Err(format!("Model '{}' not found for vendor '{}'", model_name, vendor))
            }
        } else {
            Err(format!("Vendor '{}' not found in configuration", vendor))
        }
    }

    /// Get the current active vendor and model configuration
    pub fn get_current_model(&self) -> Option<(&String, &ModelConfig)> {
        self.vendors.get(&self.current_vendor).map(|config| (&self.current_model, config))
    }

    /// Get the current vendor configuration
    pub fn get_current_vendor(&self) -> Option<&ModelConfig> {
        self.vendors.get(&self.current_vendor)
    }

    /// Get all vendors and their configurations
    pub fn get_all_vendors(&self) -> Vec<(&String, &ModelConfig)> {
        self.vendors.iter().collect()
    }

    /// Get all models across all vendors
    pub fn get_all_models(&self) -> Vec<(String, String, &ModelConfig)> {
        self.vendors
            .iter()
            .flat_map(|(vendor, config)| {
                config.models.iter().map(move |model_name| (vendor.clone(), model_name.clone(), config))
            })
            .collect()
    }

    /// Validate the configuration
    pub fn validate(&self) -> Result<(), String> {
        if self.vendors.is_empty() {
            return Err("No vendors configured".to_string());
        }

        if !self.vendors.contains_key(&self.current_vendor) {
            return Err(format!(
                "Current vendor '{}' not found in configuration",
                self.current_vendor
            ));
        }

        if let Some(current_config) = self.vendors.get(&self.current_vendor) {
            if !current_config.models.contains(&self.current_model) {
                return Err(format!(
                    "Current model '{}' not available for vendor '{}'",
                    self.current_model, self.current_vendor
                ));
            }
        }

        if !self.vendors.contains_key(&self.default_vendor) {
            return Err(format!(
                "Default vendor '{}' not found in configuration",
                self.default_vendor
            ));
        }

        if let Some(default_config) = self.vendors.get(&self.default_vendor) {
            if !default_config.models.contains(&self.default_model) {
                return Err(format!(
                    "Default model '{}' not available for vendor '{}'",
                    self.default_model, self.default_vendor
                ));
            }
        }

        Ok(())
    }
}

/// Get the default configuration file path (~/.modix/settings.json)
pub fn get_config_path() -> PathBuf {
    #[cfg(target_os = "windows")]
    {
        // On Windows, use %APPDATA%\modix\settings.json
        if let Some(app_data) = std::env::var("APPDATA")
            .ok()
            .and_then(|p| Some(PathBuf::from(p)))
        {
            app_data.join("modix").join("settings.json")
        } else {
            // Fallback to user profile
            dirs_next::home_dir()
                .unwrap_or_else(|| PathBuf::from("."))
                .join("modix")
                .join("settings.json")
        }
    }

    #[cfg(not(target_os = "windows"))]
    {
        // On Unix-like systems (macOS, Linux), use ~/.modix/settings.json
        dirs_next::home_dir()
            .unwrap_or_else(|| PathBuf::from("."))
            .join(".modix")
            .join("settings.json")
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_config_creation() {
        let config = ModixConfig::new();
        assert_eq!(config.current_vendor, "anthropic");
        assert_eq!(config.vendors.len(), 0);
    }

    #[test]
    fn test_add_vendor_and_model() {
        let mut config = ModixConfig::new();
        let model_config = ModelConfig {
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
            models: vec![],
        };

        config.add_vendor("test-vendor".to_string(), model_config);
        assert_eq!(config.vendors.len(), 1);
        assert!(config.get_vendor("test-vendor").is_some());
    }

    #[test]
    fn test_add_model_to_vendor() {
        let mut config = ModixConfig::new();
        let model_config = ModelConfig {
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
            models: vec![],
        };

        config.add_vendor("test-vendor".to_string(), model_config);
        assert!(config.add_model_to_vendor("test-vendor", "test-model".to_string()).is_ok());
        assert_eq!(config.vendors.get("test-vendor").unwrap().models.len(), 1);
    }

    #[test]
    fn test_set_current_vendor_and_model() {
        let mut config = ModixConfig::new();
        let model_config = ModelConfig {
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
            models: vec!["test-model".to_string()],
        };

        config.add_vendor("test-vendor".to_string(), model_config);
        assert!(config.set_current_vendor_and_model("test-vendor", "test-model").is_ok());
        assert_eq!(config.current_vendor, "test-vendor");
        assert_eq!(config.current_model, "test-model");
    }

    #[test]
    fn test_set_current_vendor_and_model_invalid() {
        let mut config = ModixConfig::new();
        assert!(config.set_current_vendor_and_model("nonexistent", "model").is_err());
    }

    #[test]
    fn test_get_all_models() {
        let mut config = ModixConfig::new();
        let model_config1 = ModelConfig {
            company: "TestCorp1".to_string(),
            api_endpoint: "https://api.test1.com".to_string(),
            api_key: "test-key-1".to_string(),
            models: vec!["model-1".to_string(), "model-2".to_string()],
        };
        let model_config2 = ModelConfig {
            company: "TestCorp2".to_string(),
            api_endpoint: "https://api.test2.com".to_string(),
            api_key: "test-key-2".to_string(),
            models: vec!["model-3".to_string()],
        };

        config.add_vendor("vendor1".to_string(), model_config1);
        config.add_vendor("vendor2".to_string(), model_config2);

        let all_models = config.get_all_models();
        assert_eq!(all_models.len(), 3);

        // Check that all models are present without depending on order
        let model_keys: std::collections::HashSet<_> = all_models.iter().map(|(vendor, model, _)| format!("{}:{}", vendor, model)).collect();
        assert!(model_keys.contains("vendor1:model-1"));
        assert!(model_keys.contains("vendor1:model-2"));
        assert!(model_keys.contains("vendor2:model-3"));
    }
}
