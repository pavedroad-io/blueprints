// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
//
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/google/uuid"
)

const (
	//HTTPJobType is type of Job from scheduler
	HTTPJobType string = "io.pavedraod.eventcollector.httpjob"
	//ClientTimeout in seconds to timeout client jobs
	ClientTimeout int = 30
)

type httpJob struct {
	ctx           context.Context `json:"ctx"`
	JobID         uuid.UUID       `json:"job_id"`
	Method        string          `json:"method"`
	Payload       []data          `json:"payload"`
	JobType       string          `json:"job_type"`
	client        *http.Client    `json:"client"`
	ClientTimeout int             `json:"client_timeout"`
	// TODO: FIX to errors or custom errors
	jobErrors []string  `json:"jobErrors"`
	JobURL    *url.URL  `json:"job_url"`
	Stats     httpStats `json:"stats"`
}

type httpStats struct {
	RequestTimedOut bool
	RequestTime     time.Duration
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

	// client errors are handled with errors.New()
	// so there is no defined set to check for
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

	md := j.buildMetadata(resp)
	jrsp := &httpResult{job: j,
		metaData: md,
		payload:  payload}

	return jrsp, nil
}

// buildMetadata returns a map of strings with an http.Response encoded
func (j *httpJob) buildMetadata(resp *http.Response) map[string]string {
	md := make(map[string]string)
	md["StatusCode"] = string(rune(resp.StatusCode))
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

func (j *httpJob) newJob(url url.URL) httpJob {
	newJob := httpJob{}
	pu, err := url.Parse(url.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	newJob.JobURL = pu

	// Set type and ID and http.Client
	newJob.Init()
	return newJob
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
}
