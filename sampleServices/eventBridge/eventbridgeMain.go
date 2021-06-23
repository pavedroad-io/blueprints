
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
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

// Constants to build up a k8s style URL
const (
	// EventbridgeAPIVersion Version API URL
	EventbridgeAPIVersion string = "/api/v1"

	// EventbridgeNamespaceID Prefix for namespaces
	EventbridgeNamespaceID string = "namespace"

	// EventbridgeDefaultNamespace Default namespace
	EventbridgeDefaultNamespace string = "pavedroad"

	// EventbridgeResourceType CRD Type per k8s
	EventbridgeResourceType string = "eventbridge"

	// The email or account login used by 3rd party provider
	EventbridgeKey string = "/{key}"

	// EventbridgeLivenessEndPoint
	EventbridgeLivenessEndPoint string = "liveness"

	// EventbridgeReadinessEndPoint
	EventbridgeReadinessEndPoint string = "ready"

	// EventbridgeMetricsEndPoint
	EventbridgeMetricsEndPoint string = "metrics"

	// EventbridgeManagementEndPoint
	EventbridgeManagementEndPoint string = "management"

	// EventbridgeJobsEndPoint
	EventbridgeJobsEndPoint string = "jobs"

	// EventbridgeSchedulerEndPoint
	EventbridgeSchedulerEndPoint string = "scheduler"
)

// EventbridgeApp Top level construct containing building blocks
// for this micro service
type EventbridgeApp struct {
	// Router http request router, gorilla mux for this app
	Router *mux.Router

	// Dispatcher manages jobs for workers
	Dispatcher dispatcher

	// Scheduler creates and forwards jobs to dispatcher
	Scheduler Scheduler

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

//printError
func printError(em error) {
        fmt.Println(em)
}

// main entry point for server
func main() {
	a := EventbridgeApp{}

	versionFlag := flag.Bool("v", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		printVersion()
	}

	// Setup logging
        e := openErrorLogFile(httpconf.logPath + httpconf.diagnosticsFile)
	if e != nil {
                printError(e)
                os.Exit(0)
        }
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	a.accessLog = openAccessLogFile(httpconf.logPath + httpconf.accessFile)

	a.Initialize()
	a.Run(httpconf.listenString)
}