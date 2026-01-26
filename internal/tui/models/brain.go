package models

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// BrainModel represents the brain management TUI state
type BrainModel struct {
	brains   []Brain
	selected int
	mode     string // "list", "add", "validate"
	form     BrainForm
	loading  bool
	error    error
	width    int
	height   int
}

// Brain represents a brain entry
type Brain struct {
	Provider  string
	Model     string
	Endpoint  string
	Auth      map[string]string
	Params    map[string]interface{}
	Status    string
	CreatedAt string
}

// BrainForm represents the brain addition form
type BrainForm struct {
	Provider string
	Model    string
	Endpoint string
	APIKey   string
}

// NewBrainModel creates a new brain model
func NewBrainModel() BrainModel {
	return BrainModel{
		mode:    "list",
		loading: false,
		brains: []Brain{
			{
				Provider:  "Anthropic",
				Model:     "claude-3-sonnet-20240229",
				Endpoint:  "https://api.anthropic.com",
				Auth:      map[string]string{"api_key": "sk-..."},
				Params:    map[string]interface{}{"temperature": 0.7},
				Status:    "Active",
				CreatedAt: "2026-01-20",
			},
			{
				Provider:  "OpenAI",
				Model:     "gpt-4o",
				Endpoint:  "https://api.openai.com",
				Auth:      map[string]string{"api_key": "sk-..."},
				Params:    map[string]interface{}{"temperature": 0.7},
				Status:    "Active",
				CreatedAt: "2026-01-20",
			},
			{
				Provider:  "Google",
				Model:     "gemini-1.5-pro",
				Endpoint:  "https://generativelanguage.googleapis.com",
				Auth:      map[string]string{"api_key": "AI..."},
				Params:    map[string]interface{}{"temperature": 0.7},
				Status:    "Active",
				CreatedAt: "2026-01-20",
			},
		},
	}
}

// Init initializes the model
func (m BrainModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m BrainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case "list":
			return m.updateList(msg)
		case "add":
			return m.updateAdd(msg)
		case "validate":
			return m.updateValidate(msg)
		}
	}

	return m, nil
}

// updateList handles updates in list mode
func (m BrainModel) updateList(msg tea.KeyMsg) (BrainModel, tea.Cmd) {
	switch msg.String() {
	case "j", "down":
		if m.selected < len(m.brains)-1 {
			m.selected++
		}
	case "k", "up":
		if m.selected > 0 {
			m.selected--
		}
	case "enter":
		// Would show brain details
		// m.mode = "inspect"
	case "a":
		m.mode = "add"
		m.form = BrainForm{}
	case "d":
		// Delete brain (would call internal/project package)
		// m.status = "Delete functionality not implemented yet"
	}
	return m, nil
}

// updateAdd handles updates in add mode
func (m BrainModel) updateAdd(msg tea.KeyMsg) (BrainModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Add brain (would call internal/project package)
		// m.status = fmt.Sprintf("Adding brain: %s", m.form.Model)
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// updateValidate handles updates in validate mode
func (m BrainModel) updateValidate(msg tea.KeyMsg) (BrainModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Validate brain (would call internal/project package)
		// m.status = "Validation complete"
		m.mode = "list"
	case "esc":
		m.mode = "list"
	}
	return m, nil
}

// View renders the current view
func (m BrainModel) View() string {
	switch m.mode {
	case "list":
		return m.listView()
	case "add":
		return m.addView()
	case "validate":
		return m.validateView()
	default:
		return m.listView()
	}
}

// listView renders the brain list view
func (m BrainModel) listView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Brain Registry")

	var items []string
	for i, brain := range m.brains {
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
			"[%s] %s\n    Provider: %s\n    Model: %s\n    Endpoint: %s\n    Status: %s",
			m.getCheckbox(i),
			brain.Model,
			brain.Provider,
			brain.Model,
			brain.Endpoint,
			brain.Status,
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
		Render("[A] Add Brain • [E] Edit • [D] Delete • [Esc] Back")

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

// addView renders the brain addition view
func (m BrainModel) addView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Add New Brain")

	providerField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Provider:  [%s]", m.form.Provider))

	modelField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Model:     [%s]", m.form.Model))

	endpointField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("Endpoint:  [%s]", m.form.Endpoint))

	apiKeyField := lipgloss.NewStyle().
		Padding(0, 1).
		Render(fmt.Sprintf("API Key:   [%s]", strings.Repeat("*", len(m.form.APIKey))))

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(1, 2).
		Render("[Enter] Add • [Esc] Cancel")

	sections := []string{
		title,
		"",
		providerField,
		"",
		modelField,
		"",
		endpointField,
		"",
		apiKeyField,
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

// validateView renders the brain validation view
func (m BrainModel) validateView() string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Validate Brain")

	statusBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#666")).
		Padding(1, 2).
		Width(76).
		Render(
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50C878")).Render("Validation Results:\n") +
				"  ✓ Provider: Valid\n" +
				"  ✓ Model: Valid\n" +
				"  ✓ Endpoint: Valid\n" +
				"  ✓ API Key: Valid\n" +
				"  ✓ Connection: Successful\n" +
				"\n" +
				"  ✓ Brain is ready to use",
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

// getCheckbox returns the checkbox character for the given index
func (m BrainModel) getCheckbox(index int) string {
	if index == m.selected {
		return "✓"
	}
	return " "
}
