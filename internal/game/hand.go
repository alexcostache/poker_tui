package game

// Hand holds the player's 5 cards and their hold state.
type Hand struct {
	Cards [5]Card
	Holds [5]bool
}

// ToggleHold flips the hold state at position i (0-based).
func (h *Hand) ToggleHold(i int) {
	if i >= 0 && i < 5 {
		h.Holds[i] = !h.Holds[i]
	}
}

// ClearHolds resets all hold flags.
func (h *Hand) ClearHolds() {
	h.Holds = [5]bool{}
}
