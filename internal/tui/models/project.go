package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ProjectModel represents the project management TUI state
type ProjectModel struct {
	projects []Project
	selected int
	mode     string // "list", "init", "validate", "inspect"
	form     ProjectForm
	loading  bool
	error    error
	status   string
	width    int
	height   int
}

// Project represents a modix project
type Project struct {
	Name      string
	Location  string
	Agents    int
	Runtimes  int
	CreatedAt string
}

// ProjectForm represents the project initialization form
type ProjectForm struct {
	Name     string
	Location string
	Template string
}

// NewProjectModel creates a new project model
func NewProjectModel() ProjectModel {
	return ProjectModel{
		mode:    "list",
		loading: false,
		projects: []Project{
			{
				Name:      "my-project-1",
				Location:  "/path/to/project1",
				Agents:    3,
				Runtimes:  5,
				CreatedAt: "2026-01-20",
			},
			{
				Name:      "my-project-2",
				Location:  "/path/to/project2",
				Agents:    0,
				Runtimes:  0,
				CreatedAt: "2026-01-22",
			},
		},
		form: ProjectForm{
			Template: "full",
		},
	}
}

// Init initializes the model
func (m ProjectModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m ProjectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

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
	}

	return m, nil
}

// updateList handles updates in list mode
func (m ProjectModel) updateList(msg tea.KeyMsg) (ProjectModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < len(m.projects)-1 {
			m.selected++
		}
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
	case "enter":
		m.mode = "inspect"
	case "i":
		m.mode = "init"
		m.form = ProjectForm{}
	case "v":
		m.mode = "validate"
	case "d":
		// Delete project (would call internal/project package)
		m.status = "Delete functionality not implemented yet"
	}
	return m, nil
}

// updateInit handles updates in init mode
func (m ProjectModel) updateInit(msg tea.KeyMsg) (ProjectModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Initialize project (would call internal/project package)
		m.status = fmt.Sprintf("Initializing project: %s", m.form.Name)
		m.mode = "list"
	case "tab":
		// Cycle through form fields
		// This is simplified - in real implementation would track current field
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateValidate handles updates in validate mode
func (m ProjectModel) updateValidate(msg tea.KeyMsg) (ProjectModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Validate project (would call internal/project package)
		m.status = "Validation complete"
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateInspect handles updates in inspect mode
func (m ProjectModel) updateInspect(msg tea.KeyMsg) (ProjectModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// View renders the current view
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

// listView renders the project list view
func (m ProjectModel) listView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Projects")

	var items []string
	for i, project := range m.projects {
		itemStyle := lipgloss.NewStyle().
			Padding(0, 1).
			Foreground(lipgloss.Color("#FFFFFF"))

		if i == m.selected {
			itemStyle = itemStyle.
				Bold(true).
				Foreground(lipgloss.Color("#5BCEFA")).
				Background(lipgloss.Color("#333333"))
		}

		item := fmt.Sprintf(
			"[%s] %s\n    Location: %s\n    Agents: %d, Runtimes: %d\n    Created: %s",
			m.getCheckbox(i),
			project.Name,
			project.Location,
			project.Agents,
			project.Runtimes,
			project.CreatedAt,
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
		Render("[N] New Project • [D] Delete • [Enter] Select • [V] Validate • [Esc] Back")

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

// initView renders the project initialization view
func (m ProjectModel) initView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Initialize New Project")

	nameField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Project Name: [%s]", m.form.Name))

	locationField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Location:     [%s]", m.form.Location))

	templateBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A9B8")).Render("Configuration Template:\n") +
				"  [ ] Basic (shells, brains, agents)\n" +
				"  [✓] Full (all components)\n" +
				"  [ ] Custom (select components)",
		)

	checkBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Required Tools Check:\n") +
				"  ✓ git             (version: 2.50.1)\n" +
				"  ✓ claude-code     (version: 2.1.7)\n" +
				"  ✓ codex-cli       (version: 0.58.0)\n" +
				"  ✓ gemini-cli      (version: 0.16.0)",
		)

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Initialize • [Esc] Cancel")

	sections := []string{
		title,
		"",
		nameField,
		"",
		locationField,
		"",
		templateBox,
		"",
		checkBox,
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

// validateView renders the validation view
func (m ProjectModel) validateView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Validate Configuration")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Validation Results:\n") +
				"  ✓ shells.json - Valid JSON format\n" +
				"  ✓ brains.json - Valid JSON format\n" +
				"  ✓ agents.json - Valid JSON format\n" +
				"  ✓ runtimes.json - Valid JSON format\n" +
				"  ✓ projects.json - Valid JSON format\n" +
				"  ✓ state.json - Valid JSON format\n" +
				"  ✓ version.json - Valid JSON format\n" +
				"\n" +
				"  ✓ All configuration files are valid",
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

// inspectView renders the project inspection view
func (m ProjectModel) inspectView() string {
	if m.selected >= len(m.projects) {
		return "No project selected"
	}

	project := m.projects[m.selected]

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render(fmt.Sprintf("Project: %s", project.Name))

	details := []string{
		fmt.Sprintf("Location:  %s", project.Location),
		fmt.Sprintf("Agents:    %d", project.Agents),
		fmt.Sprintf("Runtimes:  %d", project.Runtimes),
		fmt.Sprintf("Created:   %s", project.CreatedAt),
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
func (m ProjectModel) getCheckbox(index int) string {
	if index == m.selected {
		return "✓"
	}
	return " "
}
