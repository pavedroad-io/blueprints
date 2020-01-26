//
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root for full license information.
//

// User project / copyright / usage information
// Manage database of films

package main

import (
	"context"
	"database/sql"
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
func (a *FilmsApp) Initialize() {

	// Override defaults
	a.initializeEnvironment()

	// Build connection strings
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s host=%s port=%s",
		dbconf.username,
		dbconf.password,
		dbconf.database,
		dbconf.sslMode,
		dbconf.ip,
		dbconf.port)

	httpconf.listenString = fmt.Sprintf("%s:%s", httpconf.ip, httpconf.port)

	var err error
	a.DB, err = sql.Open(dbconf.dbDriver, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Start the server
func (a *FilmsApp) Run(addr string) {

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
func (a *FilmsApp) initializeEnvironment() {
	var envVar = ""

	//look for environment variables overrides
	envVar = os.Getenv("APP_DB_USERNAME")
	if envVar != "" {
		dbconf.username = envVar
	}

	envVar = os.Getenv("APP_DB_PASSWORD")
	if envVar != "" {
		dbconf.password = envVar
	}

	envVar = os.Getenv("APP_DB_NAME")
	if envVar != "" {
		dbconf.database = envVar
	}
	envVar = os.Getenv("APP_DB_SSL_MODE")
	if envVar != "" {
		dbconf.sslMode = envVar
	}

	envVar = os.Getenv("APP_DB_SQL_DRIVER")
	if envVar != "" {
		dbconf.dbDriver = envVar
	}

	envVar = os.Getenv("APP_DB_IP")
	if envVar != "" {
		dbconf.ip = envVar
	}

	envVar = os.Getenv("APP_DB_PORT")
	if envVar != "" {
		dbconf.port = envVar
	}

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

func (a *FilmsApp) initializeRoutes() {
	uri := FilmsAPIVersion + "/" + FilmsNamespaceID + "/{namespace}/" +
		FilmsResourceType + "LIST"
	a.Router.HandleFunc(uri, a.listFilms).Methods("GET")

	uri = FilmsAPIVersion + "/" + FilmsNamespaceID + "/{namespace}/" +
		FilmsResourceType + "/{key}"
	a.Router.HandleFunc(uri, a.getFilms).Methods("GET")

	uri = FilmsAPIVersion + "/" + FilmsNamespaceID + "/{namespace}/" + FilmsResourceType
	a.Router.HandleFunc(uri, a.createFilms).Methods("POST")

	uri = FilmsAPIVersion + "/" + FilmsNamespaceID + "/{namespace}/" +
		FilmsResourceType + FilmsKey
	a.Router.HandleFunc(uri, a.updateFilms).Methods("PUT")

	uri = FilmsAPIVersion + "/" + FilmsNamespaceID + "/{namespace}/" +
		FilmsResourceType + FilmsKey
	a.Router.HandleFunc(uri, a.deleteFilms).Methods("DELETE")
}

// listFilms swagger:route GET /api/v1/namespace/pavedroad.io/filmsLIST films listfilms
//
// Returns a list of films
//
// Responses:
//    default: genericError
//        200: filmsList

func (a *FilmsApp) listFilms(w http.ResponseWriter, r *http.Request) {
	films := films{}

	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	// Pre-processing hook
	a.listFilmsPreHook(w, r, count, start)

	mappings, err := films.listFilms(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Post-processing hook
	a.listFilmsPostHook(w, r)

	respondWithJSON(w, http.StatusOK, mappings)
}

// getFilms swagger:route GET /api/v1/namespace/pavedroad.io/films/{uuid} films getfilms
//
// Returns a films given a key, where key is a UUID
//
// Responses:
//    default: genericError
//        200: filmsResponse

func (a *FilmsApp) getFilms(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	films := films{}
	key := vars["key"]

	// Pre-processing hook
	a.getFilmsPreHook(w, r, key)

	//TODO: allows them to specify the column used to retrieve user
	err := films.getFilms(a.DB, key, UUID)

	if err != nil {
		errmsg := err.Error()
		errno := errmsg[0:3]
		if errno == "400" {
			respondWithError(w, http.StatusBadRequest, err.Error())
		} else {
			respondWithError(w, http.StatusNotFound, err.Error())
		}
		return
	}

	// Pre-processing hook
	a.getFilmsPostHook(w, r, key)

	respondWithJSON(w, http.StatusOK, films)
}

// createFilms swagger:route POST /api/v1/namespace/pavedroad.io/films films createfilms
//
// Create a new films
//
// Responses:
//    default: genericError
//        201: filmsResponse
//        400: genericError
func (a *FilmsApp) createFilms(w http.ResponseWriter, r *http.Request) {
	// New map structure
	films := films{}

	// Pre-processing hook
	a.createFilmsPreHook(w, r)

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	err = json.Unmarshal(htmlData, &films)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	ct := time.Now().UTC()
	films.Created = ct
	films.Updated = ct

	// Save into backend storage
	// returns the UUID if needed
	if _, err := films.createFilms(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Post-processing hook
	a.createFilmsPostHook(w, r)

	respondWithJSON(w, http.StatusCreated, films)
}

// updateFilms swagger:route PUT /api/v1/namespace/pavedroad.io/films/{key} films updatefilms
//
// Update a films specified by key, where key is a uuid
//
// Responses:
//    default: genericError
//        201: filmsResponse
//        400: genericError
func (a *FilmsApp) updateFilms(w http.ResponseWriter, r *http.Request) {
	films := films{}

	// Read URI variables
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	a.updateFilmsPreHook(w, r, key)

	htmlData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(htmlData, &films)
	if err != nil {
		log.Println(err)
		return
	}

	ct := time.Now().UTC()
	films.Updated = ct

	if err := films.updateFilms(a.DB, films.FilmsUUID); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Post-processing hook
	a.updateFilmsPostHook(w, r, key)

	respondWithJSON(w, http.StatusOK, films)
}

// deleteFilms swagger:route DELETE /api/v1/namespace/pavedroad.io/films/{key} films deletefilms
//
// Update a films specified by key, which is a uuid
//
// Responses:
//    default: genericError
//        200: filmsResponse
//        400: genericError
func (a *FilmsApp) deleteFilms(w http.ResponseWriter, r *http.Request) {
	films := films{}
	vars := mux.Vars(r)
	key := vars["key"]

	// Pre-processing hook
	a.deleteFilmsPreHook(w, r, key)

	err := films.deleteFilms(a.DB, key)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	// Post-processing hook
	a.deleteFilmsPostHook(w, r, key)

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
}

/*
func dumpFilms(m Films) {
  fmt.Println("Dump films")

}
*/
