# Modix V1 产品需求文档

## 项目概述

**项目名称**: Modix

**项目定位**: 跨平台 CLI 工具，用于管理和切换 Claude API 后端的大模型，未来可扩展到多种 Coding Agent 和 LLM。

**核心目标**:

- 简化 Claude API 使用者切换大模型的操作
- 提供跨平台（Windows、macOS、Linux）命令行工具
- 为未来扩展到多模型、多 Agent 提供技术基础

---

## 功能需求

### V1 核心功能

1. **支持的模型和服务**
   - Claude 官方 API
   - 国内兼容 Anthropic 风格 API 的大模型（包括 DeepSeek V3.2、Alibaba Qwen3 系列、ByteDance Doubao Seed Code）
   - Moonshot AI Kimi-K2
   - Kimi AI Technology KAT-Coder
   - MiniMax M2

2. **模型切换管理**
   - 通过 CLI 命令切换当前使用的模型
   - 映射 Claude 配置文件 `~/.claude/settings.json` 或等效路径
   - 支持查看当前模型信息和可用模型列表

3. **跨平台支持**
   - Windows、macOS、Linux
   - 提供原生二进制，无需额外依赖

4. **CLI 命令示例**
   - `modix list` — 列出已配置的大模型
   - `modix switch <model_name>` — 切换当前使用的大模型
   - `modix status` — 显示当前正在使用的大模型

### 非功能需求

1. **性能**: 高效读取配置文件并完成模型切换
2. **安全**: 安全管理 API Key 和配置文件，避免泄露
3. **可扩展性**: 未来可增加多 Agent 支持，如 Codex、Gemini CLI
4. **易用性**: 命令简洁明了，文档清晰

---

## 技术方案

1. **编程语言**: Rust
   - 原生跨平台，性能高，内存安全
   - 可直接生成 Windows、macOS、Linux 二进制
   - 生态成熟，支持 HTTP 请求、JSON 处理和 CLI 构建

2. **依赖库（Rust）**
   - `reqwest` — HTTP 客户端，用于调用 API
   - `serde` — JSON 解析与序列化
   - `clap` — CLI 命令解析

3. **配置文件管理**
   - 默认路径: `~/.claude/settings.json`
   - 映射模型名称到 API endpoint
   - 支持读取和修改 JSON 配置

---

## 项目名称及品牌

**最终项目名称**: Modix

**命名理由**:

- 单词简洁，易记，CLI 风格自然
- 含义明确：Model + Mix，强调多模型切换
- 科技感强，适合未来扩展和国际化

---

## 后续发展方向

1. 支持更多 Coding Agent（Codex、Gemini CLI）
2. 高性能优化，如并发请求处理和本地缓存管理
3. 增加 GUI 或 Web 管理界面（可选）
4. 多模型组合调用，动态负载均衡
5. 增加更多国产大模型支持

---

## 文档和资料

1. Claude 官方配置文档: 修改 `~/.claude/settings.json` 支持切换模型
2. 国内兼容 Anthropic API 的模型文档:
   - DeepSeek V3.2 API 文档
   - Alibaba Qwen3 系列 API 文档
   - ByteDance Doubao Seed Code 文档
   - Moonshot AI Kimi-K2 API 文档
   - Kimi AI Technology KAT-Coder 文档
   - MiniMax M2 API 文档

---

## 总结

Modix V1 是一个面向 Claude API 用户的 **跨平台 CLI 工具**，专注于 **简化大模型切换操作**，以 Rust 实现，具备高性能、安全性和未来可扩展性，为后续多 Agent、多模型的管理打下基础。

**更新**: 2025年11月19日 - 扩展支持 Moonshot AI Kimi-K2、Kimi AI Technology KAT-Coder 和 MiniMax M2 三个国内大模型。
