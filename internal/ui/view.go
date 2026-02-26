package ui

import (
	"fmt"
	"strings"

	"poker_tui/internal/engine"
	"poker_tui/internal/game"

	"github.com/charmbracelet/lipgloss"
)

// View implements tea.Model.
func (m Model) View() string {
	gs := m.gs

	switch gs.Screen {
	case engine.ScreenHelpOverlay:
		return m.viewHelp()
	case engine.ScreenHighScoreScreen:
		return m.viewStats()
	case engine.ScreenOptionsMenu:
		return m.viewOptions()
	case engine.ScreenErrorScreen:
		return m.viewError()
	case engine.ScreenGambleStage:
		return m.viewGamble()
	case engine.ScreenGambleResult:
		return m.viewGambleResult()
	}

	// Main game view (MainIdle, HandDealt, HandResolved)
	return m.viewMain()
}

// ---- main game view ---------------------------------------------------------

func (m Model) viewMain() string {
	gs := m.gs
	th := m.theme
	var sb strings.Builder

	// Top bar
	sb.WriteString(m.topBar())
	sb.WriteString("\n\n")

	// Cards area
	if gs.Screen == engine.ScreenMainIdle {
		sb.WriteString(m.emptyCardsArea())
	} else {
		var winCards [5]bool
		if gs.Screen == engine.ScreenHandResolved && gs.LastResult.IsWin {
			winCards = gs.LastResult.WinningCards
		}
		sb.WriteString(RenderCards(gs.Hand, gs.Options.CardDesign, th, winCards))
	}
	sb.WriteString("\n")

	// Result line
	if gs.Screen == engine.ScreenHandResolved {
		sb.WriteString(th.WinStyle().Render("  "+gs.LastResult.Name) + "\n\n")
	} else {
		sb.WriteString("\n")
	}

	// Message
	if gs.Message != "" {
		sb.WriteString(th.AccentStyle().Render(gs.Message) + "\n")
	}
	sb.WriteString("\n")

	// Key hints
	sb.WriteString(m.keyHints())

	return sb.String()
}

func (m Model) topBar() string {
	gs := m.gs
	th := m.theme

	roLabel := ""
	if gs.ReadOnly {
		roLabel = th.ErrorStyle().Render(" [READ-ONLY] ")
	}

	credits := fmt.Sprintf("Credits: %d", gs.Credits)
	if gs.Credits < 0 {
		credits = th.ErrorStyle().Render(credits)
	} else {
		credits = th.AccentStyle().Render(credits)
	}

	bet := th.DimStyle().Render(fmt.Sprintf("Bet: %d", gs.Bet))
	xp := th.DimStyle().Render(fmt.Sprintf("XP: %d  Lv: %d", gs.XP, gs.Level))

	return lipgloss.JoinHorizontal(lipgloss.Center,
		credits+"  "+bet+"  "+xp+roLabel,
	)
}

func (m Model) emptyCardsArea() string {
	th := m.theme
	design := m.gs.Options.CardDesign
	w := cardWidth(design)
	height := 7 // classic
	if design == engine.DesignMinimal {
		height = 5
	} else if design == engine.DesignWide {
		height = 9
	}
	emptyLine := th.DimStyle().Render(strings.Repeat("░", w))
	lines := ""
	for i := 0; i < height; i++ {
		for j := 0; j < 5; j++ {
			lines += emptyLine
			if j < 4 {
				lines += " "
			}
		}
		lines += "\n"
	}
	// hold-indicator row
	totalW := w*5 + 4
	label := th.DimStyle().Render(vcenter("[Space] to DEAL", totalW))
	lines += label + "\n"
	return lines
}

func (m Model) keyHints() string {
	gs := m.gs
	th := m.theme
	var hints []string

	switch gs.Screen {
	case engine.ScreenMainIdle:
		hints = []string{"Space", "+/-", "?", "H", "O", "^C"}
	case engine.ScreenHandDealt:
		hints = []string{"1-5", "Space", "?"}
	case engine.ScreenHandResolved:
		gg := ""
		if gs.LastResult.IsGambleEligible {
			gg = "Space gamble"
		} else {
			gg = "Space next"
		}
		hints = []string{gg, "?", "H", "O"}
	}

	parts := make([]string, len(hints))
	for i, h := range hints {
		parts[i] = th.DimStyle().Render("[" + h + "]")
	}
	return strings.Join(parts, "  ")
}

// ---- gamble screens ---------------------------------------------------------

func (m Model) viewGamble() string {
	gs := m.gs
	th := m.theme
	var sb strings.Builder

	sb.WriteString(m.topBar() + "\n\n")
	sb.WriteString(th.TitleStyle().Render("  *** GAMBLE ***") + "\n\n")

	// Current pot
	sb.WriteString(th.AccentStyle().Render(fmt.Sprintf("  Pot at risk: %d credits", gs.Gamble.CurrentPot)) + "\n\n")

	progressBlock := VerticalProgressBar(th, gs.Gamble.Stage, gs.Gamble.MaxStages)
	var cardBlock string
	if gs.Gamble.Revealed {
		cardBlock = indent(RenderGambleCard(gs.Gamble.CurrentCard, gs.Options.CardDesign, th), "  ")
	} else {
		cardBlock = indent(RenderGambleCardBack(gs.Options.CardDesign, th), "  ")
	}
	sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, progressBlock, "   ", cardBlock))
	sb.WriteString("\n")

	// History
	if len(gs.Gamble.History) > 0 {
		sb.WriteString("  Hist: ")
		for _, step := range gs.Gamble.History {
			sb.WriteString(gambleHistoryEntry(step, th) + " ")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	if gs.Message != "" {
		var msgStyle lipgloss.Style
		if gs.Gamble.Revealed {
			msgStyle = th.CorrectStyle()
		} else {
			msgStyle = th.AccentStyle()
		}
		sb.WriteString(msgStyle.Render(gs.Message) + "\n\n")
	}
	if gs.Gamble.Revealed {
		sb.WriteString(th.DimStyle().Render("[Space] Next stage") + "\n")
	} else {
		sb.WriteString(th.DimStyle().Render("[1] Red  [2] Black  [Space] Collect") + "\n")
	}
	return sb.String()
}

func (m Model) viewGambleResult() string {
	gs := m.gs
	th := m.theme
	var sb strings.Builder

	sb.WriteString(m.topBar() + "\n\n")
	sb.WriteString(th.TitleStyle().Render("  *** GAMBLE RESULT ***") + "\n\n")

	// Current pot
	sb.WriteString(th.AccentStyle().Render(fmt.Sprintf("  Pot at risk: %d credits", gs.Gamble.CurrentPot)) + "\n\n")

	progressBlock := VerticalProgressBar(th, gs.Gamble.Stage, gs.Gamble.MaxStages)
	if n := len(gs.Gamble.History); n > 0 {
		revealed := gs.Gamble.History[n-1].Card
		cardBlock := indent(RenderGambleCard(revealed, gs.Options.CardDesign, th), "  ")
		sb.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, progressBlock, "   ", cardBlock))
	} else {
		sb.WriteString(progressBlock)
	}
	sb.WriteString("\n")

	// History
	if len(gs.Gamble.History) > 0 {
		sb.WriteString("  Hist: ")
		for _, step := range gs.Gamble.History {
			sb.WriteString(gambleHistoryEntry(step, th) + " ")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	if gs.Message != "" {
		n := len(gs.Gamble.History)
		var msgStyle lipgloss.Style
		if n > 0 && gs.Gamble.History[n-1].Outcome == "lose" {
			msgStyle = th.ErrorStyle()
		} else {
			msgStyle = th.CorrectStyle()
		}
		sb.WriteString(msgStyle.Render(gs.Message) + "\n\n")
	}
	sb.WriteString(th.DimStyle().Render("[Space] next hand") + "\n")
	return sb.String()
}

// ---- overlay screens --------------------------------------------------------

func (m Model) viewHelp() string {
	th := m.theme
	var sb strings.Builder

	sb.WriteString(th.TitleStyle().Render("HELP") + "\n\n")
	sb.WriteString(th.AccentStyle().Render("Keys:") + "\n")

	keys := [][2]string{
		{"Space", "Deal / Draw / Gamble / Collect / Next hand"},
		{"1–5", "Toggle HOLD for card position"},
		{"+/-", "Increase / decrease bet"},
		{"?", "This help screen"},
		{"H", "Stats / high score"},
		{"O", "Options menu"},
		{"^C", "Quit (saves automatically)"},
	}
	for _, k := range keys {
		sb.WriteString(fmt.Sprintf("  %-6s %s\n", k[0], th.DimStyle().Render(k[1])))
	}

	sb.WriteString("\n" + th.AccentStyle().Render("Paytable (Jacks or Better):") + "\n")
	for _, row := range game.PaytableRows() {
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", row[0], th.WinStyle().Render(row[1])))
	}

	sb.WriteString("\n" + th.AccentStyle().Render("Gamble:") + "\n")
	sb.WriteString(th.DimStyle().Render("  After a win, press Space to gamble.\n"))
	sb.WriteString(th.DimStyle().Render("  Guess 1 (Red) or 2 (Black) to double or lose your winnings.\n"))
	sb.WriteString(th.DimStyle().Render("  Press Space to collect at any stage.\n"))
	sb.WriteString(th.DimStyle().Render("  Up to 5 stages (x32 max).\n"))

	sb.WriteString("\n" + th.DimStyle().Render("[ESC / Q / ?] close") + "\n")

	return th.BorderStyle().Render(sb.String())
}

func (m Model) viewStats() string {
	gs := m.gs
	th := m.theme
	s := gs.Stats
	var sb strings.Builder

	sb.WriteString(th.TitleStyle().Render("STATS") + "\n\n")
	rows := [][2]string{
		{"Hands played", fmt.Sprintf("%d", s.HandsPlayed)},
		{"Hands won", fmt.Sprintf("%d", s.HandsWon)},
		{"Hands lost", fmt.Sprintf("%d", s.HandsLost)},
		{"Biggest win", fmt.Sprintf("%d credits", s.BiggestWin)},
		{"Total wagered", fmt.Sprintf("%d credits", s.TotalWagered)},
		{"Total won", fmt.Sprintf("%d credits", s.TotalWon)},
		{"Current streak", fmt.Sprintf("%d", s.CurrentStreak)},
		{"Best streak", fmt.Sprintf("%d", s.BestStreak)},
		{"Lifetime delta", fmt.Sprintf("%+d credits", s.LifetimeDelta)},
		{"Gamble wins", fmt.Sprintf("%d", s.GambleWins)},
		{"Gamble losses", fmt.Sprintf("%d", s.GambleLosses)},
		{"XP", fmt.Sprintf("%d", gs.XP)},
		{"Level", fmt.Sprintf("%d", gs.Level)},
		{"Credits", fmt.Sprintf("%d", gs.Credits)},
	}
	for _, row := range rows {
		sb.WriteString(fmt.Sprintf("  %-20s %s\n", row[0], th.AccentStyle().Render(row[1])))
	}
	sb.WriteString("\n" + th.DimStyle().Render("[ESC / H] close") + "\n")

	return th.BorderStyle().Render(sb.String())
}

func (m Model) viewOptions() string {
	gs := m.gs
	th := m.theme
	var sb strings.Builder

	sb.WriteString(th.TitleStyle().Render("OPTIONS") + "\n\n")

	if m.optionsPrompt != "" {
		sb.WriteString(th.ErrorStyle().Render(m.optionsPrompt) + "\n")
		return th.BorderStyle().Render(sb.String())
	}

	autoHoldVal := "off"
	if gs.Options.AutoHold {
		autoHoldVal = "on"
	}
	items := []struct {
		label string
		value string
	}{
		{"Card Design", gs.Options.CardDesign.String()},
		{"Theme", gs.Options.Theme.String()},
		{"Auto-Hold", autoHoldVal},
		{"Reset Progress", ""},
		{"Back", ""},
	}

	for i, item := range items {
		cursor := "  "
		if i == m.optionsCursor {
			cursor = th.AccentStyle().Render("> ")
		}
		value := ""
		if item.value != "" {
			value = th.DimStyle().Render("  [← →] ") + th.AccentStyle().Render(item.value)
		}
		if item.label == "Reset Progress" {
			value = th.ErrorStyle().Render("  [ENTER] DANGER")
		}
		sb.WriteString(cursor + item.label + value + "\n")
	}

	sb.WriteString("\n" + th.DimStyle().Render("[↑↓] navigate  [←→/ENTER] change  [ESC/O] close") + "\n")

	return th.BorderStyle().Render(sb.String())
}

func (m Model) viewError() string {
	th := m.theme
	gs := m.gs
	var sb strings.Builder
	sb.WriteString(th.TitleStyle().Render("POKER TUI") + "\n\n")
	sb.WriteString(th.ErrorStyle().Render(gs.ErrorMessage) + "\n\n")
	sb.WriteString(th.DimStyle().Render("[Q] quit") + "\n")
	return th.BorderStyle().Render(sb.String())
}

// ---- helpers ----------------------------------------------------------------

// gambleHistoryEntry renders a single gamble history step with suit-based colour.
func gambleHistoryEntry(step engine.GambleStep, th Theme) string {
	icon := "✓"
	if step.Outcome == "lose" {
		icon = "✗"
	}
	text := fmt.Sprintf("[%s %s]", icon, step.Card.String())
	return th.CardStyle(step.Card.Suit.IsRed()).Render(text)
}

func indent(s, prefix string) string {
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
	for i, l := range lines {
		lines[i] = prefix + l
	}
	return strings.Join(lines, "\n")
}
