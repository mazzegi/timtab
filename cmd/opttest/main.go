package main

import (
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"time"

	//_ "net/http/pprof"

	"github.com/mazzegi/timtab"
)

// const profile = true
const profile = false

func main() {
	// go func() {
	// 	http.ListenAndServe("localhost:6060", nil)
	// }()
	var pw io.Writer
	if profile {
		f, err := os.Create("timtab.prof")
		if err != nil {
			panic(err)
		}
		defer f.Close()
		pw = f
	}

	fmt.Printf("start optimizing ...\n")
	cfg := simpleConfig()
	cfg.Prepare()

	t0 := time.Now()

	if profile {
		pprof.StartCPUProfile(pw)
	}
	table, rating := timtab.Optimize(cfg)
	if profile {
		pprof.StopCPUProfile()
	}

	fmt.Printf("time: %s\n", time.Since(t0))
	fmt.Printf("rating: %v\n", rating)
	fmt.Printf("schedules:\n")
	timtab.Dump(cfg, table)

}

func simpleConfig() *timtab.Configuration {
	scheds := &timtab.Schedules{}
	for day := 0; day < 5; day++ {
		for hour := 0; hour < 5; hour++ {
			scheds.Insert(timtab.Schedule{Day: timtab.Day(day), Hour: timtab.Hour(hour)})
		}
	}
	cfg := timtab.NewConfiguration(scheds)

	teach01 := &timtab.Teacher{ID: "Frank First", Name: "Frank First", Subjects: []timtab.SubjectID{"all"}}
	teach02 := &timtab.Teacher{ID: "Sabrina Second", Name: "Sabrina Second", Subjects: []timtab.SubjectID{"all"}}
	teach03 := &timtab.Teacher{ID: "Thomas Third", Name: "Thomas Third", Subjects: []timtab.SubjectID{"all"}}
	teach04 := &timtab.Teacher{ID: "Fanny Fourth", Name: "Fanny Fourth", Subjects: []timtab.SubjectID{"all"}}
	teachMusic := &timtab.Teacher{ID: "Ally Alround", Name: "Ally Alround", Subjects: []timtab.SubjectID{"music"}}
	cfg.AddTeachers(
		teach01,
		teach02,
		teach03,
		teach04,
		teachMusic,
	)

	sts01 := mkStudentsIds("01_fc_student_%d", 10)
	sts02 := mkStudentsIds("02_fc_student_%d", 10)
	sts03 := mkStudentsIds("03_fc_student_%d", 10)
	sts04 := mkStudentsIds("04_fc_student_%d", 10)

	cfg.AddClasses(
		mkClass(1, 4, "german", string(teach01.ID), sts01),
		mkClass(1, 4, "math", string(teach01.ID), sts01),
		mkClass(1, 2, "foo", string(teach01.ID), sts01),
		mkClass(1, 2, "bar", string(teach01.ID), sts01),
		mkClass(1, 3, "music", string(teachMusic.ID), sts01),

		mkClass(2, 4, "german", string(teach02.ID), sts02),
		mkClass(2, 4, "math", string(teach02.ID), sts02),
		mkClass(2, 3, "foo", string(teach02.ID), sts02),
		mkClass(2, 3, "bar", string(teach02.ID), sts02),
		mkClass(2, 3, "music", string(teachMusic.ID), sts02),

		mkClass(3, 5, "german", string(teach03.ID), sts03),
		mkClass(3, 5, "math", string(teach03.ID), sts03),
		mkClass(3, 3, "foo", string(teach03.ID), sts03),
		mkClass(3, 3, "bar", string(teach03.ID), sts03),
		mkClass(3, 4, "music", string(teachMusic.ID), sts02),

		mkClass(4, 5, "german", string(teach04.ID), sts04),
		mkClass(4, 5, "math", string(teach04.ID), sts04),
		mkClass(4, 5, "foo", string(teach04.ID), sts04),
		mkClass(4, 5, "bar", string(teach04.ID), sts04),
		mkClass(4, 4, "music", string(teachMusic.ID), sts04),
	)

	return cfg
}

func mkClass(level int, hours int, subject string, teacher string, sts []timtab.StudentID) *timtab.Class {
	name := fmt.Sprintf("%02d_%s_%s", level, subject, teacher)
	return &timtab.Class{
		//ID:           timtab.ClassID(name),
		Level:        level,
		HoursPerWeek: hours,
		Name:         name,
		Subject:      timtab.SubjectID(subject),
		Teacher:      timtab.TeacherID(teacher),
		Students:     sts,
	}
}

func mkStudentsIds(pattern string, count int) []timtab.StudentID {
	var sts []timtab.StudentID
	for i := 0; i < count; i++ {
		sts = append(sts, timtab.StudentID(fmt.Sprintf(pattern, count)))
	}
	return sts
}
