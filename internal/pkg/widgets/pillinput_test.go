package widgets_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/hedhyw/json-log-viewer/internal/pkg/widgets"
)

func TestNewPillInputModel(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level", "message", "timestamp"}
	model := widgets.NewPillInputModel(suggestions)

	assert.NotNil(t, model)

	// Verify initial state
	filterField, value := model.Value()
	assert.Empty(t, filterField, "filter field should be empty initially")
	assert.Empty(t, value, "value should be empty initially")
}

func TestPillInputModelInit(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level", "message"}
	model := widgets.NewPillInputModel(suggestions)

	cmd := model.Init()
	assert.Nil(t, cmd, "Init should return nil command")
}

func TestPillInputModelView(t *testing.T) {
	t.Parallel()

	t.Run("without_pill", func(t *testing.T) {
		t.Parallel()

		suggestions := []string{"level", "message"}
		model := widgets.NewPillInputModel(suggestions)

		view := model.View()
		assert.NotEmpty(t, view)
		// Should not contain pill styling when no field is selected
		assert.NotContains(t, view, "level:")
		assert.NotContains(t, view, "message:")
	})

	t.Run("with_pill", func(t *testing.T) {
		t.Parallel()

		suggestions := []string{"level", "message"}
		model := widgets.NewPillInputModel(suggestions)

		// Type "lev" to get "level" suggestion
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

		// Press Tab to select suggestion
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})

		view := model.View()
		assert.NotEmpty(t, view)
		// Should contain pill with selected field
		assert.Contains(t, view, "level:")
	})
}

func TestPillInputModelUpdateTabKey(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level", "message", "timestamp"}
	model := widgets.NewPillInputModel(suggestions)

	// Type "lev" to match "level"
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'v'}})

	// Press Tab to select the suggestion
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify pill is visible
	filterField, value := model.Value()
	assert.Equal(t, "level", filterField, "filter field should be set to selected suggestion")
	assert.Empty(t, value, "value should be empty after selecting field")
	assert.Nil(t, cmd, "Tab should not return a command when selecting suggestion")
}

func TestPillInputModelUpdateBackspace(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level", "message"}
	model := widgets.NewPillInputModel(suggestions)

	// Type and select a suggestion
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})

	// Verify pill is active
	filterField, _ := model.Value()
	require.Equal(t, "level", filterField)

	// Press backspace when input is empty to go back to field selection
	model, cmd := model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	// Verify we're back to field selection mode
	filterField, value := model.Value()
	assert.Empty(t, filterField, "filter field should be cleared")
	assert.Equal(t, "level", value, "value should contain the previous filter field")
	assert.Nil(t, cmd, "Backspace should not return a command when going back")
}

func TestPillInputModelUpdateTextInput(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level", "message"}
	model := widgets.NewPillInputModel(suggestions)

	// Type some text
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'s'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})

	_, value := model.Value()
	assert.Equal(t, "test", value, "value should contain typed text")
}

func TestPillInputModelValue(t *testing.T) {
	t.Parallel()

	t.Run("no_pill_selected", func(t *testing.T) {
		t.Parallel()

		suggestions := []string{"level"}
		model := widgets.NewPillInputModel(suggestions)

		// Type some text without selecting a field
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t'}})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

		filterField, value := model.Value()
		assert.Empty(t, filterField, "filter field should be empty")
		assert.Equal(t, "te", value, "value should contain typed text")
	})

	t.Run("pill_selected_with_value", func(t *testing.T) {
		t.Parallel()

		suggestions := []string{"level"}
		model := widgets.NewPillInputModel(suggestions)

		// Select field
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})

		// Type value
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

		filterField, value := model.Value()
		assert.Equal(t, "level", filterField, "filter field should be set")
		assert.Equal(t, "er", value, "value should contain typed text after pill")
	})
}

func TestPillInputModelFocus(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level"}
	model := widgets.NewPillInputModel(suggestions)

	cmd := model.Focus()
	assert.NotNil(t, cmd, "Focus should return a command")
}

func TestPillInputModelCompleteWorkflow(t *testing.T) {
	t.Parallel()

	suggestions := []string{"level", "message", "timestamp"}
	model := widgets.NewPillInputModel(suggestions)

	// Step 1: Type to filter suggestions
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'m'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})

	// Step 2: Select suggestion with Tab
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyTab})

	filterField, value := model.Value()
	assert.Equal(t, "message", filterField)
	assert.Empty(t, value)

	// Step 3: Type the filter value
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'r'}})

	filterField, value = model.Value()
	assert.Equal(t, "message", filterField)
	assert.Equal(t, "err", value)

	// Step 4: Verify view contains pill
	view := model.View()
	assert.Contains(t, view, "message:")

	// Step 5: Backspace all the way back
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	// Should still have pill
	filterField, value = model.Value()
	assert.Equal(t, "message", filterField)
	assert.Empty(t, value)

	// One more backspace should remove pill
	model, _ = model.Update(tea.KeyMsg{Type: tea.KeyBackspace})

	filterField, value = model.Value()
	assert.Empty(t, filterField)
	assert.Equal(t, "message", value)
}
