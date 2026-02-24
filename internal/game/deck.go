package game

import "math/rand"

// Deck is a 52-card deck with a cursor into the remaining cards.
type Deck struct {
	Cards  []Card
	Cursor int
}

// NewDeck returns a freshly shuffled 52-card deck.
func NewDeck(rng *rand.Rand) *Deck {
	cards := AllCards()
	rng.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return &Deck{Cards: cards, Cursor: 0}
}

// Draw returns the next card from the deck.
// It wraps around if exhausted (should not happen in normal play).
func (d *Deck) Draw() Card {
	if d.Cursor >= len(d.Cards) {
		d.Cursor = 0
	}
	c := d.Cards[d.Cursor]
	d.Cursor++
	return c
}

// Remaining returns how many cards are left.
func (d *Deck) Remaining() int {
	return len(d.Cards) - d.Cursor
}
