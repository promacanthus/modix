# Modix

一个更好管理 AI 编程助手和大模型供应商的工具

**Modix** 是一个用 Go 语言编写的命令行工具，用于统一管理和切换多个大语言模型（LLM）供应商。它简化了在不同 AI 模型之间切换的复杂性，让你可以轻松管理 Claude、DeepSeek、Qwen、Doubao、Kimi、MiniMax 等多个供应商的配置。

## 核心特性

- 🔧 **多供应商支持**: 支持 Claude、DeepSeek、Qwen、Doubao、Kimi、MiniMax、ZHIPU AI 等
- ⚡ **快速切换**: 一键切换不同的 AI 模型
- 🛡️ **安全配置**: 安全存储 API 密钥和敏感信息
- 🎨 **彩色输出**: 增强的命令行界面，提供更好的用户体验
- 🚀 **跨平台**: 支持 Windows、macOS、Linux

## 安装

### 从源码安装（推荐）

```bash
git clone https://github.com/promacanthus/modix.git
cd modix
go build -o modix ./cmd/modix
```

**前提条件**: 系统需要安装 [Go](https://go.dev/)

## 使用说明

#### 初始化配置

```bash
modix init
```

这会创建一个默认配置文件，包含预定义的模型。

#### 列出可用模型

```bash
modix list
```

显示所有已配置的模型及其状态。

#### 切换模型

```bash
# 切换到 Claude
modix switch Claude

# 切换到 DeepSeek
modix switch deepseek-reasoner
```

#### 检查配置

```bash
# 检查 Claude Code 配置
modix check claude-code

# 检查 Modix 配置
modix check modix
```

#### 查看当前状态

```bash
modix status
```

显示当前激活的模型信息。

#### 添加自定义模型

```bash
modix add my-custom-model \
  --company my-company \
  --vendor my-vendor \
  --endpoint https://api.mycustom.com \
  --api-key your-api-key
```

#### 查看配置路径

```bash
modix path
```

显示配置文件的路径。

## 配置

Modix 将配置存储在 `~/.modix/settings.json`（Windows 上为 `%APPDATA%\modix\settings.json`）。

### 支持的供应商

| 供应商 | 公司 | API 端点 | 模型 |
|--------|------|----------|------|
| `anthropic` | Anthropic | N/A (官方) | Claude |
| `deepseek` | DeepSeek | `https://api.deepseek.com/v1` | deepseek-reasoner, deepseek-chat |
| `dashscope` | Alibaba | `https://dashscope.aliyuncs.com/compatible-mode/v1` | qwen3-coder-plus, qwen3-coder-32b |
| `volcengine` | ByteDance | `https://ark.cn-beijing.volces.com/api/coding` | doubao-seed-code-preview-latest |
| `moonshot` | Moonshot AI | `https://api.moonshot.cn/anthropic` | kimi-k2-thinking-turbo |
| `StreamLake` | Kuaishou | `https://wanqing.streamlakeapi.com/api/gateway/v1/endpoints/ep-xxx-xxx/claude-code-proxy` | KAT-Coder |
| `minimax` | MiniMax | `https://api.minimaxi.com/anthropic` | MiniMax-M2 |
| `bigmodel` | ZHIPU AI | `https://open.bigmodel.cn/api/anthropic` | GLM-4.6 |

## 安全说明 ⚠️

- API 密钥当前以明文形式存储在配置文件中
- 配置文件在 Unix 系统上设置为 600 权限（仅所有者可读写）
- **生产环境建议使用环境变量或安全的密钥管理系统**
- 切勿将包含 API 密钥的配置文件提交到版本控制

## 开发

### 项目结构

```bash
cmd/modix/              # CLI 命令和主入口点
├── main.go            # 程序入口点
└── commands/          # 各个命令的实现
    ├── add.go         # 添加模型命令
    ├── list.go        # 列出模型命令
    ├── switch.go      # 切换模型命令
    └── ...
├── internal/          # 内部包
│   └── config/       # 配置管理
└── *.md              # 文档文件
```

### 构建

```bash
go build -o modix ./cmd/modix
```

### 测试

```bash
go test ./...
```

## 贡献

欢迎贡献！请查看 [Issues](https://github.com/promacanthus/modix/issues) 了解待办事项。

## 许可证

MIT License - 详情请见 [LICENSE](LICENSE) 文件。

## [English README](README.md) | 🇨🇳 [中文版 README](README_zh.md)