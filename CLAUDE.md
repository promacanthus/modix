# CLAUDE.md

---

## 项目总览

本项目是一个使用 Go 语言开发的命令行工具。主要实现多智能体协作，共同完成用户提出的问题或者任务。

项目主要分为如下的 4 个模块：

1. Project Management: 用户项目的初始化，包括检查依赖的命令行工具以及新建初始化的配置。
2. Agent Runtime Management: 这是最核心的模块，负责定义、配置、组合和校验 Agent Runtime。包括如下的 4 个子模块：
   1. Shell Registry: 负责识别系统中有哪些可用的 CLI，如 Claude code、Codex、Gemini-cli 等（这些在 `modix` 项目中称为 `shell`）并管理它们的配置入口。
   2. Brain Registry: 负责管理 Brain Profile（即 Provider + LLM Model + Params）处理 API Key、Endpoint 和 Rate Limit 等，提供这些大模型的“可用性检查“。“
   3. Agent Definition Registry: 负责定义 Agent 的身份，关联 Manifest 并给出默认的偏好（如 default shell + default brain）。
   4. Agent Runtime Composer: 负责把 Agent Definition Shell Brain 组合成一个 Agent Runtime（即 Agent Runtime = Definition × Shell × Brain × Manifest）。
3. Execution Substrate: 负责定义如何向 Shell 中发送指令，如何接收输出，这是一个抽象的执行通道，编排多智能体，根据定义相互协作完成任务。
4. State & Observability: 负责记录配置的变更、Runtime 组合历史等，保证随时可以检查项目的各种状态和历史。

### 重要注意事项

1. 项目中的多智能体之间使用 [beads](https://github.com/steveyegge/beads) 作为唯一的通信标准。
2. 每个智能体使用 FSM（Finite State Machine）描述状态：
   1. IDLE
   2. EXECUTING
   3. COMPLETED
   4. FAILED
3. 由于多智能体的存在，因此整个工作流执行的过程中会有不同的阶段，包括如下所示：
   1. INIT
   2. PLANNING
   3. EXECUTING
   4. TESTING
   5. COMPLETED
   6. FAILED

---

- **项目名（概念 / 文档 / Repo）**：`modix`
- **CLI 可执行文件名**：`mx`

所有命令统一写成：

```bash
mx <command> [subcommand] [options]
```

## 重要文件

- `README.md` - Main project documentation (keep this updated!)
- `.beads/issues.jsonl` - Current issue tracking data
- `go.mod/go.sum` - Go module and dependency definitions
- `docs/examples.md` - 项目参考示例
- `docs/fsm.md` - 有限状态机以及智能体之间交互通信步骤的说明
- `docs/manifest.md` - agent capability manifest 设计文档
- `docs/rec.md` - runtime execution context 设计、示例以及整个执行的流转过程

## Milestone

### Milestone 1

Milestone 1 遵循如下的原则：

1. 命令即系统边界：CLI 的层级 = 子系统的真实边界
2. 只允许声明，不允许执行
3. 所有命令都可被脚本调用（为后续自动化铺路）

所有的配置都落到本地的文件中，控制配置的规模和对象的数量。

一条总原则：“一个逻辑域 = 一个文件；文件内是 map / array，不是散落的对象文件。”

```shell
.modix/
├── shells.json
├── brains.json
├── agents.json
├── runtimes.json
├── projects.json
├── state.json
└── version.json
```

`cli` 和 `schema` 之间的映射关系：

`mx shell register`: shells.json
`mx shell inspect`: shells.json
`mx brain add`: brains.json
`mx agent define`: agents.json
`mx agent bind`: agents.json
`mx agent runtime compose`: runtimes.json + state.json
`mx project init`: projects.json

### Milestone 2

设计一个最小可运行闭环（MVP Loop），而不是“优雅的完整系统”。让 FSM 驱动一次真实执行，并让 beads 记住发生的一切。

**闭环标准：一次任务，从“人类给目标”到“系统给结果”，中途即使失败，也能自己知道下一步该干嘛。**

最小闭环的角色配置：

- planner: 职责是规划，把目标变成“可执行任务合同”
- executor: 职责是执行，按合同改代码
- tester：职责是验收，用事实判断“过 / 不过”

一个最小闭环目标：用一个 FSM，驱动一个 Agent，经由 beads，完成一次任务执行，并留下完整可回放的执行记录。

在这里最核心的是：REC（Runtime Execution Context）这是把所有东西连在一起的地方。

---

## Issue Tracking

We use **bd** (beads) for issue tracking instead of Markdown TODOs or external tools.

### Quick Reference

```bash
# Find ready work (no blockers)
bd ready --json

# Find ready work including future deferred issues
bd ready --include-deferred --json

# Create new issue
bd create "Issue title" -t bug|feature|task -p 0-4 -d "Description" --json

# Create issue with due date and defer (GH#820)
bd create "Task" --due=+6h              # Due in 6 hours
bd create "Task" --defer=tomorrow       # Hidden from bd ready until tomorrow
bd create "Task" --due="next monday" --defer=+1h  # Both

# Update issue status
bd update <id> --status in_progress --json

# Update issue with due/defer dates
bd update <id> --due=+2d                # Set due date
bd update <id> --defer=""               # Clear defer (show immediately)

# Link discovered work
bd dep add <discovered-id> <parent-id> --type discovered-from

# Complete work
bd close <id> --reason "Done" --json

# Show dependency tree
bd dep tree <id>

# Get issue details
bd show <id> --json

# Query issues by time-based scheduling (GH#820)
bd list --deferred              # Show issues with defer_until set
bd list --defer-before=tomorrow # Deferred before tomorrow
bd list --defer-after=+1w       # Deferred after one week from now
bd list --due-before=+2d        # Due within 2 days
bd list --due-after="next monday" # Due after next Monday
bd list --overdue               # Due date in past (not closed)
```

### Workflow

1. **Check for ready work**: Run `bd ready` to see what's unblocked
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work**: If you find bugs or TODOs, create issues:
   - `bd create "Found bug in auth" -t bug -p 1 --json`
   - Link it: `bd dep add <new-id> <current-id> --type discovered-from`
5. **Complete**: `bd close <id> --reason "Implemented"`
6. **Export**: Run `bd export -o .beads/issues.jsonl` before committing

### Issue Types

- `bug` - Something broken that needs fixing
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature composed of multiple issues
- `chore` - Maintenance work (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (nice-to-have features, minor bugs)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Dependency Types

- `blocks` - Hard dependency (issue X blocks issue Y)
- `related` - Soft relationship (issues are connected)
- `parent-child` - Epic/subtask relationship
- `discovered-from` - Track issues discovered during work

Only `blocks` dependencies affect the ready work queue.

---

## Development Guidelines

### Code Standards

- **Go version**: 1.25+
- **Linting**: `golangci-lint run ./...`
- **Testing**: All new features need tests (`go test ./...`)
- **Documentation**: Update relevant .md files
- **CLI Standards**: Use Cobra framework for commands, Viper for configuration

### File Organization

```bash
modix/
├── cmd/modix/              # CLI commands and main entry point
│   ├── main.go            # Program entry point
│   └── commands/          # Individual command implementations
├── internal/               # Internal packages
├── .beads/               # Beads issue tracking system
├── .claude/              # Claude Code configuration
├── .github/              # GitHub Actions CI/CD
└── *.md                  # Documentation files
```

### Before Committing

1. **Run tests**: `go test ./...`
2. **Run linter**: `golangci-lint run ./...`
3. **Export issues**: `bd export -o .beads/issues.jsonl`
4. **Update docs**: If you changed behavior, update README.md or other docs
5. **Git add both**: `git add .beads/issues.jsonl <your-changes>`

### Git Workflow

```bash
# Make changes
git add <files>

# Export beads issues
bd export -o .beads/issues.jsonl
git add .beads/issues.jsonl

# Commit
git commit -m "Your message"

# After pull
git pull
bd import -i .beads/issues.jsonl  # Sync SQLite cache
```

Or use the git hooks in `examples/git-hooks/` for automation.

## Questions?

- Check existing issues: `bd list`
- Look at recent commits: `git log --oneline -20`
- Read the docs: README.md, AGENTS.md
- Check CLI help: `modix --help` or `modix <command> --help`
- Create an issue if unsure: `bd create "Question: ..." -t task -p 2`

## Pro Tips for Agents

- Always use `--json` flags for programmatic use of CLI commands
- Link discoveries with `discovered-from` to maintain context
- Check `bd ready` before asking "what next?"
- Export to JSONL before committing (or use git hooks)
- Use `bd dep tree` to understand complex dependencies
- Priority 0-1 issues are usually more important than 2-4
- Modix commands support color-coded output for better UX
- Anthropic models are pre-configured in Claude Code (special handling)
- Use `modix check claude-code` to validate Claude Code integration

## Building and Testing

```bash
# Build the project
go build -o modix ./cmd/modix

# Test all packages
go test ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run locally
./modix init
./modix list
./modix status
```

## Release Process (Maintainers)

1. Update version in go.mod (if applicable)
2. Update CHANGELOG.md (if exists)
3. Run full test suite: `go test ./...`
4. Tag release: `git tag v0.x.0`
5. Push tag: `git push origin v0.x.0`
6. GitHub Actions handles CI/CD automatically

---

**Remember**:

- If you find the workflow confusing or have ideas for improvement, create an issue with your feedback.
- 如果你对某些事情不确定或缺少代码库上下文，请要求澄清，而不是编造答案。
- 在修复 Bug 时，总是在注释中简要解释你的推理。

---

## Claude Skills

项目包含自定义的 Claude Code Skills，用于辅助开发工作流。

### Requirements Gatherer Skill

**用途**: 当用户提出模糊、不明确的需求时，这个 Skill 会扮演该领域的资深专家，帮助梳理需求并生成清晰、可执行的任务清单。

**使用场景**:

- 用户描述一个功能想法但细节不完整
- 需要从模糊问题中提炼出具体任务
- 需要优先级排序和范围界定
- 需要识别隐藏的需求和边缘情况

**使用方法**:

直接向 Claude 描述你的模糊需求，例如：

```text
我想给我们的应用添加搜索功能
```

Skill 会自动：

1. 询问关键的澄清问题（谁是用户、痛点是什么、约束条件等）
2. 识别隐藏的需求（边缘情况、集成点、性能要求等）
3. 使用 RICE 或 MoSCoW 框架进行优先级排序
4. 生成具体的、可执行的任务列表

**示例输出**:

```text
作为资深搜索/UX专家，我需要先澄清几个关键问题：

## 关键问题
1. 用户搜索什么内容？（产品、文档、用户、文本内容？）
2. 当前的痛点是什么？（找不到、步骤太多、没有搜索？）
3. 需要什么搜索能力？（关键词、模糊搜索、过滤器？）

## 优先级任务清单
P0 - 必须做（3-5天）:
- 基础关键词搜索
- 移动端响应式设计

P1 - 应该做（2-3天）:
- 自动补全
- 基本过滤器

P2 - 可以做（1-2天）:
- 搜索分析
- 高级过滤器
```

**Skill 文件位置**:

- `.claude/skills/requirements-gatherer/skill.md` - Skill 定义
- `.claude/skills/requirements-gatherer/skill.json` - Skill 配置

**相关 Skills**:

- `code-simplifier` - 需求明确后用于代码简化
- `plan` - 详细的技术实现规划
- `frontend-design` - UI/UX 需求收集

Happy coding!
