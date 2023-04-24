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
		ScheduleClasses: make(map[Schedule][]ClassID),
		ClassHours:      make(map[ClassID]int),
	}
	for _, s := range cfg.Schedules.Values {
		tt.ScheduleClasses[s] = []ClassID{}
	}
	for _, cid := range cfg.ClassIDs {
		tt.ClassHours[cid] = 0
	}
	return tt
}

type Timetable struct {
	ScheduleClasses map[Schedule][]ClassID
	ClassHours      map[ClassID]int
}

func HashTimetable(cfg *Configuration, tt *Timetable) string {
	//OPT: guess len
	var sb strings.Builder
	sb.Grow(len(cfg.Schedules.Values) * (3 + len(cfg.ClassIDs)/2))
	for _, s := range cfg.Schedules.Values {
		sb.WriteString(cfg.ScheduleHashes[s])
		for _, cid := range tt.ScheduleClasses[s] {
			sb.WriteString(string(cid))
		}
	}
	return sb.String()
}

func (tt *Timetable) Clone() *Timetable {
	newtt := &Timetable{
		ScheduleClasses: make(map[Schedule][]ClassID, len(tt.ScheduleClasses)),
		ClassHours:      make(map[ClassID]int, len(tt.ClassHours)),
	}
	for s, cids := range tt.ScheduleClasses {
		newtt.ScheduleClasses[s] = slices.Clone(cids)
	}
	for cid, hs := range tt.ClassHours {
		newtt.ClassHours[cid] = hs
	}
	return newtt
}

func (tt *Timetable) AddClassSchedule(cid ClassID, sched Schedule) {
	tt.ScheduleClasses[sched] = append(tt.ScheduleClasses[sched], cid)
	tt.ClassHours[cid]++
}

func RateTimetable(cfg *Configuration, tt *Timetable) TimetableRating {
	return TimetableInvalid
}

func ClassesAt(cfg *Configuration, tt *Timetable, sched Schedule) []*Class {
	cids := tt.ScheduleClasses[sched]
	css := make([]*Class, len(cids))
	for i, cid := range cids {
		css[i] = cfg.MustClass(cid)
	}
	return css
}

func FindClassHours(cfg *Configuration, tt *Timetable, cid ClassID) int {
	return tt.ClassHours[cid]
	// var hs int
	// for _, cids := range tab.ScheduleClasses {
	// 	if slices.Contains(cids, cid) {
	// 		hs++
	// 	}
	// }
	// return hs
}
