# Agent Capability Manifest

## 1. 文档目的

用于**描述一个智能体在系统中的“能力边界、责任范围与行为约束”**。

它的目标不是告诉系统“智能体怎么做事”，而是明确：

- **可以做什么**
- **不可以做什么**
- **在什么条件下算完成**
- **失败时系统该如何处理**

Manifest 是**设计期与调度期的契约**，而不是运行时日志。

---

## 2. Manifest 的设计原则

1. **能力与实现解耦** - 不关心具体使用 `Claude Code` / `Gemini Cli` / `Codex Cli`
2. **规则前置** - 尽量在任务分配前发现问题
3. **可演进** - 允许 v2 / v3 增量扩展字段
4. **可校验** - 系统可以基于 Manifest 自动做权限与合法性检查
5. **统一性** - 所有 Agent 使用相同的结构，便于系统处理

---

## 3. MVP 版本 Manifest

MVP 版本采用**最小化、统一的 Manifest 结构**，仅包含必需字段。

### 3.1 最小化 Manifest 结构

```json
{
  "manifest_version": "1.0",

  "agent": {
    "id": "planner|executor|tester",
    "name": "Agent Name",
    "role": "task-planner|task-executor|task-tester",
    "description": "Agent description",
    "version": "0.1.0"
  },

  "runtime": {
    "fsm": "execution_fsm_v1",
    "beads": {
      "enabled": true
    }
  },

  "bindings": {
    "shell": {
      "required": boolean,
      "allowed": ["claude-code", "gemini-cli", "codex-cli"]
    },
    "brain": {
      "required": boolean,
      "capabilities": ["reasoning", "planning", "code-generation"]
    }
  },

  "execution": {
    "entrypoint": "brain|shell",
    "expects": {
      "input": "user_goal|plan.tasks|code.diff"
    },
    "produces": {
      "artifacts": ["plan.tasks|code.diff|test.report"]
    }
  },

  "handoff": {
    "next_agent": "executor|tester",
    "via_artifact": "plan.tasks|code.diff"
  },

  "failure_policy": {
    "on_shell_error": "fail|retry",
    "on_brain_error": "fail|retry",
    "on_timeout": "fail|retry",
    "max_retries": 2
  }
}
```

---

## 4. 字段说明

### 4.1 `manifest_version` (必需)

**类型**: string

**说明**: Manifest schema 版本，用于向后兼容。

**MVP 值**: "1.0"

**注意**:

- 不是 agent 的代码版本
- schema 变化才升级

---

### 4.2 `agent` (必需)

#### 4.2.1 `agent.id`

**类型**: string

**说明**: 系统内唯一标识一个智能体。

**设计要求**:

- 稳定、不随版本变化
- 不与实现工具强绑定（不要叫 `claude-planner`）

**MVP 值**: "planner" | "executor" | "tester"

---

#### 4.2.2 `agent.name`

**类型**: string

**说明**: 人类可读的角色名称，用于文档、UI、日志展示。

---

#### 4.2.3 `agent.role`

**类型**: string

**说明**: 角色类型。

**MVP 值**:

- `task-planner` - 任务规划者
- `task-executor` - 任务执行者
- `task-tester` - 任务测试者

---

#### 4.2.4 `agent.description`

**类型**: string

**说明**: 对该 agent 核心职责的简要说明，用于帮助开发者理解其定位。

---

#### 4.2.5 `agent.version`

**类型**: string

**说明**: Agent 版本号。

---

### 4.3 `runtime` (必需)

#### 4.3.1 `runtime.fsm`

**类型**: string

**说明**: 使用的 FSM ID。

**MVP 值**: "execution_fsm_v1"

---

#### 4.3.2 `runtime.beads.enabled`

**类型**: boolean

**说明**: 是否启用 beads 追踪。

**MVP 值**: true

---

### 4.4 `bindings` (必需)

#### 4.4.1 `bindings.shell`

**类型**: object

**说明**: Shell 绑定配置。

**字段**:

- `required`: boolean - 是否必需 shell
- `allowed`: string[] - 允许的 shell 列表，空数组表示不限制

---

#### 4.4.2 `bindings.brain`

**类型**: object

**说明**: Brain 绑定配置。

**字段**:

- `required`: boolean - 是否必需 brain
- `capabilities`: string[] - 要求的 brain 能力列表

**常用能力值**:

- `reasoning` - 推理能力
- `planning` - 规划能力
- `code-generation` - 代码生成能力

---

### 4.5 `execution` (必需)

#### 4.5.1 `execution.entrypoint`

**类型**: string

**说明**: 执行入口点。

**MVP 值**:

- `"brain"` - 通过 Brain 生成内容
- `"shell"` - 通过 Shell 执行命令

---

#### 4.5.2 `execution.expects.input`

**类型**: string

**说明**: 期望的输入类型。

**MVP 值**:

- `"user_goal"` - 用户目标
- `"plan.tasks"` - 任务计划
- `"code.diff"` - 代码变更

---

#### 4.5.3 `execution.produces.artifacts`

**类型**: string[]

**说明**: 期望产出的工件类型。

**MVP 值**:

- `"plan.tasks"` - 任务计划
- `"code.diff"` - 代码变更
- `"test.report"` - 测试报告
- `"agent.result"` - Agent 执行结果

---

### 4.6 `handoff` (可选，仅多 Agent 场景)

#### 4.6.1 `handoff.next_agent`

**类型**: string

**说明**: 下一个 Agent 的 ID。

---

#### 4.6.2 `handoff.via_artifact`

**类型**: string

**说明**: 通过哪个工件传递到下一个 Agent。

---

### 4.7 `failure_policy` (必需)

#### 4.7.1 `failure_policy.on_shell_error`

**类型**: string

**说明**: shell 错误处理策略。

**MVP 值**:

- `"fail"` - 直接失败
- `"retry"` - 重试

---

#### 4.7.2 `failure_policy.on_brain_error`

**类型**: string

**说明**: brain 错误处理策略。

---

#### 4.7.3 `failure_policy.on_timeout`

**类型**: string

**说明**: 超时处理策略。

---

#### 4.7.4 `failure_policy.max_retries`

**类型**: number

**说明**: 最大重试次数。

---

## 5. 完整版本字段（未来扩展）

以下字段在 MVP 版本中可以暂时省略，但在后续版本中需要添加：

### 5.1 `capabilities` (逻辑能力声明)

**说明**: 用于任务匹配与调度决策。

**示例**:

```json
"capabilities": {
  "can_plan_tasks": true,
  "can_write_code": false,
  "can_run_tests": false,
  "can_modify_docs": false
}
```

**注意**:

- 这是**系统视角的能力**
- 不等同于 LLM 技术能力
- 即使 LLM 能写代码，若此处为 false，系统应视为越权

**未来扩展方向（v2+）**:

- 能力分级（basic / advanced）
- 领域能力（frontend / backend / infra）

---

### 5.2 `inputs` (输入约束)

**说明**: 定义 Agent 允许接收的输入类型。

**示例**:

```json
"inputs": {
  "accepted_artifacts": ["epic_spec", "project_context"]
}
```

**设计价值**:

- 防止 agent 在信息不足或不适合的情况下被调用
- 提前暴露上下游依赖问题

---

### 5.3 `outputs` (输出约束)

**说明**: 定义 Agent 被期望产出的工件类型。

**示例**:

```json
"outputs": {
  "produced_artifacts": ["task_spec"]
}
```

**注意**:

- 输出类型 ≠ 输出内容
- 内容由 beads 或 artifact 系统追踪

---

### 5.4 `permissions` (权限声明)

**说明**: 系统层面的行为许可，用于判断 agent 是否越权。

**示例**:

```json
"permissions": {
  "read": ["repo_structure", "existing_docs"],
  "write": ["task_spec"]
}
```

**与 LLM 工具权限的关系**:

- LLM 工具权限：技术上“能不能”
- Manifest 权限：逻辑上“该不该”

系统应以 Manifest 为准。

---

### 5.5 `constraints` (资源与约束)

**说明**: 为调度器提供资源决策依据，用于模型选择、并发控制。

**示例**:

```json
"constraints": {
  "max_context_tokens": 32000,
  "time_budget_sec": 120,
  "cost_priority": "medium"
}
```

**说明**:

- 不是强制执行，而是调度建议
- 未来可用于自动模型切换

---

### 5.6 `quality_gates` (质量门槛)

**说明**: 定义"什么叫完成"，可由系统或其他 agent 校验。

**示例**:

```json
"quality_gates": {
  "definition_of_done": [
    "each task has clear input/output",
    "no task exceeds 1 day of work"
  ]
}
```

**价值**:

- 防止"看起来完成但不可执行"的输出

---

### 5.7 `communication` (通信约束)

**说明**: 定义"允许与谁通信"，是调度层与治理层的规则。

**示例**:

```json
"communication": {
  "protocol": "beads",
  "can_message": ["executor", "supervisor"]
}
```

**注意**:

- beads 负责记录通信事实
- Manifest 决定通信是否合法

---

## 6. 完整示例

### 6.1 Planner Agent (MVP 版本)

```json
{
  "manifest_version": "1.0",

  "agent": {
    "id": "planner",
    "name": "Planner Agent",
    "role": "task-planner",
    "description": "Transforms high-level goals into executable task plans",
    "version": "0.1.0"
  },

  "runtime": {
    "fsm": "execution_fsm_v1",
    "beads": {
      "enabled": true
    }
  },

  "bindings": {
    "shell": {
      "required": false,
      "allowed": []
    },
    "brain": {
      "required": true,
      "capabilities": ["reasoning", "planning"]
    }
  },

  "execution": {
    "entrypoint": "brain",
    "expects": {
      "input": "user_goal"
    },
    "produces": {
      "artifacts": ["plan.tasks"]
    }
  },

  "handoff": {
    "next_agent": "executor",
    "via_artifact": "plan.tasks"
  },

  "failure_policy": {
    "on_brain_error": "fail",
    "on_timeout": "fail",
    "max_retries": 1
  }
}
```

**关键设计点**:

- ❌ 不执行代码
- ❌ 不接触 repo
- ❌ 不跑 shell
- ✅ 只产出结构化任务描述

---

### 6.2 Executor Agent (MVP 版本)

```json
{
  "manifest_version": "1.0",

  "agent": {
    "id": "executor",
    "name": "Executor Agent",
    "role": "task-executor",
    "description": "Executes planned tasks using shell and brain",
    "version": "0.1.0"
  },

  "runtime": {
    "fsm": "execution_fsm_v1",
    "beads": {
      "enabled": true
    }
  },

  "bindings": {
    "shell": {
      "required": true,
      "allowed": ["claude-code", "gemini-cli", "codex-cli"]
    },
    "brain": {
      "required": true,
      "capabilities": ["code-generation", "reasoning"]
    }
  },

  "execution": {
    "entrypoint": "shell",
    "expects": {
      "input": "plan.tasks"
    },
    "produces": {
      "artifacts": ["code.diff", "agent.result"]
    }
  },

  "handoff": {
    "next_agent": "tester",
    "via_artifact": "code.diff"
  },

  "failure_policy": {
    "on_shell_error": "fail",
    "on_brain_error": "fail",
    "on_timeout": "fail",
    "max_retries": 2
  }
}
```

**关键设计点**:

- 不负责判断对不对
- 只负责按 `plan.tasks` 干活
- 把代码改出来，把痕迹留清楚

---

### 6.3 Tester Agent (MVP 版本)

```json
{
  "manifest_version": "1.0",

  "agent": {
    "id": "tester",
    "name": "Tester Agent",
    "role": "task-tester",
    "description": "Validates executor results via tests and evaluation",
    "version": "0.1.0"
  },

  "runtime": {
    "fsm": "execution_fsm_v1",
    "beads": {
      "enabled": true
    }
  },

  "bindings": {
    "shell": {
      "required": true,
      "allowed": []
    },
    "brain": {
      "required": false,
      "capabilities": []
    }
  },

  "execution": {
    "entrypoint": "shell",
    "expects": {
      "input": "code.diff"
    },
    "produces": {
      "artifacts": ["test.report"]
    }
  },

  "failure_policy": {
    "on_shell_error": "fail",
    "on_timeout": "fail",
    "max_retries": 1
  }
}
```

**关键设计点**:

1. 只跑测试
2. 引入 Brain 做失败分析
3. 引导 Executor 重试

---

### 6.4 完整版本示例（未来扩展）

以下是一个包含所有字段的完整 Manifest 示例，用于参考：

```json
{
  "manifest_version": "1.0",

  "agent": {
    "id": "planner",
    "name": "Task Planner",
    "role": "task-planner",
    "description": "Breaks epic into executable tasks and defines contracts",
    "version": "0.1.0"
  },

  "runtime": {
    "fsm": "execution_fsm_v1",
    "beads": {
      "enabled": true
    }
  },

  "bindings": {
    "shell": {
      "required": false,
      "allowed": []
    },
    "brain": {
      "required": true,
      "capabilities": ["reasoning", "planning"]
    }
  },

  "execution": {
    "entrypoint": "brain",
    "expects": {
      "input": "user_goal"
    },
    "produces": {
      "artifacts": ["plan.tasks"]
    }
  },

  "handoff": {
    "next_agent": "executor",
    "via_artifact": "plan.tasks"
  },

  "failure_policy": {
    "on_brain_error": "fail",
    "on_timeout": "fail",
    "max_retries": 1
  },

  "capabilities": {
    "can_plan_tasks": true,
    "can_write_code": false,
    "can_run_tests": false,
    "can_modify_docs": false
  },

  "inputs": {
    "accepted_artifacts": ["epic_spec", "project_context", "requirements"]
  },

  "outputs": {
    "produced_artifacts": ["task_spec"]
  },

  "permissions": {
    "read": ["repo_structure", "existing_docs", "previous_tasks"],
    "write": ["task_spec"]
  },

  "constraints": {
    "max_context_tokens": 32000,
    "time_budget_sec": 120,
    "cost_priority": "medium"
  },

  "quality_gates": {
    "definition_of_done": [
      "each task has clear input/output",
      "no task exceeds 1 day of work"
    ]
  },

  "communication": {
    "protocol": "beads",
    "can_message": ["executor", "supervisor"]
  }
}
```

---

## 7. 三 Agent 闭环执行流

```markdown
mx run pipeline planner→executor→tester

Planner
FSM: idle → executing → completed
Artifact: plan.tasks
↓
Executor
FSM: idle → executing → completed
Artifact: code.diff
↓
Tester
FSM: idle → executing → completed
Artifact: test.report
```

如果任意一步失败：

- FSM → failed
- beads 记录失败点
- Pipeline 停止

---
