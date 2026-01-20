package main

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type gotReposSuccessMsg []repo

type repo struct {
	Name string `json:"name"`
}

func getRepos() tea.Msg {
	var repos []repo
	return gotReposSuccessMsg(repos)
}

type model struct {
	textInput     textinput.Model
	help          help.Model
	keymap        keymap
	filterField   string // e.g., "level"
	isPillVisible bool   // true if a field is selected
}

type keymap struct{}

func (k keymap) ShortHelp() []key.Binding {
	return []key.Binding{
		key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "complete")),
		key.NewBinding(key.WithKeys("ctrl+n"), key.WithHelp("ctrl+n", "next")),
		key.NewBinding(key.WithKeys("ctrl+p"), key.WithHelp("ctrl+p", "prev")),
		key.NewBinding(key.WithKeys("esc"), key.WithHelp("esc", "quit")),
	}
}
func (k keymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{k.ShortHelp()}
}

func initialModel() model {
	ti := textinput.New()
	//ti.Placeholder = "Enter value..."
	//ti.Prompt = "Filter:"
	//ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	ti.Focus()
	ti.CharLimit = 50
	ti.Width = 20
	ti.ShowSuggestions = true
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Background(lipgloss.Color("245"))

	h := help.New()

	km := keymap{}

	return model{textInput: ti, help: h, keymap: km, isPillVisible: false, filterField: ""}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(getRepos, textinput.Blink)
}

var suggestions = []string{"Apples", "Ananas", "Bananas", "Oranges", "Grape"}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyTab:
			if m.textInput.ShowSuggestions && len(m.textInput.AvailableSuggestions()) > 0 {
				m.filterField = m.textInput.CurrentSuggestion()
				m.isPillVisible = true

				// 1. Reset the input buffer completely
				m.textInput.Reset()
				m.textInput.SetValue("")
				m.textInput.SetSuggestions([]string{})

				// 2. Disable suggestions for the "Value" phase
				m.textInput.ShowSuggestions = false

				// 3. IMPORTANT: Return nil to stop the tab key from
				// being passed to the textinput.Update(msg) below.
				return m, nil
			}

		case tea.KeyBackspace:
			// If input is empty and a pill is active, jump back to field selection
			if m.textInput.Value() == "" && m.isPillVisible {
				m.isPillVisible = false
				m.textInput.Reset()
				m.textInput.SetValue(m.filterField)
				m.textInput.CursorEnd()

				// Re-enable suggestions for the "Field" phase
				m.textInput.ShowSuggestions = true
				m.textInput.SetSuggestions(suggestions)
				return m, nil
			}

		case tea.KeyEnter, tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	case gotReposSuccessMsg:
		m.textInput.SetSuggestions(suggestions)
	}

	// Only update the input if we didn't handle a mode-switch above
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

var pillStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("15")).
	Background(lipgloss.Color("62")). // Purple background
	Padding(0, 1).
	MarginRight(1).
	Bold(true)

// Define a base style that clears the line by enforcing a width
var containerStyle = lipgloss.NewStyle().Width(80)

func (m model) View() string {
	var s strings.Builder

	if m.isPillVisible {
		// Render the Pill
		pill := pillStyle.Render(m.filterField + ":")
		s.WriteString(pill)
	}

	// Render the input
	s.WriteString(m.textInput.View())

	// Wrap everything in containerStyle. This ensures that when the line
	// gets shorter (pill disappears), the remaining space is overwritten with spaces.
	return containerStyle.Render(s.String()) + "\n\n" + m.help.View(m.keymap)
}
