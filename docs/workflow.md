# 多智能体工作流

## 完整工作流概览

一个完整的多智能体协作流程包含两个主要阶段：

1. **系统初始化阶段**：配置系统、加载资源、定义 Agent
2. **运行时执行阶段**：执行任务、状态转移、协作完成

---

## 第一阶段：系统初始化工作流

### 1. Initial Configuration & Rules (初始配置和规则)

**目标**：加载系统级配置，定义运行规则

**步骤**：

1. **新建配置文件**:
   1. `~/.modix/settings.json` - 全局配置文件
   2. 项目级别的配置文件：
      - `.modix/shells.json` - Shell 配置
      - `.modix/brains.json` - Brain 配置
      - `.modix/agents.json` - Agent 定义
      - `.modix/runtimes.json` - Runtime 配置
      - `.modix/projects.json` - Project 配置
      - `.modix/state.json` - 状态记录
      - `.modix/version.json` - 版本信息
2. **读取配置文件**
   1. 全局
   2. 项目级别

3. **验证配置完整性**
   - 检查必需字段是否存在
   - 验证配置格式是否正确
   - 确认配置文件版本兼容性

4. **应用默认规则**
   - 设置默认 Shell/Brain
   - 定义优先级规则
   - 配置超时和重试策略

**输出**：系统配置对象

---

### 2. Load Shell Binary (加载 Shell 二进制文件)

**目标**：识别并验证系统中可用的 Shell CLI

**步骤**：

1. **Shell 注册**
   - 扫描系统 PATH
   - 识别可用 CLI（Claude Code, Gemini-cli, Codex-cli 等）
   - 读取 shells.json 中的注册信息

2. **版本检查**
   - 检查 Shell 二进制是否存在
   - 验证版本号是否符合要求
   - 检查 Shell 能力是否匹配

3. **能力验证**
   - 测试 Shell 基本功能
   - 验证代码编辑能力
   - 确认 repo context 支持

**输出**：可用 Shell 列表及状态

---

### 3. Setting LLM Profiles (设置 LLM 配置文件)

**目标**：配置和验证 Brain（大语言模型）

**步骤**：

1. **Brain 配置加载**
   - 读取 brains.json
   - 解析 Provider、Model、Endpoint 配置
   - 加载参数（temperature, maxTokens 等）

2. **API Key 验证**
   - 检查环境变量中的 API Key
   - 验证 Key 的有效性
   - 处理 Key 过期或无效的情况

3. **模型可用性检查**
   - 测试模型连接性
   - 验证模型名称是否合法
   - 检查 Rate Limit 配置

4. **能力映射**
   - 将模型能力映射到系统能力
   - 例如：`claude-3-5-sonnet` → `["reasoning", "code-generation"]`

**输出**：可用 Brain 列表及验证状态

---

### 4. Create Agent Definitions (创建 Agent 定义)

**目标**：定义 Agent 的身份、偏好和能力

**步骤**：

1. **Agent 身份定义**
   - 从 agents.json 读取定义
   - 设置 Agent ID、名称、角色
   - 关联 Manifest 文件

2. **Manifest 关联**
   - 加载对应的 Manifest 文件
   - 验证 Manifest 格式
   - 检查能力声明是否完整

3. **默认偏好设置**
   - 设置 default shell
   - 设置 default brain
   - 定义优先级和权重

4. **能力声明验证**
   - 检查 capabilities 字段
   - 验证输入/输出约束
   - 确认权限声明

**输出**：Agent 定义列表

---

### 5. Create Finite State Machine (创建有限状态机)

**目标**：定义 Agent 的状态和状态转移规则

**步骤**：

1. **FSM 定义**
   - 读取 FSM 配置（如 execution_fsm_v1）
   - 定义状态集合
   - 设置初始状态

2. **状态定义**
   - IDLE：执行上下文已创建，但未开始
   - EXECUTING：Shell + Brain 执行中
   - COMPLETED：执行成功完成
   - FAILED：执行失败

3. **状态转移规则**
   - 定义转移事件（start, success, error）
   - 设置转移目标
   - 配置事件发射（emit）

**输出**：FSM 定义对象

---

### 6. Binding FSM to Agent (绑定 FSM 到 Agent)

**目标**：将 Agent 与 FSM 关联，形成运行时组合

**步骤**：

1. **Agent-FSM 关联**
   - 从 Agent Manifest 读取 FSM ID
   - 查找对应的 FSM 定义
   - 建立关联关系

2. **运行时组合验证**
   - 检查 Shell 绑定是否满足
   - 验证 Brain 绑定是否满足
   - 确认兼容性（shell + brain + fsm）

3. **状态初始化**
   - 设置 FSM 初始状态（idle）
   - 初始化状态历史
   - 准备状态转移记录

**输出**：Agent Runtime 组合对象

---

### 7. Create Runtime Execution Context (创建运行时执行上下文)

**目标**：构建完整的执行上下文，准备任务执行

**步骤**：

1. **REC 构建**
   - 组合 Agent、Shell、Brain、FSM
   - 设置 run_id 和时间戳
   - 初始化 artifacts 和 events

2. **配置注入**
   - 注入 Shell 配置（binary, args, env）
   - 注入 Brain 配置（model, parameters）
   - 注入 FSM 配置（状态、历史）

3. **状态初始化**
   - 设置 FSM 初始状态为 idle
   - 准备状态转移历史
   - 初始化 beads run

4. **验证准备状态**
   - 检查所有组件是否就绪
   - 验证配置完整性
   - 确认执行环境可用

**输出**：Runtime Execution Context 对象

---

## 第二阶段：运行时工作流

### 8. Modix go (执行命令)

**目标**：启动多智能体协作，完成任务

**步骤**：

1. **命令解析**
   - 解析用户输入的命令
   - 识别目标 Agent 和任务
   - 提取任务参数

2. **执行入口**
   - 调用对应的 Agent Runtime
   - 启动 FSM 状态机
   - 开始执行流程

---

## 工作流的不同阶段

- **INIT**: 只接受 Human 输入的 Epic
- **PLANNING**: Planner Agent 会产出 TaskSpec
- **EXECUTING**: Executor Agent 执行具体的动作，如修改代码等
- **TESTING**: Tester Agent 执行测试验证的动作，如允许代码的单元测试、集成测试等，如验证数据和输出格式是否正确等
- **COMPLETED**: 只读状态，表示一个任务已经完成
- **FAILED**: Agent 执行过程中遇到错误，需要 Human 介入

## 状态转移表

INIT
└─(epic_received)→ PLANNING

PLANNING
├─(task_spec_valid)→ EXECUTING
└─(task_spec_invalid)→ FAILED

EXECUTING
├─(execution_done)→ TESTING
└─(execution_failed)→ FAILED

TESTING
├─(tests_pass)→ DONE
└─(tests_fail)→ EXECUTING 这是一个受控循环

## 交互通信步骤

### Step 1：Human → System

- 人类提交 Epic
- FSM：`INIT → PLANNING`
- **没有 beads 消息**（这是系统输入）

### Step 2：Planner 工作

**Planner 读取：**

- Epic
- Project context

Planner 通过 beads 发送消息。

**Orchestrator 做三件事：**

1. 校验 TaskSpec schema
2. 校验是否符合 Planner Manifest
3. FSM：`PLANNING → EXECUTING`

### Step 3：Executor 工作

**Executor 读取：**

- TaskSpec
- Repo

Executor 通过 beads 发送消息。

FSM：`EXECUTING → TESTING`

### Step 4：Tester 工作

**Tester 读取：**

- TaskSpec
- ExecutionResult（commit ref）

Tester 通过 beads 发送消息。

FSM 分支：

- pass → DONE
- fail → EXECUTING

## 失败路径

### Planner 失败

- TaskSpec 缺字段
- Acceptance criteria 模糊

→ FSM：`PLANNING → FAILED`

→ 人类介入

### Executor 失败

- 改了不该改的文件
- 没按 TaskSpec 来

→ Orchestrator 拒绝 execution_result

→ 仍停留在 EXECUTING

### Tester 连续失败

- FSM 在 EXECUTING ↔ TESTING 循环超过 N 次

→ FSM：`TESTING → FAILED`

→ 升级给人类或 Supervisor
