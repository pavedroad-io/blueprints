{{define "templateDispatcher.go"}}package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// Map counters to JSON friendly names
const (
  dispatcherJobsSent        = "jobs_sent"
  dispatcherResultsReceived = "results_received"
)

// worker is a go worker pool pattern
type worker struct {
	currentJob   Job
	lastJob      Job
	wg           sync.WaitGroup
	jobChan      chan Job
	responseChan chan Result
	interrup     chan os.Signal
	done         chan bool
}

// Run starts listing for jobs to be processed
// Read a job from the Job channel
// Read the done channel to see if this worker should exit
// Read the interrup channel to see if this worker must exit
func (w *worker) Run() error {

	for {
		select {
		case currentJob := <-w.jobChan:
			r, e := currentJob.Run()

      if e != nil {
        // TODO: what do we want to do with an error?
        fmt.Println(e)
      }
      w.responseChan <- r
			w.lastJob = currentJob

		case done := <-w.done:
			if done {
        w.wg.Done()
        return nil
      }

		case <-w.interrup:
			w.wg.Done()
			return nil
		}
	}
	return nil
}

type dispatcher struct {
	// Points to client implemented scheduler
	// that we read Jobs from
	scheduler             *Scheduler      // Pointer to the scheduler
	schedulerCapacity     int             // buffer size of channel
	schedulerJobChan      chan Job        // Channel to read jobs from
	// TODO: Change name to schedulerResultChan
	schedulerResponseChan chan Result     // Channel to write result to
	schedulerDone         chan bool       // Shudown initiated by applicatoin
	// TODO: Fix interrupt spelling
	schedulerInterrup     chan os.Signal // Shutdown initiated by OS

	// Workers config
	wg             sync.WaitGroup
	desiredWorkers int
	currentWorkers int

	// Worker Channels
	workerJobChan     chan Job
	workerJobResponse chan Result
	workerDone        chan bool
	workerInterrup    chan os.Signal

	// Metrics counters
  metrics dispatcherMetrics
}

type dispatcherMetrics struct {
  StartTime time.Time      `json:"start_time"`
  UpTime    time.Duration  `json:"up_time"`
  Counters  map[string]int `json:"counters"`
  mux       sync.Mutex     `json:"mux"`
}

func (d *dispatcher) MetricToJSON() ([]byte, error) {
  d.metrics.mux.Lock()
  defer d.metrics.mux.Unlock()
  jb, e := json.Marshal(d.metrics)
  if e != nil {
    fmt.Println(e)
    return nil, e
  }
  return jb, nil
}

func (d *dispatcher) MetricSetStartTime() {
  d.metrics.mux.Lock()
  d.metrics.StartTime = time.Now()
  ct := time.Now()
  d.metrics.UpTime = ct.Sub(d.metrics.StartTime)
  d.metrics.mux.Unlock()
}

func (d *dispatcher) MetricUpdateUpTime() (uptime time.Duration) {
  d.metrics.mux.Lock()
  ct := time.Now()
  d.metrics.UpTime = ct.Sub(d.metrics.StartTime)
  d.metrics.mux.Unlock()
  return d.metrics.UpTime
}

func (d *dispatcher) MetricInc(key string) {
  d.metrics.mux.Lock()
  d.metrics.Counters[key]++
  d.metrics.mux.Unlock()
}

func (d *dispatcher) MetricSet(key string, value int) {
  d.metrics.mux.Lock()
  d.metrics.Counters[key] = value
  d.metrics.mux.Unlock()
}

func (d *dispatcher) MetricValue(key string) int {
  d.metrics.mux.Lock()
  defer d.metrics.mux.Unlock()
  return d.metrics.Counters[key]
}

// Init sets size of work and scheduler channels and then creates them
//   Job channels send or wait for jobs to execute
//   Done channels allow go routines to be stopped by application
//   logic
//   Inter channels handle OS interrups
func (d *dispatcher) Init(numWorkers, chanCapacity int, s Scheduler) {
	fmt.Printf("Initialize dispacther with %v workers and channel capacity of %v\n", numWorkers, chanCapacity)
	d.metrics.Counters = make(map[string]int)
	d.desiredWorkers = numWorkers
	d.schedulerCapacity = chanCapacity

	// Job channels
	// Scheduler
	d.schedulerJobChan = make(chan Job, d.schedulerCapacity)
	d.schedulerResponseChan = make(chan Result, d.schedulerCapacity)

	// Worker
	d.workerJobChan = make(chan Job, d.desiredWorkers)
	d.workerJobResponse = make(chan Result, d.desiredWorkers)

	// Done channels
	d.schedulerDone = make(chan bool)
	d.workerDone = make(chan bool)

	// Interrupt channels
	d.schedulerInterrup = make(chan os.Signal)
	d.workerInterrup = make(chan os.Signal)

  // Set scheduler internal channels
	s.SetChannels(
    d.schedulerJobChan,
    d.schedulerResponseChan,
    d.schedulerDone,
    d.schedulerInterrup)

	d.MetricSetStartTime()

	return
}

func (d dispatcher) Run() error {
	d.createWorkerPool()
	go d.Forwarder()
  go d.Responder()

	//wg.Wait()
	return nil
}

func (d dispatcher) Forwarder() {
  for {
    select {
    case currentJob := <-d.schedulerJobChan:
      //fmt.Println(currentJob.ID())
			d.MetricInc(dispatcherJobsSent)
      d.workerJobChan <- currentJob
    }
  }
}

func (d dispatcher) Responder() {
  for {
    select {
		case currentJobResponse := <-d.workerJobResponse:
			d.MetricInc(dispatcherResultsReceived)
      d.schedulerResponseChan <- currentJobResponse
    }
  }
}

func (d *dispatcher) Shutdown() error {
	return nil
}

func (d *dispatcher) createWorkerPool() error {
	for i := 0; i < d.desiredWorkers; i++ {
		newWorker := worker{wg: d.wg,
			jobChan:     d.workerJobChan,
			responseChan: d.workerJobResponse,
			interrup:    d.workerInterrup,
			done:        d.workerDone}

		d.wg.Add(1)
		go newWorker.Run()
	}
	return nil
}

func (d *dispatcher) growWorkerPool() error {
	return nil
}

func (d *dispatcher) srinkWorkerPool() error {
	return nil
}{{end}}
