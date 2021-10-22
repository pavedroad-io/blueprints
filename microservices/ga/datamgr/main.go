{{define "main.go"}}
{{.PavedroadInfo}}

// User project / copyright / usage information
// {{.ProjectInfo}}

package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"os"
	"os/user"
	"time"

	"github.com/pavedroad-io/go-core/logger"
)

// Constants to build up a k8s style URL
const (
	// {{.NameExported}}APIVersion Version API URL
	{{.NameExported}}APIVersion string = "/api/v1"

	// {{.NameExported}}NamespaceID Prefix for namespaces
	{{.NameExported}}NamespaceID string = "namespace"

	// {{.NameExported}}DefaultNamespace Default namespace
	{{.NameExported}}DefaultNamespace string = "{{.Namespace}}"

	// {{.NameExported}}ResourceType CRD Type per k8s
	{{.NameExported}}ResourceType string = "{{.Name}}"

	// The email or account login used by 3rd party provider
	{{.NameExported}}Key string = "/{key}"

	// {{.NameExported}}LivenessEndPoint
	{{.NameExported}}LivenessEndPoint string = "{{.Liveness}}"

	// {{.NameExported}}ReadinessEndPoint
	{{.NameExported}}ReadinessEndPoint string = "{{.Readiness}}"

	// {{.NameExported}}MetricsEndPoint
	{{.NameExported}}MetricsEndPoint string = "{{.Metrics}}"

	// {{.NameExported}}ManagementEndPoint
	{{.NameExported}}ManagementEndPoint string = "management"



)

// Options for looking up a user
const (
  UUID = iota
  NAME
)

// holds pointers to database and http server
type {{.NameExported}}App struct {
	Router *mux.Router
	DB     *sql.DB
	// Live http server is start
	Live bool

	// Ready once dispatcher has complete initialization
	Ready bool

	httpInterruptChan chan os.Signal

	// Logs
	accessLog *os.File
}

// both db and http configuration can be changed using environment variables
type databaseConfig struct {
  username string
  password string
  database string
  sslMode  string
  dbDriver string
  ip       string
  port     string
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
	diagnosticsFile string
	accessFile      string
}

// Global for use in the module

// Set default database configuration
var dbconf = databaseConfig{username: "root", password: "", database: "pavedroad", sslMode: "disable", dbDriver: "postgres", ip: "127.0.0.1", port: "26257"}

// Set default http configuration
var httpconf = httpConfig{ip: "{{.HTTPHost}}", port: "{{.HTTPPort}}", shutdownTimeout: 15, readTimeout: 60, writeTimeout: 60, listenString: "{{.HTTPHost}}:{{.HTTPPort}}", logPath: "logs", diagnosticsFile: "debug.log", accessFile: "access.log"}

var cuser, _ = user.Current()
var logconfig = logger.LoggerConfiguration{
		LogPackage:        logger.ZapType,
		LogLevel:          logger.InfoType,
		EnableTimeStamps:  true,
		EnableColorLevels: true,
		EnableCloudEvents: true,
		CloudEventsCfg: logger.CloudEventsConfiguration{
			SetID: logger.CEHMAC,
		},
		EnableKafka: false,
		KafkaFormat: logger.CEFormat,
		KafkaProducerCfg: logger.ProducerConfiguration{
			Brokers:       []string{"localhost:9092"},
			Topic:         "logs",
			Partition:     logger.RandomPartition,
			Key:           logger.FixedKey,
			KeyName:       cuser.Username,
			Compression:   logger.CompressionSnappy,
			AckWait:       logger.WaitForLocal,
			ProdFlushFreq: 500, // milliseconds
			EnableTLS:     false,
			EnableDebug:   false,
		},
		EnableConsole: true,
		ConsoleFormat: logger.TextFormat,
		EnableFile:    false,
		FileFormat:    logger.JSONFormat,
		FileLocation:  httpconf.logPath + httpconf.diagnosticsFile,
}


// shutdownTimeout will be initialized based on the default or HTTP_SHUTDOWN_TIMEOUT
var shutdowTimeout time.Duration

// GitTag contains current git tab for this repository
var GitTag string

// Version contains version specified in definitions file
var Version string

// Build holds latest git commit hash in short form
var Build string

// printVersion
func printVersion() {
	fmt.Printf("{\"Version\": \"%v\", \"Build\": \"%v\", \"GitTag\": \"%v\"}\n",
	Version, Build, GitTag)
	os.Exit(0)
}

//printError
func printError(em error) {
        fmt.Println(em)
}

// main entry point for server
func main() {
	a := {{.NameExported}}App{}

	versionFlag := flag.Bool("v", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		printVersion()
	}

	a.Initialize()
	a.Run(httpconf.listenString)
}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
