
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
//
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// Type of Schedulers
const (
	constantIntervaleScheduler = "Constant interval scheduler"
)

// Defaults
const (
	defaultConstantInterval = 10
	defaultResponseTimeJobs = 10
)

// Metrics constants
const (
	schedulerIterations             = "scheduler_iterations"
	jobsSent                        = "jobs_sent"
	jobListSize                     = "job_list_size"
	resultsReceived                 = "results_received"
	currentJobChannelUtilization    = "current_job_channel_utilization"
	currentJobChannelCapacity       = "current_job_channel_capacity"
	currentResultChannelUtilization = "current_result_channel_utilization"
	currentResultChannelCapacit     = "current_result_channel_capacity"
	numberOfJobTimedOut             = "number_of_jobs_sent"
	averageJobProcessingTime        = "average_job_processing_time"
)

type httpScheduler struct {
	jobList               []*httpJob
	schedulerJobChan      chan Job       // Channel to read jobs from
	schedulerResponseChan chan Result    // Channel to write repose to
	schedulerDone         chan bool      // Shutdown initiated by application
	schedulerInterrupt    chan os.Signal // Shutdown initiated by OS
	metrics               httpSchedulerMetrics
	mux                   *sync.Mutex
	schedule              httpSchedule
}

// httpSchedule holds the type of scheduler and it's configuration

type httpSchedule struct {
	ScheduleType        string `json:"schedule_type"`
	SendIntervalSeconds int64  `json:"send_interval_seconds"`
	ResponseTimeJobs    int    `json:"response_time_jobs"`
}

// httpSchedulerMetrics hold metrics about the Scheduler, Jobs, and Results
// We export attributes we want included in the JSON output
type httpSchedulerMetrics struct {
	StartTime time.Time      `json:"start_time"`
	UpTime    time.Duration  `json:"up_time"`
	Counters  map[string]int `json:"counters"`
	mux       *sync.Mutex
}

func (s *httpScheduler) MetricToJSON() ([]byte, error) {
	s.metrics.mux.Lock()
	defer s.metrics.mux.Unlock()
	jb, e := json.Marshal(s.metrics)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}
	return jb, nil
}

func (s *httpScheduler) MetricSetStartTime() {
	s.metrics.mux.Lock()
	s.metrics.StartTime = time.Now()
	s.metrics.mux.Unlock()
}

func (s *httpScheduler) MetricUpdateUpTime() (uptime time.Duration) {
	s.metrics.mux.Lock()
	ct := time.Now()
	s.metrics.UpTime = ct.Sub(s.metrics.StartTime)
	s.metrics.mux.Unlock()
	return s.metrics.UpTime
}

func (s *httpScheduler) MetricInc(key string) {
	s.metrics.mux.Lock()
	s.metrics.Counters[key]++
	s.metrics.mux.Unlock()
}

func (s *httpScheduler) MetricSet(key string, value int) {
	s.metrics.mux.Lock()
	s.metrics.Counters[key] = value
	s.metrics.mux.Unlock()
}

func (s *httpScheduler) MetricValue(key string) int {
	s.metrics.mux.Lock()
	defer s.metrics.mux.Unlock()
	return s.metrics.Counters[key]
}

// UpdateJobList to a new list safely
func (s *httpScheduler) UpdateJobList(newJobList []*httpJob) {
	s.mux.Lock()
	s.jobList = newJobList
	s.mux.Unlock()
}

// TODO: Move this to dispatcher, it is generic

// A []listJobsResponse is a single job but returned as a list
//
// swagger:response listJobResponse
type listJobsResponse struct {
	// in: body

	//id: uuid for this job
	ID string `json:"id"`
	// url for this http request
	URL string `json:"url"`

	//type: of job the represents
	Type string `json:"type"`
}

// Required object methods for interface
//
// GetScheduledJobs returns a list of job IDs and URL
func (s *httpScheduler) GetScheduledJobs() ([]byte, error) {
	var response []listJobsResponse

	for _, v := range s.jobList {
		var newRow = listJobsResponse{}
		newRow.ID = v.JobID.String()
		newRow.URL = v.JobURL.String()
		newRow.Type = v.JobType
		response = append(response, newRow)
	}

	jb, e := json.Marshal(response)
	if e != nil {
		return nil, e
	}
	return jb, nil
}

// GetScheduleJob returns a single job matching the UUID provided
func (s *httpScheduler) GetScheduleJob(UUID string) (httpStatusCode int, jsonBlob []byte, err error) {
	var newRow = listJobsResponse{}

	for _, v := range s.jobList {
		if v.ID() == UUID {
			newRow.ID = v.ID()
			newRow.URL = v.JobURL.String()
			newRow.Type = v.JobType
			break
		}
	}

	// Not found response
	if newRow.ID == "" {
		msg := fmt.Sprintf("{\"error\": \"Not found\", \"UUID\": %v}", UUID)
		return http.StatusNotFound, []byte(msg), nil
	}

	jb, e := json.Marshal(newRow)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"json.Marshal failed\", \"Error\": \"%v\"}", e.Error())
		return http.StatusInternalServerError, []byte(msg), e
	}

	return http.StatusOK, jb, nil
}

// UpdateScheduleJob decodes json data into a job and updates the jobID
// Returns httpStatusCode, JSON body, and error code
func (s *httpScheduler) UpdateScheduleJob(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error) {
	var updateData = listJobsResponse{}
	var oldJobID, newJobID string
	var newJobList []*httpJob
	foundJob := false

	e := json.Unmarshal(jsonBlob, &updateData)

	if e != nil {
		fmt.Println("Unmarshal failed", e.Error())
		msg := fmt.Sprintf("{\"error\": \"json.Unmarshal failed\", \"Error\": \"%v\"}", e.Error())
		return http.StatusBadRequest, []byte(msg), e
	}

	for _, v := range s.jobList {
		if v.ID() == updateData.ID {
			newJob := httpJob{}
			pu, err := url.Parse(updateData.URL)
			if err != nil {
		 fmt.Println(err)
                                msg := fmt.Sprintf("{\"error\": \"bad job url\", \"Error\": \"%v\"}", err)
                                return http.StatusBadRequest, []byte(msg), err
			}
			newJob.JobURL = pu
			e = newJob.Init()
                        if e != nil {
                                msg := fmt.Sprintf("{\"error\": \"job init failed\", \"Error\": \"%v\"}", e.Error())
                                return http.StatusInternalServerError, []byte(msg), e
                        }

			newJobList = append(newJobList, &newJob)
			oldJobID = v.ID()
			newJobID = newJob.ID()
			foundJob = true
			continue
		}
		newJobList = append(newJobList, v)
	}

	// Handle 404 for Job not found
	if !foundJob {
		msg := fmt.Sprintf("{\"error\": \"Not found\", \"UUID\": %v}", updateData.ID)
		return http.StatusNotFound, []byte(msg), nil
	}

	// Update job list and return
	s.UpdateJobList(newJobList)

	msg := fmt.Sprintf("{\"success\": \"Old job %v replaced by new job %v\"}",
		oldJobID, newJobID)
	return http.StatusOK, []byte(msg), nil
}

// CreateScheduleJob decodes json data into a job and inserts into jobList
// Returns httpStatusCode, JSON body, and error code
func (s *httpScheduler) CreateScheduleJob(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error) {
	var newJobType = listJobsResponse{}

	e := json.Unmarshal(jsonBlob, &newJobType)

	if e != nil {
		fmt.Println("Unmarshal failed", e.Error())
		msg := fmt.Sprintf("{\"error\": \"json.Unmarshal failed\", \"Error\": \"%v\"}", e.Error())
		return http.StatusBadRequest, []byte(msg), e
	}

	newJob := httpJob{}
	pu, err := url.Parse(newJobType.URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	newJob.JobURL = pu
	e = newJob.Init()
        if e != nil {
                msg := fmt.Sprintf("{\"error\": \"job init failed\", \"Error\": \"%v\"}", e.Error())
                return http.StatusInternalServerError, []byte(msg), e
        }

	s.jobList = append(s.jobList, &newJob)

	msg := fmt.Sprintf("{\"success\": \"new job %v added\"}", newJob.ID())
	return http.StatusCreated, []byte(msg), nil
}

// DeleteScheduleJob delete the job with ID == uuid
// Returns httpStatusCode, JSON body, and error code
func (s *httpScheduler) DeleteScheduleJob(uuid string) (httpStatusCode int, jsonb []byte, err error) {
	var newJobList []*httpJob
	var foundJob = false

	for _, v := range s.jobList {
		if v.ID() == uuid {
			foundJob = true
			continue
		}
		newJobList = append(newJobList, v)
	}

	// Handle 404 for Job not found
	if !foundJob {
		msg := fmt.Sprintf("{\"error\": \"Not found\", \"UUID\": %v}", uuid)
		return http.StatusNotFound, []byte(msg), nil
	}

	// Update job list and return
	s.UpdateJobList(newJobList)

	msg := fmt.Sprintf("{\"success\": \"Job %v deleted\"}", uuid)
	return http.StatusOK, []byte(msg), nil
}

// Object methods for schedules
func (s *httpScheduler) GetSchedule() (httpStatusCode int, jsonBlob []byte, err error) {

	jb, e := json.Marshal(s.schedule)
	if e != nil {
		msg := fmt.Sprintf("{\"json.Marsha failed\": \"%v\"}", e)
		return http.StatusInternalServerError, []byte(msg), e
	}

	return http.StatusOK, jb, nil
}

func (s *httpScheduler) UpdateSchedule(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error) {

	us := httpSchedule{}
	e := json.Unmarshal(jsonBlob, &us)
	if e != nil {
		msg := fmt.Sprintf("{\"json.Unmarshal failed\": \"%v\"}", e)
		return http.StatusInternalServerError, []byte(msg), e
	}

	s.mux.Lock()
	s.schedule.SendIntervalSeconds = us.SendIntervalSeconds
	s.mux.Unlock()

	msg := fmt.Sprintf("{\"Status\": \"Success\", \"New interval seconds\": %v}",
		s.schedule.SendIntervalSeconds)

	return http.StatusOK, []byte(msg), nil
}

// CreateSchedule replace current schdule objec
func (s *httpScheduler) CreateSchedule(jsonBlob []byte) (httpStatusCode int, jsonb []byte, err error) {

	us := httpSchedule{}
	e := json.Unmarshal(jsonBlob, &us)
	if e != nil {
		msg := fmt.Sprintf("{\"json.Unmarshal failed\": \"%v\"}", e)
		return http.StatusInternalServerError, []byte(msg), e
	}

	s.mux.Lock()
	s.schedule.ScheduleType = us.ScheduleType
	s.schedule.SendIntervalSeconds = us.SendIntervalSeconds
	s.mux.Unlock()

	msg := fmt.Sprintf("{\"Status\": \"Success\", \"New interval seconds\": %v}",
		s.schedule.SendIntervalSeconds)
	return http.StatusCreated, []byte(msg), nil
}

func (s *httpScheduler) DeleteSchedule() (httpStatusCode int, jsonb []byte, err error) {

	s.schedulerDone <- true
	msg := fmt.Sprintf("{\"Status\": \"Success scheduler stopped\"}")
	return http.StatusOK, []byte(msg), nil
}

// SetChannels initializes channels the dispatcher has created inside
// of the scheduler
func (s *httpScheduler) SetChannels(j chan Job, r chan Result, b chan bool, i chan os.Signal) {
	s.schedulerJobChan = j
	s.schedulerResponseChan = r
	s.schedulerDone = b
	s.schedulerInterrupt = i

	return
}

// Process methods
func (s *httpScheduler) Init() error {
	urlList := []string{
		"https://api.chucknorris.io/jokes/random",
		"https://swapi.dev/api/people/1/",
		"https://swapi.dev/api/people/2/",
		"https://swapi.dev/api/people/3/"}

	s.metrics.Counters = make(map[string]int)

	s.mux = new(sync.Mutex)
	s.metrics.mux = new(sync.Mutex)
	s.schedule.ResponseTimeJobs = defaultResponseTimeJobs

	s.schedule.SendIntervalSeconds = defaultConstantInterval
	s.schedule.ScheduleType = constantIntervaleScheduler

	for _, u := range urlList {
		newJob := httpJob{}
		pu, err := url.Parse(u)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		newJob.JobURL = pu

		// Set type and ID and http.Client
		em := newJob.Init()
		if em != nil {
                        return em
                }

		s.jobList = append(s.jobList, &newJob)
	}

	return nil
}

func (s *httpScheduler) Run() error {
	go s.RunScheduler()
	go s.RunResultsReader()

	return nil
}

func (s *httpScheduler) RunScheduler() error {
	s.MetricSetStartTime()
	for {
		s.MetricInc(schedulerIterations)
		for _, j := range s.jobList {
			s.schedulerJobChan <- j
			s.MetricInc(jobsSent)
			s.MetricSet(currentJobChannelCapacity, cap(s.schedulerJobChan))
			s.MetricSet(currentJobChannelUtilization, len(s.schedulerJobChan))
			s.MetricSet(jobListSize, len(s.jobList))
		}
		s.MetricUpdateUpTime()
		select {
		case <-s.schedulerDone:
			return nil
		case <-s.schedulerInterrupt:
			return nil
		default:
			time.Sleep(time.Duration(s.schedule.SendIntervalSeconds) * time.Second)
		}
	}
}

// ComputeAverageResponseTime Keep track of the last N responses
func (s *httpScheduler) ComputeAverageResponseTime(jt []int, newTime int) ([]int, int) {
	currentLength := len(jt)
	desiredLength := currentLength - 9
	if currentLength >= s.schedule.ResponseTimeJobs {
		jt = jt[desiredLength:currentLength]
	}
	jt = append(jt, newTime)
	currentLength = len(jt)

	var totalTime int = 0
	for _, t := range jt {
		totalTime += t
	}

	return jt, totalTime / currentLength
}

func (s *httpScheduler) RunResultsReader() error {
	jobTimes := make([]int, 0, s.schedule.ResponseTimeJobs)
	log.Println("Starting result reader")
	for {
		select {
		case currentResult := <-s.schedulerResponseChan:
			s.MetricInc(resultsReceived)
			s.MetricSet(currentResultChannelCapacit, cap(s.schedulerJobChan))
			s.MetricSet(currentResultChannelUtilization, len(s.schedulerJobChan))
			log.Printf("Processing response for job ID %v\n", currentResult.Job().ID())
			j := currentResult.Job()
			if j.(*httpJob).Stats.RequestTimedOut {
				s.MetricInc(numberOfJobTimedOut)
			}
			jt, avg := s.ComputeAverageResponseTime(jobTimes, int(j.(*httpJob).Stats.RequestTime))
			s.MetricSet(averageJobProcessingTime, avg)
			jobTimes = jt

		case done := <-s.schedulerDone:
			if done {
				return nil
			}

		case <-s.schedulerInterrupt:
			return nil
		}
	}
}

func (s *httpScheduler) Pause() []byte {

	return nil
}

func (s *httpScheduler) Shutdown() error {

	return nil
}

// Status methods
// DeleteSchedule stops go scheduler goroutine
func (s *httpScheduler) Metrics() []byte {
	jb, e := s.MetricToJSON()
	if e != nil {
		return nil
	}
	return jb
}