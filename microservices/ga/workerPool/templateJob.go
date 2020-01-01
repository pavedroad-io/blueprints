{{define "templateJob.go"}}package main

// Job interface abstraction for worker pools
// ID() a unique ID assigned to each job
// Type() string indicating the type of job
//   for example, "HTTP Request"
// Execute executes the job returning a Result
// Errors returns a list of errors to be logged
type Job interface {
	// Object methods
	GetJob(ID string) (jblob []byte, err error)
	UpdateJob(jblob []byte) (error)
	CreateJob(jblob []byte) (error)
	DeleteJob(ID string) (error)

	// Process methods
	ID() string
	Type() string
	Init() error
	Run() (result Result, err error)
	Pause() (status string, err error)
	Shutdown() error
	Errors() []error

	// Status methods
	Metrics() []byte
}{{end}}
