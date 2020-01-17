
// eventCollector_test.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	_ "strconv"
	"strings"
	"testing"
	"time"
)

const (
	Updated         string = "updated"
	Created         string = "created"
	Active          string = "active"
  EventCollectorURL string = "/api/v1/namespace/pavedroad.io/eventCollector/%s"
)

var newEventCollectorJSON=`{
	"eventcollectoruuid": "4241b4b6-26a7-49c7-bb4e-1078b07dd1bf",
	"id": "MdN8rIZho6rBpxN",
	"title": "7FV23nH2c4YbtVg",
	"updated": "2020-01-16T17:10:26-08:00",
	"created": "2020-01-16T17:10:26-08:00",
	"metadata": {
		"author": "QJBNZqhhmczvwSg",
		"genre": "B49JeEH3PjpLw9L",
		"rating": "bp89PHY5mswVGAh"
	}
}`
	"eventcollectoruuid": "6e3d1e83-c8f1-40a1-a805-6562269ec2c9",
	"id": "t7H7hSmX0J8YrbX",
	"title": "UXMn5RCYpoyUGQr",
	"updated": "2020-01-08T14:15:19-08:00",
	"created": "2020-01-08T14:15:19-08:00",
	"metadata": {
		"author": "9kG9C7z0chmCxRP",
		"genre": "LdHXEtDs64XChvF",
		"rating": "nDghhlloENUZD3Y"
	}
}`

var a EventCollectorApp

func TestMain(m *testing.M) {
	a = EventCollectorApp{}
	a.Initialize()

	clearDB()
	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		fmt.Println("Table check failed:", err)
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(indexCreate); err != nil {
		fmt.Println("Table check failed:", err)
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM Mirantis.EventCollector")
}

func clearDB() {
	a.DB.Exec("DROP DATABASE IF EXISTS Mirantis")
	a.DB.Exec("CREATE DATABASE Mirantis")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS Mirantis.eventCollector (
    EventCollectorUUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    eventCollector JSONB
);`


const indexCreate = `
CREATE INDEX IF NOT EXISTS eventCollectorIdx ON Mirantis.eventCollector USING GIN (eventCollector);`


func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/namespace/pavedroad.io/eventCollectorLIST", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
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

// TestGetWithBadUserUUID
// Get a users with an invalid UUID, should return 400
// and that it is an invalid UUID
//
func TestGetWithBadUserUUID(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET",
		"/api/v1/namespace/pavedroad.io/eventCollector/43ae99c9", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "400: invalid UUID: 43ae99c9" {
		t.Errorf("Expected the 'error' key of the response to be set to '400: invalid UUID: 43ae99c9'. Got '%s'", m["error"])
	}
}

// TestGetWrongUUID
// Is a valid UUID, but with leading zeros
// This will not be found and should return a 304
//
func TestGetWrongUUID(t *testing.T) {
	clearTable()
  nt := NewEventCollector()
  addEventCollector(nt)
  badUid := "00000000-d01d-4c09-a4e7-59026d143b89"

	statement := fmt.Sprintf(EventCollectorURL, badUid)

	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
// TestCreate
// Use sample data from newEventCollectorJSON) to create
// a new record.
// TODO:
//  need to assert tests for subattributes being present
//
func TestCreateEventCollector(t *testing.T) {
	clearTable()

	payload := []byte(newEventCollectorJSON)

	req, _ := http.NewRequest("POST", "/api/v1/namespace/pavedroad.io/eventCollector", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	//var md map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	//Test we can decode the data
	cs, ok := m["created"].(string)
	if ok {
		c, err := time.Parse(time.RFC3339, cs)
		if err != nil {
			t.Errorf("Parse failed on parse creataed time Got '%v'", c)
		}
	} else {
		t.Errorf("Expected creataed of string type Got '%v'", reflect.TypeOf(m["Created"]))
	}

	us, ok := m["updated"].(string)
	if ok {
		u, err := time.Parse(time.RFC3339, us)
		if err != nil {
			t.Errorf("Parse failed on parse updated time Got '%v'", u)
		}
	} else {
		t.Errorf("Expected updated of string type Got '%v'", reflect.TypeOf(m["Updated"]))
	}
}

func TestMarshallEventCollector(t *testing.T) {
	nt := NewEventCollector()
	_, err := json.Marshal(nt)
	if err != nil {
		t.Errorf("Marshal of EventCollector failed: Got '%v'", err)
	}
}

// addEventCollector
// Inserts a new user into the database and returns the UUID
// for the record that was created
//
func addEventCollector(t *eventCollector) (string) {

  statement := fmt.Sprintf("INSERT INTO Mirantis.eventCollector(eventCollector) VALUES('%s') RETURNING eventCollectorUUID", newEventCollectorJSON)
  rows, er1 := a.DB.Query(statement)

	if er1 != nil {
		log.Printf("Insert failed error %s", er1)
		return ""
	}

	defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.EventCollectorUUID)
    if err != nil {
      return ""
    }
  }

  return t.EventCollectorUUID
}

// NewEventCollector
// Create a new instance of EventCollector
// Iterate over the struct setting random values
// 
func NewEventCollector() (t *eventCollector) {
	var n eventCollector
  json.Unmarshal([]byte(newEventCollectorJSON), &n) 
	return &n
}

//test getting a eventCollector
func TestGetEventCollector(t *testing.T) {
	clearTable()
	nt := NewEventCollector()
	uid := addEventCollector(nt)
	statement := fmt.Sprintf(EventCollectorURL, uid)

	req, err := http.NewRequest("GET", statement, nil)
  if err != nil {
		panic(err)
  }

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
// TestUpdateEventCollector
func TestUpdateeventCollector(t *testing.T) {
	clearTable()
	nt := NewEventCollector()
	uid := addEventCollector(nt)

	statement := fmt.Sprintf(EventCollectorURL, uid)
	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	json.Unmarshal(response.Body.Bytes(), &nt)

	ut := nt

	//Update the new struct
	//ut.Active = "eslaf"

	jb, err := json.Marshal(ut)
	if err != nil {
		panic(err)
	}

	req, _ = http.NewRequest("PUT", statement, strings.NewReader(string(jb)))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

//	if m["active"] != "eslaf" {
//		t.Errorf("Expected active to be eslaf. Got %v", m["active"])
//	}
}

func TestDeleteeventCollector(t *testing.T) {
	clearTable()
	nt := NewEventCollector()
	uid := addEventCollector(nt)

	statement := fmt.Sprintf(EventCollectorURL, uid)
	req, _ := http.NewRequest("DELETE", statement, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", statement, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

/*
func TestDumpEventCollector(t *testing.T) {
	nt := NewEventCollector()

  err := dumpUser(*nt)

	if err != nil {
		t.Errorf("Expected erro to be nill. Got %v", err)
	}
}
*/
