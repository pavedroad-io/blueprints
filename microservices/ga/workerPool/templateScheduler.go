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
  UpdateScheduleJob(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error)
  CreateScheduleJob(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error)
  DeleteScheduleJob(UUID string) (httpStatusCode int, jsonb []byte, err error)


	// Execution methods
	Init() error
	SetChannels(chan Job, chan Result, chan bool, chan os.Signal)
	Shutdown() error
  //Status()

	// Status methods
	Metrics() []byte
}
{{end}}
