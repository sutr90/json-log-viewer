package app

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/hedhyw/json-log-viewer/internal/keymap"
	"github.com/hedhyw/json-log-viewer/internal/pkg/events"
	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

// StateFilteringModel is a state to prompt for filter term.
type StateFilteringModel struct {
	*Application

	previousState StateLoadedModel
	table         logsTableModel

	textInput widgets.PillInputModel
	keys      keymap.KeyMap
}

func newStateFiltering(
	previousState StateLoadedModel,
) StateFilteringModel {

	var s []string
	for _, f := range previousState.Config.Fields {
		s = append(s, f.Title)
	}

	textInput := widgets.NewPillInputModel(s)
	textInput.Focus()

	return StateFilteringModel{
		Application: previousState.Application,

		previousState: previousState,
		table:         previousState.table,

		textInput: textInput,
		keys:      previousState.getApplication().keys,
	}
}

// Init initializes component. It implements tea.Model.
func (s StateFilteringModel) Init() tea.Cmd {
	return s.textInput.Init()
}

// View renders component. It implements tea.Model.
func (s StateFilteringModel) View() string {
	return s.BaseStyle.Render(s.table.View()) + "\n" + s.textInput.View()
}

// Update handles events. It implements tea.Model.
func (s StateFilteringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmdBatch []tea.Cmd

	s.Application.Update(msg)

	switch msg := msg.(type) {
	case events.ErrorOccuredMsg:
		return s.handleErrorOccuredMsg(msg)
	case tea.KeyMsg:
		if mdl, cmd := s.handleKeyMsg(msg); mdl != nil {
			return mdl, cmd
		}
	default:
		s.table, cmdBatch = batched(s.table.Update(msg))(cmdBatch)
	}

	var cmd tea.Cmd
	s.textInput, cmd = s.textInput.Update(msg)
	if cmd != nil {
		cmdBatch = append(cmdBatch, cmd)
	}

	return s, tea.Batch(cmdBatch...)
}

func (s StateFilteringModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, s.keys.Back) && string(msg.Runes) != "q":
		return s.previousState.refresh()
	case key.Matches(msg, s.keys.Open):
		return s.handleEnterKeyClickedMsg()
	default:
		return nil, nil
	}
}

func (s StateFilteringModel) handleEnterKeyClickedMsg() (tea.Model, tea.Cmd) {
	filterField, input := s.textInput.Value()
	if input == "" {
		return s, events.EscKeyClicked
	}

	return initializeModel(newStateFiltered(
		s.previousState,
		input,
		filterField,
	))
}

// String implements fmt.Stringer.
func (s StateFilteringModel) String() string {
	return modelValue(s)
}
