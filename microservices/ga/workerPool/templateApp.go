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

// Initialize setups the scheduler,dispatcher, the 
// go routines, and the http server
func (a *{{.NameExported}}App) Initialize() {

        a.Router = mux.NewRouter()
        a.accessLog = openAccessLogFile(httpconf.logPath + httpconf.accessFile)
  
        a.managementOptions.Init()
        // set k8s probe
        a.LiveHTTPSever = false
        a.DispatcherReady = false
  
        // Override defaults
        a.initializeEnvironment()

	// Start the scheduler
	// Scheduler interface required for dispatcher.
        // Only interface methods, see templateScheduler,
        // available to a.Scheduler
	a.Scheduler = &{{.SchedulerName}}{}

	dConf := &dispatcherConfiguration{
		scheduler:           a.Scheduler,
		sizeOfJobChannel:    SizeOfJobChannel,
		sizeOfResultChannel: SizeOfResultChannel,
		numberOfWorkers:     NumberOfWorkers,
		gracefulShutdown:    GracefullShutdown,
		hardShutdown:        HardShutdown,
	}

	// Scheduler first
        err := a.Scheduler.Init(&a.managementOptions)
        if err != nil {
                fmt.Println(err)
                os.Exit(-1)
        }

	 //Dispatcher next, scheduler need channels
        a.Dispatcher.Init(dConf, &a.managementOptions)
        Mwg.Add(1)
        go a.Dispatcher.Run(&Mwg)

	//New: was not present in my test.
	// Start rest end points
//        httpconf.listenString = fmt.Sprintf("%s:%s", httpconf.ip, httpconf.port)

        //Make sure routes are initialize before running Scheduler
        a.initializeRoutes()

        //Now run sceduler
        Mwg.Add(1)
        go a.Scheduler.Run(&Mwg)

        a.DispatcherReady = true
	Mmutex.Lock()
        ChannelsReady = true
	Mmutex.Unlock()
}

// Run start the HTTP server for Rest endpoints
func (a *{{.NameExported}}App) Run(addr string) {

	defer Mwg.Done()
	log.Println("Listing at: " + addr)
	// Wrap router with W3C logging

	loggedRouter := handlers.LoggingHandler(a.accessLog, a.Router)
	srv := &http.Server{
		Handler:      loggedRouter,
		Addr:         addr,
		WriteTimeout: httpconf.writeTimeout * time.Second,
		ReadTimeout:  httpconf.readTimeout * time.Second,
	}

	Mwg.Add(1)

	go func() {
		defer Mwg.Done()
		//if err := srv.ListenAndServe(); err != nil {
		//	log.Println(err)
		//}
		log.Fatal(srv.ListenAndServe())
                fmt.Println("After ListenAndSerever ")
	}()

        a.LiveHTTPSever = true

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
func (a *{{.NameExported}}App) initializeEnvironment() {
	var envVar = ""
       //string environment variables map
        varStrToProc := map[string]*string{
                "HTTP_IP_ADDR": &httpconf.ip,
                "HTTP_IP_PORT": &httpconf.port,
                "HTTP_LOG":     &httpconf.logPath,
        }

        //time environment variables map
        varTimeToProc := map[string]*time.Duration{
                "HTTP_READ_TIMEOUT":     &httpconf.readTimeout,
                "HTTP_WRITE_TIMEOUT":    &httpconf.writeTimeout,
                "HTTP_SHUTDOWN_TIMEOUT": &httpconf.shutdownTimeout,
        }

	       //Environment variable list to process
        //Expand with case statement
        varsToProc := [][]string{
                []string{"HTTP_IP_ADDR", "string"},
                []string{"HTTP_IP_PORT", "string"},
                []string{"HTTP_READ_TIMEOUT", "time"},
                []string{"HTTP_WRITE_TIMEOUT", "time"},
                []string{"HTTP_SHUTDOWN_TIMEOUT", "time"},
                []string{"HTTP_LOG", "string"},
        }

                for i := 0; i < len(varsToProc); i++ {
                envVar = os.Getenv(varsToProc[i][0])
                fmt.Println(varsToProc[i][0])
                fmt.Println(envVar)
                if envVar != "" {
                        switch varTyp := varsToProc[i][1]; varTyp {
                        case "string":
                                *varStrToProc[varsToProc[i][0]] = envVar
                        case "time":
                                to, err := strconv.Atoi(envVar)
                                if err != nil {
                                        log.Printf("Failed to convert: %v, to int: %v", varsToProc[i][0], envVar)
                                } else {
                                        //toTm := (time.Duration(to) * time.Second)
                                        chV := false
                                        if to < 0 {
                                                to = to * -1
                                                chV = true
                                        }
                                        if to > MaxGraceSec {
                                                to = MaxGraceSec
                                                chV = true
                                        }
                                        if chV {
                                                envVar = strconv.Itoa(to)
                                        }
                                        // *varTimeToProc[varsToProc[i][0]] = (time.Duration(to) * time.Second)
                                        // toTm,_ := time.ParseDuration(envVar + "s")
                                        //fmt.Println(toTm.Seconds())
                                        *varTimeToProc[varsToProc[i][0]], _ = time.ParseDuration(envVar + "s")
                                }
                        default:
                                log.Printf("Env. variable type  %v, not supported.", varTyp)

                        }
                }
        }
}

{{.AllRoutesSwaggerDoc}}
func (a *{{.NameExported}}App) initializeRoutes() {

      var uriPrefix string

        uriPrefix = APIVersion + "/" + NamespaceID + "/" +
                DefaultNamespace + "/" + ResourceType + "/"

        uri := uriPrefix +
                JobsEndPoint + "LIST"

        a.Router.HandleFunc(uri, a.listJobs).Methods("GET")
        log.Println("Get: ", uri)

        uri = uriPrefix +
                SchedulerEndPoint + "LIST"
       a.Router.HandleFunc(uri, a.listSchedule).Methods("GET")
        log.Println("Get: ", uri)

        uri = uriPrefix +
                JobsEndPoint + Key
        a.Router.HandleFunc(uri, a.getJob).Methods("GET")
        log.Println("Get: ", uri)

        uri = uriPrefix +
                SchedulerEndPoint
        a.Router.HandleFunc(uri, a.getSchedule).Methods("GET")
        log.Println("Get: ", uri)

	       uri = uriPrefix +
                LivenessEndPoint
        a.Router.HandleFunc(uri, a.getLiveness).Methods("GET")
        log.Println("Get: ", uri)

        uri = uriPrefix +
                ReadinessEndPoint
        a.Router.HandleFunc(uri, a.getReadiness).Methods("GET")
        log.Println("Get: ", uri)

        uri = uriPrefix +
                MetricsEndPoint
	       a.Router.HandleFunc(uri, a.getMetrics).Methods("GET")
        log.Println("Get: ", uri)

        //uri same for next getManagement,putManagement
        uri = uriPrefix +
                ManagementEndPoint
        a.Router.HandleFunc(uri, a.getManagement).Methods("GET")
        log.Println("GET: ", uri)

        a.Router.HandleFunc(uri, a.putManagement).Methods("PUT")
        log.Println("PUT: ", uri)

        uri = uriPrefix +
               JobsEndPoint
        a.Router.HandleFunc(uri, a.createJob).Methods("POST")
        log.Println("POST :", uri)

        //uri same for deleteJob to follow
        uri = uriPrefix +
                JobsEndPoint + Key
        a.Router.HandleFunc(uri, a.updateJob).Methods("PUT")
        log.Println("PUT: ", uri)

        //uri same as above updateJob
        a.Router.HandleFunc(uri, a.deleteJob).Methods("DELETE")
        log.Println("DELETE: ", uri)

        uri = uriPrefix +
                SchedulerEndPoint
        a.Router.HandleFunc(uri, a.createSchedule).Methods("POST")
        log.Println("POST:", uri)

        //uri same for createSchedule
        a.Router.HandleFunc(uri, a.updateSchedule).Methods("PUT")
        log.Println("PUT :", uri)

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
//				200: listJobResponse
//        		500: genericError

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

// TODO: decide do kill it or do something with it
func (a *{{.NameExported}}App) listSchedule(w http.ResponseWriter, r *http.Request) {
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
//				200: listJobResponse
//				404: get404Response
//				500: genericError
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
//		200: genericResponse
//		500: genericError
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
//				200: genericResponse
//				503: genericError
func (a *{{.NameExported}}App) getLiveness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getLivenessPreHook(w, r)

	if !a.LiveHTTPSever {
		respondWithByte(w, http.StatusServiceUnavailable, []byte("{\"Live\": false}"))
		return
	}

	respondWithByte(w, http.StatusOK, []byte("{\"Live\": true}"))
}

{{.GetSwaggerDoc}}
// getReadiness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Readiness}} {{.Readiness}} get{{.Readiness}}
//
// Indicates the pod is ready to start taking traffic.
// Should return a 200 after all pod initialization has completed.
//
// Responses:
//		default: genericError
//				200: genericResponse
//				503: genericError

func (a *{{.NameExported}}App) getReadiness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
	getReadinessPreHook(w, r)

	if !a.DispatcherReady {
		respondWithByte(w, http.StatusServiceUnavailable, []byte("{\"Ready\": false}"))
		return
	}

	respondWithByte(w, http.StatusOK, []byte("{\"Ready\": true}"))
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
//				200: metricsResponse

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

// getManagement swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Management}} {{.Management}} get{{.Management}}
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: managementCommands 
func (a *{{.NameExported}}App) getManagement(w http.ResponseWriter, r *http.Request) {
	// Pre-processing hook
	getManagementPreHook(w, r)

	// Post-processing hook
	getManagementPostHook(w, r)

	respondWithJSON(w, http.StatusOK, a.managementOptions)
}

// put{{.Management}} swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Management}} {{.Management}} put{{.Management}}
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: genericResponse
//				400: genericError
//				500: genericError
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
		return
	}

	status, respBody, e := a.managementOptions.ProcessManagementRequest(requestedCommand)

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
	if requestedCommand.CommandName == "shutdown"  && requestedCommand.Resource == resDispatcher{
		//Await for a gracefulShutdown
		time.Sleep(time.Duration(a.Dispatcher.conf.gracefulShutdown) * time.Second)
		e = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if e != nil {
                        msg := fmt.Sprintf("{\"error\": \"Shutting down with\", \"Error\": \"%v\"}", e.Error())
                        log.Println("shutdown error:", msg)
                }
	}

	// Special case for hard kill
	// We've sent the reply
	if requestedCommand.CommandName == "shutdown_now" && requestedCommand.Resource == resDispatcher {
		time.Sleep(time.Duration(a.Dispatcher.conf.hardShutdown) * time.Second)
		e = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
                 if e != nil {
                        msg := fmt.Sprintf("{\"error\": \"Shutting down with\", \"Error\": \"%v\"}", e.Error())
                        log.Println("shutdown_now error:", msg)
                }

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
//				201: listJobResponse
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
//				200: listJobResponse
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
//				200: listJobResponse
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
//				201: genericResponse
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
	// We are using update that will return a 200, change status to a 201
	respondWithByte(w, http.StatusCreated, respBody)
}

{{.PutSwaggerDoc}}
// updateSchedle swagger:route PUT /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint/{key} {{.NameExported}}SchedulerEndPoint updateSchedule
//
// Update a {{.NameExported}}SchedulerEndPoint specified by key, where key is a uuid
//
// Responses:
//		default: genericError
//				200: genericResponse
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
//				200: genericResponse
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
	_,e := w.Write(payload)
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
        _, e :=	w.Write(response)
	if e != nil {
                msg := fmt.Sprintf("{\"error\": \"response error: \", \"Error\": \"%v\"}", e.Error())
                log.Println("respondWithJSON error:", msg)
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
		if err2 := lf.Close(); err2 != nil {
                        log.Println("Error closing bad access log file: os.OpenFile:", err2)
                        return nil
                } else {
                        lf, err = os.OpenFile(accesslogfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)
                        return lf
	}
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

}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
