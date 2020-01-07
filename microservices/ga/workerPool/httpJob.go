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
	clientTimeout int
	// FIX to errors or custom errors
	jobErrors []string
	jobURL    *url.URL
  Stats     httpStats
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

  j.Stats.RequestTimedOut = false

	// Set http client options
  if j.clientTimeout == 0 {
    j.clientTimeout = ClientTimeout
  }

  j.client = &http.Client{Timeout: time.Duration(j.clientTimeout) * time.Second}

	return nil
}

func (j *httpJob) New(in *chan Job, out *chan Result) error {

	return nil
}

func (j *httpJob) Run() (result Result, err error) {
	  req, err := http.NewRequest("GET", j.jobURL.String(), nil)
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
