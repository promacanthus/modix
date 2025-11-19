# Modix

A Rust-based CLI tool for managing and switching between Claude API backends and other large language models.

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

Shows all configured models with their name, vendor, and API endpoint in a clean table format.

### Switch Models

```bash
modix switch "Claude Code"
```

Switches the current active model to Claude Official API.

```bash
modix switch "DeepSeek-V3.2-Exp"
```

Switches to DeepSeek V3.2.

### Check Current Status

```bash
modix status
```

Displays information about the currently active model.

### Add Custom Model

```bash
modix add my-custom-model \
  --vendor custom \
  --endpoint https://api.mycustom.com \
  -k your-api-key
```

**Short options available:**

- `-v, --vendor` - Vendor (anthropic, deepseek, alibaba, bytedance)
- `-u, --endpoint` - API endpoint URL
- `-k, --api-key` - API key

### Remove Model

```bash
modix remove my-custom-model
```

### Show Configuration Path

```bash
modix path
```

Displays the path to the configuration file.

## Configuration

Modix stores configuration in `~/.modix/settings.json` (or `%APPDATA%\modix\settings.json` on Windows).

### Configuration Structure

```json
{
  "current_model": "Claude Code",
  "default_model": "Claude Code",
  "models": {
    "Claude Code": {
      "vendor": "anthropic",
      "api_endpoint": "https://api.anthropic.com",
      "api_key": "your-anthropic-api-key-here"
    },
    "DeepSeek-V3.2-Exp": {
      "vendor": "deepseek",
      "api_endpoint": "https://api.deepseek.com/v1",
      "api_key": "your-deepseek-api-key-here"
    },
    "Qwen3-Coder": {
      "vendor": "alibaba",
      "api_endpoint": "https://dashscope.aliyuncs.com/api/v1",
      "api_key": "your-qwen-api-key-here"
    },
    "DouBao-Seed-Code": {
      "vendor": "bytedance",
      "api_endpoint": "https://open.volcanoengine.com/api/v1",
      "api_key": "your-doubao-api-key-here"
    },
    "Kimi-K2": {
      "vendor": "moonshot",
      "api_endpoint": "https://api.moonshot.cn/v1",
      "api_key": "your-kimi-api-key-here"
    },
    "KAT-Coder": {
      "vendor": "kat",
      "api_endpoint": "",
      "api_key": "your-kat-api-key-here"
    },
    "MiniMax M2": {
      "vendor": "minimax",
      "api_endpoint": "https://api.minimax.chat/v1",
      "api_key": "your-minimax-api-key-here"
    }
  },
  "config_version": "1.0.0",
  "created_at": "2025-11-17T00:00:00Z",
  "updated_at": "2025-11-17T00:00:00Z"
}
```

### Supported Providers

- **Anthropic**: Claude Official API
- **DeepSeek**: DeepSeek V3.1
- **Alibaba**: Qwen series
- **ByteDance**: Doubao Seed Code
- **Moonshot AI**: Kimi-K2
- **Kimi AI Technology**: KAT-Coder
- **MiniMax**: MiniMax M2

## Commands Reference

| Command                | Description                      | Example                                                |
| ---------------------- | -------------------------------- | ------------------------------------------------------ |
| `modix init`           | Initialize default configuration | `modix init`                                           |
| `modix list`           | List all configured models       | `modix list`                                           |
| `modix switch Claude Code` | Switch to Claude Official API    | `modix switch "Claude Code"`                           |
| `modix switch <model>` | Switch to specified model        | `modix switch "DeepSeek-V3.2-Exp"`                     |
| `modix status`         | Show current model status        | `modix status`                                         |
| `modix add <name>`     | Add new model                    | `modix add my-model --vendor custom --endpoint https://api.example.com -k my-api-key` |
| `modix remove <name>`  | Remove model                     | `modix remove my-model`                                |
| `modix path`           | Show config file path            | `modix path`                                           |

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
- **API Integration**: Planning phase, to be implemented next

**Current Focus**: Implementing CLI commands and API integrations

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

## Getting Help

- **Issues**: [GitHub Issues](https://github.com/promacanthus/modix/issues) - Bug reports and feature requests
- **Discussions**: [GitHub Discussions](https://github.com/promacanthus/modix/discussions) - Questions and community discussions
- **Documentation**: Check the [docs/](docs/) folder for detailed technical specifications

## Changelog

### v0.1.0 (November 18, 2025)

- ✅ Initial project setup with Rust and Cargo
- ✅ Core configuration system implemented
- ✅ Model definitions and data structures completed
- ✅ Basic CLI framework structure in place
- 🔄 CLI commands under development
- 🔄 API integrations planned for next phase

## Future Plans

- 🔄 Complete CLI command implementation
- 🔄 Integrate with Claude Official API
- 🔄 Add support for DeepSeek, Qwen, Doubao, Kimi, MiniMax APIs
- 🔄 GUI configuration interface
- 🔄 Model performance benchmarking
- 🔄 Automatic model selection based on task type
- 🔄 Integration with popular IDEs and editors
