package timtab

import (
	"fmt"
	"sort"
	"strings"
)

func ClassesConflict(cs1, cs2 *Class) bool {
	//OPT: Class conflicts may be pre-compiled by config
	if cs1.Teacher == cs2.Teacher {
		return true
	}
	// if StudentsIntersect(cs1.Students, cs2.Students) {
	// 	return true
	// }
	if cs1.studentSet.Intersects(cs2.studentSet) {
		return true
	}
	return false
}

func ClassConflictsWithAny(cl *Class, css ...*Class) bool {
	for _, cs := range css {
		if cs.ID == cl.ID {
			return true
		}
		if ClassesConflict(cl, cs) {
			return true
		}
	}
	return false
}

func FindFreeSchedulesForClass(cid ClassID, cfg *Configuration, tt *Timetable) []Schedule {
	cl := cfg.MustClass(cid)
	var scheds []Schedule
	for _, sched := range cfg.Schedules.Values {
		css := ClassesAt(cfg, tt, sched)
		if !ClassConflictsWithAny(cl, css...) {
			scheds = append(scheds, sched)
		}
	}
	return scheds
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
	return "", -1, -1, false
}

func Dump(cfg *Configuration, tt *Timetable) {
	fmt.Printf("###### DUMP ######\n")
	for _, sc := range cfg.Schedules.Values {
		css := ClassesAt(cfg, tt, sc)
		var sl []string
		for _, cs := range css {
			sl = append(sl, string(cs.ID))
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

	scheds := FindFreeSchedulesForClass(cid, cfg, tt)
	if len(scheds) == 0 {
		//Dump(cfg, tt)
		return nil, false
	}
	for _, sched := range scheds {
		ctt := tt.Clone()
		ctt.AddClassSchedule(cid, sched)
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
