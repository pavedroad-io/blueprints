{{define "httpScheduler.go"}}{{.PavedroadInfo}}
package main

import (
  "encoding/json"
  "fmt"
  "net/url"
  "os"
	"sync"
  "time"
)

const (
  schedulerIterations = iota
  jobsSent
  jobListSize
  resultsRecieved
  currentJobChannelUtilization
  currentJobChannelCapacit
  currentResultChannelUtilization
  currentResultChannelCapacit
  numberOfJobTimedOut
  averageJobProcessingTime
)

type httpScheduler struct {
	jobList             []*httpJob
  sendIntervalSeconds int
  schedulerJobChan      chan Job       // Channel to read jobs from
  schedulerResponseChan chan Result    // Channel to write repose to
  schedulerDone         chan bool      // Shudown initiated by applicatoin
  schedulerInterrup     chan os.Signal // Shutdown initiated by OS
	metrics               httpSchedulerMetrics
}

// httpSchedulerMetrics hold metrics about the Scheduler, Jobs, and Results
// We export attributes we want included in the JSON output
type httpSchedulerMetrics struct {
  StartTime time.Time     `json:"start_time"`
  UpTime    time.Duration `json:"up_time"`
  Counters  map[int]int   `json:"counters"`
  mux       sync.Mutex    `json:"mux"`
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

func (s *httpScheduler) MetricSetStartTime(t time.Time) {
  s.metrics.mux.Lock()
  s.metrics.StartTime = t
  s.metrics.mux.Unlock()
}

func (s *httpScheduler) MetricUpdateUpTime() (uptime time.Duration) {
  s.metrics.mux.Lock()
  ct := time.Now()
  s.metrics.UpTime = ct.Sub(s.metrics.StartTime)
  s.metrics.mux.Unlock()
  return s.metrics.UpTime
}

func (s *httpScheduler) MetricInc(key int) {
  s.metrics.mux.Lock()
  s.metrics.Counters[key]++
  s.metrics.mux.Unlock()
}

func (s *httpScheduler) MetricSet(key, value int) {
  s.metrics.mux.Lock()
  s.metrics.Counters[key]++
  s.metrics.mux.Unlock()
}

func (s *httpScheduler) MetricValue(key int) int {
  s.metrics.mux.Lock()
  defer s.metrics.mux.Unlock()
  return s.metrics.Counters[key]
}

// Required methods

// Object methods
func (s *httpScheduler) GetSchedule() (Scheduler, error) {

	return nil, nil
}

func (s *httpScheduler) UpdateSchedule() (Scheduler, error) {

	return nil, nil
}

func (s *httpScheduler) CreateSchedule() (Scheduler, error) {

	return nil, nil
}

func (s *httpScheduler) DeleteSchedule() (Scheduler, error) {

	return nil, nil
}

// SetChannels initializes channels the dispatcher has created inside
// of the scheduler
func (s *httpScheduler) SetChannels(j chan Job, r chan Result, b chan bool, i chan os.Signal) {
  s.schedulerJobChan = j
  s.schedulerResponseChan = r
  s.schedulerDone = b
  s.schedulerInterrup = i

  return
}

// Process methods
func (s *httpScheduler) Init() error {
  urlList := []string{
    "https://api.chucknorris.io/jokes/random",
    "https://swapi.co/api/people/1/",
    "https://swapi.co/api/people/2/",
    "https://swapi.co/api/people/3/"}

  s.metrics.Counters = make(map[int]int)
  s.sendIntervalSeconds = 10

  for _, u := range urlList {
    newJob := httpJob{}
    pu, err := url.Parse(u)
    if err != nil {
      fmt.Println(err)
      os.Exit(-1)
    }
    newJob.jobURL = pu

		// Set type and ID and http.Client
    newJob.Init()
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
  fmt.Println("Running scheduler")
	s.MetricSetStartTime(time.Now())
  for {
		s.MetricInc(schedulerIterations)
    for _, j := range s.jobList {
      s.schedulerJobChan <- j
		  s.MetricInc(jobsSent)
      s.MetricSet(currentJobChannelCapacit, cap(s.schedulerJobChan))
      s.MetricSet(currentJobChannelUtilization, len(s.schedulerJobChan))
      s.MetricSet(jobListSize, len(s.jobList))
    }
    time.Sleep(time.Duration(s.sendIntervalSeconds) * time.Second)
  }

  return nil
}

// ComputeAverageResponseTime Keep track of the last N responses
func (s *httpScheduler) ComputeAverageResponseTime(jt []int, newTime int) ([]int, int) {
  currentLength := len(jt)
  desiredLength := currentLength - 9
  if currentLength >= 10 {
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
	  jobTimes := make([]int, 0, 10)
  fmt.Println("Reading job results")
  for {
    select {
    case currentResult := <-s.schedulerResponseChan:
      s.MetricInc(resultsRecieved)
      s.MetricSet(currentResultChannelCapacit, cap(s.schedulerJobChan))
      s.MetricSet(currentResultChannelUtilization, len(s.schedulerJobChan))
      fmt.Printf("Processing response for job ID %v\n", currentResult.Job().ID())
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
      } else {
        fmt.Println("Bad response on scheduler Done channel")
      }

    case <-s.schedulerInterrup:
      return nil
    }
  }

  return nil
}


func (s *httpScheduler) Pause() []byte {

	return nil
}

func (s *httpScheduler) Shutdown() error {

	return nil
}

// Status methods
func (s *httpScheduler) Metrics() []byte {
  jb, e := s.MetricToJSON()
  if e != nil {
    return nil
  }
  return jb
}{{end}}
