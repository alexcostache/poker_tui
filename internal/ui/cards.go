package ui

import (
	"fmt"
	"strings"

	"poker_tui/internal/engine"
	"poker_tui/internal/game"

	"github.com/mattn/go-runewidth"
)

// RenderCards renders all 5 hand cards side-by-side with hold indicators.
// winCards marks which positions should be highlighted green (winning hand).
func RenderCards(hand game.Hand, design engine.CardDesign, th Theme, winCards [5]bool) string {
	var cardLines [5][]string
	for i, card := range hand.Cards {
		cardLines[i] = renderCard(card, design, th, winCards[i])
	}

	// Find max height
	maxH := 0
	for _, lines := range cardLines {
		if len(lines) > maxH {
			maxH = len(lines)
		}
	}

	var sb strings.Builder

	// Card rows
	for row := 0; row < maxH; row++ {
		for i, lines := range cardLines {
			if row < len(lines) {
				sb.WriteString(lines[row])
			}
			if i < 4 {
				sb.WriteString(" ")
			}
		}
		sb.WriteString("\n")
	}

	// Hold indicators
	for i := 0; i < 5; i++ {
		label := cardHoldLabel(i, hand.Holds[i], design, th)
		sb.WriteString(label)
		if i < 4 {
			sb.WriteString(" ")
		}
	}
	sb.WriteString("\n")

	return sb.String()
}

// RenderGambleCard renders a single card (used in the gamble screen).
func RenderGambleCard(card game.Card, design engine.CardDesign, th Theme) string {
	lines := renderCard(card, design, th, false)
	return strings.Join(lines, "\n")
}

// RenderGambleCardFailed renders a card with a red X overlay, indicating a failed guess.
func RenderGambleCardFailed(card game.Card, design engine.CardDesign, th Theme) string {
	var lines []string
	switch design {
	case engine.DesignMinimal:
		lines = renderMinimalFailed(card, th)
	case engine.DesignWide:
		lines = renderWideFailed(card, th)
	default:
		lines = renderClassicFailed(card, th)
	}
	return strings.Join(lines, "\n")
}

func renderClassicFailed(card game.Card, th Theme) []string {
	rank := card.Rank.String()
	suit := card.Suit.String()
	iw := innerWidth(engine.DesignClassic) // 9
	border := strings.Repeat("─", iw)
	blank := strings.Repeat(" ", iw)
	bs := th.CardFailBorderStyle()
	cs := th.CardFailInnerStyle(card.Suit.IsRed())
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(vfillRight(rank, iw)),
		row(blank),
		row(vcenterSuit(suit, iw)),
		row(blank),
		row(vfillLeft(rank, iw)),
		bs.Render("└" + border + "┘"),
	}
}

func renderMinimalFailed(card game.Card, th Theme) []string {
	rank := card.Rank.String()
	suit := card.Suit.String()
	iw := innerWidth(engine.DesignMinimal) // 7
	border := strings.Repeat("─", iw)
	bs := th.CardFailBorderStyle()
	cs := th.CardFailInnerStyle(card.Suit.IsRed())
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(vfillRight(rank, iw)),
		row(vcenterSuit(suit, iw)),
		row(vfillLeft(rank, iw)),
		bs.Render("└" + border + "┘"),
	}
}

func renderWideFailed(card game.Card, th Theme) []string {
	rank := card.Rank.String()
	suit := card.Suit.String()
	iw := innerWidth(engine.DesignWide) // 13
	border := strings.Repeat("─", iw)
	blank := strings.Repeat(" ", iw)
	bs := th.CardFailBorderStyle()
	cs := th.CardFailInnerStyle(card.Suit.IsRed())
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(vfillRight(rank, iw)),
		row(blank),
		row(blank),
		row(vcenterSuit(suit, iw)),
		row(blank),
		row(blank),
		row(vfillLeft(rank, iw)),
		bs.Render("└" + border + "┘"),
	}
}

// RenderGambleCardBack renders a face-down card for the gamble choose phase.
func RenderGambleCardBack(design engine.CardDesign, th Theme) string {
	var lines []string
	switch design {
	case engine.DesignMinimal:
		lines = renderCardBackMinimal(th)
	case engine.DesignWide:
		lines = renderCardBackWide(th)
	default:
		lines = renderCardBackClassic(th)
	}
	return strings.Join(lines, "\n")
}

// renderCard renders one card according to design. Returns a slice of lines (equal width).
func renderCard(card game.Card, design engine.CardDesign, th Theme, highlight bool) []string {
	switch design {
	case engine.DesignMinimal:
		return renderMinimal(card, th, highlight)
	case engine.DesignWide:
		return renderWide(card, th, highlight)
	default:
		return renderClassic(card, th, highlight)
	}
}

// cardWidth returns the total visual column width of a card (including borders).
func cardWidth(design engine.CardDesign) int {
	switch design {
	case engine.DesignMinimal:
		return 9 // ┌───────┐
	case engine.DesignWide:
		return 15 // ┌─────────────┐
	default:
		return 11 // ┌─────────┐
	}
}

// innerWidth is cardWidth minus the two border chars.
func innerWidth(design engine.CardDesign) int { return cardWidth(design) - 2 }

func cardHoldLabel(i int, held bool, design engine.CardDesign, th Theme) string {
	w := cardWidth(design)
	var label string
	if held {
		label = "HOLD"
	} else {
		label = fmt.Sprintf("[%d]", i+1)
	}
	label = vcenter(label, w)
	if held {
		return th.HoldStyle().Render(label)
	}
	return th.DimStyle().Render(label)
}

// Classic design — 11 cols wide (9 inner), 7 rows tall.
//
//	┌─────────┐
//	│A        │
//	│         │
//	│    ♠    │
//	│         │
//	│        A│
//	└─────────┘
func renderClassic(card game.Card, th Theme, highlight bool) []string {
	rank := card.Rank.String()
	suit := card.Suit.String()
	iw := innerWidth(engine.DesignClassic) // 9
	border := strings.Repeat("─", iw)
	blank := strings.Repeat(" ", iw)

	cs := th.CardStyle(card.Suit.IsRed())
	bs := th.CardStyle(card.Suit.IsRed()) // border style (same colour normally)
	if highlight {
		bs = th.CardWinBorderStyle()
		cs = th.CardWinInnerStyle(card.Suit.IsRed())
	}

	// Content row: styled │ on each side, content style inside.
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}

	return []string{
		bs.Render("┌" + border + "┐"),
		row(vfillRight(rank, iw)),
		row(blank),
		row(vcenterSuit(suit, iw)),
		row(blank),
		row(vfillLeft(rank, iw)),
		bs.Render("└" + border + "┘"),
	}
}

// Minimal design — 9 cols wide (7 inner), 5 rows tall.
//
//	┌───────┐
//	│A      │
//	│   ♠   │
//	│      A│
//	└───────┘
func renderMinimal(card game.Card, th Theme, highlight bool) []string {
	rank := card.Rank.String()
	suit := card.Suit.String()
	iw := innerWidth(engine.DesignMinimal) // 7
	border := strings.Repeat("─", iw)
	cs := th.CardStyle(card.Suit.IsRed())
	bs := th.CardStyle(card.Suit.IsRed())
	if highlight {
		bs = th.CardWinBorderStyle()
		cs = th.CardWinInnerStyle(card.Suit.IsRed())
	}
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(vfillRight(rank, iw)),
		row(vcenterSuit(suit, iw)),
		row(vfillLeft(rank, iw)),
		bs.Render("└" + border + "┘"),
	}
}

// Wide design — 15 cols wide (13 inner), 9 rows tall.
//
//	┌─────────────┐
//	│A            │
//	│             │
//	│             │
//	│      ♠      │   <- single centered suit
//	│             │
//	│             │
//	│            A│
//	└─────────────┘
func renderWide(card game.Card, th Theme, highlight bool) []string {
	rank := card.Rank.String()
	suit := card.Suit.String()
	iw := innerWidth(engine.DesignWide) // 13
	border := strings.Repeat("─", iw)
	blank := strings.Repeat(" ", iw)
	cs := th.CardStyle(card.Suit.IsRed())
	bs := th.CardStyle(card.Suit.IsRed())
	if highlight {
		bs = th.CardWinBorderStyle()
		cs = th.CardWinInnerStyle(card.Suit.IsRed())
	}
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(vfillRight(rank, iw)),
		row(blank),
		row(blank),
		row(vcenterSuit(suit, iw)),
		row(blank),
		row(blank),
		row(vfillLeft(rank, iw)),
		bs.Render("└" + border + "┘"),
	}
}

func renderCardBackClassic(th Theme) []string {
	iw := innerWidth(engine.DesignClassic) // 9
	border := strings.Repeat("─", iw)
	pattern := strings.Repeat("░", iw)
	bs := th.DimStyle()
	cs := th.DimStyle()
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(pattern),
		row(pattern),
		row(pattern),
		row(pattern),
		row(pattern),
		bs.Render("└" + border + "┘"),
	}
}

func renderCardBackMinimal(th Theme) []string {
	iw := innerWidth(engine.DesignMinimal) // 7
	border := strings.Repeat("─", iw)
	pattern := strings.Repeat("░", iw)
	bs := th.DimStyle()
	cs := th.DimStyle()
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(pattern),
		row(pattern),
		row(pattern),
		bs.Render("└" + border + "┘"),
	}
}

func renderCardBackWide(th Theme) []string {
	iw := innerWidth(engine.DesignWide) // 13
	border := strings.Repeat("─", iw)
	pattern := strings.Repeat("░", iw)
	bs := th.DimStyle()
	cs := th.DimStyle()
	row := func(inner string) string {
		return bs.Render("│") + cs.Render(inner) + bs.Render("│")
	}
	return []string{
		bs.Render("┌" + border + "┐"),
		row(pattern),
		row(pattern),
		row(pattern),
		row(pattern),
		row(pattern),
		row(pattern),
		row(pattern),
		bs.Render("└" + border + "┘"),
	}
}

// --- visual-width string helpers (use terminal column counts, not rune counts) ---

// vw returns the visual terminal column width of s.
func vw(s string) int { return runewidth.StringWidth(s) }

// vfillRight pads s on the right so it occupies exactly w terminal columns.
func vfillRight(s string, w int) string {
	need := w - vw(s)
	if need <= 0 {
		return s
	}
	return s + strings.Repeat(" ", need)
}

// vfillLeft pads s on the left so it occupies exactly w terminal columns.
func vfillLeft(s string, w int) string {
	need := w - vw(s)
	if need <= 0 {
		return s
	}
	return strings.Repeat(" ", need) + s
}

// vcenter centers s in exactly w terminal columns.
// Use vcenterSuit for suit symbols (♠♥♦♣) which render as 2 columns.
func vcenter(s string, w int) string {
	need := w - vw(s)
	if need <= 0 {
		return s
	}
	left := need / 2
	right := need - left
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

// vcenterSuit centers a suit symbol in w terminal columns using runewidth.
func vcenterSuit(suit string, w int) string {
	return vcenter(suit, w)
}

// vtwoSuits places two copies of suit spaced evenly inside w columns.
// Falls back to single centered suit if they don't fit.
func vtwoSuits(suit string, w int) string {
	sw := vw(suit)
	gap := w - 2*sw - 4 // 2 cols outer padding on each side
	if gap < 1 {
		gap = 1
	}
	total := 2*sw + gap
	if total > w {
		return vcenterSuit(suit, w) // fallback
	}
	left := (w - total) / 2
	right := w - total - left
	return strings.Repeat(" ", left) + suit +
		strings.Repeat(" ", gap) + suit +
		strings.Repeat(" ", right)
}

// centerPad is kept for hold-label use (label text is ASCII-safe).
func centerPad(s string, width int) string {
	return vcenter(s, width)
}
