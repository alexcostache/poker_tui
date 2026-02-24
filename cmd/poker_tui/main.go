package main

import (
	"fmt"
	"os"

	"poker_tui/internal/engine"
	"poker_tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 1. Acquire lock (multi-instance handling).
	lock, readOnly, lockErr := engine.AcquireLock()
	if lockErr != nil {
		fmt.Fprintf(os.Stderr, "poker_tui: lock error: %v\n", lockErr)
		os.Exit(1)
	}

	// 2. Load (or create) game state.
	gs, loadErr := engine.LoadGame()
	if loadErr != nil {
		// Non-fatal; LoadGame already returns a fresh state on error.
		_ = loadErr
	}

	if readOnly {
		gs.ReadOnly = true
		ownerPid := engine.ReadLockOwner()
		msg := "Another instance of poker_tui is already running"
		if ownerPid > 0 {
			msg += fmt.Sprintf(" (PID %d)", ownerPid)
		}
		msg += ".\nRunning in READ-ONLY mode. No changes will be saved.\nPress Q to quit."
		gs.Screen = engine.ScreenMainIdle
		gs.Message = msg
	}

	// 3. Build Bubble Tea model.
	m := ui.NewModel(gs, lock)

	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// 4. Run.
	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "poker_tui: %v\n", err)
		os.Exit(1)
	}

	// 5. Save final state.
	if fm, ok := finalModel.(interface{ GameState() *engine.GameState }); ok {
		fgs := fm.GameState()
		if !fgs.ReadOnly {
			_ = engine.SaveGame(fgs)
		}
	}

	// 6. Release lock.
	lock.Release()
}
