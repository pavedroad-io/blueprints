
// eventbridge_test.go

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	_ "strconv"
	"strings"
	"testing"
	"time"
)

const (
	Updated         string = "updated"
	Created         string = "created"
	Active          string = "active"
	Namespace     string = "pavedroad"
	Service       string = "eventbridge"
	ManagementURL string = "/api/v1/namespace/" + Namespace + "/" + Service + "/management"
	JobURL        string = "/api/v1/namespace/" + Namespace + "/" + Service + "/jobs"
	JobListURL    string = "/api/v1/namespace/" + Namespace + "/" + Service + "/jobsLIST"
	ScheduleURL   string = "/api/v1/namespace/" + Namespace + "/" + Service + "/scheduler"
	ReadyURL      string = "/api/v1/namespace/" + Namespace + "/" + Service + "/ready"
	LiveURL       string = "/api/v1/namespace/" + Namespace + "/" + Service + "/liveness"
	MetricsURL    string = "/api/v1/namespace/" + Namespace + "/" + Service + "/metrics"
)

var newEventbridgeJSON=``

var a EventbridgeApp
var testSampleCode = true

func TestMain(m *testing.M) {
	a = EventbridgeApp{}

	JobTestMain()

	a.Initialize()
	go a.Run(httpconf.listenString)

	// Wait for server to start
	time.Sleep(1 * time.Second)
	defer testExit()
	defer os.Exit(m.Run())
}

func testExit() {
	data := "{\"command\": \"shutdown_now\", \"field\": \"\", \"field_value\": 0}"
	req, _ := http.NewRequest("PUT", ManagementURL, strings.NewReader(data))
	_ = executeRequest(req)
}

func TestSchedule(t *testing.T) {
	// Get the list of jobs
	req, _ := http.NewRequest("GET", ScheduleURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// This code requires knowledge about the type of scheduler
	// So make it conditional
	if testSampleCode {
		putdata := "{\"schedule_type\": \"Constant interval scheduler\", \"send_interval_seconds\": 30}"
		req, _ = http.NewRequest("PUT", ScheduleURL, strings.NewReader(putdata))
		response = executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code)

		req, _ = http.NewRequest("DELETE", ScheduleURL, strings.NewReader(putdata))
		response = executeRequest(req)
		checkResponseCode(t, http.StatusOK, response.Code)

		postdata := "{\"schedule_type\": \"Constant interval scheduler\", \"send_interval_seconds\": 5}"
		req, _ = http.NewRequest("POST", ScheduleURL, strings.NewReader(postdata))
		response = executeRequest(req)
		checkResponseCode(t, http.StatusCreated, response.Code)
	}
	return
}

func TestReady(t *testing.T) {

	if !a.Ready {
		a.Ready = true
	}

	expect := "{\"Ready\": true}"

	req, _ := http.NewRequest("GET", ReadyURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if response.Body.String() != expect {
		t.Errorf("Expected %s; Got %v\n", expect, response.Body.String())
	}

	a.Ready = false
	expect = "{\"Ready\": false}"

	req, _ = http.NewRequest("GET", ReadyURL, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusServiceUnavailable, response.Code)
	if response.Body.String() != expect {
		t.Errorf("Expected %s; Got %v\n", expect, response.Body.String())
	}

	// Reset to proper status
	a.Ready = true
	return
}

func TestLive(t *testing.T) {

	if !a.Live {
		a.Live = true
	}

	expect := "{\"Live\": true}"
	req, _ := http.NewRequest("GET", LiveURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if response.Body.String() != expect {
		t.Errorf("Expected %s; Got %v\n", expect, response.Body.String())
	}

	a.Live = false
	expect = "{\"Live\": false}"

	req, _ = http.NewRequest("GET", LiveURL, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusServiceUnavailable, response.Code)
	if response.Body.String() != expect {
		t.Errorf("Expected %s; Got %v\n", expect, response.Body.String())
	}

	// Reset to proper status
	a.Live = true
	return
}

func TestMetric(t *testing.T) {

	expect1 := "scheduler"
	expect2 := "dispatcher"

	req, _ := http.NewRequest("GET", MetricsURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	if !strings.Contains(response.Body.String(), expect1) {
		t.Errorf("Expected %s; Got %v\n", expect1, response.Body.String())
	}
	if !strings.Contains(response.Body.String(), expect2) {
		t.Errorf("Expected %s; Got %v\n", expect2, response.Body.String())
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
	os.Setenv("HTTP_IP_ADDR", "0.0.0.0")
	os.Setenv("HTTP_IP_PORT", "10000")
	os.Setenv("HTTP_READ_TIMEOUT", "60")
	os.Setenv("HTTP_WRITE_TIMEOUT", "30")
	os.Setenv("HTTP_SHUTDOWN_TIMEOUT", "30")
	os.Setenv("HTTP_LOG", "foobar.log")

	a.initializeEnvironment()

	expected := "0.0.0.0"
	if httpconf.ip != expected {
		t.Errorf("Expected IP %v; Got %v\n", expected, httpconf.ip)
	}

	expected = "10000"
	if httpconf.port != expected {
		t.Errorf("Expected PORT %v; Got %v\n", expected, httpconf.port)
	}

	expected = "foobar.log"
	if httpconf.logPath != expected {
		t.Errorf("Expected log %v; Got %v\n", expected, httpconf.logPath)
	}

	expectedInt := 60
	if int(httpconf.readTimeout.Seconds()) != expectedInt {
		t.Errorf("Expected read timeout %v; Got %v\n", expected, httpconf.port)
	}

	expectedInt = 30
	if int(httpconf.shutdownTimeout.Seconds()) != expectedInt {
		t.Errorf("Expected shutdown timeout%v; Got %v\n", expected, httpconf.port)
	}

	expectedInt = 30
	if int(httpconf.writeTimeout.Seconds()) != expectedInt {
		t.Errorf("Expected write timeout %v; Got %v\n", expected, httpconf.port)
	}

	os.Unsetenv("HTTP_IP_ADDR")
	os.Unsetenv("HTTP_IP_PORT")
	os.Unsetenv("HTTP_LOG")

	// Test error path
	os.Setenv("HTTP_READ_TIMEOUT", "A")
	os.Setenv("HTTP_WRITE_TIMEOUT", "B")
	os.Setenv("HTTP_SHUTDOWN_TIMEOUT", "C")

	a.initializeEnvironment()

	return
}

func TestJob(t *testing.T) {

	// Get the list of jobs
	req, _ := http.NewRequest("GET", JobListURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

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
		t.Errorf("jobsList return %v jobs; Wanted 4\n", len(jl))
		return
	}

	// Test getting a job
	req, _ = http.NewRequest("GET", JobURL+"/"+jl[0].ID, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

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
	nj := "{\"id\": \"\", \"url\": \"https://cat-fact.herokuapp.com/facts\", \"type\": \"\"}"
	req, _ = http.NewRequest("POST", JobURL, strings.NewReader(nj))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)

	if !strings.Contains(response.Body.String(), "success") {
		t.Errorf("Create Job: expected success; Got false\n")
	}

	// Update a job
	// It reads the UUID from Post data, may want to change that
	uj := "{\"id\": \"" + jl[0].ID + "\", \"url\": \"https://geek-jokes.sameerkumar.website/api\", \"type\": \"\"}"
	req, _ = http.NewRequest("PUT", JobURL+"/"+jl[0].ID, strings.NewReader(uj))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if !strings.Contains(response.Body.String(), "success") {
		t.Errorf("Update Job: expected success; Got false\n")
	}

	// Delete the job
	// The update above results in new ID, so use the second record
	req, _ = http.NewRequest("DELETE", JobURL+"/"+jl[1].ID, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if !strings.Contains(response.Body.String(), "success") {
		t.Errorf("Delete Job: expected success; Got false\n")
	}

}

func TestManagementGet(t *testing.T) {
	er := "{\"commands\":[{\"name\":\"set\",\"data_type\":\"int\",\"command_type\":\"config\",\"description\":\"Sets the value of a configurable field, see fields below\"},{\"name\":\"stop_scheduler\",\"data_type\":\"string\",\"command_type\":\"command\",\"description\":\"Stops the scheduler from send new jobs\"},{\"name\":\"start_scheduler\",\"data_type\":\"string\",\"command_type\":\"command\",\"description\":\"Starts the scheduler running again.  If running has no affect\"},{\"name\":\"stop_workers\",\"data_type\":\"string\",\"command_type\":\"command\",\"description\":\"Shutdown the worker pool letting jobs inflight complete\"},{\"name\":\"start_workers\",\"data_type\":\"string\",\"command_type\":\"command\",\"description\":\"Starts the worker pool if stopped\"},{\"name\":\"shutdown\",\"data_type\":\"string\",\"command_type\":\"command\",\"description\":\"Graceful shutdown\"},{\"name\":\"shutdown_now\",\"data_type\":\"string\",\"command_type\":\"command\",\"description\":\"Hard shutdown with SIGKILL\"}],\"fields\":[\"graceful_shutdown_seconds\",\"hard_shutdown_seconds\",\"number_of_workers\",\"scheduler_channel_size\",\"result_channel_size\"]}"

	req, _ := http.NewRequest("GET", ManagementURL, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != er {
		t.Errorf("Expected %s. Got %s", er, body)
	}
}

//postdata5 := "{\"command\": \"shutdown\", \"field\": \"\", \"field_value\": 0}"
//postdata6 := "{\"command\": \"shutdown_now\", \"field\": \"\", \"field_value\": 0}"
func TestManagementPutSafe(t *testing.T) {
	tc := make(map[string]string)
	tc["stop_scheduler"] = "{\"command\": \"stop_scheduler\", \"field\": \"\", \"field_value\": 0}"
	tc["start_scheduler"] = "{\"command\": \"start_scheduler\", \"field\": \"\", \"field_value\": 0}"
	tc["stop_workers"] = "{\"command\": \"stop_workers\", \"field\": \"\", \"field_value\": 0}"
	tc["start_workers"] = "{\"command\": \"start_workers\", \"field\": \"\", \"field_value\": 0}"
	tc["set_graceful_shutdown"] = "{\"command\": \"set\", \"field\": \"graceful_shutdown_seconds\", \"field_value\": 5}"
	tc["set_hard_shutdown_seconds"] = "{\"command\": \"set\", \"field\": \"hard_shutdown_seconds\", \"field_value\": 5}"
	tc["number_of_workers"] = "{\"command\": \"set\", \"field\": \"number_of_workers\", \"field_value\": 10}"
	tc["set_hard_shutdown_seconds0"] = "{\"command\": \"set\", \"field\": \"hard_shutdown_seconds\", \"field_value\": 0}"
	tc["bad_command"] = "{\"command\": \"foobar\", \"field\": \"hard_shutdown_seconds\", \"field_value\": 0}"
	tc["bad_json"] = "{\"command\": foobar, \"field\": \"hard_shutdown_seconds\", \"field_value\": 0}"

	// TODO: These are not implemented in dispatcher yet
	tc["scheduler_channel_size"] = "{\"command\": \"set\", \"field\": \"scheduler_channel_size\", \"field_value\": 10}"
	tc["result_channel_size"] = "{\"command\": \"set\", \"field\": \"result_channel_size\", \"field_value\": 10}"

	req, _ := http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["stop_scheduler"]))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["start_scheduler"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["stop_workers"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["start_workers"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// 400 if already running test that next
	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["start_workers"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	er := "{\"Status\": \"5 already running\"}"
	if body := response.Body.String(); body != er {
		t.Errorf("Expected %s. Got %s", er, body)
	}

	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["set_graceful_shutdown"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	er = "{\"Status\": \"graceful_shutdown_seconds changed from 30 to 5\"}"
	if body := response.Body.String(); body != er {
		t.Errorf("Expected %s. Got %s", er, body)
	}

	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["set_hard_shutdown_seconds"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
	er = "{\"Status\": \"hard_shutdown_seconds changed from 0 to 5\"}"
	if body := response.Body.String(); body != er {
		t.Errorf("Expected %s. Got %s", er, body)
	}

	// Set it back to 0 so we exit tests quickly
	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["set_hard_shutdown_seconds0"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// send some bad JSON to test marshal error
	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["bad_command"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	er = "{\"Status\": \"Command foobar not implemented\"}"
	if body := response.Body.String(); body != er {
		t.Errorf("Expected %s. Got %s", er, body)
	}

	// send a bad command
	req, _ = http.NewRequest("PUT", ManagementURL, strings.NewReader(tc["bad_json"]))
	response = executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)

	er = "{\"error\": \"json.Unmarshal failed\", \"Error\": \"invalid character 'o' in literal false (expecting 'a')\"}"
	if body := response.Body.String(); body != er {
		t.Errorf("Expected %s. Got %s", er, body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
