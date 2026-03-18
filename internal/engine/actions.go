package engine

import (
	"math/rand"
	"poker_tui/internal/game"
	"time"
)

var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// Deal deducts the bet and deals 5 fresh cards. Transitions to HandDealt.
func Deal(gs *GameState) {
	if gs.Screen != ScreenMainIdle {
		return
	}
	gs.Credits -= gs.Bet
	gs.Deck = *game.NewDeck(rng)
	gs.Hand.ClearHolds()
	for i := 0; i < 5; i++ {
		gs.Hand.Cards[i] = gs.Deck.Draw()
	}
	if gs.Options.AutoHold {
		applyAutoHold(&gs.Hand)
	}
	gs.Screen = ScreenHandDealt
	gs.Message = "Hold 1-5, then Space to draw."
}

// applyAutoHold holds high cards (J+) and any paired/tripped/quaded ranks.
func applyAutoHold(hand *game.Hand) {
	// Count rank frequencies
	counts := make(map[game.Rank]int, 5)
	for _, c := range hand.Cards {
		counts[c.Rank]++
	}
	for i, c := range hand.Cards {
		if counts[c.Rank] >= 2 || c.Rank >= game.Jack {
			hand.Holds[i] = true
		}
	}
}

// ToggleHold flips the hold state of card at position i (0-based). Only valid in HandDealt.
func ToggleHold(gs *GameState, i int) {
	if gs.Screen != ScreenHandDealt {
		return
	}
	gs.Hand.ToggleHold(i)
}

// Draw replaces non-held cards and evaluates the hand. Transitions to HandResolved.
func Draw(gs *GameState) {
	if gs.Screen != ScreenHandDealt {
		return
	}
	for i := 0; i < 5; i++ {
		if !gs.Hand.Holds[i] {
			gs.Hand.Cards[i] = gs.Deck.Draw()
		}
	}
	gs.Hand.ClearHolds()
	resolveHand(gs)
}

// resolveHand evaluates and pays out, updating stats and XP.
func resolveHand(gs *GameState) {
	result := game.Evaluate(gs.Hand)
	gs.LastResult = result
	gs.Stats.HandsPlayed++

	payout := gs.Bet * result.Multiplier
	if payout > 0 {
		gs.Credits += payout
		gs.Stats.HandsWon++
		gs.Stats.TotalWon += payout
		if payout > gs.Stats.BiggestWin {
			gs.Stats.BiggestWin = payout
		}
		if gs.Stats.CurrentStreak < 0 {
			gs.Stats.CurrentStreak = 1
		} else {
			gs.Stats.CurrentStreak++
		}
		if gs.Stats.CurrentStreak > gs.Stats.BestStreak {
			gs.Stats.BestStreak = gs.Stats.CurrentStreak
		}
		// XP
		gs.XP += payout
		gs.Level = xpToLevel(gs.XP)
		gs.Message = "Win: " + result.Name + " +" + creditStr(payout)
	} else {
		gs.Stats.HandsLost++
		if gs.Stats.CurrentStreak > 0 {
			gs.Stats.CurrentStreak = -1
		} else {
			gs.Stats.CurrentStreak--
		}
		gs.Message = result.Name + " (no win)"
	}
	gs.Stats.TotalWagered += gs.Bet
	gs.Stats.LifetimeDelta = gs.Credits - DefaultCredits

	gs.Screen = ScreenHandResolved
}

// NextHand resets back to MainIdle from HandResolved.
func NextHand(gs *GameState) {
	if gs.Screen != ScreenHandResolved && gs.Screen != ScreenGambleResult {
		return
	}
	if gs.Screen == ScreenGambleResult {
		gs.Hand = game.Hand{} // clear cards so the idle view shows the empty placeholder
	}
	gs.Screen = ScreenMainIdle
	gs.Message = ""
}

// IncreaseBet increases the bet by 1, capped at MaxBet. Only in MainIdle.
func IncreaseBet(gs *GameState) {
	if gs.Screen != ScreenMainIdle {
		return
	}
	if gs.Bet < MaxBet {
		gs.Bet++
	}
}

// DecreaseBet decreases the bet by 1, min MinBet. Only in MainIdle.
func DecreaseBet(gs *GameState) {
	if gs.Screen != ScreenMainIdle {
		return
	}
	if gs.Bet > MinBet {
		gs.Bet--
	}
}

// OpenHelp opens the help overlay, saving the previous screen.
func OpenHelp(gs *GameState) {
	if gs.Screen == ScreenHelpOverlay || gs.Screen == ScreenOptionsMenu || gs.Screen == ScreenHighScoreScreen {
		return
	}
	gs.PrevScreen = gs.Screen
	gs.Screen = ScreenHelpOverlay
}

// OpenOptions opens the options menu.
func OpenOptions(gs *GameState) {
	if gs.Screen == ScreenHelpOverlay || gs.Screen == ScreenOptionsMenu || gs.Screen == ScreenHighScoreScreen {
		return
	}
	gs.PrevScreen = gs.Screen
	gs.Screen = ScreenOptionsMenu
}

// OpenHighScore opens the stats screen.
func OpenHighScore(gs *GameState) {
	if gs.Screen == ScreenHelpOverlay || gs.Screen == ScreenOptionsMenu || gs.Screen == ScreenHighScoreScreen {
		return
	}
	gs.PrevScreen = gs.Screen
	gs.Screen = ScreenHighScoreScreen
}

// CloseOverlay returns to the previous game screen.
func CloseOverlay(gs *GameState) {
	gs.Screen = gs.PrevScreen
}

// creditStr formats a credit amount as a string.
func creditStr(n int) string {
	if n < 0 {
		return "-" + itoa(-n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
