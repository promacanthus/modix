package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainModel represents the main TUI application state
type MainModel struct {
	currentView string
	projectView ProjectModel
	shellView   ShellModel
	brainView   BrainModel
	agentView   AgentModel
	runtimeView RuntimeModel
	status      string
	error       error
	width       int
	height      int
}

// NewMainModel creates a new main model
func NewMainModel() MainModel {
	return MainModel{
		currentView: "dashboard",
		status:      "Ready",
		projectView: NewProjectModel(),
		shellView:   NewShellModel(),
		brainView:   NewBrainModel(),
		agentView:   NewAgentModel(),
		runtimeView: NewRuntimeModel(),
	}
}

// Init initializes the model
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.switchView()
			return m, nil

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

		case "h":
			m.currentView = "history"
			return m, nil

		case "esc":
			m.currentView = "dashboard"
			return m, nil
		}
	}

	// Update current view
	var newModel tea.Model
	switch m.currentView {
	case "project":
		newModel, cmd = m.projectView.Update(msg)
		m.projectView = newModel.(ProjectModel)
	case "shell":
		newModel, cmd = m.shellView.Update(msg)
		m.shellView = newModel.(ShellModel)
	case "brain":
		newModel, cmd = m.brainView.Update(msg)
		m.brainView = newModel.(BrainModel)
	case "agent":
		newModel, cmd = m.agentView.Update(msg)
		m.agentView = newModel.(AgentModel)
	case "runtime":
		newModel, cmd = m.runtimeView.Update(msg)
		m.runtimeView = newModel.(RuntimeModel)
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View renders the current view
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
	case "history":
		return m.historyView()
	default:
		return m.dashboardView()
	}
}

// switchView cycles through main views
func (m *MainModel) switchView() {
	views := []string{"dashboard", "project", "shell", "brain", "agent", "runtime", "history"}
	for i, v := range views {
		if v == m.currentView {
			m.currentView = views[(i+1)%len(views)]
			return
		}
	}
	m.currentView = "dashboard"
}

// dashboardView renders the main dashboard
func (m MainModel) dashboardView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Modix v1.0.0 - Multi-Agent Orchestration")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(0, 1).
		Render("Terminal User Interface")

	// Navigation bar
	navItems := []string{"[P]roject", "[S]hell", "[B]rain", "[A]gent", "[R]untime", "[H]istory"}
	navBar := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("#444")).
		Padding(0, 1).
		Render(strings.Join(navItems, "  "))

	// Status section
	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Status: ") + m.status,
		)

	// Quick actions
	actions := []string{
		"[I]nitialize Project",
		"[C]heck Dependencies",
		"[V]alidate Config",
		"[S]how Status",
	}
	actionsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A9B8")).Render("Quick Actions:\n") +
				"  " + strings.Join(actions, "    "),
		)

	// Help
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render(
			"Press [Tab] to navigate • [Enter] to select • [q] to quit",
		)

	// Assemble
	sections := []string{
		title,
		subtitle,
		"",
		navBar,
		"",
		statusBox,
		"",
		actionsBox,
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

// historyView renders the history view
func (m MainModel) historyView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Activity History")

	historyItems := []string{
		"• 2026-01-23 17:28:15 - Project initialized",
		"  Target: my-project",
		"  Status: ✓ Success",
		"",
		"• 2026-01-23 17:28:15 - Dependency check passed",
		"  Target: git, claude-code, codex-cli",
		"  Status: ✓ Success",
		"",
		"• 2026-01-23 17:28:15 - Config validated",
		"  Target: shells.json",
		"  Status: ✓ Success",
	}

	historyBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Height(15).
		Render(strings.Join(historyItems, "\n"))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[F]ilter • [E]xport • [Esc] Back")

	sections := []string{
		title,
		"",
		historyBox,
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
