package sloty

import "github.com/mazzegi/timtab/bitset"

type Constraints struct {
	CardsCanShareSlot func(i, j int) bool
}

type Card struct {
	Count int
}

func NewConfig(numSlots int, cards []Card, constraints *Constraints) *Config {
	numCards := len(cards)
	c := &Config{
		NumSlots:            numSlots,
		NumCards:            numCards,
		cards:               cards,
		constraints:         constraints,
		excludingCardMatrix: bitset.New(numCards * numCards),
	}
	for c1 := 0; c1 < numCards; c1++ {
		for c2 := 0; c2 < numCards; c2++ {
			if !constraints.CardsCanShareSlot(c1, c2) {
				c.excludingCardMatrix.Set(c1*numCards+c2, true)
			}
		}
	}
	return c
}

type Config struct {
	NumSlots            int
	NumCards            int
	cards               []Card
	constraints         *Constraints
	excludingCardMatrix bitset.Bitset
}
