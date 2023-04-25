package timtab

import (
	"fmt"
	"sort"
	"strings"

	"github.com/mazzegi/log"
)

type (
	SubjectID string
	StudentID string
	TeacherID string
	//ClassID   string
	ClassID int
	Day     int //1..5; Mo..Fr
	Hour    int //1..10
)

type Subject struct {
	ID   SubjectID
	Name string
}

type Student struct {
	ID   StudentID
	Name string
}

type Teacher struct {
	ID       TeacherID
	Name     string
	Subjects []SubjectID
}

type Class struct {
	ID           ClassID
	Level        int
	HoursPerWeek int
	Name         string
	Subject      SubjectID
	Teacher      TeacherID
	Students     []StudentID
	studentSet   *Set[StudentID]
}

type Schedule struct {
	Day  Day
	Hour Hour
}

func (s Schedule) Less(os Schedule) bool {
	if s.Day < os.Day {
		return true
	}
	if s.Day > os.Day {
		return false
	}
	return s.Hour < os.Hour
}

type Schedules struct {
	Values   []Schedule
	lastHash string
}

func (ss *Schedules) Clone() *Schedules {
	css := &Schedules{
		Values:   make([]Schedule, len(ss.Values)),
		lastHash: ss.lastHash,
	}
	copy(css.Values, ss.Values)
	return css
}

func (ss *Schedules) Insert(s Schedule) {
	ss.Values = append(ss.Values, s)
	ss.lastHash = ""
}

func (ss *Schedules) Contains(sched Schedule) bool {
	for _, s := range ss.Values {
		if s == sched {
			return true
		}
	}
	return false
}

func (s Schedule) Hash() string {
	return fmt.Sprintf("%d-%02d", s.Day, s.Hour)
}

func (ss *Schedules) Hash() string {
	if ss.lastHash != "" {
		return ss.lastHash
	}
	// sort.Slice(ss, func(i, j int) bool {
	// 	return ss[i].Less(ss[j])
	// })

	sl := make([]string, len(ss.Values))
	for i, s := range ss.Values {
		sl[i] = s.Hash()
	}
	sort.Strings(sl)
	h := strings.Join(sl, ":")
	ss.lastHash = h
	return h
}

// func (ss Schedules) Hash2() string {
// 	var b strings.Builder

// 	sl := make([]string, len(ss))
// 	for i, s := range ss {
// 		sl[i] = s.Hash()
// 	}
// 	sort.Strings(sl)
// 	return strings.Join(sl, ":")
// }

// func StudentsIntersect(st1, st2 []StudentID) bool {
// 	for _, s1 := range st1 {
// 		if slices.Contains(st2, s1) {
// 			return true
// 		}
// 	}
// 	return false
// }

// Config

func NewConfiguration(sds *Schedules) *Configuration {
	cfg := &Configuration{
		Schedules:      sds,
		ScheduleHashes: make(map[Schedule]string),
		Subjects:       make(map[SubjectID]*Subject),
		Teachers:       make(map[TeacherID]*Teacher),
		Classes:        make(map[ClassID]*Class),
	}
	for _, s := range cfg.Schedules.Values {
		cfg.ScheduleHashes[s] = s.Hash()
	}
	return cfg
}

type Configuration struct {
	Schedules                 *Schedules
	ScheduleHashes            map[Schedule]string
	Subjects                  map[SubjectID]*Subject
	Teachers                  map[TeacherID]*Teacher
	Classes                   map[ClassID]*Class
	SubjectIDs                []SubjectID
	TeacherIDs                []TeacherID
	ClassIDs                  []ClassID
	ConflictingClassPairBytes []byte
}

func (c *Configuration) Prepare() {
	c.ConflictingClassPairBytes = make([]byte, len(c.Classes)*len(c.Classes))
	for _, c1 := range c.Classes {
		for _, c2 := range c.Classes {
			if c1.ID == c2.ID ||
				c1.Teacher == c2.Teacher ||
				c1.studentSet.Intersects(c2.studentSet) {
				c.ConflictingClassPairBytes[int(c1.ID)*len(c.Classes)+int(c2.ID)] = 1
			}

		}
	}
}

func (c *Configuration) ClassesConflict(c1, c2 ClassID) bool {
	return c.ConflictingClassPairBytes[int(c1)*len(c.Classes)+int(c2)] == 1
}

func (c *Configuration) AddSubjects(sbs ...*Subject) error {
	for _, sb := range sbs {
		if _, ok := c.Subjects[sb.ID]; ok {
			return fmt.Errorf("subject with id %q already exists", sb.ID)
		}
		c.Subjects[sb.ID] = sb
		c.SubjectIDs = append(c.SubjectIDs, sb.ID)
	}
	return nil
}

func (c *Configuration) AddTeachers(ts ...*Teacher) error {
	for _, t := range ts {
		if _, ok := c.Teachers[t.ID]; ok {
			return fmt.Errorf("teacher with id %q already exists", t.ID)
		}
		c.Teachers[t.ID] = t
		c.TeacherIDs = append(c.TeacherIDs, t.ID)
	}
	return nil
}

func (c *Configuration) AddClasses(cls ...*Class) error {
	for _, cl := range cls {
		id := len(c.Classes)
		cl.ID = ClassID(id)
		cl.studentSet = NewSet(cl.Students...)

		c.Classes[cl.ID] = cl
		c.ClassIDs = append(c.ClassIDs, cl.ID)
	}
	return nil
}

func (c *Configuration) MustClass(cid ClassID) *Class {
	cl, ok := c.Classes[cid]
	if !ok {
		log.Fatalf("no such class %q", cid)
	}
	return cl
}
