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

// TODO: catch connection timeout to avoid traceback
const (
	HTTPJobType string = "httpJob"
	ClientTimeout int  = 30
)

type httpJob struct {
  ctx           context.Context `json:"ctx"`
  JobID         uuid.UUID       `json:"job_id"`
  JobType       string          `json:"job_type"`
  client        *http.Client    `json:"client"`
  ClientTimeout int             `json:"client_timeout"`
  // FIX to errors or custom errors
  jobErrors []string  `json:"job_errors"`
  JobURL    *url.URL  `json:"job_url"`
  Stats     httpStats `json:"stats"`
}

type httpStats struct {
  RequestTimedOut bool
  RequestTime     time.Duration
}

func (j *httpJob) GetJob(ID string) (jblob []byte, err error) {
	jblob, err = json.Marshal(j)
	if err != nil {
		return nil, err
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
	return j.JobID.String()
}

func (j *httpJob) Type() string {
	return HTTPJobType
}

func (j *httpJob) Init() error {

	// Generate UUID
	j.JobID = uuid.New()

	// Set job type
	j.JobType = HTTPJobType

  j.Stats.RequestTimedOut = false

	// Set http client options
  if j.ClientTimeout == 0 {
    j.ClientTimeout = ClientTimeout
  }

  j.client = &http.Client{Timeout: time.Duration(j.ClientTimeout) * time.Second}

	return nil
}

func (j *httpJob) New(in *chan Job, out *chan Result) error {

	return nil
}

func (j *httpJob) Run() (result Result, err error) {
	  req, err := http.NewRequest("GET", j.JobURL.String(), nil)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }

  start := time.Now()
  resp, err := j.client.Do(req)
  end := time.Now()
  j.Stats.RequestTime = end.Sub(start)

  if err != nil {
    j.Stats.RequestTimedOut = true
    fmt.Println(err)
    return nil, err
  }

  payload, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    fmt.Println(err)
    return nil, err
  }
  //fmt.Printf("HTTP Status: %v for %s\n", resp.StatusCode, req.URL.String())

	md := j.buildMetadata(resp)
  jrsp := &httpResult{job: j,
    metaData: md,
    payload:  payload}

  return jrsp, nil
}

// buildMetadata returns a map of strings with an http.Response encoded
func (j *httpJob) buildMetadata(resp *http.Response) map[string]string {
  md := make(map[string]string)
  md["StatusCode"] = string(resp.StatusCode)
  md["Proto"] = resp.Proto

  for n, v := range resp.Header {
    var hv string
    for _, s := range v {
      hv = hv + s + " "
    }
    md[n] = hv
  }

  md["RemoteAddr"] = resp.Request.RemoteAddr
  md["Method"] = resp.Request.Method

  return md
}

func (j *httpJob) Pause() (status string, err error) {
	return "paused", nil
}

func (j *httpJob) Shutdown() error {
	return nil
}

func (j *httpJob) Errors() []error {
	return nil

}

func (j *httpJob) Metrics() []byte {
	jblob, err := json.Marshal(j.Stats)
  if err != nil {
    return []byte("Marshal metrics failed")
  }

  return jblob
}{{end}}
