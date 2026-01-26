package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// RuntimeModel represents the runtime management TUI state
type RuntimeModel struct {
	runtimes []Runtime
	selected int
	mode     string // "list", "compose", "validate", "status"
	form     RuntimeForm
	loading  bool
	error    error
	width    int
	height   int
}

// Runtime represents a runtime entry
type Runtime struct {
	Agent      string
	Shell      string
	Brain      string
	Status     string
	Validation map[string]string
	ComposedAt string
}

// RuntimeForm represents the runtime composition form
type RuntimeForm struct {
	Agent string
	Shell string
	Brain string
}

// NewRuntimeModel creates a new runtime model
func NewRuntimeModel() RuntimeModel {
	return RuntimeModel{
		mode:    "list",
		loading: false,
		runtimes: []Runtime{
			{
				Agent:      "planner",
				Shell:      "claude-code",
				Brain:      "claude-sonnet",
				Status:     "Active",
				Validation: map[string]string{"agent": "✓", "shell": "✓", "brain": "✓"},
				ComposedAt: "2026-01-20",
			},
			{
				Agent:      "executor",
				Shell:      "codex-cli",
				Brain:      "gpt-4o",
				Status:     "Active",
				Validation: map[string]string{"agent": "✓", "shell": "✓", "brain": "✓"},
				ComposedAt: "2026-01-20",
			},
			{
				Agent:      "tester",
				Shell:      "claude-code",
				Brain:      "claude-sonnet",
				Status:     "Active",
				Validation: map[string]string{"agent": "✓", "shell": "✓", "brain": "✓"},
				ComposedAt: "2026-01-20",
			},
		},
	}
}

// Init initializes the model
func (m RuntimeModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m RuntimeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case "list":
			return m.updateList(msg)
		case "compose":
			return m.updateCompose(msg)
		case "validate":
			return m.updateValidate(msg)
		case "status":
			return m.updateStatus(msg)
		}
	}

	return m, nil
}

// updateList handles updates in list mode
func (m RuntimeModel) updateList(msg tea.KeyMsg) (RuntimeModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < len(m.runtimes)-1 {
			m.selected++
		}
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
	case "enter":
		// Would show runtime details
		// m.mode = "inspect"
	case "c":
		m.mode = "compose"
		m.form = RuntimeForm{}
	case "v":
		m.mode = "validate"
	case "s":
		m.mode = "status"
	}
	return m, nil
}

// updateCompose handles updates in compose mode
func (m RuntimeModel) updateCompose(msg tea.KeyMsg) (RuntimeModel, tea.Cmd) {
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

// updateValidate handles updates in validate mode
func (m RuntimeModel) updateValidate(msg tea.KeyMsg) (RuntimeModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Validate runtime (would call internal/project package)
		// m.status = "Validation complete"
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateStatus handles updates in status mode
func (m RuntimeModel) updateStatus(msg tea.KeyMsg) (RuntimeModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Check status (would call internal/project package)
		// m.status = "Status checked"
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// View renders the current view
func (m RuntimeModel) View() string {
	switch m.mode {
	case "list":
		return m.listView()
	case "compose":
		return m.composeView()
	case "validate":
		return m.validateView()
	case "status":
		return m.statusView()
	default:
		return m.listView()
	}
}

// listView renders the runtime list view
func (m RuntimeModel) listView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Agent Runtimes")

	var items []string
	for i, runtime := range m.runtimes {
		itemStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF"))

		if i == m.selected {
			itemStyle = itemStyle.
				Bold(true).
				Foreground(lipgloss.Color("#5BCEFA")).
				Background(lipgloss.Color("#333333"))
		}

		validation := fmt.Sprintf("agent=%s, shell=%s, brain=%s",
			runtime.Validation["agent"],
			runtime.Validation["shell"],
			runtime.Validation["brain"])

		item := fmt.Sprintf(
			"[%s] %s-runtime-%d\n    Agent: %s\n    Shell: %s\n    Brain: %s\n    Status: %s\n    Validation: %s",
			m.getCheckbox(i),
			runtime.Agent,
			i+1,
			runtime.Agent,
			runtime.Shell,
			runtime.Brain,
			runtime.Status,
			validation,
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
		Render("[C] Compose New • [V] Validate • [S] Status • [Esc] Back")

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

// composeView renders the runtime composition view
func (m RuntimeModel) composeView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Compose New Runtime")

	agentField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Agent: [%s]", m.form.Agent))

	shellField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Shell: [%s]", m.form.Shell))

	brainField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Brain: [%s]", m.form.Brain))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Compose • [Esc] Cancel")

	sections := []string{
		title,
		"",
		agentField,
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

// validateView renders the runtime validation view
func (m RuntimeModel) validateView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Validate Runtime")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Validation Results:\n") +
				"  ✓ Agent: planner (exists)\n" +
				"  ✓ Shell: claude-code (available)\n" +
				"  ✓ Brain: claude-sonnet (configured)\n" +
				"  ✓ Manifest: v1 (compatible)\n" +
				"\n" +
				"  ✓ Runtime is valid and ready to use",
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Done • [Esc] Back")

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

// statusView renders the runtime status view
func (m RuntimeModel) statusView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Runtime Status Monitor")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Active Runtimes: 3\n") +
				"  Idle: 3, Executing: 0, Completed: 0, Failed: 0\n" +
				"\n" +
				"  ┌─────────────────────────────────────────────┐\n" +
				"  │ planner-runtime-1                          │\n" +
				"  │ Status: IDLE                               │\n" +
				"  │ Last Activity: 2026-01-23 17:28:15         │\n" +
				"  └─────────────────────────────────────────────┘\n" +
				"  ┌─────────────────────────────────────────────┐\n" +
				"  │ executor-runtime-1                         │\n" +
				"  │ Status: IDLE                               │\n" +
				"  │ Last Activity: 2026-01-23 17:28:15         │\n" +
				"  └─────────────────────────────────────────────┘\n" +
				"  ┌─────────────────────────────────────────────┐\n" +
				"  │ tester-runtime-1                           │\n" +
				"  │ Status: IDLE                               │\n" +
				"  │ Last Activity: 2026-01-23 17:28:15         │\n" +
				"  └─────────────────────────────────────────────┘",
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[R] Refresh • [H] History • [Esc] Back")

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
func (m RuntimeModel) getCheckbox(index int) string {
	if index == m.selected {
		return "✓"
	}
	return " "
}
