{{define "httpJob.go"}}{{.PavedroadInfo}}
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
)

const (
	//HTTPJobType is type of Job for worker
        HTTPJobType string = "io.pavedraod.eventcollector.httpjob"
        //ClientTimeout in microseconds to timeout clients
        ClientTimeout int = 30
        //CountMaxOut is the upper limit for tracking countinuos Jobs
        CountMaxOut int = 10000
)

type httpJob struct {
        ctx           context.Context
        JobID         uuid.UUID `json:"job_id"`
        JobType       string    `json:"job_type"`
        client        *http.Client
        ClientTimeout int `json:"client_timeout"`
        // TODO: FIX to errors or custom errors
        jobErrors  []string
        JobURL     *url.URL  `json:"job_url"`
        Stats      httpStats `json:"stats"`
        ClientID   string    `json:"client_id"` //Possible business client
        sent       bool      //Internal
        continuous bool      //Internal
        count      int       //Tracking times sent for continuos Job
        startTime  time.Time
        endTime    time.Time
}


type httpStats struct {
	RequestTimedOut bool
	RequestTime			time.Duration
}

// Process methods

//ID:  returns the job ID
func (j *httpJob) ID() string {
	return j.JobID.String()
}

//UUID returns the job UUID
func (j *httpJob) UUID() uuid.UUID {
        return j.JobID
}

//Type: returns the job type
func (j *httpJob) Type() string {
	return HTTPJobType
}
//GetClientID: returns the clientID for the job
func (j *httpJob) GetClientID() string {
        return j.ClientID
}
//Excount: returns the number of time a continuous job
//executed for the current session.
//Currently will not increment past a CountMaxOut 
func (j *httpJob) ExCount() string {
        return strconv.FormatInt(int64(j.count), 10)
}
//Init: intialize a new job
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
	/*
                          defaults
                        j.sent is  false
                        j.count is 0
                        j.continuos is false

        */


	j.client = &http.Client{Timeout: time.Duration(j.ClientTimeout) * time.Second}

	return nil
}
//MarkSent: Mark job as sent before placing on channel
func (j *httpJob) MarkSent() {
        if !j.sent {
                j.sent = true
        }
        if j.continuous && j.count < CountMaxOut {
                j.count = j.count + 1

        }
}

//PausedAsSent is used to pause a continuous job.
//Job remains on jobList but is not sent.
func (j *httpJob) PausedAsSent() bool {
        if j.continuous && !j.sent {
                //After completion a continous job is
                //is marked as not sent.
                //Marking here as sent without expectation
                //for placement on the worker pool.
                j.sent = true
                return j.sent
        }
        return false
}
// GetSent: reports on the jobs status
func (j *httpJob) GetSent() bool {
        return j.sent
}

//SetContinuous: sets the job for continuous execution
func (j *httpJob) SetContinuous() {
        j.continuous = true

}
//IsContinuous: reports if the job is a continuos job
func (j *httpJob) IsContinuous() bool {
        return j.continuous

}

//Run: executes the job.
func (j *httpJob) Run(mWg *sync.WaitGroup) (result Result, err error) {
	defer mWg.Done()
	req, err := http.NewRequest("GET", j.JobURL.String(), nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//start := time.Now()
	j.startTime = time.Now()
        j.endTime = j.startTime
	resp, err := j.client.Do(req)

        j.endTime = time.Now()
        j.Stats.RequestTime = j.endTime.Sub(j.startTime)

  // client errors are handled with errors.New()
  // so there is no defined set to check for
	if err != nil {
		j.Stats.RequestTimedOut = true
		fmt.Println(err)
		return nil, err
	}

	defer resp.Body.Close()

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	md := j.buildMetadata(resp)
	jrsp := &httpResult{job: j,
		metaData: md,
		payload:	payload}

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
// Pause a job
func (j *httpJob) Pause() (status string, err error) {
	//Not implemented
	//Implementation, like PausedAsSent, with logic 
	//sent  to true and continuous to false if true
	return "paused", nil
}

//Shutdown: stop a job?
func (j *httpJob) Shutdown() error {
	//Not implemented
	return nil
}
//Errors: nill error
func (j *httpJob) Errors() []error {
	return nil

}
//Metrics: return the metrics for a job
func (j *httpJob) Metrics() []byte {
	jblob, err := json.Marshal(j.Stats)
	if err != nil {
		return []byte("Marshal metrics failed")
	}

	return jblob
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
