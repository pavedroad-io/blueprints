{{define "templateDispatcher.go"}}package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"reflect"
	"sync"
	"syscall"
	"time"
)

// Map counters to JSON friendly names
const (
	dispatcherJobsSent        string = "jobs_sent"
	dispatcherResultsReceived string = "results_received"
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
	HardShutdown int = 1

	//Maximum number of jobs to track by each worker
        MaxJobTrack int = 100

        //Maximum number of worker channels
        MaxWorkers int = 10

	 //Maximum number of seconds to wait for hard shutdown
        MaxShutSec = 60

        //Maximum number of seconds for gracefull sutdown
        MaxGraceSec int = 90
)

// Management API
const (
	gracefulShutdownSeconds string = "graceful_shutdown_seconds"
	hardShutdownSeconds     string = "hard_shutdown_seconds"
	numberOfWorkers         string = "number_of_workers"
	schedulerChannelSize    string = "scheduler_channel_size"
	resultChannelSize       string = "result_channel_size"
)

func (d *dispatcher) stopScheduler() (httpStatusCode int, jsonb []byte, err error) {
	d.metrics.mutex.Lock()
        if ChannelsReady {
		d.metrics.mutex.Unlock()
                Mwg.Add(1)
                go func() {
                        //TODO: Must update for interupt to return
                        defer Mwg.Done()
                        d.schedulerDone <- true
                }()
	        d.metrics.mutex.Lock()
                ChannelsReady = false
		d.metrics.mutex.Unlock()
                msg := fmt.Sprintf("{\"status\": \"scheduler stop initiated\"}")
                return http.StatusOK, []byte(msg), nil
        }
	d.metrics.mutex.Unlock()
        msg := fmt.Sprintf("{\"status\": \"must wait to stop scheduler\"}")
        return http.StatusForbidden, []byte(msg), nil
}

func (d *dispatcher) startScheduler() (httpStatusCode int, jsonb []byte, err error) {

        return doNotImplemented(resDispatcher)
}

func (d *dispatcher) stopWorkers() (httpStatusCode int, jsonb []byte, err error) {
	d.metrics.mutex.Lock()
        if !ChannelsReady {
                d.metrics.mutex.Unlock()
                msg := fmt.Sprintf("{\"status\": \"can't stop workers now\"}")
                return http.StatusForbidden, []byte(msg), nil
        }
        fmt.Println("in stoping wrokers")
        ChannelsReady = false
	d.metrics.mutex.Unlock()
        fmt.Println("in stoping wrokers, step 2")
        Mwg.Add(len(d.workers))
        for i := 0; i < len(d.workers); i++ {
                go func(i int) {
                        defer Mwg.Done()
                        //      d.workerDone <- true
                        //TODO: Must update for interupt for return
                        fmt.Println("stopWorkers before chan:", i)
                        d.workers[i].done <- true
                        fmt.Println("stopWorkers done:", i)
                }(i)
        }
       //TODO: More work required for the worker pool after stop
        msg := fmt.Sprintf("{\"status\": \"worker stop initiated\"}")
        return http.StatusOK, []byte(msg), nil
}

func (d *dispatcher) startWorkers() (httpStatusCode int, jsonb []byte, err error) {
        time.Sleep(5 * time.Second)
	d.metrics.mutex.Lock()
        if ChannelsReady {
                //No prior worker stop requested, so return
                msg := fmt.Sprintf("{\"status\": \" workers state %v is not expecetd\"}", ChannelsReady)
		d.metrics.mutex.Unlock()
                return http.StatusForbidden, []byte(msg), nil
        }
        if (d.conf.currentNumberOfWorkers + WorkersDown) > 0 {
                fmt.Println("startWorkers still active:", (d.conf.currentNumberOfWorkers + WorkersDown))
                msg := fmt.Sprintf("{\"status\": \" %v workers are still active\"}", (d.conf.currentNumberOfWorkers + WorkersDown))
                return http.StatusForbidden, []byte(msg), nil
        } else {
               d.metrics.mutex.Lock() 
                {

                        d.workers = make([]*worker, 0, MaxWorkers)
                        d.conf.currentNumberOfWorkers = 0
                        WorkersDown = 0
                }
               d.metrics.mutex.Unlock()
        }
        fmt.Println("startWorkers, no more workers")
        //d.createWorkerPool()
        e := d.createWorkerPool(&Mwg)
        if e != nil {
                msg := fmt.Sprintf("{\"status\": \"couldn't start worker pool: %v\"}", e)
                return http.StatusExpectationFailed, []byte(msg), nil
        }
        fmt.Println("startWorkers initated")
	d.metrics.mutex.Lock()
        ChannelsReady = true
	d.metrics.mutex.Unlock()
        msg := fmt.Sprintf("{\"status\": \"worker start initiated\"}")
        return http.StatusOK, []byte(msg), nil
}

func (d *dispatcher) shutdown() (httpStatusCode int, jsonb []byte, err error) {

        //TODO: move this logic into api Shutdown() method ?
        // Let the scheduler clean up if special logic is needed
        // d.scheduler.Shutdown()
        d.schedulerDone <- true
        d.workerDone <- true
        msg := fmt.Sprintf("{\"Status\": \"Shutdown complete\"}")
        return http.StatusOK, []byte(msg), nil

}

func (d *dispatcher) shutdownNow() (httpStatusCode int, jsonb []byte, err error) {

        d.schedulerInterrupt <- syscall.SIGINT
        d.workerInterrupt <- syscall.SIGINT
        msg := fmt.Sprintf("{\"Status\": \"Shutdown_now complete.\"}")
        return http.StatusOK, []byte(msg), nil

}

// worker is a go worker pool pattern
type worker struct {
         jobHist      map[uuid.UUID]string //Only noncontinuous jobs, showing ClientID
	currentJob   Job
	lastJob      Job
	wg           sync.WaitGroup
	jobChan      chan Job
	responseChan chan Result
	interrupt    chan os.Signal
	done         chan bool
        jobCnt       int //Number of jobs sent to  this worker chan.
        workerID     int //For monitoring work distribution
}

// Run starts listing for jobs to be processed
// Read a job from the Job channel
// Read the done channel to see if this worker should exit
// Read the interrupt channel to see if this worker must exit
func (w *worker) Run(mWg *sync.WaitGroup) error {
      defer mWg.Done()
        for {
                //fmt.Println("Woker number:", w.workerID)
                //No need to look for jobs if !ChannelsReady
                //Trying to speed up worker stop
		 Mmutex.Lock()
                if ChannelsReady {
			Mmutex.Unlock()
                        select {
                        case currentJob := <-w.jobChan:
                                {
                                        mWg.Add(1)
                                        r, e := currentJob.Run(mWg)

                                        if e != nil {
                                                log.Printf("Job: %v error: %v\n", currentJob.ID(), e.Error())
                                        }
                                        w.responseChan <- r
                                        w.lastJob = currentJob
                                        if !currentJob.IsContinuous() {
                                                w.jobHist[currentJob.UUID()] = currentJob.GetClientID()
                                        }
                                        if w.jobCnt < MaxJobTrack {
                                                w.jobCnt++
                                        }
                                }
                        case <-w.done:
                                {
                                        //TODO: More work required to clean up workers
                                        Mmutex.Lock()
                                        WorkersDown = WorkersDown - 1
                                        Mmutex.Unlock()
                                        return nil
                                }
                       case <-w.interrupt:
                                return nil
                        default:
                                {
					Mmutex.Lock()
                                        if !ChannelsReady {
          				fmt.Println("Worker waiting on channels.")

                                        }
					Mmutex.Unlock()
                                }
                        }
                } else {
			Mmutex.Unlock()
                        select {
                       case <-w.done:
                                {
                                        Mmutex.Lock()
                                        WorkersDown = WorkersDown - 1
                                        Mmutex.Unlock()
                                        return nil
                                }
                        case <-w.interrupt:
                                return nil
                        default:
				{
				 fmt.Println("Worker waiting......")
                                 time.Sleep(5 * time.Second)
                                 Mmutex.Lock()
                                if ChannelsReady {
                                       //In case worker started first
				       fmt.Println("Channels ready for workers.")
                                }
				Mmutex.Unlock()
			}
                        }
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

// SetSame Verify and set sane configuration options
// 	if not defined or exit if an option is mandatory
func (dc *dispatcherConfiguration) SetSame(d *dispatcher) {
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
	wg      sync.WaitGroup
	workers []*worker

	// Worker Channels
	workerJobChan     chan Job
	workerJobResponse chan Result
	workerDone        chan bool
	workerInterrupt   chan os.Signal

	// Management response
	managementOptions managementCommands

	// Metrics counters
	metrics dispatcherMetrics

	mutex *sync.Mutex
}

type dispatcherMetrics struct {
	StartTime time.Time      `json:"start_time"`
	UpTime    time.Duration  `json:"up_time"`
	Counters  map[string]int `json:"counters"`
	mutex     *sync.Mutex
}

func (d *dispatcher) MetricToJSON() ([]byte, error) {
	d.metrics.mutex.Lock()
	defer d.metrics.mutex.Unlock()
	jb, e := json.Marshal(d.metrics)
	if e != nil {
		fmt.Println(e)
		return nil, e
	}
	return jb, nil
}

func (d *dispatcher) MetricSetStartTime() {
	d.metrics.mutex.Lock()
	{
	d.metrics.StartTime = time.Now()
	ct := time.Now()
	d.metrics.UpTime = ct.Sub(d.metrics.StartTime)
}
	d.metrics.mutex.Unlock()
}

func (d *dispatcher) MetricUpdateUpTime() (uptime time.Duration) {
        ct := time.Now()
	if d.metrics.mutex == nil {
                return ct.Sub(time.Now())
        }
        
	d.metrics.mutex.Lock()
	d.metrics.UpTime = ct.Sub(d.metrics.StartTime)
	d.metrics.mutex.Unlock()
	return d.metrics.UpTime
}

func (d *dispatcher) MetricInc(key string) {
	d.metrics.mutex.Lock()
	d.metrics.Counters[key]++
	d.metrics.mutex.Unlock()
}

func (d *dispatcher) MetricSet(key string, value int) {
	d.metrics.mutex.Lock()
	d.metrics.Counters[key] = value
	d.metrics.mutex.Unlock()
}

func (d *dispatcher) MetricValue(key string) int {
	d.metrics.mutex.Lock()
	defer d.metrics.mutex.Unlock()
	return d.metrics.Counters[key]
}

// Init sets size of work and scheduler channels and then creates them
//	 Job channels send or wait for jobs to execute
//	 Done channels allow go routines to be stopped by application
//	 logic
//   Interrupt channels handle OS interrupts
func (d *dispatcher) Init(dc *dispatcherConfiguration, mo *managementCommands) {
	d.conf = dc
	dc.SetSame(d)

	d.metrics.Counters = make(map[string]int)
	d.workers = make([]*worker, 0, MaxWorkers)

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

	de := errors.New("Error Type")
        errorInterface = reflect.ValueOf(de).Elem()

	d.managementInit(mo)
	d.mutex = new(sync.Mutex)
        d.metrics.mutex = new(sync.Mutex)

	d.MetricSetStartTime()

	return
}

// managementInit() setup management commands and fields
func (d *dispatcher) managementInit(mo *managementCommands) {
        var funcPtr interface{}
        var parmsList []string

        //d.managementOptions.CommandNames =  make(map[string]MgtCommand)
        parmsList = make([]string, 2, 2)
        parmsList[0] = "ResourseName"
        parmsList[1] = "ResourseValue"

        funcPtr = d.SetConfigVariable

        mo.useCommand("set", "Sets the value of a configurable field, see d.managementOptions.Fields", resDispatcher, parmsList, funcPtr)

        parmsList = make([]string, 0, 5)
        funcPtr = d.stopScheduler

        mo.useCommand("stop_cheduler", "Stops the scheduler from sending new jobs.", resDispatcher, parmsList, funcPtr)

        funcPtr = d.startScheduler
        mo.useCommand("start_scheduler", "Starts the scheduler running again.  If running has no affect.", resDispatcher, parmsList, funcPtr)


        funcPtr = d.stopWorkers
        mo.useCommand("stop_workers", "Shutdown the worker pool letting jobs inflight complete", resDispatcher, parmsList, funcPtr)

        funcPtr = d.startWorkers
        mo.useCommand("start_workers", "Starts the worker pool if stopped", resDispatcher, parmsList, funcPtr)

        funcPtr = d.shutdown
        mo.useCommand("shutdown", "Graceful shutdown", resDispatcher, parmsList, funcPtr)

        funcPtr = d.shutdownNow
        mo.useCommand("shutdown_now", "Hard shutdown with SIGKILL", resDispatcher, parmsList, funcPtr)

        funcPtr = d.SetConfigVariable

        //d.managementOptions.Fields = make(map[string]ConfigValue)

        mo.setField(gracefulShutdownSeconds, "int", resDispatcher, funcPtr)
        mo.setField(hardShutdownSeconds, "int", resDispatcher, funcPtr)
        mo.setField(numberOfWorkers, "int", resDispatcher, funcPtr)
        mo.setField(schedulerChannelSize, "int", resDispatcher, funcPtr)
       mo.setField(resultChannelSize, "int", resDispatcher, funcPtr)
}

/* DONE: added hooks to allows Job and Scheduler to extend management API
*/

func (d dispatcher) Run(mWg *sync.WaitGroup) error {
	defer mWg.Done()
	e:= d.createWorkerPool(mWg)
	if e == nil {
		mWg.Add(2)
        	go d.Forwarder(mWg)
	        go d.Responder(mWg)
	        return nil
        }
	return e
}

func (d dispatcher) Forwarder(mWg *sync.WaitGroup) {
	defer mWg.Done()
	for {
		select {
		case currentJob := <-d.schedulerJobChan:
			d.MetricInc(dispatcherJobsSent)
			d.workerJobChan <- currentJob
		case <-d.workerDone:
                        return
                default:
                        time.Sleep(1 * time.Second)
                }

	}
}

func (d dispatcher) Responder(mWg *sync.WaitGroup) {
	defer mWg.Done()
	log.Println("Dispatcher result channel started worker result -> scheduler  result")
	for {
		select {
		case currentJobResponse := <-d.workerJobResponse:
			{
			d.MetricInc(dispatcherResultsReceived)
			d.schedulerResultChan <- currentJobResponse
		        }
		case <-d.workerDone:
                        return
                default:
                        time.Sleep(1 * time.Second)
		}
	}
}

func (d *dispatcher) Shutdown() error {
	return nil
}

func (d *dispatcher) doGracefulShutDown(value int) (httpStatusCode int, msg []byte, err error) {
        if value < 0 {
                value = value * -1
        }
        if value > MaxGraceSec {
                value = MaxGraceSec
        }

        old := value
        d.mutex.Lock()
        {
                old = d.conf.gracefulShutdown
                d.conf.gracefulShutdown = value
        }
        d.mutex.Unlock()
        rmsg := fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
                gracefulShutdownSeconds, old, value)
        return http.StatusOK, []byte(rmsg), nil

}

func (d *dispatcher) doHardShutdownSeconds(value int) (httpStatusCode int, msg []byte, err error) {
        if value < 0 {
                value = value * -1
        }
        if value > MaxShutSec {
                value = MaxShutSec
        }
        old := value
        d.mutex.Lock()
        {
                old = d.conf.hardShutdown
                d.conf.hardShutdown = value
        }
        d.mutex.Unlock()
        rmsg := fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
                hardShutdownSeconds, old, value)
        return http.StatusOK, []byte(rmsg), nil

}

func (d *dispatcher) doNumberOfWorkers(value int) (httpStatusCode int, msg []byte, err error) {
        if value < 0 {
                value = value * -1
        }
        if value > MaxWorkers {
                value = MaxWorkers
        }
        old := value
        d.mutex.Lock()
        {
                old = d.conf.numberOfWorkers
                d.conf.numberOfWorkers = value
        }
        d.mutex.Unlock()

        //TODO: grow or srink as necessary
        rmsg := fmt.Sprintf("{\"Status\": \"%s changed from %d to %d\"}",
                numberOfWorkers, old, value)
        return http.StatusOK, []byte(rmsg), nil
}

func doNotImplemented(resourse string) (httpStatusCode int, msg []byte, err error) {
        rmsg := fmt.Sprintf("{\"Status\": \"%s configuration is not yet implemented.\"}", resourse)
        e := errors.New("not implemented")
        return http.StatusNotImplemented, []byte(rmsg), e

}

func doUnkownCmd(resourse string) (httpStatusCode int, msg []byte, err error) {
        rmsg := fmt.Sprintf("{\"Status\": \"%s is an unknown resourse.\"}", resourse)
        e := errors.New("unknown resourse to configure")

        return http.StatusNotFound, []byte(rmsg), e

}


// SetConfigVariable changegs the value of a given field
func (d *dispatcher) SetConfigVariable(name string, value int) (httpStatusCode int, msg []byte, err error) {

	fmt.Println("In SetConfigVariable")

        switch name {
        case gracefulShutdownSeconds:
                return d.doGracefulShutDown(value)

        case hardShutdownSeconds:
                return d.doHardShutdownSeconds(value)

        case numberOfWorkers:
              return d.doNumberOfWorkers(value)

        case schedulerChannelSize:
                return doNotImplemented(name)

        case resultChannelSize:
                return doNotImplemented(name)
        default:
                return doUnkownCmd(name)
        }
}

func (d *dispatcher) createWorkerPool(mWg *sync.WaitGroup) error {
	 currNumWorkers := len(d.workers)
        if currNumWorkers > 0 {
                log.Printf("Prior worker pool with  %d workers.\n", currNumWorkers)
                return nil
        }

        d.mutex.Lock()
        {
                for i := 0; i < d.conf.numberOfWorkers; i++ {
                        newWorker := worker{
                               //                      wg:           d.wg,
                                jobHist:      make(map[uuid.UUID]string),
                                jobChan:      d.workerJobChan,
                                responseChan: d.workerJobResponse,
                                interrupt:    d.workerInterrupt,
                                done:         make(chan bool, 1),
                                workerID:     i,
                        }

                        //d.wg.Add(1)
                        mWg.Add(1)
                      // Keep track of each worker
                        d.workers = append(d.workers, &newWorker)

                        go newWorker.Run(mWg)
                }
        }
        d.mutex.Unlock()

        d.conf.currentNumberOfWorkers = d.conf.numberOfWorkers
        log.Printf("Worker pool created, %d workers\n", d.conf.numberOfWorkers)
        return nil

}

func (d *dispatcher) growWorkerPool() error {
	return nil
}

func (d *dispatcher) swrinkWorkerPool() error {
	return nil
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
