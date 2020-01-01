{{define "httpScheduler.go"}}{{.PavedroadInfo}}
package main

import (
  "fmt"
  "net/url"
  "os"
  "time"
)

type httpScheduler struct {
	jobList             []*httpJob
  sendIntervalSeconds int
  schedulerJobChan      chan Job       // Channel to read jobs from
  schedulerResponseChan chan Result    // Channel to write repose to
  schedulerDone         chan bool      // Shudown initiated by applicatoin
  schedulerInterrup     chan os.Signal // Shutdown initiated by OS
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

  s.sendIntervalSeconds = 5

  for _, u := range urlList {
    newJob := httpJob{}
    pu, err := url.Parse(u)
    if err != nil {
      fmt.Println(err)
      os.Exit(-1)
    }
    newJob.jobURL = pu
    newJob.Init()
    s.jobList = append(s.jobList, &newJob)
  }

	return nil
}

func (s *httpScheduler) Run() error {
  fmt.Println("Running scheduler")
  for {
    for _, j := range s.jobList {
      fmt.Printf("%+v\n", j.jobURL.String())
      //s.JobChanle <- j
    }
    time.Sleep(time.Duration(s.sendIntervalSeconds) * time.Second)
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
// TODO: Make this an interface
func (s *httpScheduler) Metrics() []byte {

	return nil
}{{end}}
