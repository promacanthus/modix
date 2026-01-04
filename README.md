# Modix

A better tool for managing AI coding assistants and LLM providers

**Modix** is a command-line tool written in Go that unifies and manages multiple Large Language Model (LLM) vendors. It simplifies the complexity of switching between different AI models, allowing you to easily manage configurations for Claude, DeepSeek, Qwen, Doubao, Kimi, MiniMax, and other multiple vendors.

## Key Features

- 🔧 **Multi-provider Support**: Supports Claude, DeepSeek, Qwen, Doubao, Kimi, MiniMax, ZHIPU AI, and more
- ⚡ **Quick Switching**: Switch between different AI models with one command
- 🛡️ **Secure Configuration**: Securely store API keys and sensitive information
- 🎨 **Colorful Output**: Enhanced CLI with better user experience
- 🚀 **Cross-platform**: Supports Windows, macOS, Linux

## Installation

### Build from Source (Recommended)

```bash
git clone https://github.com/promacanthus/modix.git
cd modix
go build -o modix ./cmd/modix
```

**Prerequisites**: [Go](https://go.dev/) must be installed on your system.

## Usage

### Initialize Configuration

```bash
modix init
```

This creates a default configuration file with predefined models.

### List Available Models

```bash
modix list
```

Shows all configured models with their status.

### Switch Models

```bash
# Switch to Claude
modix switch Claude

# Switch to DeepSeek
modix switch deepseek-reasoner
```

### Check Configuration

```bash
# Check Claude Code configuration
modix check claude-code

# Check Modix configuration
modix check modix
```

### View Current Status

```bash
modix status
```

Displays information about the currently active model.

### Add Custom Model

```bash
modix add my-custom-model \
  --company my-company \
  --vendor my-vendor \
  --endpoint https://api.mycustom.com \
  --api-key your-api-key
```

### View Configuration Path

```bash
modix path
```

Shows the path to the configuration file.

## Configuration

Modix stores configuration in `~/.modix/settings.json` (or `%APPDATA%\modix\settings.json` on Windows).

### Supported Vendors

| Vendor       | Company     | API Endpoint                                                                              | Models                            |
| ------------ | ----------- | ----------------------------------------------------------------------------------------- | --------------------------------- |
| `anthropic`  | Anthropic   | N/A (official)                                                                            | Claude                            |
| `deepseek`   | DeepSeek    | `https://api.deepseek.com/v1`                                                             | deepseek-reasoner, deepseek-chat  |
| `dashscope`  | Alibaba     | `https://dashscope.aliyuncs.com/compatible-mode/v1`                                       | qwen3-coder-plus, qwen3-coder-32b |
| `volcengine` | ByteDance   | `https://ark.cn-beijing.volces.com/api/coding`                                            | doubao-seed-code-preview-latest   |
| `moonshot`   | Moonshot AI | `https://api.moonshot.cn/anthropic`                                                       | kimi-k2-thinking-turbo            |
| `StreamLake` | Kuaishou    | `https://wanqing.streamlakeapi.com/api/gateway/v1/endpoints/ep-xxx-xxx/claude-code-proxy` | KAT-Coder                         |
| `minimax`    | MiniMax     | `https://api.minimaxi.com/anthropic`                                                      | MiniMax-M2                        |
| `bigmodel`   | ZHIPU AI    | `https://open.bigmodel.cn/api/anthropic`                                                  | GLM-4.6                           |

## Security Notice ⚠️

- API keys are currently stored in plain text in the configuration file
- The configuration file is set to 600 permissions (owner read/write only) on Unix-like systems
- **For production use, consider using environment variables or secure key management systems**
- Never commit configuration files containing API keys to version control

## Development

### Project Structure

```bash
cmd/modix/              # CLI commands and main entry point
├── main.go            # Program entry point
└── commands/          # Individual command implementations
    ├── add.go         # Add model command
    ├── list.go        # List models command
    ├── switch.go      # Switch model command
    └── ...
├── internal/          # Internal packages
│   └── config/       # Configuration management
└── *.md              # Documentation files
```

### Building

```bash
go build -o modix ./cmd/modix
```

### Testing

```bash
go test ./...
```

## Contributing

We welcome contributions! Please check [Issues](https://github.com/promacanthus/modix/issues) for available tasks.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## 🇨🇳 [中文版 README](README_zh.md) | [English README](README.md)
