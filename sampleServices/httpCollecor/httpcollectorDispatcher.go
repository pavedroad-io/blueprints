package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
	// NumberOfWorkers in the worker pool
	NumberOfWorkers int = 5

	// SizeOfJobChannel sets the buffers size for the jobs channel
	SizeOfJobChannel int = 5

	// SizeOfResultChannel sets the buffers size for the results channel
	SizeOfResultChannel int = 5

	// GracefullShutdown sets the number of seconds to wait for work to complete
	// durring shutdown
	GracefullShutdown int = 30

	// HardShutdown sets the number of seconds to wait for work to complete
	// durring a hard shutdown
	HardShutdown int = 0
)

// Management API
const (
	gracefulShutdownSeconds string = "graceful_shutdown_seconds"
	hardShutdownSeconds     string = "hard_shutdown_seconds"
	numberOfWorkers         string = "number_of_workers"
	schedulerChannelSize    string = "scheduler_channel_size"
	resultChannelSize       string = "result_channel_size"
)

// managementGetResponse List of available command and field options
//
// swagger:response managementGetResponse
type managementGetResponse struct {
	// The error message
	// in: body

	// Commands is a list of valid commands that can be executed
	Commands []mgtCommand `json:"commands"`
	// Fields is a list of fields that can be changed
	Fields []string `json:"fields"`
}

// get404Response Not found
//
// swagger:response get404Response
type get404Response struct {
	// The 404 error message
	// in: body

	Body struct {
		// Error message
		Error string `json:"error"`

		// UUID / key that was not found
		UUID string `json:"uuid"`
	}
}

// genericError
//
// swagger:response genericError
type genericError struct {
	// in: body
	// Error message
	Body struct {
		Error string `json:"error"`
	} `json:"body"`
}

// genericResponse
//
// swagger:response genericResponse
type genericResponse struct {
	// in: body
	Body struct {
		// JSON body
		JSONBody string `json:"json_body"`
	} `json:"body"`
}

// metricsResponse
//
// swagger:response metricsResponse
type metricsResponse struct {
	// in: body
	Body struct {
		// Error message
		SchedulerMetrics string `json:"scheduler_metrics"`
		DispatherMetrics string `json:"dispather_metrics"`
	} `json:"body"`
}

// mgtCommand List of available command and field options
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

	// Commands is a list of valid commands that can be executed
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
	wg           *sync.WaitGroup
	jobChan      chan Job
	responseChan chan Result
	interrupt    chan os.Signal
	done         chan bool
}

// Run starts listing for jobs to be processed
// Read a job from the Job channel
// Read the done channel to see if this worker should exit
// Read the interrupt channel to see if this worker must exit
func (w *worker) Run() error {

	for {
		select {
		case currentJob := <-w.jobChan:
			r, e := currentJob.Run()

			if e != nil {
				log.Printf("Job: %v error: %v\n", currentJob.ID(), e.Error())
			}
			w.responseChan <- r
			w.lastJob = currentJob

		case done := <-w.done:
			if done {
				w.wg.Done()
				return nil
			}

		case <-w.interrupt:
			w.wg.Done()
			return nil
		}
	}
}

// dispatcherConfiguration options set during Initialize
type dispatcherConfiguration struct {
	scheduler              Scheduler
	numberOfWorkers        int // target number of workers
	currentNumberOfWorkers int // current number of workers
	// this may be different during resizing
	sizeOfJobChannel    int
	sizeOfResultChannel int
	gracefulShutdown    int
	hardShutdown        int
}

// SetSane Verify and set sane configuration options
// 	if not defined or exit if an option is mandatory
func (dc *dispatcherConfiguration) SetSane(d *dispatcher) {
	if dc.scheduler == nil {
		fmt.Println("A scheduler is required")
		os.Exit(-1)
	}
	d.scheduler = dc.scheduler

	if d.conf.numberOfWorkers == 0 {
		d.conf.numberOfWorkers = NumberOfWorkers
	}

	if d.conf.sizeOfJobChannel == 0 {
		d.conf.sizeOfJobChannel = SizeOfJobChannel
	}

	if d.conf.sizeOfResultChannel == 0 {
		d.conf.sizeOfResultChannel = SizeOfResultChannel
	}

	if d.conf.gracefulShutdown == 0 {
		d.conf.gracefulShutdown = GracefullShutdown
	}

	if d.conf.hardShutdown == 0 {
		d.conf.hardShutdown = HardShutdown
	}
}

// dispatcher structure
type dispatcher struct {
	conf *dispatcherConfiguration

	// Points to client implemented scheduler
	// that we read Jobs from
	scheduler           Scheduler      // Pointer to the scheduler
	schedulerJobChan    chan Job       // Channel to read jobs from
	schedulerResultChan chan Result    // Channel to write result to
	schedulerDone       chan bool      // Shutdown initiated by applicatoin
	schedulerInterrupt  chan os.Signal // Shutdown initiated by OS

	// Workers config
	wg      *sync.WaitGroup
	workers []*worker

	// Worker Channels
	workerJobChan     chan Job
	workerJobResponse chan Result
	workerDone        chan bool
	workerInterrupt   chan os.Signal

	// Management response
	managementOptions managementGetResponse

	// Metrics counters
	metrics dispatcherMetrics

	mux *sync.Mutex
}

type dispatcherMetrics struct {
	StartTime time.Time      `json:"start_time"`
	UpTime    time.Duration  `json:"up_time"`
	Counters  map[string]int `json:"counters"`
	mux       *sync.Mutex
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
//	 Job channels send or wait for jobs to execute
//	 Done channels allow go routines to be stopped by application
//	 logic
//   Interrupt channels handle OS interrupts
func (d *dispatcher) Init(dc *dispatcherConfiguration) {
	d.conf = dc
	dc.SetSane(d)

	d.metrics.Counters = make(map[string]int)
	d.metrics.mux = &sync.Mutex{}
	d.mux = &sync.Mutex{}
	d.wg = &sync.WaitGroup{}


	// Scheduler channels
	d.schedulerJobChan = make(chan Job, d.conf.sizeOfJobChannel)
	d.schedulerResultChan = make(chan Result, d.conf.sizeOfResultChannel)

	// Worker Channels
	d.workerJobChan = make(chan Job, d.conf.numberOfWorkers)
	d.workerJobResponse = make(chan Result, d.conf.numberOfWorkers)

	// Done channels
	d.schedulerDone = make(chan bool)
	d.workerDone = make(chan bool)

	// Interrupt channels
	d.schedulerInterrupt = make(chan os.Signal)
	d.workerInterrupt = make(chan os.Signal)

	// Set scheduler internal channels
	d.scheduler.SetChannels(
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
	e := d.createWorkerPool()
	if e == nil {
                go d.Forwarder()
                go d.Responder()
                return nil
        }

	return e
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
	log.Println("Dispatcher result channel started worker result -> scheduler  result")
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

// SetConfigVariable changegs the value of a given field
func (d *dispatcher) SetConfigVariable(name string, value int) (msg []byte, err error) {
	var rmsg string

	switch name {
	case gracefulShutdownSeconds:
		d.mux.Lock()
		old := d.conf.gracefulShutdown
		d.conf.gracefulShutdown = value
		d.mux.Unlock()
		rmsg = fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
			name, old, value)
		return []byte(rmsg), nil

	case hardShutdownSeconds:
		d.mux.Lock()
		old := d.conf.hardShutdown
		d.conf.hardShutdown = value
		d.mux.Unlock()
		rmsg = fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
			name, old, value)
		return []byte(rmsg), nil

	case numberOfWorkers:
		d.mux.Lock()
		old := d.conf.numberOfWorkers
		d.conf.numberOfWorkers = value
		d.mux.Unlock()

		//TODO: grow or srink as necessary
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

// ProcessManagementRequest executes a management command
func (d *dispatcher) ProcessManagementRequest(r managementRequest) (httpStatusCode int, jsonb []byte, err error) {

	switch r.Command {
	case "set":
		msg, e := d.SetConfigVariable(r.Field, r.Value)

		if e != nil {
			return http.StatusBadRequest, []byte(msg), nil
		}
		return http.StatusOK, []byte(msg), nil

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
		d.conf.currentNumberOfWorkers = 0
		msg := fmt.Sprintf("{\"Status\": \"Worker stop initiated\"}")
		return http.StatusOK, []byte(msg), nil

	case "start_workers":
		if d.conf.currentNumberOfWorkers > 0 {
			msg := fmt.Sprintf("{\"Status\": \"%v already running\"}",
				d.conf.currentNumberOfWorkers)
			return http.StatusBadRequest, []byte(msg), nil
		}

		e := d.createWorkerPool()
		 if e != nil {
                        msg := fmt.Sprintf("{\"status\": \"couldn't start worker pool: %v\"}", e)
                        return http.StatusExpectationFailed, []byte(msg), nil
                }

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
		d.workerInterrupt <- syscall.SIGINT
		msg := fmt.Sprintf("{\"Status\": \"shutdown_now complete\"}")
		return http.StatusOK, []byte(msg), nil

	default:
		msg := fmt.Sprintf("{\"Status\": \"Command %s not implemented\"}",
			r.Command)
		return http.StatusBadRequest, []byte(msg), nil
	}

}

func (d *dispatcher) createWorkerPool() error {
	for i := 0; i < d.conf.numberOfWorkers; i++ {
		newWorker := worker{wg: d.wg,
			jobChan:      d.workerJobChan,
			responseChan: d.workerJobResponse,
			interrupt:    d.workerInterrupt,
			done:         d.workerDone}

		d.wg.Add(1)

		// Keep track of each worker
		d.workers = append(d.workers, &newWorker)

		go newWorker.Run()
	}
	d.conf.currentNumberOfWorkers = d.conf.numberOfWorkers
	log.Printf("Worker pool created, %d workers\n", d.conf.numberOfWorkers)
	return nil
}

func (d *dispatcher) growWorkerPool() error {
	return nil
}

func (d *dispatcher) srinkWorkerPool() error {
	return nil
}