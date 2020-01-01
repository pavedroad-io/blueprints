{{define "templateMain.go"}}{{.PavedroadInfo}}

// {{.ProjectInfo}}

package main

import (
        "flag"
        "fmt"
        "github.com/gorilla/mux"
        "log"
        "os"
        "time"
)

// Contants to build up a k8s style URL
const (
  // {{.NameExported}}APIVersion Version API URL
  {{.NameExported}}APIVersion string = "/api/{{.APIVersion}}"

  // {{.NameExported}}NamespaceID Prefix for namespaces
  {{.NameExported}}NamespaceID string = "namespace"

  // {{.NameExported}}DefaultNamespace Default namespace
  {{.NameExported}}DefaultNamespace string = "{{.Namespace}}"

  // {{.NameExported}}ResourceType CRD Type per k8s
  {{.NameExported}}ResourceType string = "{{.Name}}"

  // The email or account login used by 3rd parth provider
  {{.NameExported}}Key string = "/{key}"

  // {{.NameExported}}LivenessEndPoint
  {{.NameExported}}LivenessEndPoint string = "{{.Liveness}}"

  // {{.NameExported}}ReadinessEndPoint
  {{.NameExported}}ReadinessEndPoint string = "{{.Readiness}}"

  // {{.NameExported}}MetricsEndPoint
  {{.NameExported}}MetricsEndPoint string = "{{.Metrics}}"

  // {{.NameExported}}JobsEndPoint
  {{.NameExported}}JobsEndPoint string = "jobs"

  // {{.NameExported}}SchedulerEndPoint
  {{.NameExported}}SchedulerEndPoint string = "scheduler"

)

// holds pointers to database and http server
type {{.NameExported}}App struct {
  Router *mux.Router
	Dispatcher dispatcher
	//TODO: read from roadctl
  Scheduler  httpScheduler
}

// HTTP server configuration
type httpConfig struct {
  ip              string
  port            string
  shutdownTimeout time.Duration
  readTimeout     time.Duration
  writeTimeout    time.Duration
  listenString    string
  logPath         string
}

// Set default http configuration
var httpconf = httpConfig{ip: "127.0.0.1", port: "8081", shutdownTimeout: 15, readTimeout: 60, writeTimeout: 60, listenString: "127.0.0.1:8081", logPath: "logs/{{.Name}}.log"}

// shutdownTimeout will be initialized based on the default or HTTP_SHUTDOWN_TIMEOUT
var shutdowTimeout time.Duration

var GitTag string
var Version string
var Build string

// printVersion
func printVersion() {
        fmt.Printf("{\"Version\": \"%v\", \"Build\": \"%v\", \"GitTag\": \"%v\"}\n",
                Version, Build, GitTag)
        os.Exit(0)
}

// main entry point for server
func main() {

	versionFlag := flag.Bool("v", false, "Print version information")
        flag.Parse()

        if *versionFlag {
                printVersion()
        }

  // Setup loggin
  openLogFile(httpconf.logPath)
  log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
  log.Printf("Logfile opened %s", httpconf.logPath)

  a := {{.NameExported}}App{}
  a.Initialize()
  a.Run(httpconf.listenString)
}{{end}}
