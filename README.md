# Modix

A Rust-based CLI tool for managing and switching between Claude API backends and other large language models.

## 📁 Project Structure

```
modix/
├── Cargo.toml              # Rust CLI tool configuration
├── src/                    # Rust source code (modix CLI tool)
│   ├── main.rs
│   ├── config.rs
│   └── config_manager.rs
├── docs/                   # Astro documentation website
│   ├── package.json        # Astro project configuration
│   ├── astro.config.mjs
│   └── src/
│       ├── content/        # Documentation content
│       └── assets/
└── README.md               # Project documentation
```

This is a **monorepo** containing:

- **Rust CLI tool** (`src/`) - The main `modix` command-line tool
- **Documentation website** (`docs/`) - Astro-based documentation site with multi-language support

## Overview

Modix is designed to simplify the management of multiple LLM backends by providing a unified configuration system and command-line interface. Inspired by Claude's system in `~/.claude/settings.json`, Modix extends this concept to support multiple providers including Claude Official API, DeepSeek, Alibaba Qwen, and ByteDance Doubao.

## 🚀 Current Status

**Active Development** - Core infrastructure implemented, CLI commands in progress

- ✅ Basic project structure and dependencies configured
- ✅ Configuration management system implemented
- ✅ Model definitions and data structures complete
- 🔄 CLI command implementation in progress
- 🔄 API integration development upcoming

**Latest Update**: November 18, 2025 - Core configuration and model management infrastructure completed

## Features

- **Multi-Provider Support**: Manage configurations for Claude, DeepSeek, Qwen, Doubao, Kimi, MiniMax, and more
- **Easy Model Switching**: Quickly switch between different LLM backends
- **Secure Configuration**: Encrypted storage of API keys and sensitive configuration data
- **Cross-Platform**: Native binaries for Windows, macOS, and Linux
- **Colorful Output**: Enhanced CLI with color-coded output for better readability and user experience
- **Extensible Design**: Easy to add new providers and models

## Installation

### From Source (Recommended for Development)

```bash
git clone https://github.com/promacanthus/modix.git
cd modix
cargo build --release
```

The binary will be available at `target/release/modix`.

**Prerequisites**: [Rust and Cargo](https://rustup.rs/) must be installed on your system.

### Pre-built Binaries

Pre-built binaries will be available once the first stable release is published.

## Usage

**Note**: The CLI interface is currently under active development. The following commands represent the planned functionality.

### Initialize Configuration

```bash
modix init
```

This creates a default configuration file at `~/.modix/settings.json` with predefined models.

### List Available Models

```bash
modix list
```

Shows all configured models with their name, vendor, and API endpoint status in a clean table format.

**Display Legend:**

- `[ Y ]`: Configuration is present (green)
- `[ N ]`: Configuration is missing (red)
- `[ - ]`: Not applicable - pre-configured in Claude Code (blue)

For Anthropic models, both API endpoint and API key display as `[ - ]` since they are handled automatically by Claude Code.

### Check Configuration

```bash
modix check claude-code
```

Checks the configuration for different development tools. Currently supports:

- `claude-code`: Check Claude Code configuration
- `modix`: Check Modix configuration (validates vendors, endpoints, and API keys)

**Note**: Health checks for `anthropic` are automatically skipped since it's pre-configured in Claude Code and doesn't require manual API endpoint or key configuration.

```bash
modix check claude-code
modix check codex
modix check gemini-cli
```

```bash
# Switches the current active model to Claude Official API
modix switch Claude
```

```bash
# Switches to DeepSeek V3.2
modix switch deepseek-reasoner
```

### Check Current Status

Displays information about the currently active model.

```bash
modix status
```

### Add Custom Model

```bash
modix add my-custom-model \
  --company my-company \
  --vendor my-vendor \
  --endpoint https://api.mycustom.com \
  --api-key your-api-key
```

**Required parameters:**

- `<model_name>` - Model identifier/name
- `-c, --company` - Company that develops the model (e.g., Anthropic, DeepSeek, Alibaba)
- `-v, --vendor` - API provider/vendor identifier (e.g., anthropic, deepseek, dashscope, moonshot)
- `-u, --endpoint` - API endpoint URL
- `-k, --api-key` - API key

**Understanding vendor vs company:**

- **Vendor** (`-v, --vendor`): The API provider/vendor identifier. This is typically the organization that provides the API endpoint (e.g., "anthropic", "deepseek", "dashscope", "moonshot", "volcengine", "streamlake", "minimax", "bigmodel")

- **Company** (`-c, --company`): The company that develops the AI model. This is the organization that created the model (e.g., "Anthropic", "DeepSeek", "Alibaba", "Moonshot AI", "ByteDance", "Kuaishou", "MiniMax", "ZHIPU AI")

### Remove Model

```bash
modix remove my-custom-model
```

### Show Configuration Path

Displays the path to the configuration file.

```bash
modix path
```

### Show Vendor Details

Shows detailed information for a specific vendor, including all models and the API key (for debugging purposes).

```bash
modix show <vendor>
```

**Example:**

```bash
modix show deepseek
```

Output:

```shell
Vendor Details:

Vendor ID:    deepseek
Company:      DeepSeek Tech
API Endpoint: https://api.deepseek.com/v1
API Key:      sk-xxxxxxxxxxxxxxxxxxxx

Models:
  - deepseek-reasoner
  - deepseek-chat
```

## Configuration

Modix stores configuration in `~/.modix/settings.json` (or `%APPDATA%\modix\settings.json` on Windows).

### Configuration Structure

The configuration file uses a nested structure to organize models by vendor. Here's the current format:

```json
{
  "current_vendor": "anthropic",
  "current_model": "Claude",
  "default_vendor": "anthropic",
  "default_model": "Claude",
  "vendors": {
    "anthropic": {
      "company": "Anthropic",
      "api_endpoint": "",
      "api_key": "",
      "models": [
        "Claude"
      ]
    },
    "deepseek": {
      "company": "DeepSeek",
      "api_endpoint": "https://api.deepseek.com/v1",
      "api_key": "your-deepseek-api-key",
      "models": [
        "deepseek-reasoner",
        "deepseek-chat"
      ]
    },
    "bailian": {
      "company": "Alibaba",
      "api_endpoint": "https://dashscope.aliyuncs.com/compatible-mode/v1",
      "api_key": "your-qwen-api-key",
      "models": [
        "qwen3-coder-plus",
        "qwen3-coder-flash"
      ]
    },
    "volcengine": {
      "company": "ByteDance",
      "api_endpoint": "https://ark.cn-beijing.volces.com/api/coding",
      "api_key": "your-doubao-api-key",
      "models": [
        "doubao-seed-code-preview-latest"
      ]
    },
    "moonshot": {
      "company": "Moonshot AI",
      "api_endpoint": "https://api.moonshot.cn/anthropic",
      "api_key": "your-kimi-api-key",
      "models": [
        "kimi-k2-thinking-turbo"
      ]
    },
    "streamlake": {
      "company": "Kuaishou",
      "api_endpoint": "https://wanqing.streamlakeapi.com/api/gateway/v1/endpoints/ep-xxx-xxx/claude-code-proxy",
      "api_key": "your-kat-api-key",
      "models": [
        "KAT-Coder"
      ]
    },
    "minimax": {
      "company": "MiniMax",
      "api_endpoint": "https://api.minimaxi.com/anthropic",
      "api_key": "your-minimax-api-key",
      "models": [
        "MiniMax-M2"
      ]
    },
    "bigmodel": {
      "company": "ZHIPU AI",
      "api_endpoint": "https://open.bigmodel.cn/api/anthropic",
      "api_key": "your-zhipu-api-key",
      "models": [
        "GLM-4.6"
      ]
    }
  },
  "config_version": "1.0.0",
  "created_at": "2025-11-17T00:00:00Z",
  "updated_at": "2025-11-20T00:00:00Z"
}
```

Key changes from the old structure:

- **Vendor-based organization**: Models are now grouped by vendor in a `vendors` object
- **Array-based models**: Each vendor has a `models` array containing all model names
- **Shared credentials**: API endpoint and key are shared across all models in a vendor
- **Current/default tracking**: Both vendor and model are tracked for current and default selections

### Default Models and Vendors

The following table lists all default models included when running `modix init`:

| Vendor       | Company     | API Endpoint                                                                              | Models                            | Configuration                 |
| ------------ | ----------- | ----------------------------------------------------------------------------------------- | --------------------------------- | ----------------------------- |
| `anthropic`  | Anthropic   | N/A (official)                                                                            | Claude                            | Pre-configured in Claude Code |
| `deepseek`   | DeepSeek    | `https://api.deepseek.com/v1`                                                             | deepseek-reasoner, deepseek-chat  | Manual setup required         |
| `dashscope`  | Alibaba     | `https://dashscope.aliyuncs.com/compatible-mode/v1`                                       | qwen3-coder-plus, qwen3-coder-32b | Manual setup required         |
| `volcengine` | ByteDance   | `https://ark.cn-beijing.volces.com/api/coding`                                            | doubao-seed-code-preview-latest   | Manual setup required         |
| `moonshot`   | Moonshot AI | `https://api.moonshot.cn/anthropic`                                                       | kimi-k2-thinking-turbo            | Manual setup required         |
| `StreamLake` | Kuaishou    | `https://wanqing.streamlakeapi.com/api/gateway/v1/endpoints/ep-xxx-xxx/claude-code-proxy` | KAT-Coder                         | Manual setup required         |
| `minimax`    | MiniMax     | `https://api.minimaxi.com/anthropic`                                                      | MiniMax-M2                        | Manual setup required         |
| `bigmodel`   | ZHIPU AI    | `https://open.bigmodel.cn/api/anthropic`                                                  | GLM-4.6                           | Manual setup required         |

## Commands Reference

All commands provide **color-coded output** for enhanced readability:

| Command                   | Description                      | Example                                                                               |
| ------------------------- | -------------------------------- | ------------------------------------------------------------------------------------- |
| `modix init`              | Initialize default configuration | `modix init`                                                                          |
| `modix list`              | List all configured models       | `modix list`                                                                          |
| `modix check claude-code` | Check Claude Code configuration  | `modix check claude-code`                                                             |
| `modix check modix`       | Check Modix configuration health | `modix check modix`                                                                   |
| `modix switch Claude`     | Switch to Claude Official API    | `modix switch "Claude"`                                                               |
| `modix switch <model>`    | Switch to specified model        | `modix switch "deepseek-reasoner"`                                                    |
| `modix status`            | Show current model status        | `modix status`                                                                        |
| `modix add <name>`        | Add new model                    | `modix add my-model --vendor custom --endpoint https://api.example.com -k my-api-key` |
| `modix remove <name>`     | Remove model                     | `modix remove my-model`                                                               |
| `modix path`              | Show config file path            | `modix path`                                                                          |

## Development

### Project Structure

```bash
src/
├── main.rs           # CLI entry point (in development)
├── lib.rs            # Library exports
├── config.rs         # Configuration data structures
└── config_manager.rs # Configuration management logic
```

### Building

```bash
cargo build
```

### Running Tests

```bash
cargo test
```

### Running with Arguments

```bash
cargo run -- list
cargo run -- switch claude
```

### Code Formatting

```bash
cargo fmt
```

### Linting

```bash
cargo clippy
```

### Development Status

The project is currently in active development with the following components implemented:

- **Configuration System**: Complete data structures and management logic
- **Model Definitions**: All planned model types and vendor support
- **CLI Framework**: Basic structure in place, commands being implemented
- **Core Commands**: `init`, `list`, `status`, `add`, `remove`, `path`, `show`, `switch`, `update`, `check` implemented
- **Enhanced Features**:
  - Color-coded output for better user experience
  - Special handling for Anthropic (pre-configured in Claude Code)
  - Health check validation with vendor-specific logic
  - Display status indicators (`[Y]`, `[N]`, `[-]`) for configuration status
- **API Integration**: Planning phase, to be implemented next

**Current Focus**: Testing CLI commands and implementing API integrations

## Contributing

We welcome contributions! This project is currently in active development.

### Getting Started

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes
4. Add tests for new functionality
5. Run `cargo test` to ensure everything works
6. Submit a pull request

### Development Guidelines

- Follow Rust community standards and conventions
- Use `cargo fmt` to format code before submitting
- Add appropriate tests for new features
- Update documentation for any API changes
- Ensure all commits pass `cargo clippy` linting

### Current Development Priorities

- CLI command implementation
- API integration with various LLM providers
- Comprehensive test coverage
- Error handling and user feedback

## Security

⚠️ **Important Security Note**

- API keys are currently stored in plain text in the configuration file
- The configuration file is set to 600 permissions (owner read/write only) on Unix-like systems
- **For production use, consider using environment variables or secure key management systems**
- Never commit configuration files containing API keys to version control

### Recommended Security Practices

1. Use environment variables for API keys when possible
2. Set appropriate file permissions on configuration files
3. Regularly rotate API keys
4. Monitor API usage for unauthorized access
5. Consider using encrypted storage solutions for sensitive configurations

## License

MIT License - see [LICENSE](LICENSE) file for details.

## CI/CD

This project uses GitHub Actions for automated builds and releases:

- **Automatic builds**: When code is pushed to the `main` branch, binaries are automatically built for all supported platforms
- **Cross-platform support**: Linux x86_64, macOS x86_64, macOS ARM64, and Windows x86_64
- **Automatic releases**: New GitHub releases are created with all platform binaries
- **Testing**: All builds include running the complete test suite

### Build Status

[![Release](https://github.com/promacanthus/modix/actions/workflows/release.yml/badge.svg)](https://github.com/promacanthus/modix/actions/workflows/release.yml)

### Download Binaries

Pre-built binaries are automatically generated and uploaded to [GitHub Releases](https://github.com/promacanthus/modix/releases) for each commit to `main`.

### Manual Trigger

You can manually trigger a build and release by going to the [Actions tab](https://github.com/promacanthus/modix/actions) and clicking "Run workflow" on the Release workflow.

### Configuration

See [CI/CD Guide](./CI_CD_GUIDE.md) for detailed information about the CI/CD configuration and usage.

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/promacanthus/modix/issues) - Bug reports and feature requests
- **Discussions**: [GitHub Discussions](https://github.com/promacanthus/modix/discussions) - Questions and community discussions
- **Documentation**: Check the [docs/](docs/) folder for detailed technical specifications

## Changelog

### v0.1.0 (November 20, 2025)

- ✅ Initial project setup with Rust and Cargo
- ✅ Core configuration system implemented
- ✅ Model definitions and data structures completed
- ✅ Basic CLI framework structure in place
- ✅ CLI commands implemented: `init`, `list`, `status`, `add`, `remove`, `path`, `show`, `switch`, `update`, `check`
- ✅ Enhanced check command help documentation with clear parameter options
- ✅ Added colorful output for better user experience
  - Key-value pairs in status, path, show commands
  - JSON content in check command
  - Success messages with color coding
- ✅ Special handling for Anthropic vendor
  - Health checks automatically skip Anthropic validation
  - List command displays `[ - ]` for Anthropic endpoint and API key status
  - Clear documentation of pre-configured vs manual setup requirements
- 🔄 API integrations planned for next phase

## Future Plans

- 🔄 Complete CLI command implementation
- 🔄 Integrate with Claude Official API
- 🔄 Add support for DeepSeek, Qwen, Doubao, Kimi, MiniMax APIs
- 🔄 GUI configuration interface
- 🔄 Model performance benchmarking
- 🔄 Automatic model selection based on task type
- 🔄 Integration with popular IDEs and editors
