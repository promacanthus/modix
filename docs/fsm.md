# 有限状态机

FSM 管“活没活着”，Agent 管“干了什么”。

FSM(Finite State Machine)的职责只有三件事情：

1. 决定下一步该干什么（Action）
2. 决定失败时“该去哪”（Transition）
3. 记录“我现在在哪”（State）

## 状态定义

IDLE
EXECUTING
COMPLETED
FAILED

## FSM Runtime Schema

```json
{
  "id": "execution_fsm_v1",
  "states": ["idle", "executing", "completed", "failed"],
  "current_state": "idle",

  "history": [
    {
      "from": "idle",
      "to": "executing",
      "timestamp": "2026-01-10T09:32:15Z",
      "reason": "start_execution"
    }
  ]
}
```

### FSM 的设计铁律

- FSM **不包含业务逻辑**
- FSM **只记录状态变化**
- FSM **从不读取 artifacts**
- FSM **只被 Runtime 驱动**

### 在 Runtime 中，FSM 是这样被用的

- 启动任务 → `state = idle`
- 调用 shell + brain → `state = executing`
- shell 正常返回 → `state = completed`
- shell 抛错 / LLM 崩 → `state = failed`

## 示例

### execution_fsm_v1.json（完整示例）

```json
{
  "fsm_id": "execution_fsm_v1",
  "version": "1.0",

  "initial_state": "idle",
  "terminal_states": ["completed", "failed"],

  "states": {
    "idle": {
      "description": "Execution context created, not yet started",
      "on": {
        "start": {
          "to": "executing",
          "emit": ["fsm.transition"]
        }
      }
    },

    "executing": {
      "description": "Shell + Brain execution in progress",
      "on": {
        "success": {
          "to": "completed",
          "emit": ["fsm.transition"]
        },
        "error": {
          "to": "failed",
          "emit": ["fsm.transition"]
        }
      }
    },

    "completed": {
      "description": "Execution finished successfully"
    },

    "failed": {
      "description": "Execution failed"
    }
  }
}
```

`on.start / on.success / on.error` 是什么？

这是 **Runtime Driver 的事件契约**：

- Runtime 发 `start`
- Shell exit code == 0 → `success`
- Shell exit code != 0 → `error`

FSM **不判断对错，只接收事实**。

### agent-executor.manifest.json（完整示例）

下面以一个完整的 EXECUTOR Agent 为例：

```json
{
  "manifest_version": "1.0",

  "agent": {
    "id": "executor",
    "name": "Executor Agent",
    "role": "task-executor",
    "description": "Executes coding tasks using a configured shell and brain",
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
      "input": "task_description"
    },
    "produces": {
      "artifacts": ["shell.stdout", "agent.result"]
    }
  },

  "failure_policy": {
    "on_shell_error": "fail",
    "on_timeout": "fail"
  }
}
```

`bindings.shell / bindings.brain` 是关键连接点：

- Shell：外壳（Claude Code / Gemini-cli / Codex）
- Brain：大脑（Claude / DeepSeek / Grok）

Runtime 会做三件事：

1. 检查 bindings 是否满足
2. 注入具体配置
3. 绑定到 RuntimeExecutionContext

`execution.entrypoint = shell` 的深意是：

- Agent 的“动作起点”是 shell
- Brain 是 shell 的输入生成器
- 未来 Planner / Tester 可以不是 shell entrypoint

这是为 **多 Agent 类型**留的钩子。
