package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ShellModel represents the shell management TUI state
type ShellModel struct {
	shells   []Shell
	selected int
	mode     string // "list", "register", "inspect"
	form     ShellForm
	loading  bool
	error    error
	width    int
	height   int
}

// Shell represents a shell entry
type Shell struct {
	Name            string
	Binary          string
	RequiredVersion string
	Capabilities    []string
	Status          string
	RegisteredAt    string
}

// ShellForm represents the shell registration form
type ShellForm struct {
	Name            string
	Binary          string
	RequiredVersion string
	Capabilities    string
}

// NewShellModel creates a new shell model
func NewShellModel() ShellModel {
	return ShellModel{
		mode:    "list",
		loading: false,
		shells: []Shell{
			{
				Name:            "claude-code",
				Binary:          "/usr/local/bin/claude-code",
				RequiredVersion: "2.0.0",
				Capabilities:    []string{"code", "chat", "edit"},
				Status:          "Active",
				RegisteredAt:    "2026-01-20",
			},
			{
				Name:            "codex-cli",
				Binary:          "/usr/local/bin/codex",
				RequiredVersion: "0.50.0",
				Capabilities:    []string{"code", "chat"},
				Status:          "Active",
				RegisteredAt:    "2026-01-20",
			},
			{
				Name:            "gemini-cli",
				Binary:          "/usr/local/bin/gemini",
				RequiredVersion: "0.15.0",
				Capabilities:    []string{"code", "chat"},
				Status:          "Active",
				RegisteredAt:    "2026-01-20",
			},
		},
	}
}

// Init initializes the model
func (m ShellModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m ShellModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case "list":
			return m.updateList(msg)
		case "register":
			return m.updateRegister(msg)
		case "inspect":
			return m.updateInspect(msg)
		}
	}

	return m, nil
}

// updateList handles updates in list mode
func (m ShellModel) updateList(msg tea.KeyMsg) (ShellModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < len(m.shells)-1 {
			m.selected++
		}
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
	case "enter":
		m.mode = "inspect"
	case "r":
		m.mode = "register"
		m.form = ShellForm{}
	case "d":
		// Delete shell (would call internal/project package)
		// m.status = "Delete functionality not implemented yet"
	}
	return m, nil
}

// updateRegister handles updates in register mode
func (m ShellModel) updateRegister(msg tea.KeyMsg) (ShellModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Register shell (would call internal/project package)
		// m.status = fmt.Sprintf("Registering shell: %s", m.form.Name)
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateInspect handles updates in inspect mode
func (m ShellModel) updateInspect(msg tea.KeyMsg) (ShellModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// View renders the current view
func (m ShellModel) View() string {
	switch m.mode {
	case "list":
		return m.listView()
	case "register":
		return m.registerView()
	case "inspect":
		return m.inspectView()
	default:
		return m.listView()
	}
}

// listView renders the shell list view
func (m ShellModel) listView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Shell Registry")

	var items []string
	for i, shell := range m.shells {
		itemStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF"))

		if i == m.selected {
			itemStyle = itemStyle.
				Bold(true).
				Foreground(lipgloss.Color("#5BCEFA")).
				Background(lipgloss.Color("#333333"))
		}

		caps := strings.Join(shell.Capabilities, ", ")
		item := fmt.Sprintf(
			"[%s] %s\n    Binary: %s\n    Version: %s\n    Capabilities: [%s]\n    Status: %s",
			m.getCheckbox(i),
			shell.Name,
			shell.Binary,
			shell.RequiredVersion,
			caps,
			shell.Status,
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
		Render("[R] Register New • [U] Update • [D] Delete • [Enter] Select • [Esc] Back")

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

// registerView renders the shell registration view
func (m ShellModel) registerView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Register New Shell")

	nameField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Name:            [%s]", m.form.Name))

	binaryField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Binary:          [%s]", m.form.Binary))

	versionField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Required Version: [%s]", m.form.RequiredVersion))

	capsField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Capabilities:    [%s]", m.form.Capabilities))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Register • [Esc] Cancel")

	sections := []string{
		title,
		"",
		nameField,
		"",
		binaryField,
		"",
		versionField,
		"",
		capsField,
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

// inspectView renders the shell inspection view
func (m ShellModel) inspectView() string {
	if m.selected >= len(m.shells) {
		return "No shell selected"
	}

	shell := m.shells[m.selected]

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render(fmt.Sprintf("Shell: %s", shell.Name))

	details := []string{
		fmt.Sprintf("Binary:          %s", shell.Binary),
		fmt.Sprintf("Required Version: %s", shell.RequiredVersion),
		fmt.Sprintf("Capabilities:    [%s]", strings.Join(shell.Capabilities, ", ")),
		fmt.Sprintf("Status:          %s", shell.Status),
		fmt.Sprintf("Registered At:   %s", shell.RegisteredAt),
	}

	detailsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(strings.Join(details, "\n"))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Esc] Back")

	sections := []string{
		title,
		"",
		detailsBox,
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
func (m ShellModel) getCheckbox(index int) string {
	if index == m.selected {
		return "✓"
	}
	return " "
}
