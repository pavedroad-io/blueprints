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

	// Start rest end points
	httpconf.listenString = fmt.Sprintf("%s:%s", httpconf.ip, httpconf.port)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
	a.Ready = true
}

// Run start the HTTP server for Rest endpoints
func (a *{{.NameExported}}App) Run(addr string) {

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

{{.AllRoutesSwaggerDoc}}
func (a *{{.NameExported}}App) initializeRoutes() {

	uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
				{{.NameExported}}ResourceType + "/" +
				{{.NameExported}}LivenessEndPoint
	a.Router.HandleFunc(uri, a.getLiveness).Methods("GET")
	log.Println("GET: ", uri)

	uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
				{{.NameExported}}ResourceType + "/" +
				{{.NameExported}}ReadinessEndPoint
	a.Router.HandleFunc(uri, a.getReadiness).Methods("GET")
	log.Println("GET: ", uri)

	uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
				{{.NameExported}}ResourceType + "/" +
				{{.NameExported}}MetricsEndPoint
	a.Router.HandleFunc(uri, a.getMetrics).Methods("GET")
	log.Println("GET: ", uri)

	uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
				{{.NameExported}}ResourceType + "/" +
				{{.NameExported}}ManagementEndPoint
	a.Router.HandleFunc(uri, a.getManagement).Methods("GET")
	log.Println("GET: ", uri)

	uri = {{.NameExported}}APIVersion + "/" +
				{{.NameExported}}NamespaceID + "/" +
				{{.NameExported}}DefaultNamespace + "/" +
				{{.NameExported}}ResourceType + "/" +
				{{.NameExported}}ManagementEndPoint
	a.Router.HandleFunc(uri, a.putManagement).Methods("PUT")
	log.Println("PUT: ", uri)

// Generated routes
{{.ServiceRoutes}}

	return
}

{{.GetAllSwaggerDoc}}

// Generated methods
{{.ServiceMethods}}

// End generated routes

// getLiveness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}/{{.Liveness}} {{.Liveness}} get{{.Liveness}}
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

	if !a.Live {
		respondWithByte(w, http.StatusServiceUnavailable, []byte("{\"Live\": false}"))
		return
	}

	respondWithByte(w, http.StatusOK, []byte("{\"Live\": true}"))
}

{{.GetSwaggerDoc}}
// getReadiness swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}/{{.Readiness}} {{.Readiness}} get{{.Readiness}}
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

	if !a.Ready {
		respondWithByte(w, http.StatusServiceUnavailable, []byte("{\"Ready\": false}"))
		return
	}

	respondWithByte(w, http.StatusOK, []byte("{\"Ready\": true}"))
}

{{.ExplainSwaggerDoc}}
// getExplain swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}/{{.Explain}} {{.Explain}} get{{.Explain}}
//
// Return JSON/Swagger documentation for using this APi
// If the pre-processing hook is not defined return not-found
// Otherwise, return a swagger API document
//
// Responses:
//		default: genericError
//				200: genericResponse

func (a *{{.NameExported}}App) getExplain(w http.ResponseWriter, r *http.Request) {
	// Pre-processing hook
	if body, err := getExplainPreHook(w,r); err != nil {
		respondWithByte(w, http.StatusOK, []byte("{\"not-found\": true}"))
	} else {
		respondWithByte(w, http.StatusOK, body)
	}
}

{{.GetSwaggerDoc}}
// getMetrics swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}/{{.Metrics}} {{.Metrics}} getMetrics
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

// getManagement swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}/{{.Management}} {{.Management}} get{{.Management}}
//
// Returns available management commands
//
// Responses:
//		default: genericError
//				200: managementGetResponse
func (a *{{.NameExported}}App) getManagement(w http.ResponseWriter, r *http.Request) {
	// Pre-processing hook
	getManagementPreHook(w, r)

	// Post-processing hook
	getManagementPostHook(w, r)

	respondWithJSON(w, http.StatusOK, a.Dispatcher.managementOptions)
}

// put{{.Management}} swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}/{{.Management}} {{.Management}} put{{.Management}}
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

}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
