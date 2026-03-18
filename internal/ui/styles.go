package ui

import (
	"fmt"

	"poker_tui/internal/engine"

	"github.com/charmbracelet/lipgloss"
)

// Theme wraps a colour preset and exposes lipgloss style builders.
type Theme struct {
	preset engine.Theme

	// Base colours
	bg     lipgloss.Color
	fg     lipgloss.Color
	accent lipgloss.Color
	dim    lipgloss.Color
	red    lipgloss.Color
	hold   lipgloss.Color
	win    lipgloss.Color
	border lipgloss.Color
}

// ThemeFor builds a Theme for the given engine.Theme preset.
func ThemeFor(t engine.Theme) Theme {
	switch t {
	case engine.ThemeAmber:
		return Theme{
			preset: t,
			bg:     "#1A0F00", fg: "#FFBF00", accent: "#FFA500",
			dim: "#7A5C00", red: "#FF5722", hold: "#FF8C00",
			win: "#FFD700", border: "#7A5C00",
		}
	case engine.ThemeGreen:
		return Theme{
			preset: t,
			bg:     "#002200", fg: "#00FF41", accent: "#00CC33",
			dim: "#005510", red: "#CC0000", hold: "#00FF41",
			win: "#ADFF2F", border: "#005510",
		}
	case engine.ThemeMono:
		return Theme{
			preset: t,
			bg:     "#000000", fg: "#FFFFFF", accent: "#CCCCCC",
			dim: "#666666", red: "#FFFFFF", hold: "#FFFFFF",
			win: "#FFFFFF", border: "#444444",
		}
	default: // ThemeDark
		return Theme{
			preset: t,
			bg:     "#0D1117", fg: "#E6EDF3", accent: "#58A6FF",
			dim: "#484F58", red: "#FF7B72", hold: "#3FB950",
			win: "#FFD700", border: "#30363D",
		}
	}
}

func (th Theme) Base() lipgloss.Style {
	return lipgloss.NewStyle().Background(th.bg).Foreground(th.fg)
}

func (th Theme) CardStyle(isRed bool) lipgloss.Style {
	fg := th.fg
	if isRed {
		fg = th.red
	}
	return lipgloss.NewStyle().Foreground(fg).Bold(true)
}

// CardWinBorderStyle returns the bright-green bg style used for winning card borders.
func (th Theme) CardWinBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#00E676")).
		Foreground(lipgloss.Color("#000000")).
		Bold(true)
}

// CardWinInnerStyle returns the inner-cell style for a highlighted winning card.
func (th Theme) CardWinInnerStyle(isRed bool) lipgloss.Style {
	fg := th.fg
	if isRed {
		fg = th.red
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#003820")).
		Foreground(fg).
		Bold(true)
}

// CardFailBorderStyle returns the bright-red bg style used for failed-guess card borders.
func (th Theme) CardFailBorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#FF0000")).
		Foreground(lipgloss.Color("#000000")).
		Bold(true)
}

// CardFailInnerStyle returns the inner-cell style for a failed-guess card.
func (th Theme) CardFailInnerStyle(isRed bool) lipgloss.Style {
	fg := th.fg
	if isRed {
		fg = th.red
	}
	return lipgloss.NewStyle().
		Background(lipgloss.Color("#380000")).
		Foreground(fg).
		Bold(true)
}

func (th Theme) HoldStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(th.hold).Bold(true)
}

func (th Theme) DimStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(th.dim)
}

func (th Theme) AccentStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(th.accent).Bold(true)
}

func (th Theme) WinStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(th.win).Bold(true)
}

func (th Theme) BorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(th.border).
		Padding(0, 1)
}

func (th Theme) TitleStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(th.accent).
		Bold(true).
		Underline(true)
}

func (th Theme) StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(th.border).
		Foreground(th.fg).
		Padding(0, 1)
}

func (th Theme) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
}

func (th Theme) CorrectStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color("#00E676")).Bold(true)
}

// BarSegment returns a coloured bar segment character.
func (th Theme) BarSegment(filled bool) string {
	if filled {
		return th.AccentStyle().Render("█")
	}
	return th.DimStyle().Render("░")
}

// VerticalProgressBar renders a vertical bar for the gamble screen.
// height = total bar height in rows; filled = rows that are "active".
func VerticalProgressBar(th Theme, filled, total int) string {
	var rows []string
	for i := total; i >= 1; i-- {
		label := fmt.Sprintf("x%-3d", 1<<i) // x2, x4, x8 …
		if i <= filled {
			seg := th.AccentStyle().Render("███")
			rows = append(rows, label+" "+seg)
		} else {
			seg := th.DimStyle().Render("░░░")
			rows = append(rows, th.DimStyle().Render(label)+" "+seg)
		}
	}
	// Base (original bet)
	base := "x1   "
	if filled >= 0 {
		base += th.AccentStyle().Render("███")
	} else {
		base += th.DimStyle().Render("░░░")
	}
	rows = append(rows, base)
	result := ""
	for _, r := range rows {
		result += r + "\n"
	}
	return result
}
