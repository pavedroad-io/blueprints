package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"
)

const (
	ServerString string = "127.0.0.1:9000"
	EPString     string = "/endpoint"
	BADUUID      string = "00000000-0000-0000-0000-000000000000"
)

// TestMain handles setup required before tests execution
func JobTestMain() {
	// Start a server to test against
	http.HandleFunc(EPString, testEndPoint)
	go startServer()

	// Wait for server to start
	time.Sleep(1 * time.Second)

	return
}

func startServer() {
	fmt.Println("start httpJob Server")
	e := http.ListenAndServe(ServerString, nil)
	if e != nil {
		fmt.Println(e)
		os.Exit(-1)
	}

	return
}

func testEndPoint(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("true"))
	return
}

func TestHTTPJob(t *testing.T) {
	j := newJob("http://" + ServerString + EPString)
	r, e := j.Run()

	if e != nil {
		t.Errorf("j.Run() failed err = %v; Wanted nil\n", e)
	}

	if string(r.Payload()) != "true" {
		t.Errorf("r.Payload() failed body = %v; Wanted true\n", e)
	}

	if string(r.Job().ID()) == BADUUID {
		t.Errorf("r.Job.ID() failed UUID = %v; Wanted valid UUID\n", r.Job().ID())
	}

	return
}

func TestHTTPJob_Metrics(t *testing.T) {
	j := newJob("http://" + ServerString + EPString)
	_, e := j.Run()

	if e != nil {
		t.Errorf("j.Run() failed err = %v; Wanted nil\n", e)
	}

	jb := j.Metrics()

	hs := httpStats{}
	e = json.Unmarshal(jb, &hs)
	if e != nil {
		t.Errorf("json.Unmarshal() failed err = %v; Wanted nil\n", e)
	}

	if hs.RequestTimedOut {
		t.Errorf("hs.RequestTimedOut failed = %v; Wanted true\n", hs.RequestTimedOut)
	}

	if hs.RequestTime <= time.Duration(0) {
		t.Errorf("hs.RequestTime failed = %v; Wanted >0\n", hs.RequestTime)
	}

}

// newJob builds an httpJob for the specified URL
func newJob(testURL string) httpJob {
	newJob := httpJob{}
	pu, err := url.Parse(testURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	newJob.JobURL = pu

	// Set type and ID and http.Client
	newJob.Init()
	return newJob
}