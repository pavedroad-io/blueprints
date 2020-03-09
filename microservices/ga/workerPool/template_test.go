{{define "template_test.go"}}
// {{.Name}}_test.go

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

const (
	Updated         string = "updated"
	Created         string = "created"
	Active          string = "active"
	Namespace     string = "{{.Namespace}}"
	Service       string = "{{.Name}}"
	ManagementURL string = "/api/v1/namespace/" + Namespace + "/" + Service + "/{{.Management}}"
	JobURL        string = "/api/v1/namespace/" + Namespace + "/" + Service + "/jobs"
	JobListURL    string = "/api/v1/namespace/" + Namespace + "/" + Service + "/jobsLIST"
	ScheduleURL   string = "/api/v1/namespace/" + Namespace + "/" + Service + "/scheduler"
	ReadyURL      string = "/api/v1/namespace/" + Namespace + "/" + Service + "/{{.Readiness}}"
	LiveURL       string = "/api/v1/namespace/" + Namespace + "/" + Service + "/{{.Liveness}}"
	MetricsURL    string = "/api/v1/namespace/" + Namespace + "/" + Service + "/{{.Metrics}}"
)

var mo managementCommands
var mr managementRequest
var fp FieldParm

var new{{.NameExported}}JSON=`{{.PostJSON}}`

var a {{.NameExported}}App

/*
   *******
   testSampleCode should be defaulted to false.
   Only change to true if you fully understand the
   potential impact.
   *******
*/
var testSampleCode = false

func TestMain(m *testing.M) {
        Mmutex = new(sync.Mutex)
	a = {{.NameExported}}App{}

        mo = managementCommands{}
        mr = managementRequest{}
        fp = FieldParm{}

	      // Setup logging
        e := openErrorLogFile(httpconf.logPath + httpconf.diagnosticsFile)
        if e != nil {
                printError(e)
                os.Exit(0)
        }

        log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	//JobTestMain()
	a.Initialize()

	Mwg.Add(1)
	go a.Run(httpconf.listenString)

	// Wait for server to start
	time.Sleep(5 * time.Second)
	JobTestMain()
	defer os.Exit(m.Run())
	defer testExit()
	Mwg.Wait()
}

func testExit() {
        mr.CommandName = "shutdown_now"
        mr.Resource = resDispatcher
        mr.CommandType = "command"
        mr.Fields = make(map[string]FieldParm)

        mrd, e := json.Marshal(&mr)
        if e != nil {
                fmt.Println("Couldn't set command for shutdown_now")
        }
        mrds := string(mrd)

	req, _ := http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
	_ = executeRequest(req)
}

func TestSchedule(t *testing.T) {
	// Get the list of jobs
	req, _ := http.NewRequest("GET", ScheduleURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code,"Job List")

	// This code requires knowledge about the type of scheduler
	// So make it conditional
	if testSampleCode {
		putdata := "{\"schedule_type\": \"Constant interval scheduler\", \"send_interval_seconds\": 30}"
		req, _ = http.NewRequest("PUT", ScheduleURL, strings.NewReader(putdata))
		response = executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code,"interval_seconds 30")

		req, _ = http.NewRequest("DELETE", ScheduleURL, strings.NewReader(putdata))
		response = executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code,"delete shedule")

		postdata := "{\"schedule_type\": \"Constant interval scheduler\", \"send_interval_seconds\": 5}"
		req, _ = http.NewRequest("POST", ScheduleURL, strings.NewReader(postdata))
		response = executeRequest(req)
		checkResponseCode(t, http.StatusCreated, response.Code,"interval_seconds 5")
	}
	return
}

func TestReady(t *testing.T) {

	if !a.DispatcherReady {
		a.DispatcherReady = true
	}

	expect := "true"

	req, _ := http.NewRequest("GET", ReadyURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code,"dispatcher ready")
	expect2 := strings.ToLower(response.Body.String())
	if !strings.Contains(expect2, expect) {
		t.Errorf("dispatcher ready expected %s; Got %v\n", expect, expect2)
	}

	a.DispatcherReady = false
	expect = "false"

	req, _ = http.NewRequest("GET", ReadyURL, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusServiceUnavailable, response.Code,"dispatcher not ready")

 expect2 = strings.ToLower(response.Body.String())
        if !strings.Contains(expect2, expect) {
                t.Errorf("dispatcher not ready expected %s; Got %v\n", expect, expect2)
        }

	// Reset to proper status
	a.DispatcherReady = true
	return
}

func TestLive(t *testing.T) {
       if !a.LiveHTTPSever {
                a.LiveHTTPSever = true
        }

        expect := "true"
        req, _ := http.NewRequest("GET", LiveURL, nil)
        response := executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "live HTTPServer")

        expect2 := strings.ToLower(response.Body.String())
        if !strings.Contains(expect2, expect) {
               t.Errorf("live HTTPServer expected %s; Got %v\n", expect, expect2)
        }

        a.LiveHTTPSever = false
        expect = "false"

        req, _ = http.NewRequest("GET", LiveURL, nil)
        response = executeRequest(req)
        checkResponseCode(t, http.StatusServiceUnavailable, response.Code, "not live HTTPServer")
       expect2 = strings.ToLower(response.Body.String())
        if !strings.Contains(expect2, expect) {
                t.Errorf("not live HTTPServer expected %s; Got %v\n", expect, expect2)
        }

        // Reset to proper status
        a.LiveHTTPSever = true
        return
}

func TestMetric(t *testing.T) {
        expect1 := "scheduler"
        expect2 := "dispatcher"

        req, _ := http.NewRequest("GET", MetricsURL, nil)
        response := executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "metrics")

        expect3 := strings.ToLower(response.Body.String())
        if !strings.Contains(expect3, expect1) {
                t.Errorf("metrics: expected %s; Got %v\n", expect1, expect3)
        }
        if !strings.Contains(expect3, expect2) {
               t.Errorf("metrics: expected %s; Got %v\n", expect2, expect3)
        }

}

func TestOpenAccessLog(t *testing.T) {
	// Make it set default name
	f := openAccessLogFile("")
	if f == nil {
		t.Errorf("Expected *os.Files; Got %v\n", f)
	}
	_ = os.Remove("access.log")

	f = openAccessLogFile("/no permissions")
	if f != nil {
		t.Errorf("Expected nil; Got %v\n", f)
	}

	f = openAccessLogFile("/tmp/test.log")
	if f == nil {
		t.Errorf("Expected type *os.File; Got nil\n")
	}
	_ = os.Remove("/tmp/test.log")
}

func TestOpenErrorLog(t *testing.T) {
	// Make it set default name
	expect := "error.log"
	_ = openErrorLogFile("")
	if _, err := os.Stat(expect); os.IsNotExist(err) {
		t.Errorf("Expected default name error.log; Got %v\n", err)
	}

	_ = os.Remove(expect)

	e := openErrorLogFile("/no permissions")
	if e == nil {
		t.Errorf("Expected %v; Got nil\n", e)
	}

	expect = "/tmp/test.log"
	e = openErrorLogFile(expect)
	if e != nil {
		t.Errorf("Expected type nil; Got err %v\n", e)
	}
	_ = os.Remove(expect)
}

func TestRollLog(t *testing.T) {
	testFile := "/tmp/foo.log"
	_, e := os.Create(testFile)
	if e != nil {
		t.Errorf("Failed to create %v; Got err %v\n", testFile, e)
	}

	f, e := rollLogIfExists(testFile)
	if e != nil {
		t.Errorf("Failed to roll log %v; Got err %v\n", testFile, e)
	}
	_ = os.Remove(testFile)

	if f == testFile {
		t.Errorf("Expected new file name; Got %v\n", f)
	}
	_ = os.Remove(f)

	return
}

func TestInitEnv(t *testing.T) {
       //Do not change in the middle of a test.
        //These values are set to intialize the system
        //changes are not anticipated during run time

        time.Sleep(5 * time.Second)

        curIp := httpconf.ip
        curPort := httpconf.port
        curRt := httpconf.readTimeout.Seconds()
        curWt := httpconf.writeTimeout.Seconds()
        curSt := httpconf.shutdownTimeout.Seconds()
        curLp := httpconf.logPath

       curRtS := strconv.FormatFloat(curRt, 'f', 0, 64)
        curStS := strconv.FormatFloat(curSt, 'f', 0, 64)
        curWtS := strconv.FormatFloat(curWt, 'f', 0, 64)

        envIp := os.Getenv("HTTP_IP_ADDR")
        envPort := os.Getenv("HTTP_IP_PORT")
        envRto := os.Getenv("HTTP_READ_TIMEOUT")
        envWto := os.Getenv("HTTP_WRITE_TIMEOUT")
        envSto := os.Getenv("HTTP_SHUTDOWN_TIMEOUT")
        envLog := os.Getenv("HTTP_LOG")

        os.Setenv("HTTP_IP_ADDR", curIp)
       os.Setenv("HTTP_IP_PORT", curPort)
        os.Setenv("HTTP_READ_TIMEOUT", curRtS)
        os.Setenv("HTTP_WRITE_TIMEOUT", curWtS)
        os.Setenv("HTTP_SHUTDOWN_TIMEOUT", curStS)
        os.Setenv("HTTP_LOG", curLp)

        //set up same environment
        a.initializeEnvironment()

        //expected := "0.0.0.0"
        if httpconf.ip != curIp {
                t.Errorf("Expected IP %v; Got %v\n", curIp, httpconf.ip)
        }

        //expected = "10000"
        if httpconf.port != curPort {
                t.Errorf("Expected PORT %v; Got %v\n", curPort, httpconf.port)
        }

        //expected = "foobar.log"
        if httpconf.logPath != curLp {
                t.Errorf("Expected log %v; Got %v\n", curLp, httpconf.logPath)
        }

	        //expectedInt := 60
        if httpconf.readTimeout.Seconds() != curRt {
                t.Errorf("Expected read timeout %v; Got %v\n", curRt, httpconf.readTimeout.Seconds())
        }

        //expectedInt = 30
        if httpconf.shutdownTimeout.Seconds() != curSt {
                t.Errorf("Expected shutdown timeout%v; Got %v\n", curSt, httpconf.shutdownTimeout.Seconds())
        }

        //expectedInt = 30
       if httpconf.writeTimeout.Seconds() != curWt {
                t.Errorf("Expected write timeout %v; Got %v\n", curWt, httpconf.writeTimeout.Seconds())
        }

        os.Unsetenv("HTTP_IP_ADDR")
        os.Unsetenv("HTTP_IP_PORT")
        os.Unsetenv("HTTP_LOG")

        // Test error path
        os.Setenv("HTTP_READ_TIMEOUT", "A")
        os.Setenv("HTTP_WRITE_TIMEOUT", "B")
        os.Setenv("HTTP_SHUTDOWN_TIMEOUT", "C")
        //Errors from this is written to log
        a.initializeEnvironment()

        os.Unsetenv("HTTP_READ_TIMEOUT")
        os.Unsetenv("HTTP_WRITE_TIMEOUT")
        os.Unsetenv("HTTP_SHUTDOWN_TIMEOUT")

        if envIp != "" {
                os.Setenv("HTTP_IP_ADDR", envIp)

        }
        if envPort != "" {
               os.Setenv("HTTP_IP_PORT", envPort)

        }
        if envRto != "" {
                os.Setenv("HTTP_READ_TIMEOUT", envRto)

        }
        if envWto != "" {
                os.Setenv("HTTP_WRITE_TIMEOUT", envWto)

        }
        if envSto != "" {
                os.Setenv("HTTP_SHUTDOWN_TIMEOUT", envSto)
        }
        if envLog != "" {
                os.Setenv("HTTP_LOG", envLog)

        }

        return
}

func TestJob(t *testing.T) {

      // Get the list of jobs
        req, _ := http.NewRequest("GET", JobListURL, nil)
        response := executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "job list")

        var jl []listJobsResponse
        payload, e := ioutil.ReadAll(response.Body)
        if e != nil {
                t.Errorf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
        }

        e = json.Unmarshal(payload, &jl)
        if e != nil {
                t.Errorf("jobsList Unmarsahl failed for payload %v jobs; Error %v\n", payload, e)
        }

        if len(jl) <= 0 {
                t.Errorf("jobsList return %v jobs; Wanted 5\n", len(jl))
                return
        }

        // Test getting a job
        req, _ = http.NewRequest("GET", JobURL+"/"+jl[0].ID, nil)
response = executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "job by ID")

        var jd map[string]interface{}
        payload, e = ioutil.ReadAll(response.Body)
        if e != nil {
                t.Errorf("{\"error\": \"ioutil.ReadAll failed\", \"Error\": \"%v\"}", e.Error())
        }

        e = json.Unmarshal(payload, &jd)
        if e != nil {
              t.Errorf("jobsGet Unmarsahl failed for payload %v jobs; Error %v\n", payload, e)
        }

        // We know that the interface requires an id and a type
        // So we can test for those
        v, ok := jd["id"].(string)
        if !ok {
                t.Errorf("job id not found\n")
        }
        if v == "" {
               t.Errorf("job id is required but is empty\n")
        }

        v, ok = jd["type"].(string)
        if !ok {
                t.Errorf("job type not found\n")
        }
        if v == "" {
                t.Errorf("job type is required but is empty\n")
        }
       // Create a new job
        nj := "{\"id\": \"\", \"url\": \"https://cat-fact.herokuapp.com/facts\", \"type\": \"\", \"client_id\": \"\"}"
        req, _ = http.NewRequest("POST", JobURL, strings.NewReader(nj))
        response = executeRequest(req)
        checkResponseCode(t, http.StatusCreated, response.Code, "create job")

        if !strings.Contains(response.Body.String(), "success") {
                t.Errorf("Create Job: expected success; Got failed\n")
        }

        // Update a job
        // It reads the UUID from Get data, ID is used to find job
       //ID is not updated. Can update url and client_id
        uj := "{\"id\": \"" + jl[0].ID + "\", \"url\": \"https://geek-jokes.sameerkumar.website/api\", \"type\": \"\",\"client_id\": \"" + jl[0].ClientID + "\"}"
        req, _ = http.NewRequest("PUT", JobURL+"/"+jl[0].ID, strings.NewReader(uj))
        response = executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "update job")

        if !strings.Contains(response.Body.String(), "success") {
                t.Errorf("Update Job: expected success; Got %v\n", response.Body.String())
        }
      // Delete a job
        // Try another jobId
        req, _ = http.NewRequest("DELETE", JobURL+"/"+jl[1].ID, nil)
        response = executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "delete job")

        if !strings.Contains(response.Body.String(), "success") {
                t.Errorf("Delete Job: expected success; Got: %v\n", response.Body.String())
        }

}

func TestManagementGet(t *testing.T) {
      //Only testing one command and one config from management commands
        //All commands call the same routine, all config call the same
        //config routine
        oneCmd := "\"stop_job\":{\"description\":\"Stop running of a continuous job.\",\"resource\":\"Scheduler\",\"CmdParms\":{\"Job_UUID\":{\"data_type\":\"string\",\"param_order\":0}}}"

        oneFig := "\"hard_shutdown_seconds\":{\"data_type\":\"int\",\"resource\":\"Dispatcher\"}"

        //Get a list of the management commands
       req, _ := http.NewRequest("GET", ManagementURL, nil)
        response := executeRequest(req)
        checkResponseCode(t, http.StatusOK, response.Code, "management list")

        //if body := response.Body.String(); body != er {
        body := response.Body.String()
        if !strings.Contains(body, oneCmd) {
                t.Errorf("Expected %s/n. Got %s", oneCmd, body)
        }
        if !strings.Contains(body, oneFig) {
                t.Errorf("Expected %s/n. Got %s", oneFig, body)
       }

        payload, e := ioutil.ReadAll(response.Body)
        if e != nil {
                t.Errorf("{\"error\": \"ioutil. getManagement readAll failed\", \"Error\": \"%v\"}", e.Error())
        }
        //mo will have list of commands that can be tested
        //note: this is required before in TestManagementPutSafe
        e = json.Unmarshal(payload, &mo)
        if e != nil {
                t.Errorf("getManagement unmarsahl failed for payload %v ; Error %v\n", payload, e)
      }

        if len(mo.CommandNames) == 0 {
                t.Errorf("command list return %v ; more expected\n", len(mo.CommandNames))
                return
        }
        if len(mo.Fields) == 0 {
                t.Errorf("config list return %v ; more expected\n", len(mo.Fields))
                return
        }

}

func TestManagementPutSafe(t *testing.T) {

	// TODO: These are not implemented in dispatcher yet
	/*
	tc["scheduler_channel_size"] = "{\"command\": \"set\", \"field\": \"scheduler_channel_size\", \"field_value\": 10}"
	tc["result_channel_size"] = "{\"command\": \"set\", \"field\": \"result_channel_size\", \"field_value\": 10}"
*/

    //Do a bad command
        mr.CommandName = "s`top_scheduler"
        mr.Resource = resDispatcher
        mr.CommandType = "command"
        mr.Fields = make(map[string]FieldParm)

        mrd, e := json.Marshal(&mr)
       if e != nil {
                t.Errorf("Couldn't set command for: %s \n", "s`top_scheduler")
        }
        mrds := string(mrd)
        req, _ := http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
        response := executeRequest(req)
        checkResponseCode(t, http.StatusNotFound, response.Code, "s`top_scheduler")

        time.Sleep(30 * time.Second)
        _, isCn := mo.CommandNames["start_scheduler"]
        if isCn {
                mr.CommandName = "start_scheduler"
                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed for start_scheduler: %v \n", e)
                }
                mrds = string(mrd)
                fmt.Println("start_scheduler \n", mrds)
                req, _ := http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))

		                response = executeRequest(req)
                checkResponseCode(t, http.StatusNotImplemented, response.Code, "start_scheduler")
        } else {

                t.Errorf("management list doesn't have start_scheduler")

        }

        time.Sleep(30 * time.Second)

        _, isCn = mo.CommandNames["stop_workers"]
        if isCn {
               mr.CommandName = "stop_workers"
                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed for stop_workers: %v \n", e)
                }
                mrds = string(mrd)
                fmt.Println("stop_workers \n", mrds)

                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
                response = executeRequest(req)
               checkResponseCode(t, http.StatusOK, response.Code, "stop workers")
        } else {

                t.Errorf("Management list doesn't have stop_workers")

        }

        //long sleep time to see if all workers will stop
        time.Sleep(60 * time.Second)

        _, isCn = mo.CommandNames["start_workers"]
      if isCn {
                mr.CommandName = "start_workers"
                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed  fori start_workers: %v \n", e)
                }
                mrds = string(mrd)
                fmt.Println("start_workers \n", mrds)

                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
                response = executeRequest(req)
                checkResponseCode(t, http.StatusOK, response.Code, "start_workers 1")

                //long sleep time in case all workers didnt stop
                time.Sleep(30 * time.Second)

                // try start again, must do http.New
                // expecting all workers to be up, so start should be forbidden
                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
                response = executeRequest(req)
               checkResponseCode(t, http.StatusForbidden, response.Code, "start_workers again")
        } else {

                t.Errorf("management list doesn't have: start_workers")

        }

        time.Sleep(30 * time.Second)

        _, isCn = mo.Fields["graceful_shutdown_seconds"]
        if isCn {
                mr.CommandName = "set"
                mr.CommandType = "config"
                fp.FieldValue = "20"
                fp.DataType = "int"
                mr.Fields["graceful_shutdown_seconds"] = fp
                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed for graceful_shutdown_seconds: %v \n", e)
                }
                mrds = string(mrd)
                fmt.Println("graceful_shutdown_seconds \n", mrds)
                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
                response = executeRequest(req)
                checkResponseCode(t, http.StatusOK, response.Code, "graceful_shutdown_seconds")
        } else {

                t.Errorf("management list doesn't have graceful_shutdown_seconds")

        }
       time.Sleep(30 * time.Second)

        _, isCn = mo.Fields["hard_shutdown_seconds"]
        if isCn {
                mr.Fields = make(map[string]FieldParm)
                fp.FieldValue = "5"
                mr.Fields["hard_shutdown_seconds"] = fp

                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed for hard_shutdown_seconds %v \n", e)
               }
                mrds = string(mrd)
                fmt.Println("hard_shutdown_seconds 5\n", mrds)
                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))
                response = executeRequest(req)
                checkResponseCode(t, http.StatusOK, response.Code, "hard_shutdown_seconds 5")

                //TODO: fix this issue
                mr.Fields = make(map[string]FieldParm)
                fp.FieldValue = "0"
                mr.Fields["hard_shutdown_seconds"] = fp
                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed for hard_shutdown_seconds 0: %s \n", e)
                }
                mrds = string(mrd)
                fmt.Println("hard_shutdown_seconds 0\n", mrds)
                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))

                // Set it back to 0 so we exit tests quickly
                response = executeRequest(req)
               checkResponseCode(t, http.StatusOK, response.Code, "hard_shutdown_seconds 0")

                // send some bad JSON data to test marshal error
                //TODO: fix this issue
                mr.Fields = make(map[string]FieldParm)
                fp.FieldValue = "A"
                mr.Fields["hard_shutdown_seconds"] = fp

                mrd, e = json.Marshal(&mr)
                if e != nil {
                        t.Errorf("marshal failed for hard_shutdown_seconds A: %s \n", e)
                }
                mrds = string(mrd)
                fmt.Println("hard_shutdown_seconds A\n", mrds)
                req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))

                response = executeRequest(req)
                checkResponseCode(t, http.StatusBadRequest, response.Code, "hard_shutdown_seconds A")
        } else {
                t.Errorf("Management list doesn't have hard_shutdown_seconds")

        }

        // send a bad command
        mr.Fields = make(map[string]FieldParm)
        fp.FieldValue = "5"
        mr.Fields["hard_lockdown_seconds"] = fp

        mrd, e = json.Marshal(&mr)
        if e != nil {
                t.Errorf("marshal failed for hard_shutdown_seconds %v \n", e)
      }
        mrds = string(mrd)
        fmt.Println("hard_lockdown_seconds 5\n", mrds)
        req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(mrds))

        response = executeRequest(req)
        checkResponseCode(t, http.StatusNotFound, response.Code, "hard_lockdown_seconds 5")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int, testAt string) {
	if expected != actual {
	t.Errorf("%s: expected response code %d. Got %d\n", testAt, expected, actual)
	}
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
