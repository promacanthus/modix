use clap::{Parser, Subcommand};
use modix::{setup_default_models, ConfigManager, ModelConfig};
use std::process;

/// CLI tool for managing and switching between Claude API backends and other LLMs
#[derive(Parser, Debug)]
#[command(name = "modix")]
#[command(about = "CLI tool for managing and switching between Claude API backends and other LLMs")]
#[command(version = "0.1.0")]
struct Cli {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand, Debug)]
enum Commands {
    /// List configured models
    List,
    /// Switch to a different model
    Switch {
        /// Model identifier to switch to
        model_name: String,
    },
    /// Show current model status
    Status,
    /// Add a new model configuration
    Add {
        /// Model name (used as identifier)
        model_name: String,
        /// Company name
        #[arg(short = 'c', long)]
        company: String,
        /// Vendor
        #[arg(short = 'v', long)]
        vendor: String,
        /// API endpoint URL
        #[arg(short = 'u', long)]
        endpoint: String,
        /// API key
        #[arg(short = 'k', long)]
        api_key: String,
    },
    /// Remove a model configuration
    Remove {
        /// Model identifier to remove
        model_name: String,
    },
    /// Initialize default configuration
    Init,
    /// Show configuration file path
    Path,
}

fn main() {
    let cli = Cli::parse();

    match run_command(cli.command) {
        Ok(_) => {}
        Err(e) => {
            eprintln!("Error: {}", e);
            process::exit(1);
        }
    }
}

fn run_command(command: Commands) -> Result<(), Box<dyn std::error::Error>> {
    match command {
        Commands::List => cmd_list()?,
        Commands::Switch { model_name } => cmd_switch(&model_name)?,
        Commands::Status => cmd_status()?,
        Commands::Add {
            model_name,
            company,
            vendor,
            endpoint,
            api_key,
        } => cmd_add(&model_name, &company, &vendor, &endpoint, &api_key)?,
        Commands::Remove { model_name } => cmd_remove(&model_name)?,
        Commands::Init => cmd_init()?,
        Commands::Path => cmd_path()?,
    }
    Ok(())
}

fn cmd_list() -> Result<(), Box<dyn std::error::Error>> {
    let config = ConfigManager::load_config()?;

    println!("Available models:");
    println!(
        "{:<30} {:<20} {:<15} {:<30}",
        "MODEL", "COMPANY", "VENDOR", "ENDPOINT"
    );

    // Get all models and sort them by vendor and model name for consistent output
    let mut models: Vec<_> = config.get_all_models();
    models.sort_by(|a, b| {
        let vendor_cmp = a.0.cmp(&b.0);
        if vendor_cmp == std::cmp::Ordering::Equal {
            a.1.cmp(&b.1)
        } else {
            vendor_cmp
        }
    });

    for (vendor, model_name, model_config) in models {
        println!(
            "{:<30} {:<20} {:<15} {:<30}",
            model_name, model_config.company, vendor, model_config.api_endpoint
        );
    }

    Ok(())
}

fn cmd_switch(model_name: &str) -> Result<(), Box<dyn std::error::Error>> {
    let mut config = ConfigManager::load_config()?;

    // Find the model by searching all vendors for this model name
    let all_models = config.get_all_models();
    let mut found_vendor: Option<String> = None;

    for (vendor, model, _) in all_models {
        if model == model_name {
            found_vendor = Some(vendor.clone());
            break;
        }
    }

    if let Some(vendor) = found_vendor {
        config.set_current_vendor_and_model(&vendor, model_name)?;
        ConfigManager::save_config(&config)?;

        // Update Claude configuration based on the selected model
        update_claude_config_for_model(model_name, &vendor)?;

        println!("Switched to model: {}", model_name);
    } else {
        return Err(format!("Model '{}' not found in any vendor configuration", model_name).into());
    }

    Ok(())
}

/// Update Claude configuration based on the selected model
fn update_claude_config_for_model(model_name: &str, vendor: &str) -> Result<(), Box<dyn std::error::Error>> {
    // Load the current Modix configuration to get model details
    let config = ConfigManager::load_config()?;

    // Get the model configuration
    let model_config = config.get_model(vendor, model_name)
        .ok_or_else(|| format!("Model '{}' not found in configuration for vendor '{}'", model_name, vendor))?;

    // Load existing Claude configuration
    let mut claude_config = ConfigManager::load_claude_config()?;

    // If the config is empty, initialize with basic structure
    if claude_config.is_null() || claude_config.as_object().map_or(true, |obj| obj.is_empty()) {
        claude_config = serde_json::json!({
            "companyAnnouncements": ["Welcome to claude code, the configuration managed by modix."]
        });
    }

    // Check if the model is Claude (case-insensitive check for various Claude model names)
    let is_claude_model = model_name.to_lowercase().contains("claude")
        || vendor.to_lowercase() == "anthropic";

    if is_claude_model {
        // For Claude models, remove the env field to use official Claude API
        if let Some(obj) = claude_config.as_object_mut() {
            obj.remove("env");
        }
    } else {
        // For non-Claude models, set up the detailed env configuration
        let env_config = serde_json::json!({
            "ANTHROPIC_BASE_URL": model_config.api_endpoint,
            "ANTHROPIC_AUTH_TOKEN": model_config.api_key,
            "API_TIMEOUT_MS": "3000000",
            "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": 1,
            "ANTHROPIC_MODEL": model_name,
            "ANTHROPIC_SMALL_FAST_MODEL": model_name,
            "ANTHROPIC_DEFAULT_SONNET_MODEL": model_name,
            "ANTHROPIC_DEFAULT_OPUS_MODEL": model_name,
            "ANTHROPIC_DEFAULT_HAIKU_MODEL": model_name
        });

        // Update the env field with the new configuration
        if let Some(obj) = claude_config.as_object_mut() {
            obj.insert("env".to_string(), env_config);
        }
    }

    // Save the updated Claude configuration
    ConfigManager::save_claude_config(&claude_config)?;

    Ok(())
}

fn cmd_status() -> Result<(), Box<dyn std::error::Error>> {
    let config = ConfigManager::load_config()?;

    if let Some((_model_name, model_config)) = config.get_current_model() {
        println!("Current model: {}", config.current_model);
        println!("Current vendor: {}", config.current_vendor);
        println!("Company: {}", model_config.company);
        println!("API Endpoint: {}", model_config.api_endpoint);
    } else {
        println!("No current model configured");
    }

    Ok(())
}

fn cmd_add(
    model_name: &str,
    company: &str,
    vendor: &str,
    endpoint: &str,
    api_key: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut config = ConfigManager::load_config()?;

    // Check if this model name already exists in any vendor
    let all_models = config.get_all_models();
    for (_, existing_model, _) in &all_models {
        if existing_model == model_name {
            return Err(format!(
                "Model '{}' already exists. Please use a different name or remove the existing one first.",
                model_name
            ).into());
        }
    }

    // First try to add the model to an existing vendor
    if config.add_model_to_vendor(vendor, model_name.to_string()).is_ok() {
        // Update the API endpoint and key if they changed
        if let Some(vendor_config) = config.get_vendor_mut(vendor) {
            vendor_config.company = company.to_string();
            vendor_config.api_endpoint = endpoint.to_string();
            vendor_config.api_key = api_key.to_string();
        }
        ConfigManager::save_config(&config)?;
        println!("Added model '{}' to existing vendor '{}'", model_name, vendor);
    } else {
        // If vendor doesn't exist, create a new vendor config
        let model_config = ModelConfig {
            company: company.to_string(),
            api_endpoint: endpoint.to_string(),
            api_key: api_key.to_string(),
            models: vec![model_name.to_string()],
        };

        config.add_vendor(vendor.to_string(), model_config);
        ConfigManager::save_config(&config)?;

        println!("Created new vendor '{}' with model: {}", vendor, model_name);
    }

    println!("Switch to it with: modix switch {}", model_name);

    Ok(())
}

fn cmd_remove(model_name: &str) -> Result<(), Box<dyn std::error::Error>> {
    let mut config = ConfigManager::load_config()?;

    // Find which vendor this model belongs to
    let all_models = config.get_all_models();
    let mut target_vendor: Option<String> = None;
    for (vendor, model, _) in all_models {
        if model == model_name {
            target_vendor = Some(vendor.clone());
            break;
        }
    }

    if target_vendor.is_none() {
        return Err(format!("Model '{}' not found", model_name).into());
    }

    let vendor = target_vendor.unwrap();

    // Remove the model from the vendor
    if let Some(vendor_config) = config.get_vendor_mut(&vendor) {
        vendor_config.models.retain(|m| m != model_name);

        // If the vendor has no more models, remove the vendor entirely
        if vendor_config.models.is_empty() {
            config.remove_vendor(&vendor);
            println!("Vendor '{}' had no remaining models and was removed", vendor);
        }
    }

    // If we removed the current model, switch to default
    if config.current_model == model_name {
        config.current_model = config.default_model.clone();
        config.current_vendor = config.default_vendor.clone();
        println!(
            "Removed current model. Switched to default: {}@{}",
            config.default_model, config.default_vendor
        );
    }

    ConfigManager::save_config(&config)?;
    println!("Removed model: {}", model_name);

    Ok(())
}

fn cmd_init() -> Result<(), Box<dyn std::error::Error>> {
    if ConfigManager::config_exists() {
        println!(
            "Configuration already exists at: {}",
            ConfigManager::get_config_file_path()
        );
        println!("Use 'modix path' to show the configuration file path");
        return Ok(());
    }

    let config = setup_default_models();
    ConfigManager::save_config(&config)?;

    println!(
        "Initialized default configuration at: {}",
        ConfigManager::get_config_file_path()
    );
    println!("Edit the configuration file to add your API keys and enable models");
    println!("Use 'modix list' to see available models");

    Ok(())
}

fn cmd_path() -> Result<(), Box<dyn std::error::Error>> {
    println!(
        "Configuration file path: {}",
        ConfigManager::get_config_file_path()
    );
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cli_parsing() {
        // This is a basic test to ensure CLI parsing works
        let cli = Cli::parse_from(vec!["modix", "list"]);
        assert!(matches!(cli.command, Commands::List));

        let cli = Cli::parse_from(vec!["modix", "switch", "claude-official"]);
        assert!(
            matches!(cli.command, Commands::Switch { model_name } if model_name == "claude-official")
        );

        let cli = Cli::parse_from(vec!["modix", "status"]);
        assert!(matches!(cli.command, Commands::Status));

        let cli = Cli::parse_from(vec!["modix", "add", "my-model", "-c", "TestCorp", "-v", "test", "-u", "https://api.test.com", "-k", "test-key"]);
        assert!(matches!(
            cli.command,
            Commands::Add { model_name, company, vendor, endpoint, api_key }
            if model_name == "my-model" && company == "TestCorp" && vendor == "test" && endpoint == "https://api.test.com" && api_key == "test-key"
        ));
    }

    #[test]
    fn test_list_command_sorting() {
        use tempfile::NamedTempFile;

        // Create a temporary config file with models in random order
        let temp_file = NamedTempFile::new().unwrap();
        let temp_path = temp_file.path();

        // Write config with models in non-alphabetical order - using vendor structure
        let config_content = r#"{
  "current_vendor": "test",
  "current_model": "Zeta-Model",
  "default_vendor": "test",
  "default_model": "Alpha-Model",
  "vendors": {
    "test": {
      "company": "TestCorp",
      "api_endpoint": "https://api.test.com",
      "api_key": "test-key",
      "models": [
        "Zeta-Model",
        "Alpha-Model",
        "Beta-Model"
      ]
    }
  },
  "config_version": "1.0.0"
}"#;

        std::fs::write(temp_path, config_content).unwrap();

        // Load config from the temporary file
        let config = crate::ConfigManager::load_config_from_path(temp_path).unwrap();

        // Test that get_all_models returns models and we can sort them
        let mut models: Vec<_> = config.get_all_models();
        models.sort_by(|a, b| {
            let vendor_cmp = a.0.cmp(&b.0);
            if vendor_cmp == std::cmp::Ordering::Equal {
                a.1.cmp(&b.1)
            } else {
                vendor_cmp
            }
        });

        // Verify the models are sorted alphabetically by vendor then model name
        assert_eq!(models[0].1, "Alpha-Model");
        assert_eq!(models[1].1, "Beta-Model");
        assert_eq!(models[2].1, "Zeta-Model");
    }
}
