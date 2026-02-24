package game

import "fmt"

// Suit represents a card suit.
type Suit int

const (
	Clubs Suit = iota
	Diamonds
	Hearts
	Spades
)

func (s Suit) String() string {
	return [...]string{"♣", "♦", "♥", "♠"}[s]
}

func (s Suit) IsRed() bool {
	return s == Diamonds || s == Hearts
}

// Rank represents a card rank (2–Ace).
type Rank int

const (
	Two Rank = iota + 2
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

func (r Rank) String() string {
	switch r {
	case Jack:
		return "J"
	case Queen:
		return "Q"
	case King:
		return "K"
	case Ace:
		return "A"
	case Ten:
		return "10"
	default:
		return fmt.Sprintf("%d", int(r))
	}
}

// Card is a playing card.
type Card struct {
	Rank Rank
	Suit Suit
}

func (c Card) String() string {
	return c.Rank.String() + c.Suit.String()
}

// AllCards returns all 52 cards in order.
func AllCards() []Card {
	cards := make([]Card, 0, 52)
	for s := Clubs; s <= Spades; s++ {
		for r := Two; r <= Ace; r++ {
			cards = append(cards, Card{Rank: r, Suit: s})
		}
	}
	return cards
}
