# 项目参考示例

## 每个子系统的参考示例

一条总原则：每个命令都要有一个 `--json` 选项作为输出内容的格式。

1. Project Management：Project 就像一个舞台，Agent 就是其中的演员。
   1. `mx project init`
   2. `mx project inspect`
   3. `mx project config show`
   4. `mx project config validate`
2. Agent Runtime Management:
   shell 是执行外壳注册表，能看到有哪些 shell，可以注册 shell、校验 shell 是否存在，版本是否匹配等。
   1. `mx shell list`
   2. `mx shell inspect`
   3. `mx shell register`
   4. `mx shell check`
      brain 允许配置 Provider、Model、Endpoint、校验 API Key 是否存在、校验模型名是否合法等。
   5. `mx brain add`
   6. `mx brain inspect`
   7. `mx validate`
   8. `mx brain list`
      agent 定义层只涉及身份、Manifest 和默认偏好等。
   9. `mx agent define`
   10. `mx agent inspect`
   11. `mx agent list`
   12. `mx agent bind`
       agent 和 runtime 的组合
   13. `mx agent runtime compose`: 生成一个 runtime 实例
   14. `mx agent runtime status`: 是否 ready / broken
   15. `mx agent runtime validate`: 给出明确失败原因。
3. Execution Substrate：
   1. `mx run --dry`
4. State & Observability:
   查看当前配置了什么，修改了什么。
   1. `mx state status` 2.`mx state history`
5. 工程理性的守门人：把隐性失败变成显性，把“不能跑“变成“为什么不能跑“。
   doctor 是结构体检而不是 debug。
   1. `mx doctor check`
   2. `mx doctor explain`

## 每个配置文件的 schema 参考示例

`shell.json` 示例：

```json
{
  "version": "v1",
  "shells": {
    "claude-code": {
      "name": "Claude Code CLI",
      "binary": "claude",
      "requiredVersion": ">=1.0.0",
      "capabilities": ["code-edit", "repo-context"],
      "status": "unknown",
      "registeredAt": "2026-01-05T12:00:00Z"
    }
  }
}
```

`brains.json` 示例：

```json
{
  "version": "v1",
  "brains": {
    "deepseek-coder": {
      "provider": "deepseek",
      "model": "deepseek-coder",
      "endpoint": "https://api.deepseek.com/v1",
      "auth": {
        "method": "api-key",
        "env": "DEEPSEEK_API_KEY"
      },
      "params": {
        "temperature": 0.2,
        "maxTokens": 8192
      },
      "status": "unvalidated",
      "createdAt": "2026-01-05T12:10:00Z"
    }
  }
}
```

`agents.json` 示例：

```json
{
  "version": "v1",
  "agents": {
    "planner": {
      "role": "task-planner",
      "description": "High-level task planner agent",
      "manifestVersion": "v1",
      "defaults": {
        "shell": "claude-code",
        "brain": "deepseek-coder"
      },
      "capabilities": ["task-decomposition", "dependency-analysis"],
      "definedAt": "2026-01-05T12:20:00Z"
    }
  }
}
```

`runtimes.json` 示例:

```json
{
  "version": "v1",
  "runtimes": {
    "planner-runtime": {
      "agent": "planner",
      "shell": "claude-code",
      "brain": "deepseek-coder",
      "status": "ready",
      "validation": {
        "shell": "ok",
        "brain": "ok",
        "compatibility": "ok"
      },
      "composedAt": "2026-01-05T12:30:00Z"
    }
  }
}
```

`projects.json` 示例:

```json
{
  "version": "v1",
  "projects": {
    "demo-project": {
      "agents": ["planner"],
      "runtimes": ["planner-runtime"],
      "createdAt": "2026-01-05T12:40:00Z"
    }
  }
}
```

`state.json` 示例:

```json
{
  "lastUpdated": "2026-01-05T12:45:00Z",
  "counts": {
    "shells": 1,
    "brains": 1,
    "agents": 1,
    "runtimes": 1,
    "projects": 1
  },
  "history": [
    {
      "event": "agent.runtime.compose",
      "target": "planner-runtime",
      "timestamp": "2026-01-05T12:30:00Z"
    }
  ]
}
```

## Agents 参考示例

一句话定锚：

- **Planner**：把“人类意图”变成“可执行任务描述”
- **Executor**：把“任务描述”变成“代码变更”
- **Tester**：验证 Executor 的结果是否满足预期

它们**不共享内存**，只共享 **beads + artifacts**。

### Planner Manifest

Planner 是 整个闭环的“语义入口”。

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
      "required": false
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
    "on_brain_error": "fail"
  }
}
```

Planner 的关键设计点（请注意这些“克制”）

- ❌ 不执行代码
- ❌ 不接触 repo
- ❌ 不跑 shell

它只产出一个 **结构化任务描述**。

plan.tasks 的期望形态（示意）：

```json
{
  "tasks": [
    {
      "id": "fix-test-a",
      "description": "Fix failing test in test_a.py",
      "priority": 1
    }
  ]
}
```

Planner 的“成功”，不是代码好坏，而是 Executor 能不能照着。

### Executor Manifest

Executor 是 第一个真正“动手”的 Agent。

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
      "artifacts": ["shell.stdout", "code.diff", "agent.result"]
    }
  },

  "handoff": {
    "next_agent": "tester",
    "via_artifact": "code.diff"
  },

  "failure_policy": {
    "on_shell_error": "fail",
    "on_timeout": "fail"
  }
}
```

Executor 的“边界纪律”，**不负责判断对不对**，只负责：

- 按 `plan.tasks` 干活
- 把代码改出来
- 把痕迹留清楚

判断好坏，是 Tester 的事。

### Tester Manifest

Tester 是 闭环的“裁判”。

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
      "required": true
    },
    "brain": {
      "required": false
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

  "evaluation": {
    "success_condition": "all_tests_passed"
  },

  "failure_policy": {
    "on_test_failure": "fail"
  }
}
```

Tester 的关键设计点：

1. 只跑测试
2. 引入 Brain 做失败分析
3. 引导 Executor 重试

### 三 Agent 闭环执行流（文字执行图）

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
