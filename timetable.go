package timtab

import (
	"strings"
)

type TimetableRating float64

const (
	TimetableInvalid TimetableRating = 0
)

func NewTimetable(classIDs []ClassID) *Timetable {
	tt := &Timetable{
		ClassSchedules: make(map[ClassID]*Schedules),
	}
	for _, cid := range classIDs {
		tt.ClassSchedules[cid] = &Schedules{}
	}
	return tt
}

type Timetable struct {
	ClassSchedules map[ClassID]*Schedules
}

func HashTimetable(cfg *Configuration, tt *Timetable) string {
	var sb strings.Builder
	for _, cid := range cfg.ClassIDs {
		sb.WriteString(string(cid) + ":" + tt.ClassSchedules[cid].Hash() + "|")
	}
	return sb.String()

	// var sl []string
	// for _, cid := range cfg.ClassIDs {
	// 	sl = append(sl, fmt.Sprintf("%s:%s", cid, tt.ClassSchedules[cid].Hash()))
	// }

	// return strings.Join(sl, "|")
}

func (tt *Timetable) Clone() *Timetable {
	newtt := &Timetable{
		ClassSchedules: make(map[ClassID]*Schedules, len(tt.ClassSchedules)),
	}
	for cid, ss := range tt.ClassSchedules {
		newtt.ClassSchedules[cid] = ss.Clone()
	}
	return newtt
}

func (tt *Timetable) AddClassSchedule(cid ClassID, sched Schedule) {
	tt.ClassSchedules[cid].Insert(sched)
}

func RateTimetable(cfg *Configuration, tab *Timetable) TimetableRating {
	return TimetableInvalid
}

func ClassesAt(cfg *Configuration, tab *Timetable, sched Schedule) []*Class {
	var css []*Class
	for cid, scheds := range tab.ClassSchedules {
		if scheds.Contains(sched) {
			css = append(css, cfg.MustClass(cid))
		}
	}
	return css
}

func ClassHours(cfg *Configuration, tab *Timetable, cid ClassID) int {
	return len(tab.ClassSchedules[cid].Values)
}
