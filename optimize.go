package timtab

import (
	"fmt"
	"sort"
	"strings"
)

func FindFreeSchedulesForClass(cid ClassID, cfg *Configuration, tt *Timetable) []int {
	sdidxs := make([]int, 0, len(cfg.Schedules.Values))
	for sdidx := range cfg.Schedules.Values {
		if !ClassConflictsWithAnyInSched(cfg, tt, cid, sdidx) {
			sdidxs = append(sdidxs, sdidx)
		}
	}
	return sdidxs
}

func FindClassesWithToLessHours(cfg *Configuration, tt *Timetable) []ClassID {
	var css []ClassID
	for _, cid := range cfg.ClassIDs {
		if FindClassHours(cfg, tt, cid) < cfg.MustClass(cid).HoursPerWeek {
			css = append(css, cid)
		}
	}
	return css
}

func FindFirstClassWithToLessHours(cfg *Configuration, tt *Timetable) (ClassID, int, int, bool) {
	for _, cid := range cfg.ClassIDs {
		ch := FindClassHours(cfg, tt, cid)
		hpw := cfg.MustClass(cid).HoursPerWeek
		if ch < hpw {
			return cid, ch, hpw, true
		}
	}
	return -1, -1, -1, false
}

func Dump(cfg *Configuration, tt *Timetable) {
	fmt.Printf("###### DUMP ######\n")
	for sdidx, sc := range cfg.Schedules.Values {
		cids := ClassesAt(cfg, tt, sdidx)
		var sl []string
		for _, cid := range cids {
			sl = append(sl, cfg.MustClass(cid).Name)
		}
		sort.Strings(sl)
		fmt.Printf("at %d-%d: %s\n", sc.Day, sc.Hour, strings.Join(sl, ", "))
	}
}

func Step(cfg *Configuration, tt *Timetable, badCache map[string]bool, depth int, steps *int) (*Timetable, bool) {
	(*steps)++
	cid, hours, hpw, ok := FindFirstClassWithToLessHours(cfg, tt)
	if !ok {
		//this is fine - all classes have sufficient hours - done
		return tt, true
	}
	_ = hours
	_ = hpw

	sdidxs := FindFreeSchedulesForClass(cid, cfg, tt)
	if len(sdidxs) == 0 {
		//Dump(cfg, tt)
		return nil, false
	}
	for _, sdidx := range sdidxs {
		ctt := tt.Clone()
		ctt.AddClassSchedule(cfg, cid, sdidx)
		hash := HashTimetable(cfg, ctt)
		if badCache[hash] {
			continue
		}

		newtt, ok := Step(cfg, ctt, badCache, depth+1, steps)
		if ok {
			return newtt, true
		} else {
			badCache[hash] = true
		}
	}

	return nil, false
}

func Optimize(cfg *Configuration) (*Timetable, TimetableRating) {
	badCache := map[string]bool{}

	tt := NewTimetable(cfg)
	var steps int
	newtt, ret := Step(cfg, tt, badCache, 1, &steps)
	if !ret {
		return nil, TimetableInvalid
	}
	fmt.Printf("done in %d steps\n", steps)
	return newtt, 1
}
