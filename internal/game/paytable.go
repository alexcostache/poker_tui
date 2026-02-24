package game

// PaytableMultiplier returns the payout multiplier for a given HandRank
// using Jacks-or-Better 9/6 (full-pay) paytable.
func PaytableMultiplier(hr HandRank) int {
	switch hr {
	case RoyalFlush:
		return 800
	case StraightFlush:
		return 50
	case FourOfAKind:
		return 25
	case FullHouse:
		return 9
	case Flush:
		return 6
	case Straight:
		return 4
	case ThreeOfAKind:
		return 3
	case TwoPair:
		return 2
	case OnePair: // Jacks or Better
		return 1
	default: // HighCard, losing pairs
		return 0
	}
}

// PaytableRows returns a human-readable paytable for display.
func PaytableRows() [][2]string {
	return [][2]string{
		{"Royal Flush", "800x"},
		{"Straight Flush", "50x"},
		{"Four of a Kind", "25x"},
		{"Full House", "9x"},
		{"Flush", "6x"},
		{"Straight", "4x"},
		{"Three of a Kind", "3x"},
		{"Two Pair", "2x"},
		{"Jacks or Better", "1x"},
	}
}
