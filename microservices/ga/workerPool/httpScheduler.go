{{define "httpScheduler.go"}}{{.PavedroadInfo}}
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// Type of Schedulers
const (
	constantIntervaleScheduler = "Constant interval scheduler"
	//Per Worker
        trackMaxJobs string = "track_max_job_count"
)

// Defaults
const (
	defaultConstantInterval = 10
	defaultResponseTimeJobs = 10
	defaultMaxJobCount      = 20000
        defaultInitMaxJobCount  = 5000
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

var MaxJobCount int //Per Worker

//Clients for testing possible extension
var testClients = []string{
        "MT1202",
        "PR8724",
        "FB90237",
        "SC52415",
}

// SchedulerStatus tracks the scheduler and Results collector 
type SchedulerStatus struct {
        SchedulerRunning       bool
        ResultCollectorRunning bool
}


type httpScheduler struct {
	jobList               []*httpJob
	schedulerJobChan      chan Job       // Channel to read jobs from
	schedulerResponseChan chan Result    // Channel to write repose to
	schedulerDone         chan bool      // Shutdown initiated by application
	schedulerInterrupt    chan os.Signal // Shutdown initiated by OS
	metrics               httpSchedulerMetrics
	mutex                   *sync.Mutex
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
	mutex       *sync.Mutex
}

//getTestClient is temporary for business client inclusion, remove as needed
func getTestClient() (clientID string) {
        //Seed to improve selection
        rand.Seed(time.Now().UnixNano())
        return testClients[rand.Intn(len(testClients))]

}


func (s *httpScheduler) MetricToJSON() ([]byte, error) {
        if s.metrics.mutex == nil {
                return nil, nil
        }
	s.metrics.mutex.Lock()
	defer s.metrics.mutex.Unlock()
	jb, e := json.Marshal(s.metrics)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}
	return jb, nil
}

func (s *httpScheduler) MetricSetStartTime() {
	s.metrics.mutex.Lock()
	s.metrics.StartTime = time.Now()
	s.metrics.mutex.Unlock()
}

func (s *httpScheduler) MetricUpdateUpTime() (uptime time.Duration) {
	s.metrics.mutex.Lock()
	{
	ct := time.Now()
	s.metrics.UpTime = ct.Sub(s.metrics.StartTime)
}
	s.metrics.mutex.Unlock()
	return s.metrics.UpTime
}

func (s *httpScheduler) MetricInc(key string) {
	s.metrics.mutex.Lock()
	s.metrics.Counters[key]++
	s.metrics.mutex.Unlock()
}

func (s *httpScheduler) MetricSet(key string, value int) {
	s.metrics.mutex.Lock()
	s.metrics.Counters[key] = value
	s.metrics.mutex.Unlock()
}

func (s *httpScheduler) MetricValue(key string) int {
	s.metrics.mutex.Lock()
	defer s.metrics.mutex.Unlock()
	return s.metrics.Counters[key]
}

// UpdateJobList to a new list safely
// Expected internal use only.
func (s *httpScheduler) UpdateJobList(newJobList []*httpJob) {
	s.mutex.Lock()
	s.jobList = newJobList
	s.mutex.Unlock()
}

// A []listJobsResponse is a single job but returned as a list
//
// swagger:response listJobResponse
// TODO: Move this to dispatcher, it is generic
type listJobsResponse struct {
	// in: body

	//id: uuid for this job
	ID string `json:"id"`
	// url for this http request
	URL string `json:"url"`

	//type: of job the represents
	Type string `json:"type"`
	ClientID string `json:"client_id"` //For future enhancement
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
		newRow.ClientID = v.ClientID
		response = append(response, newRow)
	}

	jb, e := json.Marshal(response)
	if e != nil {
		return nil, e
	}
	return jb, nil
}

// StopContinuousJob returns the status for a single job status change request
func (s *httpScheduler) StopContinuousJob(UUID string) (httpStatusCode int, jsonBlob []byte, err error) {
        msg := ""

        for ji := range s.jobList {
                //how to lock
                if s.jobList[ji].ID() == UUID {
                        s.mutex.Lock()
                        if s.jobList[ji].PausedAsSent() {
                                msg = fmt.Sprintf("Job UUID: %v, was stopped.", UUID)
                        } else {
                                msg = fmt.Sprintf("Job UUID: %v, not stopped.", UUID)
                        }
                        s.mutex.Unlock()
                        break
                }
        }

        // Not found response
        if msg == "" {
                msg = fmt.Sprintf("Job UUID: %v, not found.", UUID)
                return http.StatusNotFound, []byte(msg), nil
        }
      return http.StatusOK, []byte(msg), nil

}


// GetScheduleJob returns a single job matching the UUID provided
func (s *httpScheduler) GetScheduleJob(UUID string) (httpStatusCode int, jsonBlob []byte, err error) {
	var newRow = listJobsResponse{}

	for _, v := range s.jobList {
		if v.ID() == UUID {
			newRow.ID = v.ID()
			newRow.URL = v.JobURL.String()
			newRow.Type = v.JobType
			newRow.ClientID = v.ClientID
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
	//var oldJobID, newJobID string
	//var newJobList []*httpJob
	foundJob := false

	e := json.Unmarshal(jsonBlob, &updateData)

	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"json.Unmarshal failed\", \"Error\": \"%v\"}", e.Error())
		return http.StatusBadRequest, []byte(msg), e
	}

	//TODO: Information on what job status can be updated?
        //Currently, job could alredy be sent and not continuous
   for ji := range s.jobList {
                if s.jobList[ji].ID() == updateData.ID {
                        //newJob := httpJob{}
                        pu, err := url.Parse(updateData.URL)
                        if err != nil {
                                fmt.Println(err)
                                msg := fmt.Sprintf("{\"error\": \"bad url: \", \"Error\": \"%v\"}", err)
                                return http.StatusBadRequest, []byte(msg), err
                        }

			//Logic adjusted to modify job url and business client,
                        //Not the JobId
                        s.mutex.Lock()
                        //newJob.JobURL = pu
                        s.jobList[ji].JobURL = pu
                        if updateData.ClientID != "" {
                                s.jobList[ji].ClientID = updateData.ClientID
                        }
                        s.mutex.Unlock()

			 //Put the udated job on the newlist
                        //newJobList = append(newJobList, &newJob)

                        foundJob = true
                        break

		}

		//put all other jobs on the newlist
                //newJobList = append(newJobList, v)

	}

	// Handle 404 for Job not found
	if !foundJob {
		msg := fmt.Sprintf("{\"error\": \"Not found\", \"UUID\": %s}", updateData.ID)
		return http.StatusNotFound, []byte(msg), nil
	}
	msg := fmt.Sprintf("{\"success\": \"job %s information was updated to Client: %s, URL: %s.\"}", updateData.ID, updateData.ClientID, updateData.URL)

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

	/* Future: could call another mircoservice to validate 
	newJobType.ClientID and return if not valid. */

	newJob := httpJob{}
	pu, err := url.Parse(newJobType.URL)
	if err != nil {
		fmt.Println(err)
		 msg := fmt.Sprintf("{\"error\": \"bad job url\", \"Error\": \"%v\"}", err)
                return http.StatusBadRequest, []byte(msg), err
	}
	e = newJob.Init()
	 if e != nil {
                msg := fmt.Sprintf("{\"error\": \"job init failed\", \"Error\": \"%v\"}", e.Error())
                return http.StatusInternalServerError, []byte(msg), e
        }
        //* See getTestClient comments
        if len(newJobType.ClientID) == 0 {
                newJobType.ClientID = getTestClient()
        }

	newJob.JobURL = pu
        newJob.ClientID = newJobType.ClientID
	newJob.SetContinuous()
        s.mutex.Lock()
        s.jobList = append(s.jobList, &newJob)
        s.mutex.Unlock()

	msg := fmt.Sprintf("{\"success\": \"new job %v added\"}", newJob.ID())
	return http.StatusCreated, []byte(msg), nil
}

// DeleteScheduleJob delete the job with ID == uuid
// Returns httpStatusCode, JSON body, and error code
func (s *httpScheduler) DeleteScheduleJob(uuid string) (httpStatusCode int, jsonb []byte, err error) {

	pstat := s.removeJob(uuid)
        msg := ""

	// Handle 404 for Job not found
	if pstat == http.StatusNotFound {
		msg = fmt.Sprintf("{\"error\": \"Not found\", \"UUID\": %v}", uuid)
	} else {
                msg = fmt.Sprintf("{\"success\": \"Job %v deleted\"}", uuid)
        }


	return pstat, []byte(msg), nil
}

// removeJob the job with ID == uuid
func (s *httpScheduler) removeJob(uuid string) (httpStatusCode int) {
        var newJobList []*httpJob
        var foundJob = false

        //TODO: Must uptimize by just adjusting the job status
        //Jobs in newJoblist might have a status change before
        //being returned to list in the incorrect status.
        for _, v := range s.jobList {
                if v.ID() == uuid {
                        foundJob = true
                        continue
               }
                newJobList = append(newJobList, v)
        }
        // Handle 404 for Job not found
        if !foundJob {
                return http.StatusNotFound
        }

        // Update job list and return
        s.UpdateJobList(newJobList)
        return http.StatusOK
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

	s.mutex.Lock()
	s.schedule.SendIntervalSeconds = us.SendIntervalSeconds
	s.mutex.Unlock()

	msg := fmt.Sprintf("{\"Status\": \"Success\", \"New interval seconds\": %v}",s.schedule.SendIntervalSeconds)

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

	s.mutex.Lock()
	{
	s.schedule.ScheduleType = us.ScheduleType
	s.schedule.SendIntervalSeconds = us.SendIntervalSeconds
}
	s.mutex.Unlock()

	msg := fmt.Sprintf("{\"Status\": \"Success\", \"New interval seconds\": %v}",s.schedule.SendIntervalSeconds)
	return http.StatusCreated, []byte(msg), nil
}

func (s *httpScheduler) DeleteSchedule() (httpStatusCode int, jsonb []byte, err error) {

	//TODO: Create a seperate response channel
	s.schedulerDone <- true
	msg := fmt.Sprintf("{\"Status\": \"Success scheduler stopped\"}")
	return http.StatusOK, []byte(msg), nil
}

// SetChannels initializes channels the dispatcher has created inside
// of the scheduler
func (s *httpScheduler) SetChannels(j chan Job, r chan Result, b chan bool, i chan os.Signal) {
	//TODO: This should be getDispatchShedChans, so sheduler can call
        //and wait if channels are not ready.
	s.schedulerJobChan = j
	s.schedulerResponseChan = r
	s.schedulerDone = b
	s.schedulerInterrupt = i

	return
}

// SetConfigVariable changes the value of a given dispatcher field
func (s *httpScheduler) SetConfigVariable(name string, value int) (httpStatusCode int, msg []byte, err error) {

        fmt.Println("In scheduler SetConfigVariable")

        switch name {
        case trackMaxJobs:
                return s.doSetMaxCount(value)
        default:
                return doUnkownCmd(name)
        }
}

//doSetMaxCount: is used to set the maximum jobs anticipated per worker
//Initial test for sheduler setConfig

func (s *httpScheduler) doSetMaxCount(value int) (httpStatusCode int, msg []byte, err error) {

        if value < 0 {
                value = value * -1
        }
        if value > defaultMaxJobCount {
                value = defaultMaxJobCount
       }

        old := value
        s.mutex.Lock()
        {
                old = MaxJobCount
                MaxJobCount = value
        }
        s.mutex.Unlock()
        rmsg := fmt.Sprintf("{\" Status\": \"%s changed from %d to %d\"}",
                trackMaxJobs, old, value)
        return http.StatusOK, []byte(rmsg), nil
}

// Process methods
func (s *httpScheduler) Init(mo *managementCommands) error {

	var funcPtr interface{}
        var parmsList []string

        MaxJobCount = defaultInitMaxJobCount

	urlList := []string{
		"https://api.chucknorris.io/jokes/random",
		"https://swapi.co/api/people/1/",
		"https://swapi.co/api/people/2/",
		"https://swapi.co/api/people/3/"}

	s.metrics.Counters = make(map[string]int)

	s.mutex = new(sync.Mutex)
	s.metrics.mutex = new(sync.Mutex)
	s.schedule.ResponseTimeJobs = defaultResponseTimeJobs

	s.schedule.SendIntervalSeconds = defaultConstantInterval
	s.schedule.ScheduleType = constantIntervaleScheduler

	//DOING: calls to set scheduler managementCommands
        parmsList = make([]string, 1, 1)
        parmsList[0] = "Job_UUID"

        funcPtr = s.StopContinuousJob

        mo.useCommand("stop_job", "Stop running of a continuous job.", resScheduler, parmsList, funcPtr)

        //DOING: calls to set scheduler managementFields
        funcPtr = s.SetConfigVariable
        mo.setField(trackMaxJobs, "int", resScheduler, funcPtr)

	for _, u := range urlList {
		newJob := httpJob{}
		pu, err := url.Parse(u)
		if err != nil {
			fmt.Println(err)
			//os.Exit(-1)
			continue
		}

		// Set type and ID and http.Client
                em := newJob.Init()
                if em != nil {
                        return em
                }


		newJob.JobURL = pu
                newJob.ClientID = getTestClient()
		newJob.SetContinuous()
		s.mutex.Lock()
		s.jobList = append(s.jobList, &newJob)
		s.mutex.Unlock()
	}

	return nil
}

func (s *httpScheduler) Run(mWg *sync.WaitGroup) error {
	mWg.Add(2)
	go s.RunScheduler(mWg)
	go s.RunResultsReader(mWg)

	return nil
}

/*
First in first out scheduler. Work
needed to optimize. Cannot lock
while ranging the list.
*/

//RunScheduler: start up the scheduler
func (s *httpScheduler) RunScheduler(mWg *sync.WaitGroup) error {
	var actSent,listLn int
        defer mWg.Done()
        //DONE: framework for sheduler management options in init
	//Which options fro list?
	s.MetricSetStartTime()
	for {
              actSent = 0
                fmt.Println("Scheduler Running..")
		s.metrics.mutex.Lock()
                if ChannelsReady {
			s.metrics.mutex.Unlock()
                        listLn = len(s.jobList)
                        if listLn != 0 {
                                s.MetricInc(schedulerIterations)
                                //Trying to block modifications to joblist
                                s.metrics.mutex.Lock()
                                for curIdx := 0; curIdx < listLn; curIdx++ {
                                        //Hold lock
                                        j := s.jobList[curIdx]
                                        sent := false
                                        if !j.GetSent() || j.IsContinuous() {
                                                sent = true
                                                s.jobList[curIdx].MarkSent()
                                        }
                                        //joblist modifications now allowed
                                        s.metrics.mutex.Unlock()

                                        if sent {
                                                //Note: Blocking
                                                j.MarkSent()
                                                s.schedulerJobChan <- j
                                                s.MetricInc(jobsSent)
                                                actSent++
                                        }
                                        //joblist could have been modified
                                        //but we can keep going
                                        //Set lock for next check on len(s.jobList)
                                        s.metrics.mutex.Lock()

                                }
                                //done with joblist
                                s.metrics.mutex.Unlock()
                                s.MetricSet(currentJobChannelCapacity, cap(s.schedulerJobChan))
                                s.MetricSet(currentJobChannelUtilization, len(s.schedulerJobChan))
                                //s.MetricSet(jobListSize, len(s.jobList))
                                s.MetricSet(jobListSize, actSent)
                        }
                } else {
			s.metrics.mutex.Unlock()
                        fmt.Println("Before worker wait sleep ..")
                        time.Sleep(5 * time.Second)
                }
                s.MetricUpdateUpTime()
                select {
                case <-s.schedulerDone:
                        return nil
                case <-s.schedulerInterrupt:
                        return nil
               default:
                        fmt.Println("Before scheduler sleep..")
                        time.Sleep(time.Duration(s.schedule.SendIntervalSeconds) * time.Second)
			break
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

func (s *httpScheduler) RunResultsReader(mWg *sync.WaitGroup) error {
	defer mWg.Done()
	jobTimes := make([]int, 0, s.schedule.ResponseTimeJobs)
	log.Println("Starting result reader")
	for {
                        select {
                        case currentResult := <-s.schedulerResponseChan:
                                {
                                        log.Println("Processing response")
                                        s.MetricInc(resultsReceived)
                                        s.MetricSet(currentResultChannelCapacit, cap(s.schedulerJobChan))
                                        s.MetricSet(currentResultChannelUtilization, len(s.schedulerJobChan))
                                        log.Printf("Processing response for job ID %v\n", currentResult.Job().ID())
                                       log.Println("Processing response after Printf")
                                        j := currentResult.Job()
                                        if j.(*httpJob).Stats.RequestTimedOut {
                                                s.MetricInc(numberOfJobTimedOut)
                                        }
                                        jt, avg := s.ComputeAverageResponseTime(jobTimes, int(j.(*httpJob).Stats.RequestTime))
                                        s.MetricSet(averageJobProcessingTime, avg)
                                        jobTimes = jt
                                        //DOING: take job off s.jobList
                                        if !j.(*httpJob).IsContinuous() {
                                                st := s.removeJob(j.(*httpJob).ID())
                                                if st == http.StatusNotFound {
                                                        log.Printf("could not remove non-continuos job ID %v, %v\n", currentResult.Job().ID(), "not found")
                                                }
                                        }
                                }
                        case <-s.schedulerDone:
                                return nil
                        case <-s.schedulerInterrupt:
                                return nil
                        default:
                                {
                                        fmt.Println("Result reader Before Sleep..")
                                        time.Sleep(1 * time.Second)
					break
                                }

                                //select
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
// Metrics return the computed sheduler metrics
func (s *httpScheduler) Metrics() []byte {
	jb, e := s.MetricToJSON()
	if e != nil {
		return nil
	}
	return jb
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
