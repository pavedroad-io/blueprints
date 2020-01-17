{{define "templateApp.go"}}{{.PavedroadInfo}}

// User project / copyright / usage information
// {{.ProjectInfo}}

package main

import (
  "context"
  "database/sql"
  "encoding/json"
  "fmt"
  "github.com/gorilla/mux"
  _ "github.com/lib/pq"
  "io/ioutil"
  "log"
  "net/http"
  "os"
  "strconv"
  "time"
)

// Initialize setups database connection object and the http server
//
func (a *{{.NameExported}}App) Initialize() {

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

// Run starts the server
//
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
  if  err := srv.Shutdown(ctx); err != http.ErrServerClosed {
	  log.Printf("HTTP server shut down: %v", err)
  }
  log.Println("shutting down")
  os.Exit(0)
}

// Get for ennvironment variable overrides
func (a *{{.NameExported}}App) initializeEnvironment() {
  var envVar = ""

  //Look for environment variables overrides
  //You should set the evironment variables as needed
  //Defaults are hard coded in the Main package
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

{{.AllRoutesSwaggerDoc}}
func (a *{{.NameExported}}App) initializeRoutes() {
  uri := {{.NameExported}}APIVersion + "/" + {{.NameExported}}NamespaceID + "/{namespace}/" +
    {{.NameExported}}ResourceType + "LIST"
  a.Router.HandleFunc(uri, a.list{{.NameExported}}).Methods("GET")

  uri = {{.NameExported}}APIVersion + "/" + {{.NameExported}}NamespaceID + "/{namespace}/" +
    {{.NameExported}}ResourceType + "/{key}"
  a.Router.HandleFunc(uri, a.get{{.NameExported}}).Methods("GET")

  uri = {{.NameExported}}APIVersion + "/" + {{.NameExported}}NamespaceID + "/{namespace}/" + {{.NameExported}}ResourceType
  a.Router.HandleFunc(uri, a.create{{.NameExported}}).Methods("POST")

  uri = {{.NameExported}}APIVersion + "/" + {{.NameExported}}NamespaceID + "/{namespace}/" +
    {{.NameExported}}ResourceType + {{.NameExported}}Key
  a.Router.HandleFunc(uri, a.update{{.NameExported}}).Methods("PUT")

  uri = {{.NameExported}}APIVersion + "/" + {{.NameExported}}NamespaceID + "/{namespace}/" +
    {{.NameExported}}ResourceType + {{.NameExported}}Key
  a.Router.HandleFunc(uri, a.delete{{.NameExported}}).Methods("DELETE")
}

{{.GetAllSwaggerDoc}}
// list{{.NameExported}} swagger:route GET /api/v1/namespace/pavedroad.io/{{.Name}}LIST {{.Name}} list{{.Name}}
//
// Returns a list of {{.Name}}
//
// Responses:
//    default: genericError
//        200: {{.Name}}List

func (a *{{.NameExported}}App) list{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
  {{.Name}} := {{.Name}}{}

  count, _ := strconv.Atoi(r.FormValue("count"))
  start, _ := strconv.Atoi(r.FormValue("start"))

  if count > 10 || count < 1 {
    count = 10
  }
  if start < 0 {
    start = 0
  }

  // Pre-processing hook
  a.list{{.NameExported}}PreHook(w, r, count, start)

  mappings, err := {{.Name}}.list{{.NameExported}}(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

	// Post-processing hook
  a.list{{.NameExported}}PostHook(w, r)

  respondWithJSON(w, http.StatusOK, mappings)
}

{{.GetSwaggerDoc}}
// get{{.NameExported}} swagger:route GET /api/v1/namespace/pavedroad.io/{{.Name}}/{uuid} {{.Name}} get{{.Name}}
//
// Returns a {{.Name}} given a key, where key is a UUID
//
// Responses:
//    default: genericError
//        200: {{.Name}}Response

func (a *{{.NameExported}}App) get{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  {{.Name}} := {{.Name}}{}
	key := vars["key"]

	// Pre-processing hook
  a.get{{.NameExported}}PreHook(w, r, key)

  //TODO: allows them to specify the column used to retrieve user
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

  // Pre-processing hook
  a.get{{.NameExported}}PostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.PostSwaggerDoc}}
// create{{.NameExported}} swagger:route POST /api/v1/namespace/pavedroad.io/{{.Name}} {{.Name}} create{{.Name}}
//
// Create a new {{.Name}}
//
// Responses:
//    default: genericError
//        201: {{.Name}}Response
//        400: genericError
func (a *{{.NameExported}}App) create{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
  // New map structure
  {{.Name}} := {{.Name}}{}

  // Pre-processing hook
  a.create{{.NameExported}}PreHook(w, r)

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

 // Post-processing hook
  a.create{{.NameExported}}PostHook(w, r)

  respondWithJSON(w, http.StatusCreated, {{.Name}})
}

{{.PutSwaggerDoc}}
// update{{.NameExported}} swagger:route PUT /api/v1/namespace/pavedroad.io/{{.Name}}/{key} {{.Name}} update{{.Name}}
//
// Update a {{.Name}} specified by key, where key is a uuid
//
// Responses:
//    default: genericError
//        201: {{.Name}}Response
//        400: genericError
func (a *{{.NameExported}}App) update{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
  {{.Name}} := {{.Name}}{}

  // Read URI variables
  vars := mux.Vars(r)
  key := vars["key"]

  // Pre-processing hook
  a.update{{.NameExported}}PreHook(w, r, key)

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

  // Post-processing hook
  a.update{{.NameExported}}PostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, {{.Name}})
}

{{.DeleteSwaggerDoc}}
// delete{{.NameExported}} swagger:route DELETE /api/v1/namespace/pavedroad.io/{{.Name}}/{key} {{.Name}} delete{{.Name}}
//
// Update a {{.Name}} specified by key, which is a uuid
//
// Responses:
//    default: genericError
//        200: {{.Name}}Response
//        400: genericError
func (a *{{.NameExported}}App) delete{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
  {{.Name}} := {{.Name}}{}
  vars := mux.Vars(r)
	key := vars["key"]

  // Pre-processing hook
  a.delete{{.NameExported}}PreHook(w, r, key)

  err := {{.Name}}.delete{{.NameExported}}(a.DB, key)
  if err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
    return
  }

  // Post-processing hook
  a.delete{{.NameExported}}PostHook(w, r, key)

  respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
  respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
  response, _ := json.Marshal(payload)

  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(code)
  _, err := w.Write(response)
  if err != nil {
	  log.Printf("Response errror: %s", err)
  }
}

func logRequest(handler http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
    handler.ServeHTTP(w, r)
  })
}

func openLogFile(logfile string) {
  if logfile != "" {
    lf, err := os.OpenFile(logfile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0600)

    if err != nil {
      log.Fatal("OpenLogfile: os.OpenFile:", err)
    }
    log.SetOutput(lf)
  }
}

/*
func dump{{.NameExported}}(m {{.NameExported}}) {
  fmt.Println("Dump {{.Name}}")
  {{.DumpStructs}}
}
*/{{end}}
