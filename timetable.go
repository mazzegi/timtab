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
		ScheduleClassesBitset: NewBitset(schedCount * len(cfg.ClassIDs)),
		ClassHoursBytes:       make([]byte, len(cfg.ClassIDs)),
	}
	return tt
}

type Timetable struct {
	ScheduleClassesBitset Bitset
	ClassHoursBytes       []byte
}

func HashTimetable(cfg *Configuration, tt *Timetable) string {
	return string(tt.ScheduleClassesBitset)
}

func (tt *Timetable) Clone() *Timetable {
	newtt := &Timetable{
		ScheduleClassesBitset: tt.ScheduleClassesBitset.Clone(),
		ClassHoursBytes:       bytes.Clone(tt.ClassHoursBytes),
	}
	return newtt
}

func (tt *Timetable) AddClassSchedule(cfg *Configuration, cid ClassID, sdidx int) {
	tt.ScheduleClassesBitset.Set(sdidx*len(cfg.ClassIDs)+int(cid), true)
	tt.ClassHoursBytes[cid]++
}

func RateTimetable(cfg *Configuration, tt *Timetable) TimetableRating {
	return TimetableInvalid
}

func ClassesAt(cfg *Configuration, tt *Timetable, sdidx int) []ClassID {
	css := make([]ClassID, 0, len(cfg.ClassIDs))
	off := sdidx * len(cfg.ClassIDs)
	for i := 0; i < len(cfg.ClassIDs); i++ {
		if tt.ScheduleClassesBitset.Get(off + i) {
			css = append(css, ClassID(i))
		}
	}
	return css
}

func ClassConflictsWithAnyInSched(cfg *Configuration, tt *Timetable, cid ClassID, sdidx int) bool {
	off := sdidx * len(cfg.ClassIDs)
	for i := 0; i < len(cfg.ClassIDs); i++ {
		if tt.ScheduleClassesBitset.Get(off+i) &&
			cfg.ClassesConflict(cid, ClassID(i)) {
			return true
		}
	}
	return false
}

func FindClassHours(cfg *Configuration, tt *Timetable, cid ClassID) int {
	return int(tt.ClassHoursBytes[cid])
}
