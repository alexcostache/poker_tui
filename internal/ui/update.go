package ui

import (
	"poker_tui/internal/engine"

	tea "github.com/charmbracelet/bubbletea"
)

// saveCmd triggers an autosave (fire-and-forget side effect).
type saveMsg struct{}

func doSave(gs *engine.GameState) tea.Cmd {
	return func() tea.Msg {
		_ = engine.SaveGame(gs)
		return saveMsg{}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case saveMsg:
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	gs := m.gs
	key := msg.String()

	// Global quit
	if key == "ctrl+c" {
		_ = engine.SaveGame(gs)
		return m, tea.Quit
	}

	// Read-only banner — only allow q to quit.
	if gs.ReadOnly && gs.Screen == engine.ScreenMainIdle {
		if key == "q" || key == "Q" {
			_ = engine.SaveGame(gs)
			return m, tea.Quit
		}
		return m, nil
	}

	// ---- Error screen -------------------------------------------------------
	if gs.Screen == engine.ScreenErrorScreen {
		switch key {
		case "q", "Q", "esc":
			return m, tea.Quit
		}
		return m, nil
	}

	// ---- Overlay / menu screens — handle close first ------------------------
	switch gs.Screen {
	case engine.ScreenHelpOverlay, engine.ScreenHighScoreScreen:
		if key == "esc" || key == "q" || key == "?" || key == "h" {
			engine.CloseOverlay(gs)
			return m, doSave(gs)
		}
		return m, nil

	case engine.ScreenOptionsMenu:
		return m.handleOptions(key)
	}

	// ---- Gameplay screens ---------------------------------------------------
	switch gs.Screen {
	case engine.ScreenMainIdle:
		return m.handleMainIdle(key)
	case engine.ScreenHandDealt:
		return m.handleHandDealt(key)
	case engine.ScreenHandResolved:
		return m.handleHandResolved(key)
	case engine.ScreenGambleStage:
		return m.handleGamble(key)
	case engine.ScreenGambleResult:
		return m.handleGambleResult(key)
	}

	return m, nil
}

// ---- per-screen handlers ----------------------------------------------------

func (m Model) handleMainIdle(key string) (tea.Model, tea.Cmd) {
	gs := m.gs
	switch key {
	case " ":
		engine.Deal(gs)
		return m, doSave(gs)
	case "+", "=":
		engine.IncreaseBet(gs)
		return m, doSave(gs)
	case "-", "_":
		engine.DecreaseBet(gs)
		return m, doSave(gs)
	case "?":
		engine.OpenHelp(gs)
		return m, nil
	case "h", "H":
		engine.OpenHighScore(gs)
		return m, nil
	case "o", "O":
		engine.OpenOptions(gs)
		return m, nil
	case "esc":
		_ = engine.SaveGame(gs)
		return m, tea.Quit
	}
	return m, nil
}

func (m Model) handleHandDealt(key string) (tea.Model, tea.Cmd) {
	gs := m.gs
	switch key {
	case "1":
		engine.ToggleHold(gs, 0)
	case "2":
		engine.ToggleHold(gs, 1)
	case "3":
		engine.ToggleHold(gs, 2)
	case "4":
		engine.ToggleHold(gs, 3)
	case "5":
		engine.ToggleHold(gs, 4)
	case " ":
		engine.Draw(gs)
		return m, doSave(gs)
	case "?":
		engine.OpenHelp(gs)
		return m, nil
	case "o", "O":
		engine.OpenOptions(gs)
		return m, nil
	case "h", "H":
		engine.OpenHighScore(gs)
		return m, nil
	}
	return m, nil
}

func (m Model) handleHandResolved(key string) (tea.Model, tea.Cmd) {
	gs := m.gs
	switch key {
	case " ":
		engine.NextHand(gs)
		return m, doSave(gs)
	case "q", "Q":
		if gs.LastResult.IsGambleEligible {
			engine.StartGamble(gs)
			return m, doSave(gs)
		}
	case "?":
		engine.OpenHelp(gs)
		return m, nil
	case "h", "H":
		engine.OpenHighScore(gs)
		return m, nil
	case "o", "O":
		engine.OpenOptions(gs)
		return m, nil
	case "esc":
		engine.NextHand(gs)
		return m, doSave(gs)
	}
	return m, nil
}

func (m Model) handleGamble(key string) (tea.Model, tea.Cmd) {
	gs := m.gs
	switch key {
	case "1":
		engine.GambleGuess(gs, "red")
		return m, doSave(gs)
	case "2":
		engine.GambleGuess(gs, "black")
		return m, doSave(gs)
	case "c", "C":
		engine.CollectGamble(gs)
		return m, doSave(gs)
	}
	return m, nil
}

func (m Model) handleGambleResult(key string) (tea.Model, tea.Cmd) {
	gs := m.gs
	switch key {
	case "q", "Q", "esc", "enter", " ":
		engine.NextHand(gs)
		return m, doSave(gs)
	}
	return m, nil
}

// ---- options menu -----------------------------------------------------------

const (
	optDesign = iota
	optTheme
	optAutoHold
	optReset
	optBack
	optCount
)

func (m Model) handleOptions(key string) (tea.Model, tea.Cmd) {
	gs := m.gs

	// Reset confirmation prompt
	if m.optionsPrompt != "" {
		switch key {
		case "y", "Y":
			_ = engine.DeleteSave()
			*gs = *engine.NewGameState()
			m.optionsPrompt = ""
			m.optionsCursor = 0
			m.theme = ThemeFor(gs.Options.Theme)
			engine.CloseOverlay(gs)
			return m, nil
		default:
			m.optionsPrompt = ""
			return m, nil
		}
	}

	switch key {
	case "up", "k":
		if m.optionsCursor > 0 {
			m.optionsCursor--
		}
	case "down", "j":
		if m.optionsCursor < optCount-1 {
			m.optionsCursor++
		}
	case "enter", " ", "right", "l":
		switch m.optionsCursor {
		case optDesign:
			gs.Options.CardDesign = (gs.Options.CardDesign + 1) % 3
			return m, doSave(gs)
		case optTheme:
			gs.Options.Theme = (gs.Options.Theme + 1) % 4
			m.theme = ThemeFor(gs.Options.Theme)
			return m, doSave(gs)
		case optAutoHold:
			gs.Options.AutoHold = !gs.Options.AutoHold
			return m, doSave(gs)
		case optReset:
			m.optionsPrompt = "Reset all progress? (Y/N)"
			return m, nil
		case optBack:
			engine.CloseOverlay(gs)
			return m, nil
		}
	case "left", "h":
		switch m.optionsCursor {
		case optDesign:
			gs.Options.CardDesign = (gs.Options.CardDesign + 2) % 3
			return m, doSave(gs)
		case optTheme:
			gs.Options.Theme = (gs.Options.Theme + 3) % 4
			m.theme = ThemeFor(gs.Options.Theme)
			return m, doSave(gs)
		case optAutoHold:
			gs.Options.AutoHold = !gs.Options.AutoHold
			return m, doSave(gs)
		}
	case "esc", "q", "o":
		engine.CloseOverlay(gs)
		return m, nil
	case "?":
		engine.CloseOverlay(gs)
		engine.OpenHelp(gs)
		return m, nil
	}
	return m, nil
}
