{{define "templateDispatcher.go"}}package main

import (
	"fmt"
	"os"
	"sync"
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
}

// Init sets size of work and scheduler channels and then creates them
//   Job channels send or wait for jobs to execute
//   Done channels allow go routines to be stopped by application
//   logic
//   Inter channels handle OS interrups
func (d *dispatcher) Init(numWorkers, chanCapacity int, s Scheduler) {
	fmt.Printf("Initialize dispacther with %v workers and channel capacity of %v\n", numWorkers, chanCapacity)
	d.desiredWorkers = numWorkers
	d.schedulerCapacity = chanCapacity

	// Job channels
	d.schedulerJobChan = make(chan Job, d.schedulerCapacity)
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
      fmt.Println(currentJob.ID())
      d.workerJobChan <- currentJob
    }
  }
}

func (d dispatcher) Responder() {
  for {
    select {
    case currentJob := <-d.workerJobResponse:
      //fmt.Println(currentJob.ID())
      d.schedulerResponseChan <- currentJob
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
