package timtab

import (
	"bytes"
)

type TimetableRating float64

const (
	TimetableInvalid TimetableRating = 0
)

func NewTimetable(cfg *Configuration) *Timetable {
	schedCount := len(cfg.Schedules.Values)

	tt := &Timetable{
		ScheduleClassesBytes: make([]byte, schedCount*len(cfg.ClassIDs)),
		//ClassHours:           make(map[ClassID]int),
		ClassHoursBytes: make([]byte, len(cfg.ClassIDs)),
	}
	// s
	// for _, cid := range cfg.ClassIDs {
	// 	tt.ClassHours[cid] = 0
	// }
	return tt
}

type Timetable struct {
	ScheduleClassesBytes []byte
	//ClassHours           map[ClassID]int
	ClassHoursBytes []byte
}

func HashTimetable(cfg *Configuration, tt *Timetable) string {
	return string(tt.ScheduleClassesBytes)
}

func (tt *Timetable) Clone() *Timetable {
	newtt := &Timetable{
		ScheduleClassesBytes: bytes.Clone(tt.ScheduleClassesBytes),
		ClassHoursBytes:      bytes.Clone(tt.ClassHoursBytes),
		//ClassHours:           make(map[ClassID]int, len(tt.ClassHours)),
	}
	// for cid, hs := range tt.ClassHours {
	// 	newtt.ClassHours[cid] = hs
	// }
	return newtt
}

func (tt *Timetable) AddClassSchedule(cfg *Configuration, cid ClassID, sdidx int) {
	tt.ScheduleClassesBytes[sdidx*len(cfg.ClassIDs)+int(cid)] = 1
	tt.ClassHoursBytes[cid]++
}

func RateTimetable(cfg *Configuration, tt *Timetable) TimetableRating {
	return TimetableInvalid
}

func ClassesAt(cfg *Configuration, tt *Timetable, sdidx int) []ClassID {
	css := make([]ClassID, 0, len(cfg.ClassIDs))
	off := sdidx * len(cfg.ClassIDs)
	for i := 0; i < len(cfg.ClassIDs); i++ {
		if tt.ScheduleClassesBytes[off+i] > 0 {
			css = append(css, ClassID(i))
		}
	}
	return css
}

func FindClassHours(cfg *Configuration, tt *Timetable, cid ClassID) int {
	return int(tt.ClassHoursBytes[cid])
}
