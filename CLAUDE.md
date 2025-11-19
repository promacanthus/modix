# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Modix** is a Rust-based CLI tool for managing and switching between Claude API backends and other large language models. The project is in early development stage with comprehensive planning documentation but minimal implementation.

## Current State

This repository is in initial development phase:

- Only contains documentation and configuration files
- No actual Rust source code or Cargo.toml exists yet
- Product requirements are well-defined in `docs/modix_product_requirements.md`

## Build System & Development Commands

Since this is a Rust project in planning phase, the following commands will be used once implementation begins:

### Initial Setup

```bash
# Initialize new Rust project
cargo init

# Add required dependencies
cargo add reqwest serde clap
```

### Standard Rust Commands

```bash
# Build the project
cargo build

# Run the project
cargo run

# Run with arguments
cargo run -- [args]

# Run tests
cargo test

# Check code formatting
cargo fmt

# Check for issues
cargo clippy

# Build release version
cargo build --release

# Run specific test
cargo test test_name

# Run specific binary
cargo run --bin binary_name
```

## Architecture & Key Components

### Planned Technical Stack

- **Language**: Rust (native cross-platform binaries)
- **HTTP Client**: `reqwest` for API calls to various LLM backends
- **CLI Parsing**: `clap` for command-line argument handling
- **JSON Handling**: `serde` for configuration file management

### Core Functionality

The tool will manage `~/.claude/settings.json` and provide these commands:

- `modix list` - List configured models
- `modix switch <model_name>` - Switch current model
- `modix status` - Show current model

### Supported Models (Planned)

- Claude Official API
- DeepSeek V3.1
- Alibaba Qwen series
- ByteDance Doubao Seed Code
- Moonshot AI Kimi-K2
- Kimi AI Technology KAT-Coder
- MiniMax M2

### Configuration Management

- Default config path: `~/.claude/settings.json`
- Maps model names to API endpoints
- Handles API keys and authentication securely

## File Structure (Planned)

```
modix/
├── Cargo.toml          # Rust project configuration
├── src/
│   ├── main.rs         # CLI entry point
│   ├── config.rs       # Configuration management
│   ├── models.rs       # Model definitions and switching logic
│   └── api.rs          # HTTP client and API integration
├── docs/
│   └── modix_product_requirements.md  # Product requirements
└── .gitignore          # Rust/Cargo patterns
```

## Development Workflow

### Starting Implementation

1. Create `Cargo.toml` with planned dependencies
2. Set up basic CLI structure with `clap`
3. Implement configuration file management
4. Add model switching functionality
5. Integrate with various API backends

### Key Development Areas

- CLI argument parsing and validation
- JSON configuration file reading/writing
- HTTP API integration with multiple providers
- Cross-platform path handling
- Error handling and user feedback

## Important Notes

- This project is currently in planning phase with no actual implementation
- The comprehensive product requirements document (`docs/modix_product_requirements.md`) contains the complete technical specification
- The `.claude/settings.local.json` file shows the current user's Claude configuration structure
- Focus on building a minimal viable product first, then expanding to additional models
- Security is important when handling API keys and configuration files
