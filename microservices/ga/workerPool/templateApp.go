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
  "time"

  "github.com/gorilla/mux"
  _ "github.com/lib/pq"
)

// Initialize setups database connection object and the http server
//
func (a *{{.NameExported}}App) Initialize() {

  // Override defaults
  a.initializeEnvironment()

	// Start the Dispatcher
	// TODO: make workers and channel buffer size configurable
  a.Dispatcher.Init(5, 5, &a.Scheduler)
  go a.Dispatcher.Run()

  // Scheduler
  err := a.Scheduler.Init()
  if err != nil {
    fmt.Println(err)
    os.Exit(-1)
  }
  go a.Scheduler.Run()

	// Start rest end points
  httpconf.listenString = fmt.Sprintf("%s:%s", httpconf.ip, httpconf.port)
  a.Router = mux.NewRouter()
  a.initializeRoutes()
}

// Start the server
func (a *{{.NameExported}}App) Run(addr string) {

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

  // Listen for SIGHUP
  c := make(chan os.Signal, 1)
  <-c

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

{{.AllRoutesSwaggerDoc}}
func (a *{{.NameExported}}App) initializeRoutes() {

  uri := {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorJobsEndPoint + "LIST"
  a.Router.HandleFunc(uri, a.listJobs).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorSchedulerEndPoint + "LIST"
  a.Router.HandleFunc(uri, a.listSchedule).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorJobsEndPoint + EventCollectorKey
  a.Router.HandleFunc(uri, a.getJob).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorSchedulerEndPoint + EventCollectorKey
  a.Router.HandleFunc(uri, a.getSchedule).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorLivenessEndPoint
  a.Router.HandleFunc(uri, a.getLiveness).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorReadinessEndPoint
  a.Router.HandleFunc(uri, a.getReadiness).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorMetricsEndPoint
  a.Router.HandleFunc(uri, a.getMetrics).Methods("GET")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorJobsEndPoint + EventCollectorKey
  a.Router.HandleFunc(uri, a.updateJob).Methods("PUT")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorJobsEndPoint + EventCollectorKey
  a.Router.HandleFunc(uri, a.deleteJob).Methods("DELETE")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorJobsEndPoint
  a.Router.HandleFunc(uri, a.createJob).Methods("POST")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorSchedulerEndPoint + EventCollectorKey
  a.Router.HandleFunc(uri, a.updateSchedule).Methods("PUT")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorSchedulerEndPoint + EventCollectorKey
  a.Router.HandleFunc(uri, a.deleteSchedule).Methods("DELETE")
	fmt.Println(uri)

  uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
    		{{.NameExported}}ResourceType + "/" +
				EventCollectorSchedulerEndPoint
  a.Router.HandleFunc(uri, a.createSchedule).Methods("POST")
	fmt.Println(uri)

	return
}

{{.GetAllSwaggerDoc}}
// listJobs swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/EventCollectorJobsEndPointLIST jobs listjobs
//
// Returns a list of Jobs
//
// Responses:
//    default: genericError
//        200: jobsList

func (a *{{.NameExported}}App) listJobs(w http.ResponseWriter, r *http.Request) {
  //{{.Name}} := {{.Name}}{}

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

	/*
	New logic here
  mappings, err := {{.Name}}.list{{.NameExported}}(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }
	*/

	// Post-processing hook
  listJobsPostHook(w, r)

  respondWithJSON(w, http.StatusOK, "{}")
}

{{.GetAllSwaggerDoc}}
// listSchedule swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/EventCollectorSchedulerEndPointLIST schedules listschedule
//
// Returns a list of schedules
//
// Responses:
//    default: genericError
//        200: scheduleList

func (a *{{.NameExported}}App) listSchedule(w http.ResponseWriter, r *http.Request) {
  //{{.Name}} := {{.Name}}{}

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
	New logic here
  mappings, err := {{.Name}}.list{{.NameExported}}(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }
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
//    default: genericError
//        200: jobResponse

func (a *{{.NameExported}}App) getJob(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  //{{.Name}} := {{.Name}}{}
	key := vars["key"]

	// Pre-processing hook
  getJobPreHook(w, r, key)

	/*
	New logic here
  err := {{.Name}}.get{{.NameExported}}(a.DB, key, UUID)

  if err != nil {
    errmsg := err.Error()
    errno :=  errmsg[0:3]
    if errno == "400" {
      respondWithError(w, http.StatusBadRequest, err.Error())
    } else {
      respondWithError(w, http.StatusNotFound, err.Error())
    }
    return
  }
	*/

  // Pre-processing hook
  getJobPostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.GetSwaggerDoc}}
// getSchedule swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/schedule/{key} schedule getschedule
//
// Returns a schedule given a key, where key is a UUID
//
// Responses:
//    default: genericError
//        200: scheduleResponse

func (a *{{.NameExported}}App) getSchedule(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  //{{.Name}} := {{.Name}}{}
	key := vars["key"]

	// Pre-processing hook
  getSchedulePreHook(w, r, key)

	/*
	New logic here
  err := {{.Name}}.get{{.NameExported}}(a.DB, key, UUID)

  if err != nil {
    errmsg := err.Error()
    errno :=  errmsg[0:3]
    if errno == "400" {
      respondWithError(w, http.StatusBadRequest, err.Error())
    } else {
      respondWithError(w, http.StatusNotFound, err.Error())
    }
    return
  }
	*/

  // Pre-processing hook
  getSchedulePostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, "{}")
}

{{.GetSwaggerDoc}}
// getLiveness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Liveness}} {{.Liveness}} get{{.Liveness}}
//
// A HTTP response status code between 200-400 indicates the pod is alive.
// Any other status code will cause kubelet to restart the pod.
//
// Responses:
//    default: genericError
//        200: livenessResponse

func (a *{{.NameExported}}App) getLiveness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
  getLivenessPreHook(w, r)

	/*
	New logic here
	*/

  // Pre-processing hook
  getLivenessPostHook(w, r)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.GetSwaggerDoc}}
// getReadiness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Readiness}} {{.Readiness}} get{{.Readiness}}
//
// Indicates the pod is ready to start taking traffic.
// Should return a 200 after all pod initialization has completed.
//
// Responses:
//    default: genericError
//        200: readinessResponse

func (a *{{.NameExported}}App) getReadiness(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
  getReadinessPreHook(w, r)

	/*
	New logic here
	*/

  // Pre-processing hook
  getReadinessPostHook(w, r)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.GetSwaggerDoc}}
// getMetrics swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.Metrics}} {{.Metrics}} getMetrics
//
// Returns metrics for {{.Name}} service
// Metrics should include:
//   - Scheduler
//   - Dispatcher
//   - Workers
//   - Jobs
//
// Responses:
//    default: genericError
//        200: readinessResponse

func (a *{{.NameExported}}App) getMetrics(w http.ResponseWriter, r *http.Request) {

	// Pre-processing hook
  getMetricsPreHook(w, r)

	/*
	New logic here
	*/

  // Pre-processing hook
  getMetricsPostHook(w, r)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.PostSwaggerDoc}}
// createJob swagger:route POST /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}JobsEndPoint {{.NameExported}}JobsEndPoint createJob
//
// Create a new Job
//
// Responses:
//    default: genericError
//        201: jobResponse
//        400: genericError
func (a *{{.NameExported}}App) createJob(w http.ResponseWriter, r *http.Request) {
  // New job structure
  //{{.Name}} := {{.Name}}{}

  // Pre-processing hook
  createJobPreHook(w, r)

	/*
  htmlData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }

  err = json.Unmarshal(htmlData, &{{.Name}})
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }

  ct := time.Now().UTC()
  {{.Name}}.Created = ct
  {{.Name}}.Updated = ct

  // Save into backend storage
  // returns the UUID if needed
  if _, err := {{.Name}}.create{{.NameExported}}(a.DB); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

	*/
 // Post-processing hook
  createJobPostHook(w, r)

  respondWithJSON(w, http.StatusCreated, {{.Name}})
}

{{.PutSwaggerDoc}}
// updateJob swagger:route PUT /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}JobsEndPoint/{key} {{.NameExported}}SchedulerEndPoint updateJob
//
// Update a {{.NameExported}}JobsEndPoint specified by key, where key is a uuid
//
// Responses:
//    default: genericError
//        201: jobResponse
//        400: genericError
func (a *{{.NameExported}}App) updateJob(w http.ResponseWriter, r *http.Request) {
	// {{.Name}} := {{.Name}}{}

  // Read URI variables
  vars := mux.Vars(r)
  key := vars["key"]

  // Pre-processing hook
  updateJobPreHook(w, r, key)

	/*
  htmlData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    return
  }

  err = json.Unmarshal(htmlData, &{{.Name}})
  if err != nil {
    log.Println(err)
    return
  }

  ct := time.Now().UTC()
  {{.Name}}.Updated = ct

  if err := {{.Name}}.update{{.NameExported}}(a.DB, {{.Name}}.{{.NameExported}}UUID); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
	*/

  // Post-processing hook
  updateJobPostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.DeleteSwaggerDoc}}
// deleteJob swagger:route DELETE /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}JobsEndPoint/{key} {{.NameExported}}JobsEndPoint deleteJobs
//
// Update a job specified by key, which is a uuid
//
// Responses:
//    default: genericError
//        200: jobResponse
//        400: genericError
func (a *{{.NameExported}}App) deleteJob(w http.ResponseWriter, r *http.Request) {
  //{{.Name}} := {{.Name}}{}
  vars := mux.Vars(r)
	key := vars["key"]

  // Pre-processing hook
  deleteJobPreHook(w, r, key)

	/*
  err := {{.Name}}.delete{{.NameExported}}(a.DB, key)
  if err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
    return
  }
	*/

  // Post-processing hook
  deleteJobPostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

{{.PostSwaggerDoc}}
// createSchedule swagger:route POST /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint {{.NameExported}}SchedulerEndPoint createSchedule
//
// Create a new scheduler
//
// Responses:
//    default: genericError
//        201: schedulerResponse
//        400: genericError
func (a *{{.NameExported}}App) createSchedule(w http.ResponseWriter, r *http.Request) {
  // New scheduler structure
  //{{.Name}} := {{.Name}}{}

  // Pre-processing hook
  createSchedulePreHook(w, r)

	/*
  htmlData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }

  err = json.Unmarshal(htmlData, &{{.Name}})
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }

  ct := time.Now().UTC()
  {{.Name}}.Created = ct
  {{.Name}}.Updated = ct

  // Save into backend storage
  // returns the UUID if needed
  if _, err := {{.Name}}.create{{.NameExported}}(a.DB); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

	*/
 // Post-processing hook
  createSchedulePostHook(w, r)

  respondWithJSON(w, http.StatusCreated, {{.Name}})
}

{{.PutSwaggerDoc}}
// updateSchedle swagger:route PUT /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint/{key} {{.NameExported}}SchedulerEndPoint updateSchedule
//
// Update a {{.NameExported}}SchedulerEndPoint specified by key, where key is a uuid
//
// Responses:
//    default: genericError
//        200: schedulerResponse
//        400: genericError
func (a *{{.NameExported}}App) updateSchedule(w http.ResponseWriter, r *http.Request) {
	// {{.Name}} := {{.Name}}{}

  // Read URI variables
  vars := mux.Vars(r)
  key := vars["key"]

  // Pre-processing hook
  updateSchedulePreHook(w, r, key)

	/*
  htmlData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    return
  }

  err = json.Unmarshal(htmlData, &{{.Name}})
  if err != nil {
    log.Println(err)
    return
  }

  ct := time.Now().UTC()
  {{.Name}}.Updated = ct

  if err := {{.Name}}.update{{.NameExported}}(a.DB, {{.Name}}.{{.NameExported}}UUID); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }
	*/

  // Post-processing hook
  updateSchedulePostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.DeleteSwaggerDoc}}
// deleteSchedule swagger:route DELETE /api/v1/namespace/{{.Namespace}}/{{.Name}}/{{.NameExported}}SchedulerEndPoint/{key} {{.NameExported}}SchedulerEndPoint deleteSchudler
//
// Delete a job specified by key, which is a uuid
//
// Responses:
//    default: genericError
//        200: schedulerResponse
//        400: genericError
func (a *{{.NameExported}}App) deleteSchedule(w http.ResponseWriter, r *http.Request) {
  //{{.Name}} := {{.Name}}{}
  vars := mux.Vars(r)
	key := vars["key"]

  // Pre-processing hook
  deleteSchedulePreHook(w, r, key)

	/*
  err := {{.Name}}.delete{{.NameExported}}(a.DB, key)
  if err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
    return
  }
	*/

  // Post-processing hook
  deleteSchedulePostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  w.Write(response)
}

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
}{{end}}
