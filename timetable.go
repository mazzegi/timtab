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
	}
	for _, s := range cfg.Schedules.Values {
		tt.ScheduleClasses[s] = []ClassID{}
	}
	return tt
}

type Timetable struct {
	ScheduleClasses map[Schedule][]ClassID
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
	}
	for s, cids := range tt.ScheduleClasses {
		newtt.ScheduleClasses[s] = slices.Clone(cids)
	}
	return newtt
}

func (tt *Timetable) AddClassSchedule(cid ClassID, sched Schedule) {
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
	return css
}

func ClassHours(cfg *Configuration, tab *Timetable, cid ClassID) int {
	var hs int
	for _, cids := range tab.ScheduleClasses {
		if slices.Contains(cids, cid) {
			hs++
		}
	}
	return hs
}
