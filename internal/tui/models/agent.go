package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// AgentModel represents the agent management TUI state
type AgentModel struct {
	agents   []Agent
	selected int
	mode     string // "list", "define", "bind", "compose"
	form     AgentForm
	loading  bool
	error    error
	width    int
	height   int
}

// Agent represents an agent entry
type Agent struct {
	Role            string
	Description     string
	ManifestVersion string
	Defaults        map[string]string
	Capabilities    []string
	DefinedAt       string
}

// AgentForm represents the agent definition form
type AgentForm struct {
	Role        string
	Description string
	Shell       string
	Brain       string
}

// NewAgentModel creates a new agent model
func NewAgentModel() AgentModel {
	return AgentModel{
		mode:    "list",
		loading: false,
		agents: []Agent{
			{
				Role:            "planner",
				Description:     "Planning Agent",
				ManifestVersion: "v1",
				Defaults:        map[string]string{"shell": "claude-code", "brain": "claude-sonnet"},
				Capabilities:    []string{"plan", "analyze"},
				DefinedAt:       "2026-01-20",
			},
			{
				Role:            "executor",
				Description:     "Execution Agent",
				ManifestVersion: "v1",
				Defaults:        map[string]string{"shell": "codex-cli", "brain": "gpt-4o"},
				Capabilities:    []string{"code", "edit", "test"},
				DefinedAt:       "2026-01-20",
			},
			{
				Role:            "tester",
				Description:     "Testing Agent",
				ManifestVersion: "v1",
				Defaults:        map[string]string{"shell": "claude-code", "brain": "claude-sonnet"},
				Capabilities:    []string{"test", "validate"},
				DefinedAt:       "2026-01-20",
			},
		},
	}
}

// Init initializes the model
func (m AgentModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m AgentModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case "list":
			return m.updateList(msg)
		case "define":
			return m.updateDefine(msg)
		case "bind":
			return m.updateBind(msg)
		case "compose":
			return m.updateCompose(msg)
		}
	}

	return m, nil
}

// updateList handles updates in list mode
func (m AgentModel) updateList(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < len(m.agents)-1 {
			m.selected++
		}
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
	case "enter":
		// Would show agent details
		// m.mode = "inspect"
	case "d":
		m.mode = "define"
		m.form = AgentForm{}
	case "b":
		m.mode = "bind"
	case "c":
		m.mode = "compose"
	}
	return m, nil
}

// updateDefine handles updates in define mode
func (m AgentModel) updateDefine(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Define agent (would call internal/project package)
		// m.status = fmt.Sprintf("Defining agent: %s", m.form.Role)
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateBind handles updates in bind mode
func (m AgentModel) updateBind(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Bind agent (would call internal/project package)
		// m.status = "Agent bound successfully"
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateCompose handles updates in compose mode
func (m AgentModel) updateCompose(msg tea.KeyMsg) (AgentModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Compose runtime (would call internal/project package)
		// m.status = "Runtime composed successfully"
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// View renders the current view
func (m AgentModel) View() string {
	switch m.mode {
	case "list":
		return m.listView()
	case "define":
		return m.defineView()
	case "bind":
		return m.bindView()
	case "compose":
		return m.composeView()
	default:
		return m.listView()
	}
}

// listView renders the agent list view
func (m AgentModel) listView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Agent Definitions")

	var items []string
	for i, agent := range m.agents {
		itemStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF"))

		if i == m.selected {
			itemStyle = itemStyle.
				Bold(true).
				Foreground(lipgloss.Color("#5BCEFA")).
				Background(lipgloss.Color("#333333"))
		}

		caps := strings.Join(agent.Capabilities, ", ")
		item := fmt.Sprintf(
			"[%s] %s\n    Role: %s\n    Default Shell: %s\n    Default Brain: %s\n    Capabilities: [%s]",
			m.getCheckbox(i),
			agent.Role,
			agent.Description,
			agent.Defaults["shell"],
			agent.Defaults["brain"],
			caps,
		)

		items = append(items, itemStyle.Render(item))
	}

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Height(12).
		Render(strings.Join(items, "\n\n"))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[D] Define New • [B] Bind • [C] Compose Runtime • [Esc] Back")

	sections := []string{
		title,
		"",
		listBox,
		"",
		help,
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5BCEFA")).
		Padding(1, 2).
		Width(m.width - 4).
		Height(m.height - 2).
		Render(strings.Join(sections, "\n"))
}

// defineView renders the agent definition view
func (m AgentModel) defineView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Define New Agent")

	roleField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Role:        [%s]", m.form.Role))

	descField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Description: [%s]", m.form.Description))

	shellField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Shell:       [%s]", m.form.Shell))

	brainField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Brain:       [%s]", m.form.Brain))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Define • [Esc] Cancel")

	sections := []string{
		title,
		"",
		roleField,
		"",
		descField,
		"",
		shellField,
		"",
		brainField,
		"",
		help,
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5BCEFA")).
		Padding(1, 2).
		Width(m.width - 4).
		Height(m.height - 2).
		Render(strings.Join(sections, "\n"))
}

// bindView renders the agent binding view
func (m AgentModel) bindView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Bind Agent to Shell/Brain")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Binding Options:\n") +
				"  [ ] planner → claude-code + claude-sonnet\n" +
				"  [✓] executor → codex-cli + gpt-4o\n" +
				"  [ ] tester → claude-code + claude-sonnet\n" +
				"\n" +
				"  Select agent and configure shell/brain",
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Bind • [Esc] Back")

	sections := []string{
		title,
		"",
		statusBox,
		"",
		help,
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5BCEFA")).
		Padding(1, 2).
		Width(m.width - 4).
		Height(m.height - 2).
		Render(strings.Join(sections, "\n"))
}

// composeView renders the runtime composition view
func (m AgentModel) composeView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Compose Agent Runtime")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A9B8")).Render("Runtime Composition:\n") +
				"  Agent:    [planner]\n" +
				"  Shell:    [claude-code]\n" +
				"  Brain:    [claude-sonnet]\n" +
				"  Manifest: [v1]\n" +
				"\n" +
				"  Composition: planner × claude-code × claude-sonnet",
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Compose • [Esc] Back")

	sections := []string{
		title,
		"",
		statusBox,
		"",
		help,
	}

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#5BCEFA")).
		Padding(1, 2).
		Width(m.width - 4).
		Height(m.height - 2).
		Render(strings.Join(sections, "\n"))
}

// getCheckbox returns the checkbox character for the given index
func (m AgentModel) getCheckbox(index int) string {
	if index == m.selected {
		return "✓"
	}
	return " "
}
