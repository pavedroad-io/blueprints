{{define "templateScheduler.go"}}{{.PavedroadInfo}}
package main

import "os"

// Schedule
type Scheduler interface {
	// Data methods
	GetSchedule() (Scheduler, error)
	UpdateSchedule() (Scheduler, error)
	CreateSchedule() (Scheduler, error)
	DeleteSchedule() (Scheduler, error)

	// Execution methods
	Init() error
	SetChannels(chan Job, chan Result, chan bool, chan os.Signal)
	Shutdown() error
  //Status()

	// Status methods
	Metrics() []byte
}
{{end}}
