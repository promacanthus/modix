# Modix TUI 设计文档

## 概述

Modix TUI 是一个基于 Bubbletea 框架的终端用户界面应用，提供交互式的多智能体协作管理体验。

## 设计原则

### 1. 保持 CLI 兼容性
- 所有 TUI 功能都必须支持 CLI 命令调用
- TUI 是 CLI 的增强，不是替代
- 脚本调用仍然使用 CLI 模式

### 2. 渐进式体验
- TUI 提供交互式引导
- CLI 提供快速自动化
- 用户可以根据场景选择

### 3. 状态驱动
- 使用 FSM（有限状态机）管理界面状态
- 每个状态对应一个 TUI view
- 状态转换清晰可见

## 架构设计

### 1. 应用结构

```
modix/
├── cmd/modix/
│   ├── main.go                    # 程序入口（支持 CLI 和 TUI 模式）
│   └── commands/                  # CLI 命令（保持不变）
│       ├── root.go
│       └── project/
│           ├── init.go
│           ├── check.go
│           ├── validate.go
│           └── inspect.go
├── internal/
│   ├── tui/                       # 新增：TUI 相关代码
│   │   ├── models/               # Bubbletea models
│   │   │   ├── main.go           # 主应用 model
│   │   │   ├── project.go        # 项目管理 model
│   │   │   ├── shell.go          # Shell 管理 model
│   │   │   ├── brain.go          # Brain 管理 model
│   │   │   ├── agent.go          # Agent 管理 model
│   │   │   └── runtime.go        # Runtime 管理 model
│   │   ├── views/                # TUI views
│   │   │   ├── main.go           # 主界面
│   │   │   ├── project.go        # 项目视图
│   │   │   ├── shell.go          # Shell 视图
│   │   │   ├── brain.go          # Brain 视图
│   │   │   ├── agent.go          # Agent 视图
│   │   │   └── runtime.go        # Runtime 视图
│   │   ├── components/           # 可复用 UI 组件
│   │   │   ├── list.go           # 列表组件
│   │   │   ├── form.go           # 表单组件
│   │   │   ├── table.go          # 表格组件
│   │   │   └── status.go         # 状态组件
│   │   └── styles/               # 样式定义
│   │       ├── colors.go         # 颜色主题
│   │       └── layout.go         # 布局定义
│   └── project/                  # 业务逻辑（保持不变）
│       ├── init.go
│       ├── check.go
│       ├── validate.go
│       └── inspect.go
```

### 2. 运行模式

#### CLI 模式（默认）
```bash
mx project init
mx project check --format json
```

#### TUI 模式
```bash
mx tui                    # 启动 TUI 应用
mx                        # 如果没有子命令，也启动 TUI
```

#### 混合模式
```bash
mx project init --tui     # 在 TUI 中执行初始化
```

## TUI 界面设计

### 1. 主界面（Dashboard）

```
┌─────────────────────────────────────────────────────────┐
│  Modix v1.0.0 - Multi-Agent Orchestration              │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  [Project]  [Shells]  [Brains]  [Agents]  [Runtimes]   │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Project Status: Active                          │ │
│  │  Location: /Users/zhubowen/project               │ │
│  │  Last Updated: 2026-01-23 17:28:15               │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Quick Actions:                                   │ │
│  │  [I] Initialize Project    [C] Check Dependencies │ │
│  │  [V] Validate Config       [S] View Status        │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Recent Activity:                                 │ │
│  │  • 2026-01-23 17:28:15 - Project initialized      │ │
│  │  • 2026-01-23 17:28:15 - Dependency check passed  │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  Press [Tab] to navigate, [Enter] to select, [q] quit │
└─────────────────────────────────────────────────────────┘
```

### 2. 项目管理界面

#### 项目列表视图
```
┌─────────────────────────────────────────────────────────┐
│  Projects                                               │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] my-project-1                                │ │
│  │      Location: /path/to/project1                  │ │
│  │      Agents: 3, Runtimes: 5                      │ │
│  │      Created: 2026-01-20                         │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [ ] my-project-2                                │ │
│  │      Location: /path/to/project2                  │ │
│  │      Agents: 0, Runtimes: 0                      │ │
│  │      Created: 2026-01-22                         │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [N] New Project    [D] Delete    [Enter] Select      │
└─────────────────────────────────────────────────────────┘
```

#### 项目初始化向导
```
┌─────────────────────────────────────────────────────────┐
│  Initialize New Project                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  Project Name: [___________________________]           │
│                                                         │
│  Location:     [~/Documents/projects/______]           │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Configuration Template:                          │ │
│  │  [ ] Basic (shells, brains, agents)              │ │
│  │  [✓] Full (all components)                       │ │
│  │  [ ] Custom (select components)                   │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Required Tools Check:                            │ │
│  │  ✓ git             (version: 2.50.1)             │ │
│  │  ✓ claude-code     (version: 2.1.7)              │ │
│  │  ✓ codex-cli       (version: 0.58.0)             │ │
│  │  ✓ gemini-cli      (version: 0.16.0)             │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [Enter] Initialize    [Esc] Cancel                   │
└─────────────────────────────────────────────────────────┘
```

### 3. Shell 管理界面

```
┌─────────────────────────────────────────────────────────┐
│  Shell Registry                                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] claude-code                                │ │
│  │      Binary: /usr/local/bin/claude-code         │ │
│  │      Version: 2.1.7                             │ │
│  │      Capabilities: [code, chat, edit]           │ │
│  │      Status: Active                             │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] codex-cli                                  │ │
│  │      Binary: /usr/local/bin/codex               │ │
│  │      Version: 0.58.0                            │ │
│  │      Capabilities: [code, chat]                 │ │
│  │      Status: Active                             │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [R] Register New    [U] Update    [D] Delete        │
└─────────────────────────────────────────────────────────┘
```

### 4. Brain 管理界面

```
┌─────────────────────────────────────────────────────────┐
│  Brain Registry                                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] claude-sonnet                              │ │
│  │      Provider: Anthropic                         │ │
│  │      Model: claude-3-sonnet-20240229            │ │
│  │      Endpoint: https://api.anthropic.com        │ │
│  │      Status: Active                             │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] gpt-4o                                     │ │
│  │      Provider: OpenAI                            │ │
│  │      Model: gpt-4o                              │ │
│  │      Endpoint: https://api.openai.com           │ │
│  │      Status: Active                             │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [A] Add Brain     [E] Edit    [D] Delete            │
└─────────────────────────────────────────────────────────┘
```

### 5. Agent 管理界面

```
┌─────────────────────────────────────────────────────────┐
│  Agent Definitions                                     │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] planner                                     │ │
│  │      Role: Planning Agent                        │ │
│  │      Default Shell: claude-code                  │ │
│  │      Default Brain: claude-sonnet                │ │
│  │      Capabilities: [plan, analyze]               │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] executor                                    │ │
│  │      Role: Execution Agent                       │ │
│  │      Default Shell: codex-cli                    │ │
│  │      Default Brain: gpt-4o                       │ │
│  │      Capabilities: [code, edit, test]            │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [D] Define New    [B] Bind    [C] Compose Runtime    │
└─────────────────────────────────────────────────────────┘
```

### 6. Runtime 管理界面

```
┌─────────────────────────────────────────────────────────┐
│  Agent Runtimes                                        │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] planner-runtime-1                          │ │
│  │      Agent: planner                              │ │
│  │      Shell: claude-code                          │ │
│  │      Brain: claude-sonnet                        │ │
│  │      Status: Active                              │ │
│  │      Validation: ✓ Passed                        │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  [✓] executor-runtime-1                         │ │
│  │      Agent: executor                             │ │
│  │      Shell: codex-cli                            │ │
│  │      Brain: gpt-4o                               │ │
│  │      Status: Active                              │ │
│  │      Validation: ✓ Passed                        │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [C] Compose New    [V] Validate    [S] Status        │
└─────────────────────────────────────────────────────────┘
```

### 7. 状态监控界面

```
┌─────────────────────────────────────────────────────────┐
│  Runtime Status Monitor                                │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Active Runtimes: 2                              │ │
│  │  Idle: 2, Executing: 0, Completed: 0, Failed: 0 │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  Runtime Details:                                 │ │
│  │  ┌─────────────────────────────────────────────┐ │ │
│  │  │ planner-runtime-1                          │ │ │
│  │  │ Status: IDLE                               │ │ │
│  │  │ Last Activity: 2026-01-23 17:28:15         │ │ │
│  │  └─────────────────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────────────────┐ │ │
│  │  │ executor-runtime-1                         │ │ │
│  │  │ Status: IDLE                               │ │ │
│  │  │ Last Activity: 2026-01-23 17:28:15         │ │ │
│  │  └─────────────────────────────────────────────┘ │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [R] Refresh    [H] History    [Esc] Back             │
└─────────────────────────────────────────────────────────┘
```

### 8. 历史记录界面

```
┌─────────────────────────────────────────────────────────┐
│  Activity History                                      │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  2026-01-23 17:28:15 - Project initialized       │ │
│  │      Target: my-project                          │ │
│  │      Status: ✓ Success                           │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  2026-01-23 17:28:15 - Dependency check passed   │ │
│  │      Target: git, claude-code, codex-cli         │ │
│  │      Status: ✓ Success                           │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  ┌───────────────────────────────────────────────────┐ │
│  │  2026-01-23 17:28:15 - Config validated          │ │
│  │      Target: shells.json                         │ │
│  │      Status: ✓ Success                           │ │
│  └───────────────────────────────────────────────────┘ │
│                                                         │
│  [F] Filter    [E] Export    [Esc] Back               │
└─────────────────────────────────────────────────────────┘
```

## Bubbletea Models 设计

### 1. Main Model (主应用模型)

```go
package models

type MainModel struct {
    currentView string
    projectView ProjectModel
    shellView   ShellModel
    brainView   BrainModel
    agentView   AgentModel
    runtimeView RuntimeModel
    status      string
    error       error
}

func (m MainModel) Init() tea.Cmd {
    return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "q", "ctrl+c":
            return m, tea.Quit
        case "tab":
            return m, m.switchView()
        case "p":
            m.currentView = "project"
            return m, nil
        case "s":
            m.currentView = "shell"
            return m, nil
        case "b":
            m.currentView = "brain"
            return m, nil
        case "a":
            m.currentView = "agent"
            return m, nil
        case "r":
            m.currentView = "runtime"
            return m, nil
        }
    }
    return m, nil
}

func (m MainModel) View() string {
    switch m.currentView {
    case "project":
        return m.projectView.View()
    case "shell":
        return m.shellView.View()
    case "brain":
        return m.brainView.View()
    case "agent":
        return m.agentView.View()
    case "runtime":
        return m.runtimeView.View()
    default:
        return m.dashboardView()
    }
}
```

### 2. Project Model (项目管理模型)

```go
package models

type ProjectModel struct {
    projects     []Project
    selected     int
    mode         string // "list", "init", "validate", "inspect"
    form         ProjectForm
    loading      bool
    error        error
}

type Project struct {
    Name      string
    Location  string
    Agents    int
    Runtimes  int
    CreatedAt string
}

type ProjectForm struct {
    Name     string
    Location string
    Template string
}

func (m ProjectModel) Init() tea.Cmd {
    return m.loadProjects()
}

func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch m.mode {
        case "list":
            return m.updateList(msg)
        case "init":
            return m.updateInit(msg)
        case "validate":
            return m.updateValidate(msg)
        case "inspect":
            return m.updateInspect(msg)
        }
    case projectsLoadedMsg:
        m.projects = msg.projects
        m.loading = false
        return m, nil
    case errorMsg:
        m.error = msg.err
        return m, nil
    }
    return m, nil
}

func (m ProjectModel) View() string {
    switch m.mode {
    case "list":
        return m.listView()
    case "init":
        return m.initView()
    case "validate":
        return m.validateView()
    case "inspect":
        return m.inspectView()
    default:
        return m.listView()
    }
}
```

### 3. Shell Model (Shell 管理模型)

```go
package models

type ShellModel struct {
    shells      []Shell
    selected    int
    mode        string // "list", "register", "inspect"
    form        ShellForm
    loading     bool
    error       error
}

type Shell struct {
    Name            string
    Binary          string
    RequiredVersion string
    Capabilities    []string
    Status          string
    RegisteredAt    string
}

type ShellForm struct {
    Name            string
    Binary          string
    RequiredVersion string
    Capabilities    string
}

func (m ShellModel) Init() tea.Cmd {
    return m.loadShells()
}

func (m ShellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Similar pattern to ProjectModel
    return m, nil
}

func (m ShellModel) View() string {
    // Render shell list/register/inspect views
    return ""
}
```

### 4. Brain Model (Brain 管理模型)

```go
package models

type BrainModel struct {
    brains      []Brain
    selected    int
    mode        string // "list", "add", "validate"
    form        BrainForm
    loading     bool
    error       error
}

type Brain struct {
    Provider  string
    Model     string
    Endpoint  string
    Auth      map[string]string
    Params    map[string]interface{}
    Status    string
    CreatedAt string
}

type BrainForm struct {
    Provider string
    Model    string
    Endpoint string
    APIKey   string
}

func (m BrainModel) Init() tea.Cmd {
    return m.loadBrains()
}

func (m BrainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Similar pattern
    return m, nil
}

func (m BrainModel) View() string {
    // Render brain list/add/validate views
    return ""
}
```

### 5. Agent Model (Agent 管理模型)

```go
package models

type AgentModel struct {
    agents      []Agent
    selected    int
    mode        string // "list", "define", "bind", "compose"
    form        AgentForm
    loading     bool
    error       error
}

type Agent struct {
    Role            string
    Description     string
    ManifestVersion string
    Defaults        map[string]string
    Capabilities    []string
    DefinedAt       string
}

type AgentForm struct {
    Role        string
    Description string
    Shell       string
    Brain       string
}

func (m AgentModel) Init() tea.Cmd {
    return m.loadAgents()
}

func (m AgentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Similar pattern
    return m, nil
}

func (m AgentModel) View() string {
    // Render agent list/define/bind/compose views
    return ""
}
```

### 6. Runtime Model (Runtime 管理模型)

```go
package models

type RuntimeModel struct {
    runtimes    []Runtime
    selected    int
    mode        string // "list", "compose", "validate", "status"
    form        RuntimeForm
    loading     bool
    error       error
}

type Runtime struct {
    Agent      string
    Shell      string
    Brain      string
    Status     string
    Validation map[string]string
    ComposedAt string
}

type RuntimeForm struct {
    Agent string
    Shell string
    Brain string
}

func (m RuntimeModel) Init() tea.Cmd {
    return m.loadRuntimes()
}

func (m RuntimeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // Similar pattern
    return m, nil
}

func (m RuntimeModel) View() string {
    // Render runtime list/compose/validate/status views
    return ""
}
```

## Components (可复用组件)

### 1. List Component

```go
package components

type List struct {
    Items       []string
    Selected    int
    Title       string
    PageSize    int
    ShowIndex   bool
}

func (l List) View() string {
    // Render a scrollable list with selection highlight
    return ""
}
```

### 2. Form Component

```go
package components

type Form struct {
    Fields      []FormField
    Title       string
    SubmitLabel string
}

type FormField struct {
    Label    string
    Value    string
    Required bool
    Type     string // "text", "password", "select"
    Options  []string
}

func (f Form) View() string {
    // Render a form with input fields
    return ""
}
```

### 3. Table Component

```go
package components

type Table struct {
    Headers []string
    Rows    [][]string
    Title   string
}

func (t Table) View() string {
    // Render a table with aligned columns
    return ""
}
```

### 4. Status Component

```go
package components

type Status struct {
    Title   string
    Message string
    Type    string // "info", "success", "warning", "error"
    Loading bool
}

func (s Status) View() string {
    // Render a status message with appropriate styling
    return ""
}
```

## Styles (样式定义)

### 1. Colors (颜色主题)

```go
package styles

var (
    Primary   = lipgloss.Color("#5BCEFA")
    Secondary = lipgloss.Color("#F5A9B8")
    Success   = lipgloss.Color("#50C878")
    Warning   = lipgloss.Color("#FFD700")
    Error     = lipgloss.Color("#FF6B6B")
    Info      = lipgloss.Color("#87CEEB")
    Text      = lipgloss.Color("#FFFFFF")
    Border    = lipgloss.Color("#666666")
)
```

### 2. Layout (布局定义)

```go
package styles

var (
    // Main container
    Container = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(Border).
        Padding(1, 2).
        Width(80).
        Height(24)

    // Header
    Header = lipgloss.NewStyle().
        Bold(true).
        Foreground(Primary).
        Padding(0, 1)

    // List item
    ListItem = lipgloss.NewStyle().
        Padding(0, 1).
        Foreground(Text)

    // Selected item
    SelectedItem = lipgloss.NewStyle().
        Padding(0, 1).
        Bold(true).
        Foreground(Primary).
        Background(lipgloss.Color("#333333"))

    // Form field
    FormField = lipgloss.NewStyle().
        Padding(0, 1).
        Foreground(Text)

    // Status message
    StatusInfo = lipgloss.NewStyle().
        Padding(0, 1).
        Foreground(Info)

    StatusSuccess = lipgloss.NewStyle().
        Padding(0, 1).
        Foreground(Success)

    StatusWarning = lipgloss.NewStyle().
        Padding(0, 1).
        Foreground(Warning)

    StatusError = lipgloss.NewStyle().
        Padding(0, 1).
        Foreground(Error)
)
```

## CLI 兼容性

### 1. 命令行参数

```bash
# CLI 模式（默认）
mx project init
mx project check --format json

# TUI 模式
mx tui                    # 启动 TUI 应用
mx                        # 如果没有子命令，也启动 TUI

# 混合模式
mx project init --tui     # 在 TUI 中执行初始化
```

### 2. 环境变量

```bash
# 强制使用 CLI 模式
MODIX_MODE=cli mx project init

# 强制使用 TUI 模式
MODIX_MODE=tui mx project init

# TUI 配置
MODIX_TUI_THEME=dark        # 主题：dark, light, mono
MODIX_TUI_WIDTH=120         # 界面宽度
MODIX_TUI_HEIGHT=40         # 界面高度
```

### 3. 配置文件

```toml
# ~/.config/modix/config.toml
[ui]
mode = "auto"  # auto, cli, tui
theme = "dark"
width = 120
height = 40

[cli]
format = "human"  # human, json
color = true
```

## 实现步骤

### Phase 1: 基础架构 (1-2 天)
1. 添加 Bubbletea 依赖
2. 创建 TUI 目录结构
3. 实现 MainModel 和主界面
4. 添加基础样式和组件

### Phase 2: 项目管理 TUI (2-3 天)
1. 实现 ProjectModel
2. 创建项目列表视图
3. 创建项目初始化向导
4. 创建项目验证和检查视图

### Phase 3: Shell/Brain/Agent TUI (3-4 天)
1. 实现 ShellModel 和视图
2. 实现 BrainModel 和视图
3. 实现 AgentModel 和视图
4. 实现 RuntimeModel 和视图

### Phase 4: 高级功能 (2-3 天)
1. 状态监控界面
2. 历史记录界面
3. 搜索和过滤功能
4. 导出和导入功能

### Phase 5: 测试和优化 (1-2 天)
1. 功能测试
2. 性能优化
3. 错误处理
4. 文档更新

## 参考资源

- [Bubbletea 官方文档](https://github.com/charmbracelet/bubbletea)
- [Lipgloss 样式库](https://github.com/charmbracelet/lipgloss)
- [Charmbracelet 示例](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [TUI 设计模式](https://github.com/charmbracelet/bubbletea/discussions)

## 注意事项

1. **保持 CLI 兼容性**：所有功能必须支持 CLI 调用
2. **错误处理**：TUI 必须优雅地处理错误
3. **性能**：TUI 应该响应迅速，避免卡顿
4. **可访问性**：考虑键盘导航和屏幕阅读器支持
5. **国际化**：未来可能需要支持多语言
