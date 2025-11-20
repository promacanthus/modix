use crate::config::{get_config_path, ModixConfig};
use anyhow::{Context, Result};
use serde_json;
use std::fs;
use std::path::{Path, PathBuf};

/// Configuration manager for handling JSON file operations
pub struct ConfigManager;

impl ConfigManager {
    /// Load configuration from the default path (~/.modix/settings.json)
    pub fn load_config() -> Result<ModixConfig> {
        let config_path = get_config_path();

        if !config_path.exists() {
            // Return default configuration if file doesn't exist
            return Ok(ModixConfig::new());
        }

        let config_content =
            fs::read_to_string(&config_path).context("Failed to read configuration file")?;

        if config_content.trim().is_empty() {
            // Return default configuration if file is empty
            return Ok(ModixConfig::new());
        }

        let config: ModixConfig =
            serde_json::from_str(&config_content).context("Failed to parse configuration JSON")?;

        // Validate the loaded configuration
        config
            .validate()
            .map_err(|e| anyhow::anyhow!("Invalid configuration: {}", e))?;

        Ok(config)
    }

    /// Load configuration from the Claude settings file (~/.claude/settings.json)
    pub fn load_claude_config() -> Result<serde_json::Value> {
        let claude_config_path = dirs_next::home_dir()
            .unwrap_or_else(|| PathBuf::from("."))
            .join(".claude")
            .join("settings.json");

        if !claude_config_path.exists() {
            // Return empty JSON object if file doesn't exist
            return Ok(serde_json::json!({}));
        }

        let config_content = fs::read_to_string(&claude_config_path)
            .context("Failed to read Claude configuration file")?;

        if config_content.trim().is_empty() {
            // Return empty JSON object if file is empty
            return Ok(serde_json::json!({}));
        }

        let config: serde_json::Value = serde_json::from_str(&config_content)
            .context("Failed to parse Claude configuration JSON")?;

        Ok(config)
    }

    /// Save configuration to the Claude settings file (~/.claude/settings.json)
    pub fn save_claude_config(config: &serde_json::Value) -> Result<()> {
        let claude_config_path = dirs_next::home_dir()
            .unwrap_or_else(|| PathBuf::from("."))
            .join(".claude")
            .join("settings.json");

        // Create parent directory if it doesn't exist
        if let Some(parent) = claude_config_path.parent() {
            fs::create_dir_all(parent)
                .context("Failed to create Claude configuration directory")?;
        }

        let config_json = serde_json::to_string_pretty(config)
            .context("Failed to serialize Claude configuration to JSON")?;

        fs::write(&claude_config_path, config_json)
            .context("Failed to write Claude configuration file")?;

        // Set appropriate file permissions (Unix-like systems only)
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let permissions = std::fs::Permissions::from_mode(0o600); // Owner read/write only
            fs::set_permissions(&claude_config_path, permissions)
                .context("Failed to set Claude configuration file permissions")?;
        }

        Ok(())
    }

    /// Switch Claude configuration to use official backend (remove env field)
    pub fn switch_to_claude_official() -> Result<()> {
        let mut claude_config = Self::load_claude_config()?;

        // Remove the env field to switch back to official Claude
        if let Some(obj) = claude_config.as_object_mut() {
            obj.remove("env");
        }

        Self::save_claude_config(&claude_config)?;
        Ok(())
    }

    /// Switch Claude configuration to use specified environment
    pub fn switch_to_claude_env(env: &str) -> Result<()> {
        let mut claude_config = Self::load_claude_config()?;

        // Set the env field
        if claude_config.is_null() {
            claude_config = serde_json::json!({});
        }

        if let Some(obj) = claude_config.as_object_mut() {
            obj.insert(
                "env".to_string(),
                serde_json::Value::String(env.to_string()),
            );
        }

        Self::save_claude_config(&claude_config)?;
        Ok(())
    }

    /// Load configuration from a custom path
    pub fn load_config_from_path<P: AsRef<Path>>(path: P) -> Result<ModixConfig> {
        let config_content =
            fs::read_to_string(&path).context("Failed to read configuration file")?;

        let config: ModixConfig =
            serde_json::from_str(&config_content).context("Failed to parse configuration JSON")?;

        // Validate the loaded configuration
        config
            .validate()
            .map_err(|e| anyhow::anyhow!("Invalid configuration: {}", e))?;

        Ok(config)
    }

    /// Save configuration to the default path (~/.modix/settings.json)
    pub fn save_config(config: &ModixConfig) -> Result<()> {
        let config_path = get_config_path();

        // Create parent directory if it doesn't exist
        if let Some(parent) = config_path.parent() {
            fs::create_dir_all(parent).context("Failed to create configuration directory")?;
        }

        // Update timestamp
        let mut config_to_save = config.clone();
        let current_time = chrono::Utc::now().to_rfc3339();
        config_to_save.updated_at = Some(current_time.clone());

        // Set created_at if this is a new configuration
        if config_to_save.created_at.is_none() {
            config_to_save.created_at = Some(current_time);
        }

        let config_json = serde_json::to_string_pretty(&config_to_save)
            .context("Failed to serialize configuration to JSON")?;

        fs::write(&config_path, config_json).context("Failed to write configuration file")?;

        // Set appropriate file permissions (Unix-like systems only)
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let permissions = std::fs::Permissions::from_mode(0o600); // Owner read/write only
            fs::set_permissions(&config_path, permissions)
                .context("Failed to set configuration file permissions")?;
        }

        Ok(())
    }

    /// Save configuration to a custom path
    pub fn save_config_to_path<P: AsRef<Path>>(config: &ModixConfig, path: P) -> Result<()> {
        // Create parent directory if it doesn't exist
        if let Some(parent) = path.as_ref().parent() {
            fs::create_dir_all(parent).context("Failed to create configuration directory")?;
        }

        // Update timestamp
        let mut config_to_save = config.clone();
        let current_time = chrono::Utc::now().to_rfc3339();
        config_to_save.updated_at = Some(current_time.clone());

        // Set created_at if this is a new configuration
        if config_to_save.created_at.is_none() {
            config_to_save.created_at = Some(current_time);
        }

        let config_json = serde_json::to_string_pretty(&config_to_save)
            .context("Failed to serialize configuration to JSON")?;

        fs::write(&path, config_json).context("Failed to write configuration file")?;

        // Set appropriate file permissions (Unix-like systems only)
        #[cfg(unix)]
        {
            use std::os::unix::fs::PermissionsExt;
            let permissions = std::fs::Permissions::from_mode(0o600); // Owner read/write only
            fs::set_permissions(&path, permissions)
                .context("Failed to set configuration file permissions")?;
        }

        Ok(())
    }

    /// Check if configuration file exists
    pub fn config_exists() -> bool {
        get_config_path().exists()
    }

    /// Get the path to the configuration file
    pub fn get_config_file_path() -> String {
        get_config_path().to_string_lossy().to_string()
    }

    /// Reset configuration to default values
    pub fn reset_config() -> Result<()> {
        let default_config = ModixConfig::new();
        Self::save_config(&default_config)
    }

    /// Import configuration from another file
    pub fn import_config<P: AsRef<Path>>(source_path: P) -> Result<()> {
        let config = Self::load_config_from_path(source_path)?;
        Self::save_config(&config)
    }

    /// Export configuration to another file
    pub fn export_config<P: AsRef<Path>>(target_path: P) -> Result<()> {
        let config = Self::load_config()?;
        Self::save_config_to_path(&config, target_path)
    }
}

/// Initialize default models for common vendors
pub fn setup_default_models() -> ModixConfig {
    let mut config = ModixConfig::new();

    // Anthropic
    let claude_config = crate::config::ModelConfig {
        company: "Anthropic".to_string(),
        api_endpoint: "".to_string(),
        api_key: "".to_string(),
        models: vec!["Claude".to_string()],
    };
    config.add_vendor("anthropic".to_string(), claude_config);

    // DeepSeek
    let deepseek_config = crate::config::ModelConfig {
        company: "DeepSeek".to_string(),
        api_endpoint: "https://api.deepseek.com/v1".to_string(),
        api_key: "".to_string(),
        models: vec!["deepseek-reasoner".to_string(), "deepseek-chat".to_string()],
    };
    config.add_vendor("deepseek".to_string(), deepseek_config);

    // Alibaba
    let qwen_config = crate::config::ModelConfig {
        company: "Alibaba".to_string(),
        api_endpoint: "https://dashscope.aliyuncs.com/compatible-mode/v1".to_string(),
        api_key: "".to_string(),
        models: vec![
            "qwen3-coder-plus".to_string(),
            "qwen3-coder-flash".to_string(),
        ],
    };
    config.add_vendor("bailian".to_string(), qwen_config);

    // ByteDance
    let doubao_config = crate::config::ModelConfig {
        company: "ByteDance".to_string(),
        api_endpoint: "https://ark.cn-beijing.volces.com/api/coding".to_string(),
        api_key: "".to_string(),
        models: vec!["doubao-seed-code-preview-latest".to_string()],
    };
    config.add_vendor("volcengine".to_string(), doubao_config);

    // Moonshot AI
    let kimi_config = crate::config::ModelConfig {
        company: "Moonshot AI".to_string(),
        api_endpoint: "https://api.moonshot.cn/anthropic".to_string(),
        api_key: "".to_string(),
        models: vec!["kimi-k2-thinking-turbo".to_string()],
    };
    config.add_vendor("moonshot".to_string(), kimi_config);

    // KuaiShou
    let kat_config = crate::config::ModelConfig {
        company: "Kuaishou".to_string(),
        api_endpoint: "https://wanqing.streamlakeapi.com/api/gateway/v1/endpoints/ep-xxx-xxx/claude-code-proxy".to_string(),
        api_key: "".to_string(),
        models: vec!["KAT-Coder".to_string()],
    };
    config.add_vendor("streamlake".to_string(), kat_config);

    // MiniMax
    let minimax_config = crate::config::ModelConfig {
        company: "MiniMax".to_string(),
        api_endpoint: "https://api.minimaxi.com/anthropic".to_string(),
        api_key: "".to_string(),
        models: vec!["MiniMax-M2".to_string()],
    };
    config.add_vendor("minimax".to_string(), minimax_config);

    // ZHIPU AI
    let zhipu_config = crate::config::ModelConfig {
        company: "ZHIPU AI".to_string(),
        api_endpoint: "https://open.bigmodel.cn/api/anthropic".to_string(),
        api_key: "".to_string(),
        models: vec!["GLM-4.6".to_string()],
    };
    config.add_vendor("bigmodel".to_string(), zhipu_config);

    // Set current and default vendor/model to the first enabled model
    config.current_vendor = "anthropic".to_string();
    config.current_model = "Claude".to_string();
    config.default_vendor = "anthropic".to_string();
    config.default_model = "Claude".to_string();

    config
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::NamedTempFile;

    #[test]
    fn test_load_config_nonexistent() {
        // Create a temporary directory to avoid conflicts with real config file
        let temp_dir = tempfile::tempdir().unwrap();
        let temp_config_path = temp_dir.path().join(".modix").join("settings.json");

        // Ensure the temporary directory structure exists
        std::fs::create_dir_all(temp_config_path.parent().unwrap()).unwrap();

        // Create a test that uses the temporary path
        let config_result = ConfigManager::load_config_from_path(&temp_config_path);

        match config_result {
            Ok(config) => {
                assert_eq!(config.vendors.len(), 0);
            }
            Err(_) => {
                // If file doesn't exist, it should return empty config
                let empty_config = ModixConfig::new();
                assert_eq!(empty_config.vendors.len(), 0);
            }
        }
    }

    #[test]
    fn test_save_and_load_config() {
        let temp_file = NamedTempFile::new().unwrap();
        let temp_path = temp_file.path();

        let mut config = ModixConfig::new();
        let model_config = crate::config::ModelConfig {
            company: "TestCorp".to_string(),
            api_endpoint: "https://api.test.com".to_string(),
            api_key: "test-key".to_string(),
            models: vec!["test-model".to_string()],
        };
        config.add_vendor("test-vendor".to_string(), model_config);
        config.current_vendor = "test-vendor".to_string();
        config.current_model = "test-model".to_string();
        config.default_vendor = "test-vendor".to_string();
        config.default_model = "test-model".to_string();

        // Save config
        ConfigManager::save_config_to_path(&config, temp_path).unwrap();

        // Load config back
        let loaded_config = ConfigManager::load_config_from_path(temp_path).unwrap();

        assert_eq!(config.vendors.len(), loaded_config.vendors.len());
        assert!(loaded_config.get_vendor("test-vendor").is_some());
    }

    #[test]
    fn test_setup_default_models() {
        let config = setup_default_models();
        assert_eq!(config.vendors.len(), 8);
        assert!(config.get_vendor("anthropic").is_some());
        assert!(config.get_vendor("deepseek").is_some());
        assert!(config.get_vendor("alibaba").is_some());
        assert!(config.get_vendor("bytedance").is_some());
        assert!(config.get_vendor("moonshot").is_some());
        assert!(config.get_vendor("kat").is_some());
        assert!(config.get_vendor("minimax").is_some());
        assert!(config.get_vendor("zhipu").is_some());

        // Verify models are present
        assert!(config.get_model("anthropic", "Claude").is_some());
        assert!(config.get_model("deepseek", "deepseek-reasoner").is_some());
        assert!(config.get_model("alibaba", "qwen3-coder-plus").is_some());
    }
}
