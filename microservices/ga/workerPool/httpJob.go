{{define "httpJob.go"}}{{.PavedroadInfo}}
package main

import (
	"context"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "net/http"
  "net/url"
  "time"

	"github.com/google/uuid"
)

const (
	HTTPJobType string = "httpJob"
	ClientTimeout int    = 5
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
  if j.clientTimeout == 0 {
    j.clientTimeout = ClientTimeout
  }

  j.client = &http.Client{Timeout: time.Duration(j.clientTimeout) * time.Second}

	return nil
}

func (j *httpWorker) New(in *chan Job, out *chan Result) error {

	return nil
}

func (j *httpWorker) Run() (result Result, err error) {
	  req, err := http.NewRequest("GET", j.jobURL.String(), nil)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  resp, err := j.client.Do(req)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  payload, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  fmt.Printf("HTTP Status: %v for %s\n", resp.StatusCode, req.URL.String())

  jrsp := &httpResult{job: j,
    metaData: nil,
    payload:  payload}

  return jrsp, nil
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
