package main

import (
	"github.com/KingsEpic/kinglib"
	// "fmt"
	"gopkg.in/v0/qml"
)

type JobModel struct {
	Len  int
	Jobs []*kinglib.JobStatus
}

func (jm *JobModel) Update() {
	qml.Lock()

	jm.Len = len(w.Jobs)

	jm.Jobs = make([]*kinglib.JobStatus, jm.Len)
	count := 0
	for _, js := range w.Jobs {
		jm.Jobs[count] = js
		count++
	}
	// fmt.Printf("JobModel Length: %d\n", jm.Len)
	qml.Unlock()
	qml.Changed(jm, &jm.Len)
}

func (jm *JobModel) Get(index int) *kinglib.JobStatus {
	return jm.Jobs[index]
}
