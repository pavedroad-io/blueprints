{{define "templateDispatcher.go"}}package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"syscall"
	"time"
)

// Map counters to JSON friendly names
const (
  dispatcherJobsSent        = "jobs_sent"
  dispatcherResultsReceived = "results_received"
)

// Defaults for dispatcher and workers
const (
	NUMBEROFWORKERS     int = 5
	SIZEOFJOBCHANNEL    int = 5
	SIZEOFRESULTCHANNEL int = 5
	GRACEFULLSHUTDOWN   int = 30
	HARDSHUTDOWN        int = 0
)

// Management API
const (
	gracefulShutdownSeconds string = "graceful_shutdown_seconds"
	hardShutdownSeconds     string = "hard_shutdown_seconds"
	numberOfWorkers         string = "number_of_workers"
	schedulerChannelSize    string = "scheduler_channel_size"
	resultChannelSize       string = "result_channel_size"
)

// managementGetResponse List of valaible command and field options
//
// swagger:response managementGetResponse
type managementGetResponse struct {
	// The error message
	// in: body

	// Commands is a list of valide commands that can be executed
	Commands []mgtCommand `json:"commands"`
	// Fields is a list of fields that can be changed
	Fields []string `json:"fields"`
}

// mgtCommand List of valaible command and field options
//
type mgtCommand struct {
	// Name of the command
	Name string `json:"name"`

	// Go data type
	DataType string `json:"data_type"`

	// Go data type
	CommandType string `json:"command_type"`

	// Description of what this command does
	Description string `json:"description"`
}

// managementRequest user request to execute a management command
//
// swagger:response managementRequest
type managementRequest struct {
	// The error message
	// in: body

	// Commands is a list of valide commands that can be executed
	Command string `json:"command"`

	// Field to set
	Field string `json:"field"`

	// Value for field
	Value int `json:"field_value"`
}

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
	schedulerResultChan chan Result     // Channel to write result to
	schedulerDone         chan bool       // Shudown initiated by applicatoin
	schedulerInterrupt     chan os.Signal // Shutdown initiated by OS
	schedulerInterrupt  chan os.Signal // Shutdown initiated by OS

	// Workers config
	wg             sync.WaitGroup
	desiredWorkers int
	currentWorkers int
	workers        []*worker

	// Worker Channels
	workerJobChan     chan Job
	workerJobResponse chan Result
	workerDone        chan bool
	workerInterrup    chan os.Signal

	// Management response
	managementOptions managementGetResponse

	// Metrics counters
  metrics dispatcherMetrics

	// Shutdown options
	gracefulShutdown int //Seconds to wait
	hardShutdown     int
	mux              sync.Mutex
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
	d.gracefulShutdown = GRACEFULLSHUTDOWN
	d.hardShutdown = HARDSHUTDOWN

	if numWorkers == 0 {
		d.desiredWorkers = NUMBEROFWORKERS
	} else {
		d.desiredWorkers = numWorkers
	}

	if chanCapacity == 0 {
		d.schedulerCapacity = SIZEOFJOBCHANNEL
		d.resultsCapacity = SIZEOFRESULTCHANNEL
	} else {
		d.schedulerCapacity = chanCapacity
		d.resultsCapacity = chanCapacity
	}

	d.metrics.Counters = make(map[string]int)
	d.desiredWorkers = numWorkers
	d.schedulerCapacity = chanCapacity

	// Job channels
	// Scheduler
	d.schedulerJobChan = make(chan Job, d.schedulerCapacity)
	d.schedulerResultChan = make(chan Result, d.schedulerCapacity)

	// Worker
	d.workerJobChan = make(chan Job, d.desiredWorkers)
	d.workerJobResponse = make(chan Result, d.desiredWorkers)

	// Done channels
	d.schedulerDone = make(chan bool)
	d.workerDone = make(chan bool)

	// Interrupt channels
	d.schedulerInterrupt = make(chan os.Signal)
	d.workerInterrup = make(chan os.Signal)

  // Set scheduler internal channels
	s.SetChannels(
    d.schedulerJobChan,
    d.schedulerResultChan,
    d.schedulerDone,
    d.schedulerInterrupt)

	d.managementInit()
	d.MetricSetStartTime()

	return
}

// managementInit() setup management commands and fields
func (d *dispatcher) managementInit() {

	newCMD := mgtCommand{Name: "set", DataType: "int", CommandType: "config",
		Description: "Sets the value of a configurable field, see fields below"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	newCMD = mgtCommand{Name: "stop_scheduler", DataType: "string",
		CommandType: "command",
		Description: "Stops the scheduler from send new jobs"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	newCMD = mgtCommand{Name: "start_scheduler", DataType: "string",
		CommandType: "command",
		Description: "Starts the scheduler running again.  If running has no affect"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	newCMD = mgtCommand{Name: "stop_workers", DataType: "string",
		CommandType: "command",
		Description: "Shutdown the worker pool letting jobs inflight complete"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	newCMD = mgtCommand{Name: "start_workers", DataType: "string",
		CommandType: "command",
		Description: "Starts the worker pool if stopped"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	newCMD = mgtCommand{Name: "shutdown", DataType: "string",
		CommandType: "command",
		Description: "Graceful shutdown"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	newCMD = mgtCommand{Name: "shutdown_now", DataType: "string",
		CommandType: "command",
		Description: "Hard shutdown with SIGKILL"}
	d.managementOptions.Commands = append(d.managementOptions.Commands, newCMD)

	d.managementOptions.Fields = append(d.managementOptions.Fields,
		gracefulShutdownSeconds,
		hardShutdownSeconds,
		numberOfWorkers,
		schedulerChannelSize,
		resultChannelSize)

	/* TODO: add hooks to allows Job and Scheduler to extend management API
	d.managementOptions.Commands = append(d.managementOption.Command, s.AddSchedulerCommands())
	d.managementOptions.Fields = append(d.managementOption.Fields, s.AddSchedulerFields())

	d.managementOptions.Commands = append(d.managementOption.Command, s.AddJobCommands())
	d.managementOptions.Fields = append(d.managementOption.Fields, s.AddJobFields())
	*/
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
      d.schedulerResultChan <- currentJobResponse
    }
  }
}

func (d *dispatcher) Shutdown() error {
	return nil
}

// SetConfigVariable
func (d *dispatcher) SetConfigVariable(name string, value int) (msg []byte, err error) {
	var rmsg string

	switch name {
	case gracefulShutdownSeconds:
		d.mux.Lock()
		old := d.gracefulShutdown
		d.gracefulShutdown = value
		d.mux.Unlock()
		rmsg = fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
			name, old, value)
		return []byte(rmsg), nil

	case hardShutdownSeconds:
		d.mux.Lock()
		old := d.hardShutdown
		d.hardShutdown = value
		d.mux.Unlock()
		rmsg = fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
			name, old, value)
		return []byte(rmsg), nil

	case numberOfWorkers:
		d.mux.Lock()
		old := d.desiredWorkers
		d.desiredWorkers = value
		d.mux.Unlock()

		//TODO: grow or srinkt as necessary
		rmsg = fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
			name, old, value)
		return []byte(rmsg), nil

	case schedulerChannelSize:
		rmsg = fmt.Sprintf("{\"Status\": \"%s not implemented\"}", name)
		e := errors.New("not implemented")
		return []byte(rmsg), e

	case resultChannelSize:
		rmsg = fmt.Sprintf("{\"Status\": \"%s not implemented\"}", name)
		e := errors.New("not implemented")
		return []byte(rmsg), e
	default:
		rmsg = fmt.Sprintf("{\"Status\": \"%s unknown\"}", name)
		e := errors.New("unknown set field")
		return []byte(rmsg), e
	}

}

func (d *dispatcher) ProcessManagementRequest(r managementRequest) (httpStatusCode int, jsonb []byte, err error) {

	switch r.Command {
	case "set":
		msg, e := d.SetConfigVariable(r.Field, r.Value)

		if e != nil {
			return http.StatusBadRequest, []byte(msg), nil
		} else {
			return http.StatusOK, []byte(msg), nil
		}

	case "stop_scheduler":
		d.schedulerDone <- true
		msg := fmt.Sprintf("{\"Status\": \"Scheduler stop initiated\"}")
		return http.StatusOK, []byte(msg), nil

	case "start_scheduler":
		// TODO: this will require changes to the scheduler interface
		msg := fmt.Sprintf("{\"Status\": \"Scheduler start initiated\"}")
		return http.StatusOK, []byte(msg), nil

	case "stop_workers":
		d.workerDone <- true
		msg := fmt.Sprintf("{\"Status\": \"Worker stop initiated\"}")
		return http.StatusOK, []byte(msg), nil

	case "start_workers":
		// TODO: make sure it isn't running first
		d.createWorkerPool()
		msg := fmt.Sprintf("{\"Status\": \"Worker stop initiated\"}")
		return http.StatusOK, []byte(msg), nil

		//TODO: move this logic into Shutdown() method
	case "shutdown":
		// Let the scheduler clean up if special logic is needed
		// d.scheduler.Shutdown()
		d.schedulerDone <- true
		d.workerDone <- true
		msg := fmt.Sprintf("{\"Status\": \"Shutdown complete\"}")
		return http.StatusOK, []byte(msg), nil

	case "shutdown_now":
		d.schedulerInterrupt <- syscall.SIGINT
		// TODO: fix this spelling
		d.workerInterrup <- syscall.SIGINT
		msg := fmt.Sprintf("{\"Status\": \"shutdown_now complete\"}")
		return http.StatusOK, []byte(msg), nil

	default:
		msg := fmt.Sprintf("{\"Status\": \"Command: %s not implemeted\"}",
			r.Command)
		return http.StatusOK, []byte(msg), nil
	}

	return 0, nil, nil
}

// TODO: keep a list of points to workes so we can call Shutdown()
//       Errors(), etc
func (d *dispatcher) createWorkerPool() error {
	for i := 0; i < d.desiredWorkers; i++ {
		newWorker := worker{wg: d.wg,
			jobChan:     d.workerJobChan,
			responseChan: d.workerJobResponse,
			interrup:    d.workerInterrup,
			done:        d.workerDone}

		d.wg.Add(1)
		d.workers = append(d.workers, &newWorker)
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
