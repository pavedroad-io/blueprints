{{define "templateApp.go"}}{{.PavedroadInfo}}

// User project / copyright / usage information
// {{.ProjectInfo}}

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
func (a *{{.NameExported}}App) Initialize() {

	// set k8s probe
	a.Live = false
	a.Ready = false

	// Override defaults
	a.initializeEnvironment()

	// Start the Dispatcher
	// TOOD generate this next line from roadctl
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
func (a *{{.NameExported}}App) Run(addr string) {

	log.Println("Listing at: " + addr)
	// Wrap router with w3C logging

	lf, _ := os.OpenFile("logs/access.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

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
	srv.Shutdown(ctx)
	log.Println("shutting down")
	os.Exit(0)
}

// Get for ennvironment variable overrides
func (a *{{.NameExported}}App) initializeEnvironment() {
	var envVar = ""

        //string environment variables map 
        varStrToProc := map[string]*string{
                "HTTP_IP_ADDR": &httpconf.ip,
                "HTTP_IP_PORT": &httpconf.port,
                "HTTP_LOG": &httpconf.logPath,
        }

        //time environment variables map
        varTimeToProc := map[string]*int{
                "HTTP_READ_TIMEOUT": &httpconf.readTimeout,
                "HTTP_WRITE_TIMEOUT": &httpconf.writeTimeout,
                "HTTP_SHUTDOWN_TIMEOUT": &httpconf.shutdownTimeout,
        }

        //Environment variable list to process
        //Expand with case statement
        varsToProc := [][]string{
                []string{"HTTP_IP_ADDR","string"},
                []string{"HTTP_IP_PORT","string"},
                []string{"HTTP_READ_TIMEOUT","time"},
                []string{"HTTP_WRITE_TIMEOUT","time"},
                []string{"HTTP_SHUTDOWN_TIMEOUT","time"},
                []string{"HTTP_LOG","string"},
        }

	for i := 0; i < len(varsToProc); i++ {
                envVar = os.Getenv(varToProc[i][0])
                if envVar != "" {
                   switch varTyp := varToProc[i][1]; varTyp {
                   case "string" :
                         *varStrToProc[varToProc[i][0]] = envVar
                   case "time" :
                         to, err := strconv.Atoi(envVar)
                         if err == nil {
                        log.Printf("failed to convert %s : %s to int",varToProc[i,0] ,envVar)
                        } else {
                          to = time.Duration(to) * time.Second
                          *varTimeToProc[varToProc[i][0]] = to
                        log.Printf("%s : %d", varToProc[i,0], to)
                        }
                   default:
                        log.Printf("Env variable type  %s, not supporte",varTyp)

                }
              }
        }

}

{{.AllRoutesSwaggerDoc}}
func (a *{{.NameExported}}App) initializeRoutes() {

	var uriPrefix = ""

        uriPrefix =  {{.NameExported}}APIVersion + "/" +
                                {{.NameExported}}NamespaceID + "/" +
                                {{.NameExported}}DefaultNamespace + "/" +
                                {{.NameExported}}ResourceType + "/"


        uri := uriPrefix +
               EventCollectorJobsEndPoint + "LIST"

        a.Router.HandleFunc(uri, a.listJobs).Methods("GET")
	log.Println("Get: ",uri)

        uri = uriPrefix +
              EventCollectorSchedulerEndPoint + "LIST"
        a.Router.HandleFunc(uri, a.listSchedule).Methods("GET")
        log.Println("Get: ",uri)
        
	uri = uriPrefix +
              EventCollectorJobsEndPoint + EventCollectorKey
        a.Router.HandleFunc(uri, a.getJob).Methods("GET")
        log.Println("Get: ",uri)

        uri = uriPrefix +
              EventCollectorSchedulerEndPoint 
        a.Router.HandleFunc(uri, a.getSchedule).Methods("GET")
        log.Println("Get: ",uri)

        uri = uriPrefix +
              EventCollectorLivenessEndPoint
        a.Router.HandleFunc(uri, a.getLiveness).Methods("GET")
        log.Println("Get: ",uri)

        uri = uriPrefix +
              EventCollectorReadinessEndPoint
        a.Router.HandleFunc(uri, a.getReadiness).Methods("GET")
        log.Println("Get: ",uri)

        uri = uriPrefix +
              EventCollectorMetricsEndPoint
        a.Router.HandleFunc(uri, a.getMetrics).Methods("GET")
        log.Println("Get: ",uri)

        //uri same for next getManagement,putManagement
        uri = uriPrefix +
              EventCollectorManagementEndPoint
        a.Router.HandleFunc(uri, a.getManagement).Methods("GET")
        log.Println("GET: ", uri)

        a.Router.HandleFunc(uri, a.putManagement).Methods("PUT")
        log.Println("PUT: ", uri)


	uri = uriPrefix +
              EventCollectorJobsEndPoint
        a.Router.HandleFunc(uri, a.createJob).Methods("POST")
        log.Println("POST :" ,uri)


        //uri same for deleteJob to follow
	uri = uriPrefix +
        EventCollectorJobsEndPoint + EventCollectorKey
        a.Router.HandleFunc(uri, a.updateJob).Methods("PUT")
        log.Println(("PUT: ",uri)

	//uri same as above updateJob
        a.Router.HandleFunc(uri, a.deleteJob).Methods("DELETE")
        log.Println("DELETE: ", uri)
        

	uri = uriPrefix +
        EventCollectorSchedulerEndPoint
        a.Router.HandleFunc(uri, a.createSchedule).Methods("POST")
        log.Println("POST:", uri)


	//uri same for createSchedule
        a.Router.HandleFunc(uri, a.updateSchedule).Methods("PUT")
        log.Println(("PUT :" ,uri)

        //uri same for createSchedule above
        a.Router.HandleFunc(uri, a.deleteSchedule).Methods("DELETE")
        log.Println("Delete:", uri)

	return
}

{{.GetAllSwaggerDoc}}
// listJobs swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/EventCollectorJobsEndPointLIST jobs listjobs
//
// Returns a list of Jobs
//
// Responses:
//		default: genericError
//				200: jobsList

func (a *{{.NameExported}}App) listJobs(w http.ResponseWriter, r *http.Request) {
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

{{.GetAllSwaggerDoc}}
// listSchedule swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/EventCollectorSchedulerEndPointLIST schedules listschedule
//
// Returns a list of schedules
//
// Responses:
//		default: genericError
//				200: scheduleList

// TODO: decide do kill it or do something with it
func (a *{{.NameExported}}App) listSchedule(w http.ResponseWriter, r *http.Request) {

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

{{.GetSwaggerDoc}}
// getJob swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/jobs/{key} job getjob
//
// Returns a job given a key, where key is a UUID
//
// Responses:
//		default: genericError
//				200: jobResponse

func (a *{{.NameExported}}App) getJob(w http.ResponseWriter, r *http.Request) {
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

{{.GetSwaggerDoc}}
// getSchedule swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/schedule/{key} schedule getschedule
//
// Returns a schedule given a key, where key is a UUID
//
// Responses:
//		default: genericError
//				200: scheduleResponse

func (a *{{.NameExported}}App) getSchedule(w http.ResponseWriter, r *http.Request) {
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

{{.GetSwaggerDoc}}
// getLiveness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Liveness}} {{.Liveness}} get{{.Liveness}}
//
// A HTTP response status code between 200-400 indicates the pod is alive.
// Any other status code will cause kubelet to restart the pod.
//
// Responses:
//		default: genericError
//				200: livenessResponse

func (a *{{.NameExported}}App) getLiveness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getLivenessPreHook(w, r)

	if !a.Live {
		respondWithError(w, http.StatusServiceUnavailable, "{\"Live\": false}")
	}

	respondWithJSON(w, http.StatusOK, a.Live)
}

{{.GetSwaggerDoc}}
// getReadiness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Readiness}} {{.Readiness}} get{{.Readiness}}
//
// Indicates the pod is ready to start taking traffic.
// Should return a 200 after all pod initialization has completed.
//
// Responses:
//		default: genericError
//				200: readinessResponse

func (a *{{.NameExported}}App) getReadiness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getReadinessPreHook(w, r)

	if !a.Ready {
		respondWithError(w, http.StatusServiceUnavailable, "{\"Ready\": false}")
	}

	respondWithJSON(w, http.StatusOK, a.Ready)
}

{{.GetSwaggerDoc}}
// getMetrics swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Metrics}} {{.Metrics}} getMetrics
//
// Returns metrics for {{.Name}} service
// Metrics should include:
//	 - Scheduler
//	 - Dispatcher
//	 - Workers
//	 - Jobs
//
// Responses:
//		default: genericError
//				200: readinessResponse

func (a *{{.NameExported}}App) getMetrics(w http.ResponseWriter, r *http.Request) {
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
func (a *{{.NameExported}}App) getManagement(w http.ResponseWriter, r *http.Request) {
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
func (a *{{.NameExported}}App) putManagement(w http.ResponseWriter, r *http.Request) {

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
		time.Sleep(time.Duration(a.Dispatcher.conf.gracefulShutdown) * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	// Special case for hard kill
	// We've sent the reply
	if requestedCommand.Command == "shutdown_now" {
		// Give it 1 second to be sent
		time.Sleep(time.Duration(a.Dispatcher.conf.hardShutdown) * time.Second)
		syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}

	return
}

{{.PostSwaggerDoc}}
// createJob swagger:route POST /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}JobsEndPoint {{.NameExported}}JobsEndPoint createJob
//
// Create a new Job
//
// Responses:
//		default: genericError
//				201: jobResponse
//				400: genericError
func (a *{{.NameExported}}App) createJob(w http.ResponseWriter, r *http.Request) {

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

{{.PutSwaggerDoc}}
// updateJob swagger:route PUT /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}JobsEndPoint/{key} {{.NameExported}}SchedulerEndPoint updateJob
//
// Update a {{.NameExported}}JobsEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: jobResponse
//				400: genericError
//				404: genericError
func (a *{{.NameExported}}App) updateJob(w http.ResponseWriter, r *http.Request) {
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

{{.DeleteSwaggerDoc}}
// deleteJob swagger:route DELETE /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}JobsEndPoint/{key} {{.NameExported}}JobsEndPoint deleteJobs
//
// Delete a job specified by key, which is a uuid
//
// Responses:
//		default: genericError
//				200: jobResponse
//				400: genericError
func (a *{{.NameExported}}App) deleteJob(w http.ResponseWriter, r *http.Request) {
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

{{.PostSwaggerDoc}}
// createSchedule swagger:route POST /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint {{.NameExported}}SchedulerEndPoint createSchedule
//
// Create a new scheduler
//
// Responses:
//		default: genericError
//				201: schedulerResponse
//				400: genericError
func (a *{{.NameExported}}App) createSchedule(w http.ResponseWriter, r *http.Request) {

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

	respondWithByte(w, status, respBody)
}

{{.PutSwaggerDoc}}
// updateSchedle swagger:route PUT /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint/{key} {{.NameExported}}SchedulerEndPoint updateSchedule
//
// Update a {{.NameExported}}SchedulerEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: schedulerResponse
//				400: genericError
func (a *{{.NameExported}}App) updateSchedule(w http.ResponseWriter, r *http.Request) {
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

{{.DeleteSwaggerDoc}}
// deleteSchedule swagger:route DELETE /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint/{key} {{.NameExported}}SchedulerEndPoint deleteSchudler
//
// Delete a job specified by key, which is a uuid
//
// Responses:
//		default: genericError
//				200: schedulerResponse
//				400: genericError
func (a *{{.NameExported}}App) deleteSchedule(w http.ResponseWriter, r *http.Request) {
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
	w.Write(payload)
}

// respondWithJSON will Marshal the payload
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func openAccessLogFile(accesslogfile string) *os.File {
	var lf *os.File
	var err error

	if accesslogfile == "" {
		accesslogfile = "access.log"
		log.Println("Access log file name not declared using errors.log")
	}

	rollLogIfExists(accesslogfile)

	lf, err = os.OpenFile(accesslogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

	if err != nil {
		log.Fatal("Error opening access log file: os.OpenFile:", err)
		return nil
	}

	return lf
}

func openErrorLogFile(errorlogfile string) {
	if errorlogfile == "" {
		errorlogfile = "errors.log"
		log.Println("Error log file name not declared using errors.log")
	}

	rollLogIfExists(errorlogfile)

	lf, err := os.OpenFile(errorlogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0640)

	if err != nil {
		log.Fatal("Error opening error log file: os.OpenFile:", err)
	}
	log.SetOutput(lf)
}

func rollLogIfExists(logfilename string) {

	if _, err := os.Stat(logfilename); os.IsNotExist(err) {
		return
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
		log.Printf("Rename logfile %v to %v failed with error %v\n",
			logfilename, newFileName, err.Error())
	}

	return

}{{end}}
