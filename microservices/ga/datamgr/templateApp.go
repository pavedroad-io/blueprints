{{define "templateApp.go"}}
// Pavedroad license / copyright information
{{.pavedroad-info}}

// User project / copyright / usage information
{{.project-info}}

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
func (a *{{.name}}App) Initialize() {

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
func (a *{{.name}}App) Run(addr string) {

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
func (a *prUserIdMapperApp) initializeEnvironment() {
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

{{.all-routes-swagger-doc}}
func (a *{{.name}}App) initializeRoutes() {
  uri := {{.name}}APIVersion + "/" + {{.name}}NamespaceID + "/{namespace}/" +
    {{.name}}ResourceType + "LIST"
  a.Router.HandleFunc(uri, a.get{{.name-exported}}s).Methods("GET")

  uri = {{.name}}APIVersion + "/" + {{.name}}NamespaceID + "/{namespace}/" +
    {{.name}}ResourceType + "/{key}"
  a.Router.HandleFunc(uri, a.get{{.name-exported}}).Methods("GET")

  uri = {{.name}}APIVersion + "/" + {{.name}}NamespaceID + "/{namespace}/" + {{.name}}ResourceType
  a.Router.HandleFunc(uri, a.create{{.name-exported}}).Methods("POST")

  uri = {{.name}}APIVersion + "/" + {{.name}}NamespaceID + "/{namespace}/" +
    {{.name}}ResourceType + {{.name-exported}}Key
  a.Router.HandleFunc(uri, a.update{{.name-exported}}).Methods("PUT")

  uri = {{.name}}APIVersion + "/" + {{.name}}NamespaceID + "/{namespace}/" +
    {{.name}}ResourceType + {{.name-exported}}Key
  a.Router.HandleFunc(uri, a.delete{{.name-exported}}).Methods("DELETE")
}

{{.get-all-swagger-doc}}
func (a *{{.name}}App) get{{.-name-exported}}s(w http.ResponseWriter, r *http.Request) {
  {{.-name-exported}} := {{.name}}{}

  //vars := mux.Vars(r)
  //fmt.Println("list tokens: ", vars)

  count, _ := strconv.Atoi(r.FormValue("count"))
  start, _ := strconv.Atoi(r.FormValue("start"))

  if count > 10 || count < 1 {
    count = 10
  }
  if start < 0 {
    start = 0
  }

  mappings, err := {{.-name-exported}}.get{{.-name-exported}}s(a.DB, start, count)
  if err != nil {
    respondWithError(w, http.StatusInternalServerError, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, mappings)
}

{{.get-swagger-doc}}
func (a *{{.name-exported}}App) get{{.name-exported}}(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r)
  {{.name-exported}} := {{.name-exported}}{}

  err := {{.name-exported}}.get{{.name-exported}}(a.DB, vars["key"])
  if err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
    return
  }

  respondWithJSON(w, http.StatusOK, {{.name-exported}})
}

{{.post-swagger-doc}}
func (a *{{.name-exported}}App) create{{.name-exported}}(w http.ResponseWriter, r *http.Request) {
  // New map structure
  {{.name}} := {{.name-exported}}{}

  htmlData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }

  err = json.Unmarshal(htmlData, &{{.name}})
  if err != nil {
    log.Println(err)
    os.Exit(1)
  }

  ct := time.Now().UTC()
  {{.name}}.Created = ct.Format(time.RFC3339)
  {{.name}}.Updated = ct.Format(time.RFC3339)
  {{.name}}.LoginCount = 1

  // Save into backend storage
  if err := {{.name}}.create{{.name-exported}}(a.DB); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  respondWithJSON(w, http.StatusCreated, {{.name}})
}

{{.put-swagger-doc}}
func (a *{{.name-exported}}App) update{{.name-exported}}(w http.ResponseWriter, r *http.Request) {
  {{.name}} := {{.name-exported}}{}

  // Read URI variables
  // vars := mux.Vars(r)

  htmlData, err := ioutil.ReadAll(r.Body)
  if err != nil {
    log.Println(err)
    return
  }

  err = json.Unmarshal(htmlData, &{{.name}})
  if err != nil {
    log.Println(err)
    return
  }

  ct := time.Now().UTC()
  {{.name}}.Updated = ct.Format(time.RFC3339)
  {{.name}}.LoginCount += 1

  if err := {{.name}}.update{{.name-exported}}(a.DB); err != nil {
    respondWithError(w, http.StatusBadRequest, "Invalid request payload")
    return
  }

  respondWithJSON(w, http.StatusOK, {{.name}})
}

{{.delete-swagger-doc}}
func (a *{{.name-exported}}App) delete{{.name-exported}}(w http.ResponseWriter, r *http.Request) {
  {{.name-exported}} := {{.name-exported}}{}
  vars := mux.Vars(r)

  err := {{.name-exported}}.delete{{.name-exported}}(a.DB, vars["key"])
  if err != nil {
    respondWithError(w, http.StatusNotFound, err.Error())
    return
  }

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

func dump{{.name-exported}}(m {{.name-exported}}) {
  fmt.Println("Dump {{.name}}")
  {{.dump-structs}}
}
{{end}}

