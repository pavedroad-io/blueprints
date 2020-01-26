
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
//
package main

import "os"

// Scheduler defines the interfaces a scheduler must implement
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
	Run() error
	//RestartScheduler() error
	//RestartResultsCollector() error

	// Status methods
	Metrics() []byte
	//Status()
	//RestartScheduler() error
	//RestartResultsCollector() error
}

// SchedulerStatus tracks if the scheduler and Results collectors are running
type SchedulerStatus struct {
	SchedulerRunning       bool
	ResultCollectorRunning bool
}