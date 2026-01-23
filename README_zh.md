# Modix

**Declarative Multi-Agent Runtime for Real-World Engineering.**

Modix 是一个基于 Go 语言实现的、多智能体运行时（Multi-Agent Runtime），用于在真实工程项目中**声明、编排、运行和演进 AI Agent**。

Modix 的核心目标是：**让多 Agent 协作像配置系统一样可靠、可复现、可观测。**

## 为什么需要 Modix

随着 AI Agent 从“单轮对话”走向“工程协作”，问题开始发生变化：

- Agent 不再是一次性调用，而是长期运行的工作单元
- 多个 Agent 需要明确分工、边界和协作关系
- 执行过程需要可追踪、可恢复、可重放
- Agent 可能运行在本地、远程或集群中
- 工程团队希望**通过配置使用 Agent，而不是重新造轮子**

Modix 正是为这一阶段而设计的。

## 核心理念

### 1. 声明式优先（Declarative First）

在 Modix 中，你不需要“写代码组织 Agent”。

你通过声明式配置描述：

- Agent 是什么
- Agent 能做什么
- Agent 如何协作
- 工作流如何推进
- 在什么状态下发生什么行为

**配置即系统真实结构。**

### 2. Agent = 外壳 + 大脑

在 Modix 中，一个可运行的 Agent 明确由两部分组成：

- **Shell（外壳）**
  如 Claude Code、Gemini CLI、Codex CLI
  负责执行、文件操作、命令调用等
- **Brain（大脑）**
  任意 LLM 后端（Claude / GPT / Gemini / DeepSeek / Qwen 等）
  通过统一配置进行管理和切换

这使得：

- Agent 与具体模型解耦
- 模型切换不影响工程结构
- 本地 / 云端 / 私有模型统一管理

### 3. FSM 驱动的 Agent 生命周期

Modix 使用**有限状态机（FSM）**描述 Agent 和工作流的生命周期：

- 每个 Agent 都处在明确的状态中
- 状态迁移是显式、可观测、可恢复的
- 失败不是异常，而是状态的一部分

这让 Modix 天然适合：

- 长任务
- 不稳定执行环境
- 远程 Agent
- 人机混合流程

### 4. 基于事实的 Agent 通信（beads）

Agent 之间不通过隐式调用通信，而是通过**事实记录（beads）**：

- 输入是事实
- 输出是事实
- 所有行为可追踪、可重放

这让系统具备：

- 可审计性
- 可回放性
- 可恢复性
- 去中心化执行能力

### 5. 执行位置无关

在 Modix 中：

- Agent 不等同于本地进程
- CLI 不是“控制执行”，而是“提交意图”
- 执行地点是运行时细节

同一套配置可以运行在：

- 本地开发环境
- 远程服务器
- CI / CD 环境
- Kubernetes 集群

无需改变 Agent 定义。

## 项目结构（概念层）

Modix 在当前阶段主要由以下几个子系统组成：

- **CLI（mx）**
  管理配置、Agent、LLM、项目初始化和运行控制
- **Agent Runtime**
  负责加载 Manifest、驱动 FSM、执行 Agent
- **Brain Manager**
  管理和切换不同 LLM 后端
- **Shell Adapter**
  对接不同 Agent 外壳（Claude Code / Gemini CLI / Codex CLI）
- **Fact Store（beads）**
  记录和传递 Agent 之间的事实与状态

## 使用方式（高层）

Modix 面向工程用户，而不是框架开发者。

典型使用流程是：

1. 初始化项目
2. 声明 Agent（Manifest）
3. 声明工作流（FSM）
4. 配置 LLM 和 Shell
5. 运行、观察、迭代

几乎不需要编写 glue code。

## 当前状态

Modix 仍处于早期阶段，但已经明确：

- 架构方向
- 抽象边界
- 演进路径

当前重点在于：

- 稳定 Runtime Schema
- 打磨 CLI 体验
- 构建最小可运行闭环

---

## 长期愿景

Modix 希望成为：

> AI 时代的工程级 Agent Runtime
>
> 像 Terraform 管理基础设施一样，
>
> 像 Kubernetes 编排容器一样，
>
> 用配置管理 Agent。
