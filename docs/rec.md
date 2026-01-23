# Runtime Execution Context & workflow

## Runtime Execution Context 字段详细设计

这是最小集合：

- `agent_id`
- `brain_id`
- `shell_id`
- `fsm_id`
- `task`（用户输入 or planner 产物）
- `state`（FSM 当前状态）
- `beads_run_id`
- `artifacts`（中间产物）
- `events`（状态变化、工具调用）

可以理解为一句话：REC = Config + FSM State + beads Run + Memory

```json
{
  "run_id": "bd-run-uuid",
  "created_at": "2026-01-10T09:32:11Z",

  "agent": {
    "id": "executor",
    "version": "v1",
    "manifest_version": "1.0"
  },

  "brain": {
    "id": "claude-3-5-sonnet",
    "provider": "anthropic",
    "config_ref": "brains.json#claude-3-5-sonnet"
  },

  "shell": {
    "id": "claude-code",
    "type": "cli",
    "binary": "claude",
    "config_ref": "shells.json#claude-code"
  },

  "fsm": {
    "id": "execution_fsm_v1",
    "current_state": "idle",
    "history": []
  },

  "task": {
    "input": "Fix failing tests in repo",
    "metadata": {
      "source": "cli",
      "invocation": "mx run agent executor"
    }
  },

  "artifacts": {},
  "events": [],

  "status": "running"
}
```

## Shell × Brain：执行层字段模型

这是“Agent 外壳 + 大脑”的第一次真正协作。

### Brain Invocation Schema

```json
{
  "brain_id": "claude-3-5-sonnet",
  "prompt": {
    "system": "You are a coding agent.",
    "instruction": "Fix failing tests in repo",
    "constraints": ["Do not change public APIs"]
  },
  "parameters": {
    "temperature": 0.2,
    "max_tokens": 4096
  }
}
```

注意：Brain 只关心 **“我要生成什么”**，完全不关心 FSM / beads / shell。

### Shell Execution Schema

```json
{
  "shell_id": "claude-code",
  "command": "claude",
  "args": ["--project", ".", "--input", "prompt.txt"],
  "env": {
    "ANTHROPIC_API_KEY": "****"
  },
  "working_dir": "/repo"
}
```

Shell 是纯执行器：不理解 Agent，不理解 FSM，不理解 beads。

## 执行流

### 把所有东西连成一次真实执行（文字时序图）

```markdown
mx run agent executor
↓
load agents.json / brains.json / shells.json
↓
create RuntimeExecutionContext
↓
beads.create
↓
FSM: idle → executing
↓
brain.generate_prompt
↓
shell.execute
↓
shell.exit (success)
↓
artifact: stdout / result
↓
FSM: executing → completed
↓
beads.finalize
```

### Step-by-step 极简版

1. `mx run agent foo --task "do X"`
2. CLI：
   - 读取 `agents.json`
   - 读取 `brains.json`
   - 读取 `shells.json`
3. 构建 Runtime Execution Context
4. 创建 beads 的 epic 和 task
5. FSM：`idle → executing`
6. Shell 启动（Claude Code / Gemini CLI / Codex CLI）
7. Brain 生成 prompt / policy
8. Shell 执行
9. 结果写入 beads
10. FSM：
    - success → `completed`
    - error → `failed`
11. CLI 返回结果 / run_id
