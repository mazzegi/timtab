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

func Step(cfg *Config, t *Table, badCache map[string]bool) (*Table, bool) {
	cid := -1
	for ic := 0; ic < cfg.NumCards; ic++ {
		if t.CardCounts[ic] < byte(cfg.cards[ic].Count) {
			cid = ic
			break
		}
	}
	if cid < 0 {
		// all classes assigned
		return t, true
	}
	var freeSlots []int
	for is := 0; is < cfg.NumSlots; is++ {
		off := is * cfg.NumCards
		free := true
		//cas := t.CardsAtSlot(cfg, is)
		for ic := 0; ic < cfg.NumCards; ic++ {
			if t.SlotCards.Get(off+ic) && cfg.excludingCardMatrix.Get(cid*cfg.NumCards+ic) {
				free = false
				break
			}
		}
		//_ = cas
		if free {
			freeSlots = append(freeSlots, is)
		}
	}
	if len(freeSlots) == 0 {
		//RawDump(cfg, t)
		return nil, false
	}
	for _, is := range freeSlots {
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
