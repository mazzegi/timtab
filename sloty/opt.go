package sloty

import "fmt"

func Optimize(cfg *Config) (*Table, error) {
	badCache := map[string]bool{}

	t := NewTable(cfg)
	optt, ok := Step(cfg, t, badCache)
	if !ok {
		return nil, fmt.Errorf("unable to create table")
	}
	return optt, nil
}

func nextCard(cfg *Config, t *Table) int {
	for ic := 0; ic < cfg.NumCards; ic++ {
		if t.CardCounts[ic] < byte(cfg.cards[ic].Count) {
			return ic
		}
	}
	return -1
}

func freeSlots(cfg *Config, t *Table, cid int) []int {
	//fs := make([]int, 0, cfg.NumSlots)
	var fs []int
	for is := 0; is < cfg.NumSlots; is++ {
		off := is * cfg.NumCards
		free := true
		for ic := 0; ic < cfg.NumCards; ic++ {
			if t.SlotCards.Get(off+ic) && cfg.excludingCardMatrix.Get(cid*cfg.NumCards+ic) {
				free = false
				break
			}
		}
		if free {
			fs = append(fs, is)
		}
	}
	return fs
}

func Step(cfg *Config, t *Table, badCache map[string]bool) (*Table, bool) {
	cid := nextCard(cfg, t)
	if cid < 0 {
		// all classes assigned
		return t, true
	}
	free := freeSlots(cfg, t, cid)
	if len(free) == 0 {
		//RawDump(cfg, t)
		return nil, false
	}
	for _, is := range free {
		ct := t.Clone()
		ct.AddCardToSlot(cfg, cid, is)
		hash := ct.Hash()
		if badCache[hash] {
			continue
		}
		newt, ok := Step(cfg, ct, badCache)
		if ok {
			return newt, true
		} else {
			badCache[hash] = true
		}
	}
	return nil, false
}

func RawDump(cfg *Config, table *Table) {
	fmt.Printf("### RAW ###\n")
	for slot := 0; slot < cfg.NumSlots; slot++ {
		cs := table.CardsAtSlot(cfg, slot)
		fmt.Printf("%02d: %v\n", slot, cs)
	}
}
