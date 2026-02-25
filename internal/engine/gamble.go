package engine

import "poker_tui/internal/game"

// StartGamble takes the most recent win amount back from credits to put it at risk.
func StartGamble(gs *GameState) {
	if gs.Screen != ScreenHandResolved {
		return
	}
	if !gs.LastResult.IsGambleEligible {
		return
	}
	pot := gs.Bet * gs.LastResult.Multiplier
	if pot <= 0 {
		return
	}
	gs.Credits -= pot
	gs.Gamble = GambleState{
		Active:     true,
		Stage:      0,
		MaxStages:  MaxGambleStages,
		CurrentPot: pot,
	}
	gs.Screen = ScreenGambleStage
	drawGambleCard(gs)
	gs.Message = "Pick 1=Red, 2=Black, Space=Collect."
}

func drawGambleCard(gs *GameState) {
	gs.Gamble.CurrentCard = gs.Deck.Draw()
	gs.Gamble.Revealed = false
}

func GambleGuess(gs *GameState, choice string) {
	if gs.Screen != ScreenGambleStage {
		return
	}
	card := gs.Gamble.CurrentCard
	correct := (choice == "red" && card.Suit.IsRed()) ||
		(choice == "black" && !card.Suit.IsRed())
	step := GambleStep{
		Card:      card,
		Choice:    choice,
		PotBefore: gs.Gamble.CurrentPot,
	}
	gs.Gamble.Revealed = true
	gs.Gamble.Stage++
	if correct {
		gs.Gamble.CurrentPot *= 2
		step.Outcome = "win"
		step.PotAfter = gs.Gamble.CurrentPot
		gs.Stats.GambleWins++
		gs.Gamble.History = append(gs.Gamble.History, step)
		if gs.Gamble.Stage >= gs.Gamble.MaxStages {
			collectGamble(gs)
			return
		}
		drawGambleCard(gs)
		gs.Message = "Correct. Pot " + creditStr(gs.Gamble.CurrentPot) + ". 1/2 or Space."
	} else {
		step.Outcome = "lose"
		step.PotAfter = 0
		gs.Gamble.CurrentPot = 0
		gs.Stats.GambleLosses++
		gs.Gamble.History = append(gs.Gamble.History, step)
		gs.Gamble.Active = false
		gs.Screen = ScreenGambleResult
		gs.Message = "Wrong. Pot lost. Space for next hand."
	}
}

func CollectGamble(gs *GameState) {
	if gs.Screen != ScreenGambleStage {
		return
	}
	collectGamble(gs)
}

func collectGamble(gs *GameState) {
	gs.Credits += gs.Gamble.CurrentPot
	profit := gs.Gamble.CurrentPot - (gs.Bet * gs.LastResult.Multiplier)
	if profit > 0 {
		gs.XP += profit
		gs.Level = xpToLevel(gs.XP)
	}
	gs.Stats.LifetimeDelta = gs.Credits - DefaultCredits
	gs.Gamble.Active = false
	gs.Screen = ScreenGambleResult
	gs.Message = "Collected " + creditStr(gs.Gamble.CurrentPot) + ". Space for next hand."
}

func GambleCurrentCard(gs *GameState) game.Card {
	return gs.Gamble.CurrentCard
}
