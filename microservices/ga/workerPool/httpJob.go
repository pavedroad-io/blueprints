{{define "httpJob.go"}}{{.PavedroadInfo}}
package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

const (
	HTTPJobType string = "httpJob"
)

type httpJob struct {
	ctx		  context.Context
	id      uuid.UUID
	jobType string
	client  *http.Client
	// FIX to errors or custom errors
	jobErrors []string
	jobURL    *url.URL
}

type httpWorker struct {
	currentJob  *httpJob
	previousJob *httpJob
}

func (j *httpJob) GetJob(ID string) (jblob []byte, err error) {
	jblob, err = json.Marshal(j)
	if err != nil {
		return err
	}

	return jblob, nil
}

func (j *httpJob) UpdateJob(jblob []byte) error {
	err := json.Unmarshal(jblob, j)
	if err != nil {
		return err
	}

	return nil
}

func (j *httpJob) CreateJob(jblob []byte) error {
	err := json.Unmarshal(jblob, j)
	if err != nil {
		return err
	}

	return nil
}

func (j *httpJob) DeleteJob(ID string) error {
	// Cancel request
	return nil

}

// Process methods
func (j *httpJob) ID() string {
	return j.id.String()
}

func (j *httpJob) Type() string {
	return HTTPJobType
}

func (j *httpJob) Init() error {

	// Generate UUID
	j.id = uuid.New()

	// Set job type
	j.jobType = HTTPJobType

	// Set http client options
}

func (j *httpJob) Errors() []error {
	return nil
}

func (j *httpWorker) New(in *chan Job, out *chan Result) error {

	return nil
}

func (j *httpWorker) Run() (result Result, err error) {
	return nil, nil
}

func (j *httpWorker) Pause() (status string, err error) {
	return "paused", nil
}

func (j *httpWorker) Shutdown() error {
	return nil

}

func (j *httpWorker) Errors() []error {
	return nil

}

func (j *httpWorker) Metrics() []byte {

	return nil
}{{end}}
