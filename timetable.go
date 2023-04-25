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
		//ScheduleClasses:      make(map[Schedule][]ClassID),
		ScheduleClassesBytes: make([]byte, schedCount*len(cfg.ClassIDs)),
		ClassHours:           make(map[ClassID]int),
	}
	// s
	for _, cid := range cfg.ClassIDs {
		tt.ClassHours[cid] = 0
	}
	return tt
}

type Timetable struct {
	//ScheduleClasses      map[Schedule][]ClassID
	ScheduleClassesBytes []byte
	ClassHours           map[ClassID]int
}

func HashTimetable(cfg *Configuration, tt *Timetable) string {
	return string(tt.ScheduleClassesBytes)
	//OPT: guess len
	// var sb strings.Builder
	// sb.Grow(len(cfg.Schedules.Values) * (3 + len(cfg.ClassIDs)))
	// for _, s := range cfg.Schedules.Values {
	// 	sb.WriteString(cfg.ScheduleHashes[s])
	// 	// -- write classids is the killer
	// 	// - IDEA: make scheduleclasses a flat thing (one byte array) with constant lookup time and cheap copying resp. hashing
	// 	sb.WriteString(fmt.Sprint(tt.ScheduleClasses[s]))
	// 	// for _, cid := range tt.ScheduleClasses[s] {
	// 	// 	sb.WriteString(string(cid))
	// 	// }
	// }
	// return sb.String()
}

func (tt *Timetable) Clone() *Timetable {
	newtt := &Timetable{
		ScheduleClassesBytes: bytes.Clone(tt.ScheduleClassesBytes),
		//ScheduleClasses: make(map[Schedule][]ClassID, len(tt.ScheduleClasses)),
		ClassHours: make(map[ClassID]int, len(tt.ClassHours)),
	}
	// for s, cids := range tt.ScheduleClasses {
	// 	newtt.ScheduleClasses[s] = slices.Clone(cids)
	// }
	for cid, hs := range tt.ClassHours {
		newtt.ClassHours[cid] = hs
	}
	return newtt
}

func (tt *Timetable) AddClassSchedule(cfg *Configuration, cid ClassID, sdidx int) {
	tt.ScheduleClassesBytes[sdidx*len(cfg.ClassIDs)+int(cid)] = 1

	//tt.ScheduleClasses[sched] = append(tt.ScheduleClasses[sched], cid)
	tt.ClassHours[cid]++
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
	//return tt.ScheduleClasses[sched]
}

func FindClassHours(cfg *Configuration, tt *Timetable, cid ClassID) int {
	return tt.ClassHours[cid]
}
