package ui

import (
	"poker_tui/internal/engine"

	tea "github.com/charmbracelet/bubbletea"
)

// Model is the Bubble Tea model for the entire application.
type Model struct {
	gs     *engine.GameState
	lock   *engine.Lock
	theme  Theme
	width  int
	height int

	// Options sub-state
	optionsCursor int
	optionsPrompt string // set when confirming reset

	// Error / read-only
	errorText string
}

// NewModel creates a Model from a loaded GameState and optional lock.
func NewModel(gs *engine.GameState, lock *engine.Lock) Model {
	return Model{
		gs:    gs,
		lock:  lock,
		theme: ThemeFor(gs.Options.Theme),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// GameState exposes the inner state (used by main for save-on-quit).
func (m *Model) GameState() *engine.GameState {
	return m.gs
}

// Lock exposes the file lock (used by main for cleanup).
func (m *Model) Lock() *engine.Lock {
	return m.lock
}
