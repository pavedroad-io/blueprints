{{define "templateJob.go"}}package main

import (
        "github.com/google/uuid"
        "sync"
)


// Job interface abstraction for worker pools
// ID() a unique ID assigned to each job
// Type() string indicating the type of job
//   for example, "HTTP Request"
// ClientID() Client for the job (future usage)
// Execute executes the job returning a Result
// Errors returns a list of errors to be logged
type Job interface {
	// Process methods
	ID() string
	UUID() uuid.UUID
	Type() string
	GetClientID() string
	Init() error
	Run(mWg *sync.WaitGroup) (result Result, err error)
	Pause() (status string, err error)
	Shutdown() error
	Errors() []error
	GetSent() bool
        MarkSent()
        SetContinuous()
        ExCount() string
        PausedAsSent() bool
        IsContinuous() bool


	// Status methods
	Metrics() []byte
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
