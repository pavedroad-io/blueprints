
//
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root for full license information.
//

// Allocate jobs to workers in a pool

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gorilla/mux"
)

// Contants to build up a k8s style URL
const (
	// EventCollectorAPIVersion Version API URL
	EventCollectorAPIVersion string = "/api/v1"

	// EventCollectorNamespaceID Prefix for namespaces
	EventCollectorNamespaceID string = "namespace"

	// EventCollectorDefaultNamespace Default namespace
	EventCollectorDefaultNamespace string = "mirantis"

	// EventCollectorResourceType CRD Type per k8s
	EventCollectorResourceType string = "eventCollector"

	// The email or account login used by 3rd parth provider
	EventCollectorKey string = "/{key}"

	// EventCollectorLivenessEndPoint
	EventCollectorLivenessEndPoint string = "liveness"

	// EventCollectorReadinessEndPoint
	EventCollectorReadinessEndPoint string = "ready"

	// EventCollectorMetricsEndPoint
	EventCollectorMetricsEndPoint string = "metrics"

	// EventCollectorManagementEndPoint
	EventCollectorManagementEndPoint string = "management"

	// EventCollectorJobsEndPoint
	EventCollectorJobsEndPoint string = "jobs"

	// EventCollectorSchedulerEndPoint
	EventCollectorSchedulerEndPoint string = "scheduler"
)

// EventCollectorApp Top level construct containing building blockes
// for this microservice
type EventCollectorApp struct {
	// Router http request router, gorilla mux for this app
	Router *mux.Router

	// Dispatcher manages jobs for workers
	Dispatcher dispatcher

	// Scheduler creates and forwards jobs to dispatcher
	Scheduler  Scheduler

	// Live http server is start
	Live bool

	// Ready once dispatcher has complete initialization
	Ready bool

	httpInterruptChan chan os.Signal

	// Logs
	accessLog *os.File
}

// HTTP server configuration
type httpConfig struct {
	ip							string
	port						string
	shutdownTimeout time.Duration
	readTimeout			time.Duration
	writeTimeout		time.Duration
	listenString		string
	logPath					string
	diagnosticsFile string
	accessFile      string
}

// Set default http configuration
var httpconf = httpConfig{ip: "127.0.0.1", port: "8081", shutdownTimeout: 15, readTimeout: 60, writeTimeout: 60, listenString: "127.0.0.1:8081", logPath: "logs/", diagnosticsFile: "diagnostics.log", accessFile: "access.log"}

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

// main entry point for server
func main() {
	a := EventCollectorApp{}

	versionFlag := flag.Bool("v", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		printVersion()
	}

	// Setup logging
	openErrorLogFile(httpconf.logPath + httpconf.diagnosticsFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	a.accessLog = openAccessLogFile(httpconf.logPath + httpconf.accessFile)

	a.Initialize()
	a.Run(httpconf.listenString)
}