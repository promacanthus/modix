package models

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InputModel represents the input-based TUI state
type InputModel struct {
	input      string
	history    []string
	historyPos int
	focus      bool
	width      int
	height     int
}

// NewInputModel creates a new input model
func NewInputModel() InputModel {
	return InputModel{
		input:      "",
		history:    []string{},
		historyPos: -1,
		focus:      true,
	}
}

// Init initializes the model
func (m InputModel) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if m.input != "" {
				m.history = append(m.history, m.input)
				m.historyPos = len(m.history)
				m.input = ""
			}
			return m, nil

		case "up":
			if m.historyPos > 0 {
				m.historyPos--
				m.input = m.history[m.historyPos]
			}
			return m, nil

		case "down":
			if m.historyPos < len(m.history)-1 {
				m.historyPos++
				m.input = m.history[m.historyPos]
			} else {
				m.historyPos = len(m.history)
				m.input = ""
			}
			return m, nil

		case "tab":
			// Auto-complete for slash commands
			m.autoComplete()
			return m, nil

		case "backspace", "ctrl+h":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil

		case "ctrl+u":
			m.input = ""
			return m, nil

		case "ctrl+w":
			// Delete last word
			if len(m.input) > 0 {
				lastSpace := strings.LastIndex(m.input, " ")
				if lastSpace == -1 {
					m.input = ""
				} else {
					m.input = m.input[:lastSpace]
				}
			}
			return m, nil

		default:
			// Handle regular character input
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
			return m, nil
		}
	}

	return m, nil
}

// autoComplete provides tab completion for slash commands
func (m *InputModel) autoComplete() {
	slashCommands := []string{
		"/project",
		"/shell",
		"/brain",
		"/agent",
		"/runtime",
		"/help",
		"/status",
	}

	if strings.HasPrefix(m.input, "/") {
		for _, cmd := range slashCommands {
			if strings.HasPrefix(cmd, m.input) {
				m.input = cmd + " "
				return
			}
		}
	}
}

// View renders the current view
func (m InputModel) View() string {
	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#5BCEFA")).
		Padding(0, 1).
		Render("Modix v1.0.0 - Multi-Agent Orchestration")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(0, 1).
		Render("Terminal User Interface - Type /help for commands")

	// Calculate available width for text inside the input box
	// Width minus padding (2 spaces) and prompt "mx> " (4 chars)
	availableWidth := m.width - 6
	if availableWidth < 1 {
		availableWidth = 1
	}

	// Wrap input text to available width
	wrappedInput := wrapText(m.input, availableWidth)

	// Input prompt
	inputPrompt := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5BCEFA")).
		Bold(true).
		Render("mx> ")

	// Highlight the slash command
	var inputContent string
	if strings.HasPrefix(m.input, "/") {
		parts := strings.Fields(m.input)
		if len(parts) > 0 {
			cmd := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#F5A9B8")).
				Bold(true).
				Render(parts[0])
			args := ""
			if len(parts) > 1 {
				args = " " + strings.Join(parts[1:], " ")
			}
			inputContent = cmd + lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(args)
		} else {
			inputContent = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(m.input)
		}
	} else {
		inputContent = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Render(m.input)
	}

	// Cursor - positioned based on wrapped input
	cursor := ""
	if m.focus {
		cursor = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5BCEFA")).
			Bold(true).
			Render("▌")
	}

	// For wrapped input, we need to handle cursor positioning differently
	// For now, we'll append cursor at the end (simple approach)
	// In a more advanced implementation, we'd calculate cursor position
	var displayContent string
	if wrappedInput != "" {
		displayContent = wrappedInput + cursor
	} else {
		displayContent = inputContent + cursor
	}

	// Create horizontal lines for input box
	// Extend to full terminal width (minus 2 spaces for padding before prompt)
	// Ensure width is at least 10 to avoid negative repeat count
	lineWidth := m.width - 2
	if lineWidth < 10 {
		lineWidth = 10
	}
	horizontalLine := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#5BCEFA")).
		Render(strings.Repeat("─", lineWidth))

	// Assemble input box with horizontal lines
	inputBox := strings.Join([]string{
		horizontalLine,
		"  " + inputPrompt + displayContent,
		horizontalLine,
	}, "\n")

	// Available commands - only show when user is typing a slash command
	var commandsBox string
	if strings.HasPrefix(m.input, "/") {
		commands := []string{
			"/project  - Manage projects",
			"/shell    - Manage shells",
			"/brain    - Manage brains/models",
			"/agent    - Manage agents",
			"/runtime  - Manage runtimes",
			"/status   - Show system status",
			"/help     - Show all commands",
		}

		commandsBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#666")).
			Padding(1, 2).
			Width(m.width - 4).
			Render(
				lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#F5A9B8")).Render("Available Commands:\n") +
					"  " + strings.Join(commands, "\n  "),
			)
	}

	// Help
	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888")).
		Padding(0, 1).
		Render(
			"Press [Enter] to execute • [Tab] to autocomplete • [↑/↓] history • [Ctrl+U] clear • [q] quit",
		)

	// Assemble
	sections := []string{
		title,
		subtitle,
		"",
		inputBox,
	}
	if commandsBox != "" {
		sections = append(sections, "", commandsBox)
	}
	sections = append(sections, "", help)

	return strings.Join(sections, "\n")
}

// wrapText wraps text to fit within the specified width
func wrapText(text string, width int) string {
	if width <= 0 || text == "" {
		return text
	}

	var result strings.Builder
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	currentLine := ""
	for _, word := range words {
		// Check if adding this word would exceed the width
		testLine := currentLine
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if len(testLine) <= width {
			// Word fits on current line
			if currentLine != "" {
				currentLine += " "
			}
			currentLine += word
		} else {
			// Word doesn't fit, start new line
			if currentLine != "" {
				result.WriteString(currentLine)
				result.WriteString("\n")
			}
			currentLine = word
		}
	}

	// Add the last line
	if currentLine != "" {
		result.WriteString(currentLine)
	}

	return result.String()
}

// GetInput returns the current input
func (m InputModel) GetInput() string {
	return m.input
}

// SetInput sets the input
func (m *InputModel) SetInput(input string) {
	m.input = input
}
