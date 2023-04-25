package timtab

import (
	"fmt"
	"sort"
	"strings"
)

// func ClassesConflict(cs1, cs2 *Class) bool {
// 	//OPT: Class conflicts may be pre-compiled by config
// 	if cs1.Teacher == cs2.Teacher {
// 		return true
// 	}
// 	// if StudentsIntersect(cs1.Students, cs2.Students) {
// 	// 	return true
// 	// }
// 	if cs1.studentSet.Intersects(cs2.studentSet) {
// 		return true
// 	}
// 	return false
// }

func ClassConflictsWithAnyOf(cid ClassID, cfg *Configuration, cids ...ClassID) bool {
	for _, acid := range cids {
		if cfg.ClassesConflict(cid, acid) {
			return true
		}
	}
	return false
}

func FindFreeSchedulesForClass(cid ClassID, cfg *Configuration, tt *Timetable) []int {
	sdidxs := make([]int, 0, len(cfg.Schedules.Values))
	for sdidx := range cfg.Schedules.Values {
		css := ClassesAt(cfg, tt, sdidx)
		if !ClassConflictsWithAnyOf(cid, cfg, css...) {
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

func Step(cfg *Configuration, tt *Timetable, badCache map[string]bool, depth int) (*Timetable, bool) {
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

		newtt, ok := Step(cfg, ctt, badCache, depth+1)
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
	newtt, ret := Step(cfg, tt, badCache, 1)
	if !ret {
		return nil, TimetableInvalid
	}
	return newtt, 1
}
