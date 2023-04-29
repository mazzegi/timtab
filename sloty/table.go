package sloty

import (
	"bytes"

	"github.com/mazzegi/timtab/bitset"
)

func NewTable(cfg *Config) *Table {
	t := &Table{
		SlotCards:  bitset.New(cfg.NumSlots * cfg.NumCards),
		CardCounts: make([]byte, cfg.NumCards),
		//CardAssignment: bitset.New(cfg.numCards),
	}
	return t
}

type Table struct {
	SlotCards  bitset.Bitset
	CardCounts []byte
	//CardAssignment bitset.Bitset
}

func (t *Table) Hash() string {
	return string(t.SlotCards)
}

func (t *Table) Clone() *Table {
	return &Table{
		SlotCards:  t.SlotCards.Clone(),
		CardCounts: bytes.Clone(t.CardCounts),
		//CardAssignment: t.CardAssignment.Clone(),
	}
}

func (t *Table) AddCardToSlot(cfg *Config, card int, slot int) {
	t.SlotCards.Set(slot*cfg.NumCards+card, true)
	t.CardCounts[card]++
	//t.CardAssignment.Set(card, true)
}

func (t *Table) CardsAtSlot(cfg *Config, slot int) []int {
	cards := []int{}
	off := slot * cfg.NumCards
	for i := 0; i < cfg.NumCards; i++ {
		if t.SlotCards.Get(off + i) {
			cards = append(cards, i)
		}
	}
	return cards
}
