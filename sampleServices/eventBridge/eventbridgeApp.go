// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
//

// User project / copyright / usage information
// Allocate jobs to workers in a pool

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Initialize setups database connection object and the http server
func (a *EventbridgeApp) Initialize() {

	// set k8s probe
	a.Live = false
	a.Ready = false

	// Override defaults
	a.initializeEnvironment()

	// Start the Dispatcher
	//a.Scheduler = &bridgeScheduler{}
	a.Scheduler = &httpScheduler{}

	dConf := &dispatcherConfiguration{
		scheduler:           a.Scheduler,
		sizeOfJobChannel:    SizeOfJobChannel,
		sizeOfResultChannel: SizeOfResultChannel,
		numberOfWorkers:     NumberOfWorkers,
		gracefulShutdown:    GracefullShutdown,
		hardShutdown:        HardShutdown,
	}

	a.Dispatcher.Init(dConf)
	go a.Dispatcher.Run()

	// Scheduler
	err := a.Scheduler.Init()
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	go a.Scheduler.Run()

	a.Ready = true
	// Start rest end points
	httpconf.listenString = fmt.Sprintf("%s:%s", httpconf.ip, httpconf.port)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run start the HTTP server for Rest endpoints
func (a *EventbridgeApp) Run(addr string) {

	log.Println("Listing at: " + addr)
	// Wrap router with W3C logging

	lf, _ := os.OpenFile("logs/access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	loggedRouter := handlers.LoggingHandler(lf, a.Router)
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         addr,
		WriteTimeout: httpconf.writeTimeout * time.Second,
		ReadTimeout:  httpconf.readTimeout * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	a.Live = true

	// Listen for SIGHUP
	a.httpInterruptChan = make(chan os.Signal, 1)

	<-a.httpInterruptChan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), httpconf.shutdownTimeout)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	err := srv.Shutdown(ctx)
	if err == nil {
		log.Println("Shutting Down...")
	} else {
		log.Println("Shutting Down...", err)
	}

	os.Exit(0)
}

// Get for environment variable overrides
func (a *EventbridgeApp) initializeEnvironment() {
	var envVar = ""

	envVar = os.Getenv("HTTP_IP_ADDR")
	if envVar != "" {
		httpconf.ip = envVar
	}

	envVar = os.Getenv("HTTP_IP_PORT")
	if envVar != "" {
		httpconf.port = envVar
	}

	envVar = os.Getenv("HTTP_READ_TIMEOUT")
	if envVar != "" {
		to, err := strconv.Atoi(envVar)
		if err != nil {
			log.Printf("failed to convert HTTP_READ_TIMEOUT: %s to int", envVar)
		} else {
			httpconf.readTimeout = (time.Second * time.Duration(to))
		}
		log.Printf("Read timeout: %d", httpconf.readTimeout)
	}

	envVar = os.Getenv("HTTP_WRITE_TIMEOUT")
	if envVar != "" {
		to, err := strconv.Atoi(envVar)
		if err != nil {
			log.Printf("failed to convert HTTP_READ_TIMEOUT: %s to int", envVar)
		} else {
			httpconf.writeTimeout = time.Duration(to) * time.Second
		}
		log.Printf("Write timeout: %d", httpconf.writeTimeout)
	}

	envVar = os.Getenv("HTTP_SHUTDOWN_TIMEOUT")
	if envVar != "" {
		if envVar != "" {
			to, err := strconv.Atoi(envVar)
			if err != nil {
				log.Printf("failed to convert HTTP_SHUTDOWN_TIMEOUT: %s to int", envVar)
			} else {
				httpconf.shutdownTimeout = time.Second * time.Duration(to)
			}
			log.Println("Shutdown timeout", httpconf.shutdownTimeout)
		}
	}

	envVar = os.Getenv("HTTP_LOG")
	if envVar != "" {
		httpconf.logPath = envVar
	}

}

func (a *EventbridgeApp) initializeRoutes() {

	uri := EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeJobsEndPoint + "LIST"
	a.Router.HandleFunc(uri, a.listJobs).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeSchedulerEndPoint + "LIST"
	a.Router.HandleFunc(uri, a.listSchedule).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeJobsEndPoint + EventbridgeKey
	a.Router.HandleFunc(uri, a.getJob).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeSchedulerEndPoint
	a.Router.HandleFunc(uri, a.getSchedule).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeLivenessEndPoint
	a.Router.HandleFunc(uri, a.getLiveness).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeReadinessEndPoint
	a.Router.HandleFunc(uri, a.getReadiness).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeMetricsEndPoint
	a.Router.HandleFunc(uri, a.getMetrics).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeManagementEndPoint
	a.Router.HandleFunc(uri, a.getManagement).Methods("GET")
	log.Println("GET: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeManagementEndPoint
	a.Router.HandleFunc(uri, a.putManagement).Methods("PUT")
	log.Println("PUT: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeJobsEndPoint + EventbridgeKey
	a.Router.HandleFunc(uri, a.updateJob).Methods("PUT")
	log.Println("PUT: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeJobsEndPoint + EventbridgeKey
	a.Router.HandleFunc(uri, a.deleteJob).Methods("DELETE")
	log.Println("DELETE: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeJobsEndPoint
	a.Router.HandleFunc(uri, a.createJob).Methods("POST")
	log.Println("POST: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeSchedulerEndPoint
	a.Router.HandleFunc(uri, a.updateSchedule).Methods("PUT")
	log.Println("PUT: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeSchedulerEndPoint
	a.Router.HandleFunc(uri, a.deleteSchedule).Methods("DELETE")
	log.Println("DELETE: ", uri)

	uri = EventbridgeAPIVersion + "/" +
		EventbridgeNamespaceID + "/" +
		EventbridgeDefaultNamespace + "/" +
		EventbridgeResourceType + "/" +
		EventbridgeSchedulerEndPoint
	a.Router.HandleFunc(uri, a.createSchedule).Methods("POST")
	log.Println("POST: ", uri)

	return
}

// listJobs swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/EventbridgeJobsEndPointLIST jobs listjobs
//
// Returns a list of Jobs
//
// Responses:
//				200: listJobResponse
//        		500: genericError

func (a *EventbridgeApp) listJobs(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	// Pre-processing hook
	listJobsPreHook(w, r, count, start)

	jl, e := a.Scheduler.GetScheduledJobs()

	if e != nil {
		respondWithError(w, http.StatusInternalServerError, e.Error())
	}

	// Post-processing hook
	listJobsPostHook(w, r)

	respondWithByte(w, http.StatusOK, jl)
}

// TODO: decide do kill it or do something with it

// listSchedule swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/EventbridgeSchedulerEndPointLIST schedules listschedule
//
// Returns a list of schedules
//
// Responses:
//		default: genericError

func (a *EventbridgeApp) listSchedule(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, "{}")
}

// getJob swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/jobs/{key} job getjob
//
// Returns a job given a key, where key is a UUID
//
// Responses:
//		default: genericError
//				200: jobResponse
//				200: listJobResponse
//				404: get404Response
//				500: genericError
func (a *EventbridgeApp) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	getJobPreHook(w, r, key)

	status, jb, e := a.Scheduler.GetScheduleJob(key)

	if e != nil {
		log.Println(e)
	}

	// Pre-processing hook
	getJobPostHook(w, r, key)

	respondWithByte(w, status, jb)
}

// getSchedule swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/schedule/{key} schedule getschedule
//
// Returns a schedule given a key, where key is a UUID
//
// Responses:
//		default: genericError
//		200: genericResponse
//		500: genericError
func (a *EventbridgeApp) getSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	getSchedulePreHook(w, r, key)

	status, respBody, e := a.Scheduler.GetSchedule()
	if e != nil {
		log.Println(e)
	}

	// Pre-processing hook
	getSchedulePostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// getLiveness swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/liveness liveness getliveness
//
// A HTTP response status code between 200-400 indicates the pod is alive.
// Any other status code will cause kubelet to restart the pod.
//
// Responses:
//		default: genericError
//				200: genericResponse
//				503: genericError
func (a *EventbridgeApp) getLiveness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getLivenessPreHook(w, r)

	if !a.Live {
		respondWithByte(w, http.StatusServiceUnavailable, []byte("{\"Live\": false}"))
		return
	}

	respondWithByte(w, http.StatusOK, []byte("{\"Live\": true}"))
}

// getReadiness swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/ready ready getready
//
// Indicates the pod is ready to start taking traffic.
// Should return a 200 after all pod initialization has completed.
//
// Responses:
//		default: genericError
//				200: genericResponse
//				503: genericError

func (a *EventbridgeApp) getReadiness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getReadinessPreHook(w, r)

	if !a.Ready {
		respondWithByte(w, http.StatusServiceUnavailable, []byte("{\"Ready\": false}"))
		return
	}

	respondWithByte(w, http.StatusOK, []byte("{\"Ready\": true}"))
}

// getMetrics swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/metrics metrics getMetrics
//
// Returns metrics for eventbridge service
// Metrics should include:
//	 - Scheduler
//	 - Dispatcher
//	 - Workers
//	 - Jobs
//
// Responses:
//		default: genericError
//				200: metricsResponse

func (a *EventbridgeApp) getMetrics(w http.ResponseWriter, r *http.Request) {
	var combinedJSON string = "{"

	// Pre-processing hook
	getMetricsPreHook(w, r)

	sm := a.Scheduler.Metrics()
	if sm != nil {
		combinedJSON += `"scheduler":`
		combinedJSON += string(sm)
	}

	a.Dispatcher.MetricUpdateUpTime()
	dm, _ := a.Dispatcher.MetricToJSON()
	if dm != nil {
		combinedJSON += `,"dispatcher":`
		combinedJSON += string(dm)
	}

	combinedJSON += "}"

	// Post-processing hook
	getMetricsPostHook(w, r)

	respondWithByte(w, http.StatusOK, []byte(combinedJSON))
}

// getManagement swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/management management getmanagement
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: managementGetResponse
func (a *EventbridgeApp) getManagement(w http.ResponseWriter, r *http.Request) {
	// Pre-processing hook
	getManagementPreHook(w, r)

	// Post-processing hook
	getManagementPostHook(w, r)

	respondWithJSON(w, http.StatusOK, a.Dispatcher.managementOptions)
}

// putmanagement swagger:route GET /api/v1/namespace/pavedroad/Eventbridge/management management putmanagement
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: genericResponse
//				400: genericError
//				500: genericError
func (a *EventbridgeApp) putManagement(w http.ResponseWriter, r *http.Request) {

	var requestedCommand managementRequest

	// Pre-processing hook
	putManagementPreHook(w, r)

	payload, e := ioutil.ReadAll(r.Body)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, http.StatusBadRequest, []byte(msg))
		return
	}

	e = json.Unmarshal(payload, &requestedCommand)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"json.Unmarshal failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, http.StatusBadRequest, []byte(msg))
		return
	}

	status, respBody, e := a.Dispatcher.ProcessManagementRequest(requestedCommand)

	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"Management command failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, status, []byte(msg))
		return
	}

	// Post-processing hook
	putManagementPostHook(w, r)

	// TODO: find a way to flush this out
	respondWithByte(w, status, respBody)

	// Special case for shutting down
	if requestedCommand.Command == "shutdown" {
		time.Sleep(time.Duration(a.Dispatcher.conf.gracefulShutdown) * time.Second)
		e = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if e != nil {
			msg := fmt.Sprintf("{\"error\": \"Shutting down with\", \"Error\": \"%v\"}", e.Error())
			log.Println("shutdown error:", msg)
		}

	}

	// Special case for hard kill
	// We've sent the reply
	if requestedCommand.Command == "shutdown_now" {
		time.Sleep(time.Duration(a.Dispatcher.conf.hardShutdown) * time.Second)
		e = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if e != nil {
			msg := fmt.Sprintf("{\"error\": \"Shutting down with\", \"Error\": \"%v\"}", e.Error())
			log.Println("shutdown_now error:", msg)
		}

	}

	return
}

// createJob swagger:route POST /api/v1/namespace/pavedroad/Eventbridge/EventbridgeJobsEndPoint EventbridgeJobsEndPoint createJob
//
// Create a new Job
//
// Responses:
//		default: genericError
//				201: listJobResponse
//				400: genericError
func (a *EventbridgeApp) createJob(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	createJobPreHook(w, r)

	payload, e := ioutil.ReadAll(r.Body)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, http.StatusBadRequest, []byte(msg))
		return
	}

	status, respBody, e := a.Scheduler.CreateScheduleJob(payload)

	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"CreateScheduleJob failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, status, []byte(msg))
		return
	}

	// Post-processing hook
	createJobPostHook(w, r)

	respondWithByte(w, status, respBody)
}

// updateJob swagger:route PUT /api/v1/namespace/pavedroad/Eventbridge/EventbridgeJobsEndPoint/{key} EventbridgeSchedulerEndPoint updateJob
//
// Update a EventbridgeJobsEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: listJobResponse
//				400: genericError
//				404: genericError
func (a *EventbridgeApp) updateJob(w http.ResponseWriter, r *http.Request) {
	// Read URI variables
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	updateJobPreHook(w, r, key)

	payload, e := ioutil.ReadAll(r.Body)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, http.StatusBadRequest, []byte(msg))
		return
	}

	status, respBody, e := a.Scheduler.UpdateScheduleJob(payload)

	if e != nil {
		log.Printf("UpdateScheduleJob error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	updateJobPostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// deleteJob swagger:route DELETE /api/v1/namespace/pavedroad/Eventbridge/EventbridgeJobsEndPoint/{key} EventbridgeJobsEndPoint deleteJobs
//
// Delete a job specified by key, which is a uuid
//
// Responses:
//		default: genericError
//				200: listJobResponse
//				400: genericError
func (a *EventbridgeApp) deleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	deleteJobPreHook(w, r, key)

	status, respBody, e := a.Scheduler.DeleteScheduleJob(key)

	if e != nil {
		log.Printf("DeleteScheduleJob error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	deleteJobPostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// createSchedule swagger:route POST /api/v1/namespace/pavedroad/Eventbridge/EventbridgeSchedulerEndPoint EventbridgeSchedulerEndPoint createSchedule
//
// Create a new scheduler
//
// Responses:
//		default: genericError
//				201: genericResponse
//				400: genericError
func (a *EventbridgeApp) createSchedule(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	createSchedulePreHook(w, r)

	payload, e := ioutil.ReadAll(r.Body)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, http.StatusBadRequest, []byte(msg))
		return
	}

	status, respBody, e := a.Scheduler.UpdateSchedule(payload)

	if e != nil {
		log.Printf("CreateSchedule error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	createSchedulePostHook(w, r)
	// We are using update that will return a 200, change status to a 201
	respondWithByte(w, http.StatusCreated, respBody)
}

// updateSchedle swagger:route PUT /api/v1/namespace/pavedroad/Eventbridge/EventbridgeSchedulerEndPoint/{key} EventbridgeSchedulerEndPoint updateSchedule
//
// Update a EventbridgeSchedulerEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: genericResponse
//				400: genericError
func (a *EventbridgeApp) updateSchedule(w http.ResponseWriter, r *http.Request) {
	// Read URI variables
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	updateSchedulePreHook(w, r, key)

	payload, e := ioutil.ReadAll(r.Body)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
		respondWithByte(w, http.StatusBadRequest, []byte(msg))
		return
	}

	status, respBody, e := a.Scheduler.UpdateSchedule(payload)

	if e != nil {
		log.Printf("updateSchedule error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	updateSchedulePostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// deleteSchedule swagger:route DELETE /api/v1/namespace/pavedroad/Eventbridge/EventbridgeSchedulerEndPoint/{key} EventbridgeSchedulerEndPoint deleteSchudler
//
// Delete a job specified by key, which is a uuid
//
// Responses:
//		default: genericError
//				200: genericResponse
//				400: genericError
func (a *EventbridgeApp) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	// Read URI variables
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	deleteSchedulePreHook(w, r, key)

	status, respBody, e := a.Scheduler.DeleteSchedule()

	if e != nil {
		log.Printf("DeleteSchedule error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	deleteSchedulePostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

// respondWithByte not need to Marshal the JSON
func respondWithByte(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, e := w.Write(payload)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"payload error: \", \"Error\": \"%v\"}", e.Error())
		log.Println("respondWithByte error:", msg)
	}
}

// respondWithJSON will Marshal the payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, e := w.Write(response)
	if e != nil {
		msg := fmt.Sprintf("{\"error\": \"payload error: \", \"Error\": \"%v\"}", e.Error())
		log.Println("respondWithByte error:", msg)
	}

}

func openAccessLogFile(accesslogfile string) *os.File {
	var lf *os.File
	var err error

	if accesslogfile == "" {
		accesslogfile = "access.log"
		log.Println("Access log file name not declared using access.log")
	}

	_, _ = rollLogIfExists(accesslogfile)

	lf, err = os.OpenFile(accesslogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		log.Println("Error opening access log file: os.OpenFile:", err)
		return nil
	}

	return lf
}

func openErrorLogFile(errorlogfile string) error {
	if errorlogfile == "" {
		errorlogfile = "error.log"
		log.Println("Error log file name not declared using errors.log")
	}

	_, _ = rollLogIfExists(errorlogfile)

	lf, err := os.OpenFile(errorlogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

	if err != nil {
		log.Println("Error opening error log file: os.OpenFile:", err)
		return err
	}
	log.SetOutput(lf)
	return nil
}

func rollLogIfExists(logfilename string) (string, error) {

	if _, err := os.Stat(logfilename); os.IsNotExist(err) {
		return "", err
	}
	var newFileName string
	tn := time.Now()
	endsWithDotLogIdx := strings.LastIndex(logfilename, ".log")
	if endsWithDotLogIdx == -1 {
		newFileName = logfilename + tn.Format(time.RFC3339)
	} else {
		newFileName = logfilename[0:endsWithDotLogIdx] +
			tn.Format(time.RFC3339) + ".log"
	}

	err := os.Rename(logfilename, newFileName)
	if err != nil {
		msg := fmt.Sprintf("Rename logfile %v to %v failed with error %v\n",
			logfilename, newFileName, err.Error())
		log.Printf(msg)
		return "", err
	}

	return newFileName, nil

}
