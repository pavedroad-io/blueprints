{{define "templateMain.go"}}{{.PavedroadInfo}}

// {{.ProjectInfo}}

package main

import (
       "flag"
        "fmt"
        "github.com/gorilla/mux"
        "log"
        "net/http"
        "os"
        "reflect"
        "strconv"
        "strings"
        "sync"
        "time"
)

// Constants to build up a k8s style URL
const (
	// {{.NameExported}}APIVersion Version API URL
	APIVersion string = "/api/{{.APIVersion}}"

	// {{.NameExported}}NamespaceID Prefix for namespaces
	NamespaceID string = "namespace"

	// {{.NameExported}}DefaultNamespace Default namespace
	DefaultNamespace string = "{{.Namespace}}"

	// {{.NameExported}}ResourceType CRD Type per k8s
	ResourceType string = "{{.Name}}"

	// The email or account login used by 3rd party provider
	Key string = "/{key}"

	// {{.NameExported}}LivenessEndPoint
	LivenessEndPoint string = "{{.Liveness}}"

	// {{.NameExported}}ReadinessEndPoint
	ReadinessEndPoint string = "{{.Readiness}}"

	// {{.NameExported}}MetricsEndPoint
	MetricsEndPoint string = "{{.Metrics}}"

	// EventCollectorManagementEndPoint
	ManagementEndPoint string = "management"

	// {{.NameExported}}JobsEndPoint
	JobsEndPoint string = "jobs"

	// {{.NameExported}}SchedulerEndPoint
	SchedulerEndPoint string = "scheduler"
)

var (
        //Resources management Option
        // managementOptions   managementCommand

        // shutdownTimeout will be initialized based on the
        //default or HTTP_SHUTDOWN_TIMEOUT
        shutdowTimeout time.Duration

        // GitTag contains current git tab for this repository
        GitTag string

	       // Version contains version specified in definitions file
        Version string

        // Build holds latest git commit hash in short form
        Build string

        //Mwg is used to support go routines
        Mwg sync.WaitGroup

        //Mmutex is the main control on critical code sections
        Mmutex *sync.Mutex

        //ChannelsReady is used to manage channels blocking
        ChannelsReady bool

        //WorkersDown is a Temp fix to allow worker start and stop
        WorkersDown int
)

// {{.NameExported}}App Top level construct containing building blocks
// for this micro service
type {{.NameExported}}App struct {
	// Router http request router, gorilla mux for this app
	Router *mux.Router

	// Dispatcher manages jobs for workers
	Dispatcher dispatcher

	// Scheduler creates and forwards jobs to dispatcher
	Scheduler Scheduler

	// Live http server is start
	LiveHTTPSever bool

	// Ready once dispatcher has complete initialization
	DispatcherReady bool

	httpInterruptChan chan os.Signal

	// Logs
	accessLog *os.File

	//Resources management Option
        managementOptions managementCommands
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
var httpconf = httpConfig{
        ip:              "127.0.0.1",
        port:            "8081",
        shutdownTimeout: 15 * time.Second,
        readTimeout:     time.Minute,
        writeTimeout:    time.Minute,
        listenString:    "127.0.0.1:8081",
        logPath:         "logs/",
        diagnosticsFile: "diagnostics.log",
        accessFile:      "access.log",
}

// Options
const (
        runConfig     string = "config"     //managementRequest.CommandType
        runCommand    string = "command"    //managementRequest.CommandType
        resDispatcher string = "Dispatcher" //Resource Type
        resScheduler  string = "Scheduler"  //Resource Type
        resWorker     string = "Worker"     //Resource Type
        resJob        string = "Job"        //Resource Type
        MaxCmdParms   int    = 2            //Maximum parameters in runCommand
)

var errorInterface interface{}

// CmdParm has parametr information for the functions in managementCommands
type CmdParm struct {
        //Go parametr data type
        DataType string `json:"data_type"`
        //Function parameter order
        ParamOrder int `json:"param_order"`
}

// ConfigValue information for the keyed Field in managementCommands
//
type ConfigValue struct {
        //Go data type
        DataType string `json:"data_type"`

        //Should be Resourse type option
        Resource string `json:"resource"`

        //Function to set the keyed field
        funcName interface{} `json:"func_name"`
}

//MgtCommand has information CommandNames maps in managementCommands
type MgtCommand struct {
        // Description of what this command does
        Description string `json:"description"`

        //Indicate service command will impact.
        //Should be from  Resource type option
        Resource string `json:"resource"`

        //Function to process the command
        funcName interface{} `json:"func_name"`

      //Values to set for command
        CmdParms map[string]CmdParm
}

// managementCommands is a list of available commanads and field options
//
// swagger:response managementCommands
type managementCommands struct {
        // The error message
        // in: body

        // CommandNamess is a map of valid command names that can be executed
        // Map key is the name of the command
        CommandNames map[string]MgtCommand `json:"commands"`

       // Fields is a map of fields that can be changed
        // CommandNames and Fields are not related
        //But Fields can be checked when a 'config' command type is called.
        //Map key is the name of the field that can be set.
        //value and type for the field is in ConfigValue
        Fields map[string]ConfigValue `json:"fields"`
}

// get404Response Not found
//
// swagger:response get404Response
type get404Response struct {
        // The 404 error message
        // in: body

        Body struct {
                // Error message
                Error string `json:"error"`

                // UUID / key that was not found
                UUID string `json:"uuid"`
        }
}

// genericError
//
// swagger:response genericError
type genericError struct {
        // in: body
        // Error message
        Body struct {
                Error string `json:"error"`
        } `json:"body"`
}


// genericResponse
//
// swagger:response genericResponse
type genericResponse struct {
        // in: body
        Body struct {
                // JSON body
                JSONBody string `json:"json_body"`
        } `json:"body"`
}



//
// swagger:response metricsResponse
type metricsResponse struct {
        // in: body
        Body struct {
                // Error message
                SchedulerMetrics string `json:"scheduler_metrics"`
                DispatherMetrics string `json:"dispather_metrics"`
        } `json:"body"`
}

func (mc *managementCommands) Init() {
        mc.CommandNames = make(map[string]MgtCommand)
        mc.Fields = make(map[string]ConfigValue)
}

//Only processing for functions with a specific return
//signature. Cannot be used if panic must be avoided.
//func (mc *managementCommands) callFunc(fn interface{}, params ...interface{}) (result []reflect.Value) {
func (mc *managementCommands) callFunc(fn interface{}, params ...interface{}) (httpStatusCode int, jsonb []byte, err error) {
        //result:= make([]reflect.Value,3)
        f := reflect.ValueOf(fn)
        //ft := f.Elem()
        if f.Kind() != reflect.Func {
                log.Printf("Function interface expected but %v was provide. \n", f.Kind())
                panic("Unexpected function call.")
        }

        if f.Type().NumIn() != len(params) {
                //TODO: How to print function name?
                log.Printf("Incorrect parametres list for function call.")
                panic("Incorrect number of function parameters")
        }

        //Three sepcific return parameters expected
        if f.Type().NumOut() != 3 {
                log.Printf("Operation has specific function return requirements.")

               panic("Incorrect number of expected function return parameters")
        }

        if f.Type().Out(0).Kind() != reflect.Int {
                log.Printf("Must return a httpStatusCode")
                panic("Incorrect type for expected httpStatusCode.")
        }
        //Specific check for slice of bytes can be done.
        if f.Type().Out(1).Kind() != reflect.Slice {
               log.Printf("Must return a byte slice payload")

                panic("Incorrect type for expected payload.")
        }
        //if !f.Type().Out(2).Implements(errorInterface) {
        if f.Type().Out(2).Kind() != reflect.Interface {
                log.Printf("Must return an error ")
                panic("Incorrect, error return expecetd.")
        }

        inputs := make([]reflect.Value, len(params))
       for k, in := range params {
                inputs[k] = reflect.ValueOf(in)
        }
        result := f.Call(inputs)
        return int(result[0].Int()), result[1].Bytes(), nil
}

//FieldParm is the information for Fields in managementCommands
type FieldParm struct {
        // Value for field
        FieldValue string `json:"field_value"`

        // Golang Datatype for field
        // Supported types not specified, must review
        DataType string `json:"data_type"`
}


// managementRequest user request to execute a management command
// swagger:request managementRequest
type managementRequest struct {
        // The error message
        // in: body

        // CommandName is a valid Name that can be executed
        // from a.managementOptions
        // Can be blank/nil if CommandType is config
        // Must be available if CommandType is runCommand
        CommandName string `json:"command"`
        //Fields to set if CommandType is config
        // Must have Fileds[0] if config
        // Can have 0 or mutiple fields if runCommand
        Fields map[string]FieldParm `json:"fields"`

        //Resource the command is acting on
        // Dispatcher, Scheduler, Jobs, Workers
        Resource string `json:"resource"`

        //Type of request: config or command
        CommandType string `json:"command_type"`
}

// useCommand() sets up a new managementOptions.CommandNames[] entry if it does not exist
//Name:  Short name for the command to execute
//Desc:  Description of what the command does
//Resourse: Resourse the command will affect
//ParamsName: Descriptive ordered name list for the parameters to the command
//  i.e. paramsName[0] is the descriptive name of the first command.
//  paramsName of length 0 indicates that there are no parameters
//ValPtr:  Function to call to execute the command
func (mc *managementCommands) useCommand(name, desc, resourse string, paramsName []string, valPtr interface{}) {
        var oldCmd MgtCommand
        var exist bool


        if valPtr == nil {
                log.Printf("Command %s needs a processing function. \n", name)
                return

        }

        f := reflect.ValueOf(valPtr)
        if f.Kind() != reflect.Func {
                log.Printf("Command %s needs a processing function. \n", name)
                return

        }

       aLen := f.Type().NumIn()
        if aLen > MaxCmdParms {
                log.Printf("Function supporting command: %s has more parameters than expecetd. \n", name)
                return

        }

        Mmutex.Lock()
        defer Mmutex.Unlock()
        oldCmd, exist = mc.CommandNames[name]
        if !exist {
                pLen := len(paramsName)

                if pLen != aLen {
                        log.Printf("Command %s should have %v parameter name(s). \n", name, aLen)
                        return
                }

                newParm := make(map[string]CmdParm)
                if aLen > 0 {

                        for i := 0; i < pLen; i++ {
                               newCmdParm := CmdParm{
                                        DataType:   f.Type().In(i).Name(),
                                        ParamOrder: i,
                                }
                                newParm[paramsName[i]] = newCmdParm
                        }
                }

                newCMD := MgtCommand{
                        Description: desc,
                        Resource:    resourse,
                        funcName:    valPtr,
                       CmdParms:    newParm,
                }

                mc.CommandNames[name] = newCMD
        } else {
                log.Printf("Command key taken for: %s with: %+v \n", name, oldCmd)

        }
}

// setField() sets up a new managementOptions.Field[] if it does not exist
// A simplified version of useCommand.
//Name:  Name of resourse option to change.
//DataType: The Go type of the Name (string or int are only type
//curenntly implemented)
//Resourse: Resourse the command will affect (see resResourse constants)
//ValPtr:  Function to call to process the change

func (mc *managementCommands) setField(name, datatype, resourse string, valPtr interface{}) {

        Mmutex.Lock()
       defer Mmutex.Unlock()

        oldVal, exist := mc.Fields[name]

        if !exist {
                newVal := ConfigValue{
                        DataType: datatype,
                        Resource: resourse,
                        funcName: valPtr,
                }

                mc.Fields[name] = newVal
        } else {
                log.Printf("Filed key taken for: %s with: %+v \n", name, oldVal)

        }
}

/*
  cmdProcessManagementRequest only supports options
  from managementOptions.CommandNames[]
*/

func (mc *managementCommands) cmdProcessManagementRequest(r *managementRequest) (httpStatusCode int, jsonb []byte, err error) {
        name := r.CommandName
        msg := ""
        i := len(name)
        if i == 0 {
                msg = fmt.Sprintf("{\"Status\": \"Command name is required.\"}")
               return http.StatusOK, []byte(msg), nil

        }

        reqCmd, isSup := mc.CommandNames[name]

        if !isSup {

                msg = fmt.Sprintf("{\"Status\": \"Command %s is not supported.\"}", r.CommandName)
                return http.StatusNotFound, []byte(msg), nil
        }
      i = len(reqCmd.CmdParms)

        if i != len(r.Fields) {
                msg = fmt.Sprintf("{\"Status\": \"Command %s requires %v parameters.\"}", r.CommandName, i)
                return http.StatusBadRequest, []byte(msg), nil

        }
        execFunc := reqCmd.funcName
        if reflect.ValueOf(execFunc).IsNil() {
                msg = fmt.Sprintf("{\"status\": \"unexpected resourse : %v .\"}", name)

              return http.StatusInternalServerError, []byte(msg), nil
        }

        //Just a cmmand without any parameters
        if i == 0 {
                return mc.callFunc(execFunc)

        }
        //Now process some limited supported parameter

        //keys := make([]string, i)
      parms := make([]interface{}, i)

        for m := range reqCmd.CmdParms {
                //Only int and string currently supported
                if reqCmd.CmdParms[m].DataType != r.Fields[m].DataType {
                        //Params datatype should match
                        msg := fmt.Sprintf("{\"Status\": \"Unexpected data type : %v .\"}", r.Fields[m].DataType)
                        return http.StatusBadRequest, []byte(msg), nil

                }

                switch strings.ToUpper(reqCmd.CmdParms[m].DataType) {
               case "STRING":
                        parms[reqCmd.CmdParms[m].ParamOrder] = r.Fields[m].FieldValue
                case "INT":
                        //TODO: Convert fldVal to int
                        int6FldVal, err := strconv.ParseInt(r.Fields[m].FieldValue, 0, 16)
                        if err != nil {
                                msg := fmt.Sprintf("{\"Status\": \"Unexpected data value : %v .\"}", err)
                                return http.StatusBadRequest, []byte(msg), nil
                       }
                        intFldVal := int(int6FldVal)
                        intStrVal := strconv.Itoa(intFldVal)
                        if intStrVal != r.Fields[m].FieldValue {
                                msg := fmt.Sprintf("{\"Status\": \"Unexpected int value : %v .\"}", int6FldVal)
                                return http.StatusBadRequest, []byte(msg), nil
                        }
                        parms[reqCmd.CmdParms[m].ParamOrder] = intFldVal

                default:
                        msg := fmt.Sprintf("{\"Status\": \"Not a supported data type : %v .\"}", r.Fields[m].DataType)
                       return http.StatusBadRequest, []byte(msg), nil

                        //switch over datatype
                }
                //Range over keys
        }

        return mc.callFunc(execFunc, parms)
}

/*
  configProcessManagementRequest only supports options
  from managementOptions.Fields[]
*/
func (mc *managementCommands) configProcessManagementRequest(r *managementRequest) (httpStatusCode int, jsonb []byte, err error) {
        i := len(r.Fields)
        msg := ""
        if i != 1 {
                msg = fmt.Sprintf("{\"status\": \"command is implemented for only one resource name.\"}")
                return http.StatusBadRequest, []byte(msg), nil
       }
        keys := make([]string, i)
        i = 0
        for k := range r.Fields {
                keys[i] = k
                i++
        }
        field := keys[0]
        reqCfig, isSup := mc.Fields[field]
        if !isSup {
                msg = fmt.Sprintf("{\"status\": \"not a supported resourse : %v .\"}", field)
               return http.StatusNotFound, []byte(msg), nil
        }
        execFunc := reqCfig.funcName
        if reflect.ValueOf(execFunc).IsNil() {
                msg = fmt.Sprintf("{\"status\": \"unexpected resourse : %v .\"}", field)
                return http.StatusInternalServerError, []byte(msg), nil
        }
        fldVal := r.Fields[field].FieldValue
        fldType := mc.Fields[field].DataType

	/*
           Note:
           keys[0]
           field to set
           mc.managementOptions.Fields[keys[0]].FuncName
           function to use to set the field
           r.FIELDS[

        */

       if strings.ToUpper(fldType) == "STRING" {
                return mc.callFunc(execFunc, field, fldVal)
        }
        if strings.ToUpper(fldType) == "INT" {
                int6FldVal, err := strconv.ParseInt(fldVal, 0, 16)
                if err != nil {
                        msg = fmt.Sprintf("{\"Status\": \"Unexpected data type : %v .\"}", err)
                        return http.StatusBadRequest, []byte(msg), nil

                }
                intFldVal := int(int6FldVal)
                intStrVal := strconv.Itoa(intFldVal)

               if intStrVal != fldVal {
                        msg = fmt.Sprintf("{\"Status\": \"Unexpected int value : %v .\"}", int6FldVal)
                        return http.StatusBadRequest, []byte(msg), nil
                }

                return mc.callFunc(execFunc, field, intFldVal)
        }

        msg = fmt.Sprintf("{\"Status\": \"Command is not implemented for field type : %v .\"}", fldType)
        return http.StatusBadRequest, []byte(msg), nil
}


// ProcessManagementRequest executes a management command
func (mc *managementCommands) ProcessManagementRequest(r managementRequest) (httpStatusCode int, jsonb []byte, err error) {

        switch r.CommandType {
        case runConfig:
                return mc.configProcessManagementRequest(&r)

        case runCommand:
                return mc.cmdProcessManagementRequest(&r)

        default:
                msg := fmt.Sprintf("{\"Status\": \"Command: %s not implemeted\"}",
                        r.CommandName)
                return http.StatusBadRequest, []byte(msg), nil
        }

}


// printVersion
func printVersion() {
	fmt.Printf("{\"Version\": \"%v\", \"Build\": \"%v\", \"GitTag\": \"%v\"}\n", Version, Build, GitTag)
}

func printError(em error) {
        fmt.Println(em)
}


// main entry point for server
func main() {
	Mmutex = new(sync.Mutex)
	a := {{.NameExported}}App{}
        WorkersDown = 0


	versionFlag := flag.Bool("v", false, "Print version information")
	flag.Parse()

	if *versionFlag {
		printVersion()
		os.Exit(0)
	}

	// Setup logging
	e := openErrorLogFile(httpconf.logPath + httpconf.diagnosticsFile)
	if e != nil {
                printError(e)
                os.Exit(0)
        }

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	a.Initialize()
	Mwg.Add(1)
	a.Run(httpconf.listenString)
	Mwg.Wait()
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
