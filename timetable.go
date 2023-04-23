package timtab

import (
	"strings"

	"golang.org/x/exp/slices"
)

type TimetableRating float64

const (
	TimetableInvalid TimetableRating = 0
)

func NewTimetable(cfg *Configuration) *Timetable {
	tt := &Timetable{
		//ClassSchedules: make(map[ClassID]*Schedules),
		ScheduleClasses: make(map[Schedule][]ClassID),
	}
	// for _, cid := range classIDs {
	// 	tt.ClassSchedules[cid] = &Schedules{}
	// }
	for _, s := range cfg.Schedules.Values {
		tt.ScheduleClasses[s] = []ClassID{}
	}
	return tt
}

type Timetable struct {
	//ClassSchedules  map[ClassID]*Schedules
	ScheduleClasses map[Schedule][]ClassID
}

// func HashTimetable(cfg *Configuration, tt *Timetable) string {
// 	var sb strings.Builder
// 	for _, cid := range cfg.ClassIDs {
// 		sb.WriteString(string(cid) + ":" + tt.ClassSchedules[cid].Hash() + "|")
// 	}
// 	return sb.String()
// }

func joinClassIDs(cids []ClassID, sep string) string {
	var sb strings.Builder
	for _, cid := range cids {
		sb.WriteString(string(cid) + sep)
	}
	return sb.String()
}

// func HashTimetable(cfg *Configuration, tt *Timetable) string {
// 	var sb strings.Builder
// 	for _, s := range cfg.Schedules.Values {
// 		sb.WriteString(string(s.Hash()) + ":" + joinClassIDs(tt.ScheduleClasses[s], ":"))
// 	}
// 	return sb.String()
// }

func HashTimetable(cfg *Configuration, tt *Timetable) string {
	var sb strings.Builder
	for _, s := range cfg.Schedules.Values {
		sb.WriteString(cfg.ScheduleHashes[s])
		for _, cid := range tt.ScheduleClasses[s] {
			sb.WriteString(string(cid))
		}

		//sb.WriteString(string(s.Hash()) + ":" + joinClassIDs(tt.ScheduleClasses[s], ":"))
	}
	return sb.String()
}

func (tt *Timetable) Clone() *Timetable {
	newtt := &Timetable{
		//ClassSchedules:  make(map[ClassID]*Schedules, len(tt.ClassSchedules)),
		ScheduleClasses: make(map[Schedule][]ClassID, len(tt.ScheduleClasses)),
	}
	// for cid, ss := range tt.ClassSchedules {
	// 	newtt.ClassSchedules[cid] = ss.Clone()
	// }
	for s, cids := range tt.ScheduleClasses {
		newtt.ScheduleClasses[s] = slices.Clone(cids)
	}
	return newtt
}

func (tt *Timetable) AddClassSchedule(cid ClassID, sched Schedule) {
	//tt.ClassSchedules[cid].Insert(sched)
	tt.ScheduleClasses[sched] = append(tt.ScheduleClasses[sched], cid)
}

func RateTimetable(cfg *Configuration, tab *Timetable) TimetableRating {
	return TimetableInvalid
}

func ClassesAt(cfg *Configuration, tab *Timetable, sched Schedule) []*Class {
	var css []*Class
	for _, cid := range tab.ScheduleClasses[sched] {
		css = append(css, cfg.MustClass(cid))
	}

	// for cid, scheds := range tab.ClassSchedules {
	// 	if scheds.Contains(sched) {
	// 		css = append(css, cfg.MustClass(cid))
	// 	}
	// }
	return css
}

func ClassHours(cfg *Configuration, tab *Timetable, cid ClassID) int {
	//return len(tab.ClassSchedules[cid].Values)
	var hs int
	for _, cids := range tab.ScheduleClasses {
		if slices.Contains(cids, cid) {
			hs++
		}
	}
	return hs
}
