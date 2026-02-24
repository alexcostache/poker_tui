package game

import "sort"

// HandRank is the strength category of a poker hand.
type HandRank int

const (
	HighCard HandRank = iota
	OnePair
	TwoPair
	ThreeOfAKind
	Straight
	Flush
	FullHouse
	FourOfAKind
	StraightFlush
	RoyalFlush
)

func (hr HandRank) String() string {
	return [...]string{
		"High Card",
		"One Pair",
		"Two Pair",
		"Three of a Kind",
		"Straight",
		"Flush",
		"Full House",
		"Four of a Kind",
		"Straight Flush",
		"Royal Flush",
	}[hr]
}

// HandResult is the evaluated outcome of a 5-card hand.
type HandResult struct {
	Rank             HandRank
	Name             string
	Multiplier       int     // applied to bet
	IsWin            bool    // multiplier > 0
	IsGambleEligible bool    // win that allows the gamble feature
	WinningCards     [5]bool // which original card positions form the winning hand
}

// Evaluate assesses a 5-card hand and returns the HandResult.
func Evaluate(hand Hand) HandResult {
	cards := hand.Cards[:]
	sorted := make([]Card, 5)
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Rank < sorted[j].Rank
	})

	isFlush := checkFlush(sorted)
	isStraight, isRoyal := checkStraight(sorted)

	switch {
	case isFlush && isRoyal:
		result := makeResult(RoyalFlush)
		result.WinningCards = allFive()
		return result
	case isFlush && isStraight:
		result := makeResult(StraightFlush)
		result.WinningCards = allFive()
		return result
	}

	counts := rankCounts(sorted)
	if hr, ok := checkCounts(counts); ok {
		result := makeResult(hr)
		result.WinningCards = winCardsForCounts(hr, hand.Cards[:], counts)
		return result
	}

	switch {
	case isFlush:
		result := makeResult(Flush)
		result.WinningCards = allFive()
		return result
	case isStraight:
		result := makeResult(Straight)
		result.WinningCards = allFive()
		return result
	}

	// Pair check — only Jacks or Better pays
	if hr, ok := checkPair(sorted, counts); ok {
		result := makeResult(hr)
		result.WinningCards = winCardsForCounts(hr, hand.Cards[:], counts)
		return result
	}

	return makeResult(HighCard)
}

func allFive() [5]bool { return [5]bool{true, true, true, true, true} }

// winCardsForCounts marks positions in original cards that participate in the winning combination.
func winCardsForCounts(hr HandRank, cards []Card, counts map[Rank]int) [5]bool {
	var win [5]bool
	switch hr {
	case FourOfAKind:
		for i, c := range cards {
			if counts[c.Rank] == 4 {
				win[i] = true
			}
		}
	case FullHouse:
		return allFive()
	case ThreeOfAKind:
		for i, c := range cards {
			if counts[c.Rank] == 3 {
				win[i] = true
			}
		}
	case TwoPair:
		// Mark cards belonging to either pair rank.
		var pairRanks []Rank
		for rank, cnt := range counts {
			if cnt == 2 {
				pairRanks = append(pairRanks, rank)
			}
		}
		for i, c := range cards {
			for _, pr := range pairRanks {
				if c.Rank == pr {
					win[i] = true
				}
			}
		}
	case OnePair:
		for i, c := range cards {
			if counts[c.Rank] == 2 && c.Rank >= Jack {
				win[i] = true
			}
		}
	}
	return win
}

// ---- helpers ----------------------------------------------------------------

func makeResult(hr HandRank) HandResult {
	mult := PaytableMultiplier(hr)
	return HandResult{
		Rank:             hr,
		Name:             hr.String(),
		Multiplier:       mult,
		IsWin:            mult > 0,
		IsGambleEligible: mult > 0,
	}
}

func checkFlush(sorted []Card) bool {
	s := sorted[0].Suit
	for _, c := range sorted[1:] {
		if c.Suit != s {
			return false
		}
	}
	return true
}

// checkStraight returns (isStraight, isRoyal).
func checkStraight(sorted []Card) (bool, bool) {
	// Normal straight
	normal := true
	for i := 1; i < 5; i++ {
		if int(sorted[i].Rank) != int(sorted[i-1].Rank)+1 {
			normal = false
			break
		}
	}
	if normal {
		isRoyal := sorted[0].Rank == Ten && sorted[4].Rank == Ace
		return true, isRoyal
	}
	// Ace-low straight: A 2 3 4 5
	// sorted might be [2,3,4,5,A] because Ace rank=14 is highest
	if sorted[0].Rank == Two && sorted[1].Rank == Three &&
		sorted[2].Rank == Four && sorted[3].Rank == Five &&
		sorted[4].Rank == Ace {
		return true, false
	}
	return false, false
}

func rankCounts(sorted []Card) map[Rank]int {
	m := make(map[Rank]int)
	for _, c := range sorted {
		m[c.Rank]++
	}
	return m
}

func checkCounts(counts map[Rank]int) (HandRank, bool) {
	var pairs, trips, quads int
	for _, cnt := range counts {
		switch cnt {
		case 2:
			pairs++
		case 3:
			trips++
		case 4:
			quads++
		}
	}
	switch {
	case quads == 1:
		return FourOfAKind, true
	case trips == 1 && pairs == 1:
		return FullHouse, true
	case trips == 1:
		return ThreeOfAKind, true
	case pairs == 2:
		return TwoPair, true
	}
	return HighCard, false
}

// checkPair returns OnePair only when the pair is Jacks or Better (Jacks, Queens, Kings, Aces).
func checkPair(sorted []Card, counts map[Rank]int) (HandRank, bool) {
	for rank, cnt := range counts {
		if cnt == 2 && rank >= Jack {
			return OnePair, true
		}
	}
	return HighCard, false
}
