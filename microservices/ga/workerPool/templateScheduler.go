{{define "templateScheduler.go"}}{{.PavedroadInfo}}
package main

import "os"

// Schedule
type Scheduler interface {
	// Data methods
	// For schedulers
	GetSchedule() (Scheduler, error)
	UpdateSchedule() (Scheduler, error)
	CreateSchedule() (Scheduler, error)
	DeleteSchedule() (Scheduler, error)

	 // For jobs
  GetScheduledJobs() ([]byte, error)
  GetScheduleJob(UUID string) (httpStatusCode int, jsonBlob []byte, err error)
  //CreateScheduleJob(jb []byte)(error)
  //UpddateScheduleJob(jb []byte)(error)
  //DeleteScheduleJob(UUID string)(error)

	// Execution methods
	Init() error
	SetChannels(chan Job, chan Result, chan bool, chan os.Signal)
	Shutdown() error
  //Status()

	// Status methods
	Metrics() []byte
}
{{end}}
