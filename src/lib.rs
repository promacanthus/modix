pub mod config;
pub mod config_manager;

pub use config::{get_config_path, ModelConfig, ModixConfig};
pub use config_manager::{setup_default_models, ConfigManager};
