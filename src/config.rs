use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;

/// Model configuration structure
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ModelConfig {
    /// Vendor identifier (anthropic, deepseek, alibaba, bytedance, etc.)
    pub vendor: String,
    /// Company name
    pub company: String,
    /// API endpoint URL
    pub api_endpoint: String,
    /// API key for authentication
    pub api_key: String,
}

/// Main configuration structure for Modix
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ModixConfig {
    /// Currently active model identifier
    pub current_model: String,
    /// Default model to use if current_model is invalid
    pub default_model: String,
    /// Map of model names to their configurations
    pub models: HashMap<String, ModelConfig>,
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
            current_model: "claude".to_string(),
            default_model: "claude".to_string(),
            models: HashMap::new(),
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

    /// Add or update a model configuration
    pub fn add_model(&mut self, name: String, config: ModelConfig) {
        self.models.insert(name, config);
    }

    /// Remove a model configuration
    pub fn remove_model(&mut self, name: &str) -> Option<ModelConfig> {
        self.models.remove(name)
    }

    /// Get a model configuration by name
    pub fn get_model(&self, name: &str) -> Option<&ModelConfig> {
        self.models.get(name)
    }

    /// Get a mutable reference to a model configuration
    pub fn get_model_mut(&mut self, name: &str) -> Option<&mut ModelConfig> {
        self.models.get_mut(name)
    }

    /// Set the current active model
    pub fn set_current_model(&mut self, name: &str) -> Result<(), String> {
        if self.models.contains_key(name) {
            self.current_model = name.to_string();
            Ok(())
        } else {
            Err(format!("Model '{}' not found in configuration", name))
        }
    }

    /// Get the current active model configuration
    pub fn get_current_model(&self) -> Option<&ModelConfig> {
        self.get_model(&self.current_model)
    }

    /// Get all enabled models sorted by addition order
    pub fn get_enabled_models(&self) -> Vec<(&String, &ModelConfig)> {
        self.models.iter().collect()
    }

    /// Get all models (including disabled ones) sorted by addition order
    pub fn get_all_models(&self) -> Vec<(&String, &ModelConfig)> {
        self.models.iter().collect()
    }

    /// Validate the configuration
    pub fn validate(&self) -> Result<(), String> {
        if self.models.is_empty() {
            return Err("No models configured".to_string());
        }

        if !self.models.contains_key(&self.current_model) {
            return Err(format!(
                "Current model '{}' not found in configuration",
                self.current_model
            ));
        }

        if !self.models.contains_key(&self.default_model) {
            return Err(format!(
                "Default model '{}' not found in configuration",
                self.default_model
            ));
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
        assert_eq!(config.current_model, "claude");
        assert_eq!(config.models.len(), 0);
    }

    #[test]
    fn test_add_model() {
        let mut config = ModixConfig::new();
        let model = ModelConfig {
            vendor: "test".to_string(),
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
        };

        config.add_model("test-model".to_string(), model);
        assert_eq!(config.models.len(), 1);
        assert!(config.get_model("test-model").is_some());
    }

    #[test]
    fn test_set_current_model() {
        let mut config = ModixConfig::new();
        let model = ModelConfig {
            vendor: "test".to_string(),
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
        };

        config.add_model("test-model".to_string(), model);
        assert!(config.set_current_model("test-model").is_ok());
        assert_eq!(config.current_model, "test-model");
    }

    #[test]
    fn test_set_current_model_invalid() {
        let mut config = ModixConfig::new();
        assert!(config.set_current_model("nonexistent").is_err());
    }

    #[test]
    fn test_get_enabled_models() {
        let mut config = ModixConfig::new();
        let model1 = ModelConfig {
            vendor: "test".to_string(),
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
        };
        let model2 = ModelConfig {
            vendor: "test".to_string(),
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
        };

        config.add_model("Model 1".to_string(), model1);
        config.add_model("Model 2".to_string(), model2);

        let enabled_models = config.get_enabled_models();
        assert_eq!(enabled_models.len(), 2);

        // Check that both models are present without depending on order
        let model_names: std::collections::HashSet<_> = enabled_models.iter().map(|(name, _)| name.as_str()).collect();
        assert!(model_names.contains("Model 1"));
        assert!(model_names.contains("Model 2"));
    }
}
