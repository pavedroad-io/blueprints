{{define "templateScheduler.go"}}{{.PavedroadInfo}}
package main

import "os"

// Schedule
type Scheduler interface {
	// Data methods
	// For schedulers
	GetSchedule() (httpStatusCode int, jsonBlob []byte, err error)
	UpdateSchedule(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error)
	CreateSchedule(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error)
	DeleteSchedule() (httpStatusCode int, jsonb []byte, err error)

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
