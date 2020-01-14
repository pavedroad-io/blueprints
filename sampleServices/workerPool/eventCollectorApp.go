//
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root for full license information.
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
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Initialize setups database connection object and the http server
func (a *EventCollectorApp) Initialize() {

	// set k8s probe
	a.Live = false
	a.Ready = false

	// Override defaults
	a.initializeEnvironment()

	// Start the Dispatcher
	// TODO: make size of result channel configurable
	// TODO: change to take a configuration structure
	a.Dispatcher.Init(NUMBEROFWORKERS, SIZEOFJOBCHANNEL, &a.Scheduler)
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
func (a *EventCollectorApp) Run(addr string) {

	log.Println("Listing at: " + addr)
	srv := &http.Server{
		Handler:      a.Router,
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
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// Get for ennvironment variable overrides
func (a *EventCollectorApp) initializeEnvironment() {
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
		if err == nil {
			log.Printf("failed to convert HTTP_READ_TIMEOUT: %s to int", envVar)
		} else {
			httpconf.readTimeout = time.Duration(to) * time.Second
		}
		log.Printf("Read timeout: %d", httpconf.readTimeout)
	}

	envVar = os.Getenv("HTTP_WRITE_TIMEOUT")
	if envVar != "" {
		to, err := strconv.Atoi(envVar)
		if err == nil {
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
				httpconf.shutdownTimeout = time.Second * time.Duration(to)
			} else {
				httpconf.shutdownTimeout = time.Second * httpconf.shutdownTimeout
			}
			log.Println("Shutdown timeout", httpconf.shutdownTimeout)
		}
	}

	envVar = os.Getenv("HTTP_LOG")
	if envVar != "" {
		httpconf.logPath = envVar
	}

}

func (a *EventCollectorApp) initializeRoutes() {

	uri := EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorJobsEndPoint + "LIST"
	a.Router.HandleFunc(uri, a.listJobs).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorSchedulerEndPoint + "LIST"
	a.Router.HandleFunc(uri, a.listSchedule).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorJobsEndPoint + EventCollectorKey
	a.Router.HandleFunc(uri, a.getJob).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorSchedulerEndPoint
	a.Router.HandleFunc(uri, a.getSchedule).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorLivenessEndPoint
	a.Router.HandleFunc(uri, a.getLiveness).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorReadinessEndPoint
	a.Router.HandleFunc(uri, a.getReadiness).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorMetricsEndPoint
	a.Router.HandleFunc(uri, a.getMetrics).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorManagementEndPoint
	a.Router.HandleFunc(uri, a.getManagement).Methods("GET")
	fmt.Println("GET: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorManagementEndPoint
	a.Router.HandleFunc(uri, a.putManagement).Methods("PUT")
	fmt.Println("PUT: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorJobsEndPoint + EventCollectorKey
	a.Router.HandleFunc(uri, a.updateJob).Methods("PUT")
	fmt.Println("PUT: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorJobsEndPoint + EventCollectorKey
	a.Router.HandleFunc(uri, a.deleteJob).Methods("DELETE")
	fmt.Println("DELETE: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorJobsEndPoint
	a.Router.HandleFunc(uri, a.createJob).Methods("POST")
	fmt.Println("POST: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorSchedulerEndPoint
	a.Router.HandleFunc(uri, a.updateSchedule).Methods("PUT")
	fmt.Println("PUT: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorSchedulerEndPoint
	a.Router.HandleFunc(uri, a.deleteSchedule).Methods("DELETE")
	fmt.Println("DELETE: ", uri)

	uri = EventCollectorAPIVersion + "/" +
		EventCollectorNamespaceID + "/" +
		EventCollectorDefaultNamespace + "/" +
		EventCollectorResourceType + "/" +
		EventCollectorSchedulerEndPoint
	a.Router.HandleFunc(uri, a.createSchedule).Methods("POST")
	fmt.Println("POST: ", uri)

	return
}

// listJobs swagger:route GET /api/v1/namespace/mirantis/eventCollector/EventCollectorJobsEndPointLIST jobs listjobs
//
// Returns a list of Jobs
//
// Responses:
//		default: genericError
//				200: jobsList

func (a *EventCollectorApp) listJobs(w http.ResponseWriter, r *http.Request) {

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

// listSchedule swagger:route GET /api/v1/namespace/mirantis/eventCollector/EventCollectorSchedulerEndPointLIST schedules listschedule
//
// Returns a list of schedules
//
// Responses:
//		default: genericError
//				200: scheduleList

// TODO: decide do kill it or do something with it
func (a *EventCollectorApp) listSchedule(w http.ResponseWriter, r *http.Request) {
	//eventCollector := eventCollector{}

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	// Pre-processing hook
	listSchedulePreHook(w, r, count, start)

	/*
	 jl, e := a.Scheduler.GetSchedule()
	*/

	// Post-processing hook
	listSchedulePostHook(w, r)

	respondWithJSON(w, http.StatusOK, "{}")
}

// getJob swagger:route GET /api/v1/namespace/mirantis/eventCollector/jobs/{key} job getjob
//
// Returns a job given a key, where key is a UUID
//
// Responses:
//		default: genericError
//				200: jobResponse
func (a *EventCollectorApp) getJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	getJobPreHook(w, r, key)

	status, jb, e := a.Scheduler.GetScheduleJob(key)

	if e != nil {
		// TODO: log for internal errors
		fmt.Println(e)
	}

	// Pre-processing hook
	getJobPostHook(w, r, key)

	respondWithByte(w, status, jb)
}

// getSchedule swagger:route GET /api/v1/namespace/mirantis/eventCollector/schedule/{key} schedule getschedule
//
// Returns a schedule given a key, where key is a UUID
//
// Responses:
//		default: genericError
//				200: scheduleResponse

func (a *EventCollectorApp) getSchedule(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	getSchedulePreHook(w, r, key)

	status, respBody, e := a.Scheduler.GetSchedule()

	if e != nil {
		// TODO: log for internal errors
		fmt.Println(e)
	}

	// Pre-processing hook
	getSchedulePostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// getLiveness swagger:route GET /api/v1/namespace/mirantis/eventCollector/liveness liveness getliveness
//
// A HTTP response status code between 200-400 indicates the pod is alive.
// Any other status code will cause kubelet to restart the pod.
//
// Responses:
//		default: genericError
//				200: livenessResponse

func (a *EventCollectorApp) getLiveness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getLivenessPreHook(w, r)

	if !a.Live {
		respondWithError(w, http.StatusServiceUnavailable, "{\"Live\": false}")
	}

	respondWithJSON(w, http.StatusOK, a.Live)
}

// getReadiness swagger:route GET /api/v1/namespace/mirantis/eventCollector/ready ready getready
//
// Indicates the pod is ready to start taking traffic.
// Should return a 200 after all pod initialization has completed.
//
// Responses:
//		default: genericError
//				200: readinessResponse

func (a *EventCollectorApp) getReadiness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getReadinessPreHook(w, r)

	if !a.Ready {
		respondWithError(w, http.StatusServiceUnavailable, "{\"Ready\": false}")
	}

	respondWithJSON(w, http.StatusOK, a.Ready)
}

// getMetrics swagger:route GET /api/v1/namespace/mirantis/eventCollector/metrics metrics getMetrics
//
// Returns metrics for eventCollector service
// Metrics should include:
//	 - Scheduler
//	 - Dispatcher
//	 - Workers
//	 - Jobs
//
// Responses:
//		default: genericError
//				200: readinessResponse

func (a *EventCollectorApp) getMetrics(w http.ResponseWriter, r *http.Request) {
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

// getManagement swagger:route GET /api/v1/namespace/mirantis/eventCollector/management management getManagement
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: managementResponse
func (a *EventCollectorApp) getManagement(w http.ResponseWriter, r *http.Request) {
	// Pre-processing hook
	getManagementPreHook(w, r)

	// Post-processing hook
	getManagementPostHook(w, r)

	respondWithJSON(w, http.StatusOK, a.Dispatcher.managementOptions)
}

// putManagement swagger:route PUT /api/v1/namespace/mirantis/eventCollector/management management putManagement
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: managementResponse
func (a *EventCollectorApp) putManagement(w http.ResponseWriter, r *http.Request) {

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
		respondWithByte(w, http.StatusInternalServerError, []byte(msg))
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
		// Give it 1 second to be sent
		time.Sleep(time.Duration(a.Dispatcher.gracefulShutdown) * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	// Special case for hard kill
	// We've sent the reply
	if requestedCommand.Command == "shutdown_now" {
		// Give it 1 second to be sent
		time.Sleep(time.Duration(a.Dispatcher.hardShutdown) * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	return
}

// createJob swagger:route POST /api/v1/namespace/mirantis/eventCollector/EventCollectorJobsEndPoint EventCollectorJobsEndPoint createJob
//
// Create a new Job
//
// Responses:
//		default: genericError
//				201: jobResponse
//				400: genericError
func (a *EventCollectorApp) createJob(w http.ResponseWriter, r *http.Request) {

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

// updateJob swagger:route PUT /api/v1/namespace/mirantis/eventCollector/EventCollectorJobsEndPoint/{key} EventCollectorSchedulerEndPoint updateJob
//
// Update a EventCollectorJobsEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: jobResponse
//				400: genericError
//				404: genericError
func (a *EventCollectorApp) updateJob(w http.ResponseWriter, r *http.Request) {
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
		// TODO: log this
		fmt.Printf("UpdateScheduleJob error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	updateJobPostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// deleteJob swagger:route DELETE /api/v1/namespace/mirantis/eventCollector/EventCollectorJobsEndPoint/{key} EventCollectorJobsEndPoint deleteJobs
//
// Delete a job specified by key, which is a uuid
//
// Responses:
//		default: genericError
//				200: jobResponse
//				400: genericError
func (a *EventCollectorApp) deleteJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	deleteJobPreHook(w, r, key)

	status, respBody, e := a.Scheduler.DeleteScheduleJob(key)

	if e != nil {
		// TODO: log this
		fmt.Printf("DeleteScheduleJob error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	deleteJobPostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// createSchedule swagger:route POST /api/v1/namespace/mirantis/eventCollector/EventCollectorSchedulerEndPoint EventCollectorSchedulerEndPoint createSchedule
//
// Create a new scheduler
//
// Responses:
//		default: genericError
//				201: schedulerResponse
//				400: genericError
func (a *EventCollectorApp) createSchedule(w http.ResponseWriter, r *http.Request) {

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
		// TODO: log this
		fmt.Printf("CreateSchedule error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	createSchedulePostHook(w, r)

	respondWithByte(w, status, respBody)
}

// updateSchedle swagger:route PUT /api/v1/namespace/mirantis/eventCollector/EventCollectorSchedulerEndPoint/{key} EventCollectorSchedulerEndPoint updateSchedule
//
// Update a EventCollectorSchedulerEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: schedulerResponse
//				400: genericError
func (a *EventCollectorApp) updateSchedule(w http.ResponseWriter, r *http.Request) {
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
		// TODO: log this
		fmt.Printf("updateSchedule error: %v status %v", e.Error(), status)
	}

	// Post-processing hook
	updateSchedulePostHook(w, r, key)

	respondWithByte(w, status, respBody)
}

// deleteSchedule swagger:route DELETE /api/v1/namespace/mirantis/eventCollector/EventCollectorSchedulerEndPoint/{key} EventCollectorSchedulerEndPoint deleteSchudler
//
// Delete the current schedule by setting interval to 272 billion years
//
// Responses:
//		default: genericError
//				200: schedulerResponse
//				400: genericError
func (a *EventCollectorApp) deleteSchedule(w http.ResponseWriter, r *http.Request) {
	// Read URI variables
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	deleteSchedulePreHook(w, r, key)

	status, respBody, e := a.Scheduler.DeleteSchedule()

	if e != nil {
		// TODO: log this
		fmt.Printf("DeleteSchedule error: %v status %v", e.Error(), status)
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
	w.Write(payload)
}

// respondWithJSON will Marshal the payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// TODO: Figure out why this isn't logging each request, aka access log
// TODO: two logs one for diag and one for access
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func openLogFile(logfile string) {
	if logfile != "" {
		lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

		if err != nil {
			log.Fatal("OpenLogfile: os.OpenFile:", err)
		}
		log.SetOutput(lf)
	}
}
