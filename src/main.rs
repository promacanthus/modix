use clap::{Parser, Subcommand};
use colored::*;
use modix::{setup_default_models, ConfigManager, ModelConfig};
use std::process;

// Force colored output for all terminals
fn init_colors() {
    control::set_override(true);
}

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
    /// Add a new model configuration
    Add {
        /// Model name (used as identifier)
        model_name: String,
        /// Company that develops the model (e.g., "Anthropic", "DeepSeek", "Alibaba")
        #[arg(short = 'c', long, help = "Company that develops the model")]
        company: String,
        /// API provider/vendor (e.g., "anthropic", "deepseek", "alibaba", "moonshot")
        #[arg(short = 'v', long, help = "API provider/vendor identifier")]
        vendor: String,
        /// API endpoint URL
        #[arg(short = 'u', long)]
        endpoint: String,
        /// API key
        #[arg(short = 'k', long)]
        api_key: String,
    },
    /// Check configuration for different tools
    ///
    /// Available tools:
    ///   - claude-code: Check Claude Code configuration
    ///   - modix: Check Modix configuration
    ///   - codex: Check Codex configuration
    ///   - gemini-cli: Check Gemini CLI configuration
    Check {
        /// Tool to check configuration for (claude-code, modix, codex, gemini-cli)
        tool: String,
    },
    /// Initialize default configuration
    Init,
    /// List configured models
    List,
    /// Show configuration file path
    Path,
    /// Remove a model configuration
    Remove {
        /// Model identifier to remove
        model_name: String,
    },
    /// Show details for a specific vendor (including API key)
    Show {
        /// Vendor identifier to display
        vendor: String,
    },
    /// Show current model status
    Status,
    /// Switch to a different model
    Switch {
        /// Model identifier to switch to
        model_name: String,
    },
    /// Update an existing vendor configuration (model, company, api_endpoint, or api_key)
    Update {
        /// The vendor identifier to update
        vendor: String,
        /// Add a model to the vendor
        #[arg(short = 'm', long)]
        add_model: Option<String>,
        /// Update company name
        #[arg(short = 'c', long)]
        company: Option<String>,
        /// Update API endpoint URL
        #[arg(short = 'u', long)]
        endpoint: Option<String>,
        /// Update API key
        #[arg(short = 'k', long)]
        api_key: Option<String>,
    },
}

fn main() {
    init_colors();
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
        Commands::Add {
            model_name,
            company,
            vendor,
            endpoint,
            api_key,
        } => cmd_add(&model_name, &company, &vendor, &endpoint, &api_key)?,
        Commands::Check { tool } => cmd_check(&tool)?,
        Commands::Init => cmd_init()?,
        Commands::List => cmd_list()?,
        Commands::Path => cmd_path()?,
        Commands::Remove { model_name } => cmd_remove(&model_name)?,
        Commands::Show { vendor } => cmd_show(&vendor)?,
        Commands::Status => cmd_status()?,
        Commands::Switch { model_name } => cmd_switch(&model_name)?,
        Commands::Update {
            vendor,
            add_model,
            company,
            endpoint,
            api_key,
        } => cmd_update(
            &vendor,
            add_model.as_deref(),
            company.as_deref(),
            endpoint.as_deref(),
            api_key.as_deref(),
        )?,
    }
    Ok(())
}

fn cmd_list() -> Result<(), Box<dyn std::error::Error>> {
    let config = ConfigManager::load_config()?;

    println!(
        "{:<35} {:<15} {:<15} {:<10} {:<10}",
        "MODEL".bright_cyan().bold(),
        "COMPANY".bright_cyan().bold(),
        "VENDOR".bright_cyan().bold(),
        "ENDPOINT".bright_cyan().bold(),
        "API_KEY".bright_cyan().bold(),
    );
    println!(
        "{} {} {} {} {}",
        "-".repeat(35).bright_black(),
        "-".repeat(15).bright_black(),
        "-".repeat(15).bright_black(),
        "-".repeat(10).bright_black(),
        "-".repeat(10).bright_black(),
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

    for (vendor, model_name, model_config) in &models {
        // Show endpoint status - special handling for Anthropic
        let endpoint_display = if vendor.to_lowercase() == "anthropic" {
            "[ - ]".bright_blue().bold()
        } else if model_config.api_endpoint.is_empty() {
            "[ N ]".bright_red().bold()
        } else {
            "[ Y ]".bright_green().bold()
        };

        // Show API key status - special handling for Anthropic
        let api_key_display = if vendor.to_lowercase() == "anthropic" {
            "[ - ]".bright_blue().bold()
        } else if model_config.api_key.is_empty() {
            "[ N ]".bright_red().bold()
        } else {
            "[ Y ]".bright_green().bold()
        };

        // Highlight current model
        let model_display =
            if config.current_vendor == *vendor && config.current_model == *model_name {
                model_name.bright_yellow().bold()
            } else {
                model_name.bright_white()
            };

        // Highlight current vendor
        let vendor_display = if config.current_vendor == *vendor {
            vendor.bright_yellow().bold()
        } else {
            vendor.bright_cyan()
        };

        println!(
            "{:<35} {:<15} {:<15} {:<10} {:<10}",
            model_display,
            model_config.company.bright_blue(),
            vendor_display,
            endpoint_display,
            api_key_display
        );
    }

    // Show summary information
    let total_models = models.len();
    let configured_models: Vec<_> = models
        .iter()
        .filter(|(_, _, config)| !config.api_endpoint.is_empty() && !config.api_key.is_empty())
        .collect();
    let current_model_info = if let Some((_, model_name, _)) =
        models.iter().find(|(vendor, model_name, _)| {
            config.current_vendor == *vendor && config.current_model == *model_name
        }) {
        format!("{}@{}", config.current_vendor, model_name)
    } else {
        "None".to_string()
    };

    println!();
    println!("{}", "--- Summary ---".bright_cyan().bold());
    println!(
        "{}: {}",
        "Total models".bright_blue(),
        total_models.to_string().bright_yellow()
    );
    println!(
        "{}: {}",
        "Configured models".bright_blue(),
        configured_models.len().to_string().bright_green()
    );
    println!(
        "{}: {}",
        "Current model".bright_blue(),
        current_model_info.bright_yellow()
    );

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

        println!(
            "{}: {}",
            "Switched to model".bright_green().bold(),
            model_name.bright_yellow()
        );
    } else {
        return Err(format!(
            "Model '{}' not found in any vendor configuration",
            model_name
        )
        .into());
    }

    Ok(())
}

/// Update Claude configuration based on the selected model
fn update_claude_config_for_model(
    model_name: &str,
    vendor: &str,
) -> Result<(), Box<dyn std::error::Error>> {
    // Load the current Modix configuration to get model details
    let config = ConfigManager::load_config()?;

    // Get the model configuration
    let model_config = config.get_model(vendor, model_name).ok_or_else(|| {
        format!(
            "Model '{}' not found in configuration for vendor '{}'",
            model_name, vendor
        )
    })?;

    // Load existing Claude configuration
    let mut claude_config = ConfigManager::load_claude_config()?;

    // If the config is empty, initialize with basic structure
    if claude_config.is_null() || claude_config.as_object().map_or(true, |obj| obj.is_empty()) {
        claude_config = serde_json::json!({
            "companyAnnouncements": ["Welcome to Claude code, the configuration managed by modix."]
        });
    }

    // Check if the model is Claude (case-insensitive check for various Claude model names)
    let is_claude_model =
        model_name.to_lowercase().contains("claude") || vendor.to_lowercase() == "anthropic";

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
        println!(
            "{}: {}",
            "Current model".bright_cyan().bold(),
            config.current_model.bright_green()
        );
        println!(
            "{}: {}",
            "Current vendor".bright_cyan().bold(),
            config.current_vendor.bright_green()
        );
        println!(
            "{}: {}",
            "Company".bright_cyan().bold(),
            model_config.company.bright_yellow()
        );
        println!(
            "{}: {}",
            "API Endpoint".bright_cyan().bold(),
            model_config.api_endpoint.bright_blue()
        );
    } else {
        println!("{}", "No current model configured".bright_red().bold());
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
    if config
        .add_model_to_vendor(vendor, model_name.to_string())
        .is_ok()
    {
        // Update the API endpoint and key if they changed
        if let Some(vendor_config) = config.get_vendor_mut(vendor) {
            vendor_config.company = company.to_string();
            vendor_config.api_endpoint = endpoint.to_string();
            vendor_config.api_key = api_key.to_string();
        }
        ConfigManager::save_config(&config)?;
        println!(
            "{}: '{}' {}",
            "Added model".bright_green().bold(),
            model_name.bright_yellow(),
            format!("to existing vendor '{}'", vendor).bright_cyan()
        );
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

        println!(
            "{}: '{}' {}",
            "Created new vendor".bright_green().bold(),
            vendor.bright_yellow(),
            format!("with model: {}", model_name).bright_cyan()
        );
    }

    println!(
        "{}: modix switch {}",
        "Switch to it with".bright_blue().bold(),
        model_name.bright_yellow()
    );

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
            println!(
                "Vendor '{}' had no remaining models and was removed",
                vendor
            );
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
    let path = ConfigManager::get_config_file_path();
    println!(
        "{}: {}",
        "Configuration file path".bright_cyan().bold(),
        path.bright_yellow()
    );
    Ok(())
}

fn cmd_update(
    vendor: &str,
    add_model: Option<&str>,
    company: Option<&str>,
    endpoint: Option<&str>,
    api_key: Option<&str>,
) -> Result<(), Box<dyn std::error::Error>> {
    let mut config = ConfigManager::load_config()?;

    // Check if vendor exists
    if config.get_vendor(vendor).is_none() {
        return Err(format!(
            "Vendor '{}' not found. Use 'modix add' to create a new vendor first.",
            vendor
        )
        .into());
    }

    let mut updates: Vec<String> = Vec::new();

    // Update company if provided
    if let Some(company_name) = company {
        if let Some(vendor_config) = config.get_vendor_mut(vendor) {
            vendor_config.company = company_name.to_string();
            updates.push(format!("Company: {}", company_name));
        }
    }

    // Update endpoint if provided
    if let Some(endpoint_url) = endpoint {
        if let Some(vendor_config) = config.get_vendor_mut(vendor) {
            vendor_config.api_endpoint = endpoint_url.to_string();
            updates.push(format!("API Endpoint: {}", endpoint_url));
        }
    }

    // Update API key if provided
    if let Some(key) = api_key {
        if let Some(vendor_config) = config.get_vendor_mut(vendor) {
            vendor_config.api_key = key.to_string();
            updates.push("API Key: [updated]".to_string());
        }
    }

    // Add model if provided
    if let Some(model_name) = add_model {
        if config
            .add_model_to_vendor(vendor, model_name.to_string())
            .is_ok()
        {
            updates.push(format!("Added model: {}", model_name));
        } else {
            return Err(format!(
                "Model '{}' already exists in vendor '{}'",
                model_name, vendor
            )
            .into());
        }
    }

    // Check if any updates were made
    if updates.is_empty() {
        println!("No updates were specified. Use --help to see available options.");
        return Ok(());
    }

    // Save the configuration
    ConfigManager::save_config(&config)?;

    // Display what was updated
    println!("Updated vendor '{}' with:", vendor);
    for update in updates {
        println!("  - {}", update);
    }

    Ok(())
}

fn cmd_check(tool: &str) -> Result<(), Box<dyn std::error::Error>> {
    match tool {
        "claude-code" => check_claude_code()?,
        "modix" => check_modix_config()?,
        "codex" => {
            println!("Codex configuration check is not yet implemented");
        }
        "gemini-cli" => {
            println!("Gemini CLI configuration check is not yet implemented");
        }
        _ => {
            return Err(format!(
                "Unknown tool '{}'. Supported tools: claude-code, modix, codex, gemini-cli",
                tool
            )
            .into());
        }
    }
    Ok(())
}

fn check_claude_code() -> Result<(), Box<dyn std::error::Error>> {
    let home_dir = dirs::home_dir().ok_or("Could not determine home directory")?;
    let claude_config_path = home_dir.join(".claude").join("settings.json");

    println!(
        "{}",
        "Checking Claude Code configuration..."
            .bright_yellow()
            .bold()
    );
    println!(
        "{}: {}",
        "Config file path".bright_cyan().bold(),
        claude_config_path.display().to_string().bright_yellow()
    );
    println!();

    if !claude_config_path.exists() {
        println!(
            "{}: {}",
            "Configuration file not found".bright_red().bold(),
            claude_config_path.display().to_string().bright_yellow()
        );
        return Ok(());
    }

    let contents = std::fs::read_to_string(&claude_config_path)
        .map_err(|e| format!("Failed to read configuration file: {}", e))?;

    println!(
        "{}",
        "=== Claude Code Configuration ===".bright_cyan().bold()
    );
    println!();

    // Pretty print JSON with syntax highlighting
    if let Ok(json_value) = serde_json::from_str::<serde_json::Value>(&contents) {
        let pretty_json = serde_json::to_string_pretty(&json_value)?;
        println!("{}", pretty_json.bright_white());
    } else {
        // If not valid JSON, just print the raw content
        println!("{}", contents.bright_white());
    }

    Ok(())
}

fn check_modix_config() -> Result<(), Box<dyn std::error::Error>> {
    let home_dir = dirs::home_dir().ok_or("Could not determine home directory")?;
    let modix_config_path = home_dir.join(".modix").join("settings.json");

    println!(
        "{}",
        "Checking Modix configuration...".bright_yellow().bold()
    );
    println!(
        "{}: {}",
        "Config file path".bright_cyan().bold(),
        modix_config_path.display().to_string().bright_yellow()
    );
    println!();

    if !modix_config_path.exists() {
        println!(
            "{}: {}",
            "Configuration file not found".bright_red().bold(),
            modix_config_path.display().to_string().bright_yellow()
        );
        println!(
            "{}",
            "Run 'modix init' to create a default configuration".bright_blue()
        );
        return Ok(());
    }

    let contents = std::fs::read_to_string(&modix_config_path)
        .map_err(|e| format!("Failed to read configuration file: {}", e))?;

    println!("{}", "=== Modix Configuration ===".bright_cyan().bold());
    println!();

    // Parse and validate JSON structure
    match serde_json::from_str::<serde_json::Value>(&contents) {
        Ok(json_value) => {
            // Pretty print JSON with syntax highlighting
            let pretty_json = serde_json::to_string_pretty(&json_value)?;
            println!("{}", pretty_json.bright_white());

            // Show configuration summary
            println!();
            println!("{}", "--- Configuration Summary ---".bright_cyan().bold());

            if let Some(vendors) = json_value.get("vendors").and_then(|v| v.as_object()) {
                let total_vendors = vendors.len();
                let total_models = vendors
                    .values()
                    .map(|v| {
                        v.get("models")
                            .and_then(|m| m.as_array())
                            .map_or(0, |arr| arr.len())
                    })
                    .sum::<usize>();

                println!(
                    "{}: {}",
                    "Total vendors".bright_blue(),
                    total_vendors.to_string().bright_yellow()
                );
                println!(
                    "{}: {}",
                    "Total models".bright_blue(),
                    total_models.to_string().bright_yellow()
                );

                if let (Some(current_vendor), Some(current_model)) = (
                    json_value.get("current_vendor").and_then(|v| v.as_str()),
                    json_value.get("current_model").and_then(|m| m.as_str()),
                ) {
                    println!(
                        "{}: {}@{}",
                        "Current selection".bright_blue(),
                        current_vendor.bright_green().bold(),
                        current_model.bright_green().bold()
                    );
                }
            }

            // Check for common configuration issues
            println!();
            println!(
                "{}",
                "--- Configuration Health Check ---".bright_cyan().bold()
            );
            check_config_health(&json_value)?;
        }
        Err(e) => {
            // If not valid JSON, just print the raw content and show error
            println!("{}", contents.bright_white());
            println!();
            println!(
                "{}: {}",
                "JSON Parse Error".bright_red().bold(),
                e.to_string().bright_yellow()
            );
        }
    }

    Ok(())
}

fn check_config_health(config: &serde_json::Value) -> Result<(), Box<dyn std::error::Error>> {
    let mut issues = Vec::new();

    // Check if vendors section exists
    if config.get("vendors").is_none() {
        issues.push("Missing 'vendors' section".to_string());
    }

    // Check for vendors with empty endpoints or keys
    if let Some(vendors) = config.get("vendors").and_then(|v| v.as_object()) {
        for (vendor_name, vendor_config) in vendors {
            // Skip health check for Anthropic since it's pre-configured in Claude-code
            if vendor_name.to_lowercase() == "anthropic" {
                continue;
            }

            if let Some(endpoint) = vendor_config.get("api_endpoint").and_then(|e| e.as_str()) {
                if endpoint.is_empty() {
                    issues.push(format!("Vendor '{}' has empty API endpoint", vendor_name));
                }
            } else {
                issues.push(format!("Vendor '{}' missing API endpoint", vendor_name));
            }

            if let Some(api_key) = vendor_config.get("api_key").and_then(|k| k.as_str()) {
                if api_key.is_empty() {
                    issues.push(format!("Vendor '{}' has empty API key", vendor_name));
                }
            } else {
                issues.push(format!("Vendor '{}' missing API key", vendor_name));
            }
        }
    }

    // Check current model selection
    if let (Some(current_vendor), Some(current_model)) = (
        config.get("current_vendor").and_then(|v| v.as_str()),
        config.get("current_model").and_then(|m| m.as_str()),
    ) {
        if let Some(vendors) = config.get("vendors").and_then(|v| v.as_object()) {
            if let Some(vendor_config) = vendors.get(current_vendor) {
                if let Some(models) = vendor_config.get("models").and_then(|m| m.as_array()) {
                    let model_names: Vec<String> = models
                        .iter()
                        .filter_map(|m| m.as_str())
                        .map(|s| s.to_string())
                        .collect();

                    if !model_names.contains(&current_model.to_string()) {
                        issues.push(format!(
                            "Current model '{}' not found in vendor '{}' models",
                            current_model, current_vendor
                        ));
                    }
                }
            } else {
                issues.push(format!(
                    "Current vendor '{}' not found in configuration",
                    current_vendor
                ));
            }
        }
    }

    // Report issues
    if issues.is_empty() {
        println!(
            "{}: All checks passed! 🎉",
            "Health Check".bright_green().bold()
        );
    } else {
        println!(
            "{}: Found {} issue(s)",
            "Health Check".bright_red().bold(),
            issues.len().to_string().bright_yellow()
        );
        for issue in issues {
            println!("  - {}", issue.bright_red());
        }
    }

    Ok(())
}

fn cmd_show(vendor: &str) -> Result<(), Box<dyn std::error::Error>> {
    let config = ConfigManager::load_config()?;

    // Get the vendor configuration
    let vendor_config = config
        .get_vendor(vendor)
        .ok_or_else(|| format!("Vendor '{}' not found", vendor))?;

    println!("{}", "Vendor Details:".bright_cyan().bold());
    println!();
    println!(
        "{}: {}",
        "Vendor ID".bright_cyan().bold(),
        vendor.bright_green()
    );
    println!(
        "{}: {}",
        "Company".bright_cyan().bold(),
        vendor_config.company.bright_yellow()
    );
    println!(
        "{}: {}",
        "API Endpoint".bright_cyan().bold(),
        vendor_config.api_endpoint.bright_blue()
    );
    println!(
        "{}: {}",
        "API Key".bright_cyan().bold(),
        vendor_config.api_key.bright_red()
    );
    println!();
    println!("{}", "Models:".bright_cyan().bold());
    for model in &vendor_config.models {
        let current_marker = if config.current_vendor == vendor && config.current_model == *model {
            " (current)".bright_green().bold()
        } else {
            "".normal()
        };
        println!("  - {}{}", model.bright_magenta(), current_marker);
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_cli_parsing() {
        // This is a basic test to ensure CLI parsing works
        let cli = Cli::parse_from(vec![
            "modix",
            "add",
            "my-model",
            "-c",
            "TestCorp",
            "-v",
            "test",
            "-u",
            "https://api.test.com",
            "-k",
            "test-key",
        ]);
        assert!(matches!(
            cli.command,
            Commands::Add { model_name, company, vendor, endpoint, api_key }
            if model_name == "my-model" && company == "TestCorp" && vendor == "test" && endpoint == "https://api.test.com" && api_key == "test-key"
        ));

        let cli = Cli::parse_from(vec!["modix", "init"]);
        assert!(matches!(cli.command, Commands::Init));

        let cli = Cli::parse_from(vec!["modix", "list"]);
        assert!(matches!(cli.command, Commands::List));

        let cli = Cli::parse_from(vec!["modix", "path"]);
        assert!(matches!(cli.command, Commands::Path));

        let cli = Cli::parse_from(vec!["modix", "remove", "my-model"]);
        assert!(matches!(cli.command, Commands::Remove { model_name } if model_name == "my-model"));

        let cli = Cli::parse_from(vec!["modix", "show", "test"]);
        assert!(matches!(cli.command, Commands::Show { vendor } if vendor == "test"));

        let cli = Cli::parse_from(vec!["modix", "status"]);
        assert!(matches!(cli.command, Commands::Status));

        let cli = Cli::parse_from(vec!["modix", "switch", "claude-official"]);
        assert!(
            matches!(cli.command, Commands::Switch { model_name } if model_name == "claude-official")
        );

        let cli = Cli::parse_from(vec!["modix", "update", "test"]);
        assert!(matches!(cli.command, Commands::Update { vendor, .. } if vendor == "test"));
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
