package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainModel represents the main TUI application state
type MainModel struct {
	currentView string
	inputView   InputModel
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
		currentView: "input",
		status:      "Ready",
		inputView:   NewInputModel(),
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

		case "esc":
			// Return to input view from any sub-view
			m.currentView = "input"
			return m, nil
		}
	}

	// Update current view
	var newModel tea.Model
	switch m.currentView {
	case "input":
		newModel, cmd = m.inputView.Update(msg)
		m.inputView = newModel.(InputModel)

		// Check if enter was pressed and process slash commands
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "enter" {
			input := m.inputView.GetInput()
			if strings.HasPrefix(input, "/") {
				m.processSlashCommand(input)
			}
		}

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
	case "history":
		// History view doesn't have its own model, just handle key events
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
			m.currentView = "input"
		}

	case "status":
		// Status view doesn't have its own model, just handle key events
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "esc" {
				m.currentView = "input"
			} else if keyMsg.String() == "r" {
				// Refresh status
				m.status = "Ready"
			}
		}

	case "help":
		// Help view doesn't have its own model, just handle key events
		if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
			m.currentView = "input"
		}
	}

	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// processSlashCommand processes slash commands from the input view
func (m *MainModel) processSlashCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := strings.ToLower(parts[0])

	switch command {
	case "/project":
		m.currentView = "project"
	case "/shell":
		m.currentView = "shell"
	case "/brain":
		m.currentView = "brain"
	case "/agent":
		m.currentView = "agent"
	case "/runtime":
		m.currentView = "runtime"
	case "/help", "/h":
		m.currentView = "help"
	case "/status", "/s":
		m.currentView = "status"
	case "/history", "/hist":
		m.currentView = "history"
	default:
		// Unknown command - stay in input view
		m.currentView = "input"
	}
}

// View renders the current view
func (m MainModel) View() string {
	switch m.currentView {
	case "input":
		// Show initial welcome page if no input has been entered yet
		if m.inputView.GetInput() == "" && len(m.inputView.history) == 0 {
			return m.initialView()
		}
		return m.inputView.View()
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
	case "help":
		return m.helpView()
	case "status":
		return m.statusView()
	case "history":
		return m.historyView()
	default:
		return m.inputView.View()
	}
}

// initialView renders the initial welcome page for mx command
func (m MainModel) initialView() string {
	// Title with gradient colors
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("╔════════════════════════════════════════════════════════════════════════════╗")

	subtitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#F5A9B8")).
		Padding(0, 1).
		Render("║                           Modix v1.0.0                                    ║")

	subtitle2 := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(0, 1).
		Render("║                  Multi-Agent Orchestration Platform                       ║")

	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("╠════════════════════════════════════════════════════════════════════════════╣")

	// Features section
	features := []string{
		"  • Project Management: Initialize, validate, and inspect projects",
		"  • Shell Registry: Register and manage CLI tools (Claude, Codex, etc.)",
		"  • Brain Registry: Configure LLM providers and models",
		"  • Agent Definition: Define agent roles and capabilities",
		"  • Runtime Composer: Compose and orchestrate agent runtimes",
		"  • Execution Substrate: Execute multi-agent workflows",
	}

	featuresBox := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Render("║ " + strings.Join(features, "\\n║ ") + " ")

	// Quick Start section
	quickStartContent := []string{
		"  Quick Start:",
		"  • Type /help to see all available commands",
		"  • Type /project init <name> to create a new project",
		"  • Type /shell register <name> to register a CLI tool",
		"  • Type /agent define <role> to define an agent",
		"  • Type /status to view system status",
	}

	quickStartBox := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF")).
		Padding(0, 1).
		Render("║ " + strings.Join(quickStartContent, "\\n║ ") + " ")

	// Input section
	inputPrompt := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("╠════════════════════════════════════════════════════════════════════════════╣")

	inputLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(0, 1).
		Render("║ Type a command or press [Enter] to continue...                            ║")

	// Bottom border
	bottomBorder := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("╚════════════════════════════════════════════════════════════════════════════╝")

	// Keyboard shortcuts
	shortcuts := []string{
		"[Enter] Execute",
		"[Tab] Autocomplete",
		"[↑/↓] History",
		"[Ctrl+U] Clear",
		"[q] Quit",
	}

	shortcutsLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(0, 1).
		Render("║ " + strings.Join(shortcuts, "  │  ") + strings.Repeat(" ", m.width-2-len(strings.Join(shortcuts, "  │  "))) + " ║")

	// Assemble all sections
	sections := []string{
		title,
		subtitle,
		subtitle2,
		separator,
		featuresBox,
		separator,
		quickStartBox,
		inputPrompt,
		inputLine,
		bottomBorder,
		"",
		shortcutsLine,
	}

	return strings.Join(sections, "\n")
}

// helpView renders the help view
func (m MainModel) helpView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Help - Available Commands")

	commands := []string{
		"/project  - Manage projects (list, init, validate, inspect)",
		"/shell    - Manage shells (list, register, inspect)",
		"/brain    - Manage brains/models (list, add, validate)",
		"/agent    - Manage agents (list, define, bind, compose)",
		"/runtime  - Manage runtimes (list, compose, validate, status)",
		"/status   - Show system status and statistics",
		"/history  - View activity history",
		"/help     - Show this help message",
	}

	commandsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A9B8")).Render("Slash Commands:\n") +
				"  " + strings.Join(commands, "\n  "),
		)

	shortcuts := []string{
		"Navigation:",
		"  • Type /command and press [Enter] to execute",
		"  • Use [Tab] for autocomplete",
		"  • Use [↑/↓] to navigate command history",
		"  • Press [Ctrl+U] to clear input",
		"  • Press [Esc] to return to input view",
		"  • Press [q] or [Ctrl+C] to quit",
		"",
		"Sub-view Navigation:",
		"  • [j]/[↓] - Move down",
		"  • [k]/[↑] - Move up",
		"  • [Enter] - Select/Confirm",
		"  • [Esc] - Go back",
	}

	shortcutsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Keyboard Shortcuts:\n") +
				"  " + strings.Join(shortcuts, "\n  "),
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Esc] Back to input view")

	sections := []string{
		title,
		"",
		commandsBox,
		"",
		shortcutsBox,
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

// statusView renders the system status view
func (m MainModel) statusView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("System Status")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("System Overview:\n") +
				"  ✓ Modix v1.0.0\n" +
				"  ✓ TUI: Active\n" +
				"  ✓ Status: " + m.status + "\n" +
				"",
		)

	componentsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A9B8")).Render("Component Status:\n") +
				"  • Projects: 2 registered\n" +
				"  • Shells: 3 registered\n" +
				"  • Brains: 3 configured\n" +
				"  • Agents: 3 defined\n" +
				"  • Runtimes: 3 composed\n" +
				"",
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[R] Refresh • [Esc] Back")

	sections := []string{
		title,
		"",
		statusBox,
		"",
		componentsBox,
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
