
// films_test.go

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
  FilmsURL string = "/api/v1/namespace/pavedroad.io/films/%s"
)

var newFilmsJSON=`{
	"filmsuuid": "7545e80d-8758-41b0-98a8-4476da49b5ac",
	"id": "gWAfKJvIvCYgRfD",
	"title": "WCzeGnMuETejBkC",
	"updated": "2019-12-18T09:33:44-08:00",
	"created": "2019-12-18T09:33:44-08:00",
	"metadata": {
		"author": "pVncA5odR19wAdJ",
		"genre": "ncKGOa8vkUFmoaG",
		"rating": "Gz3QVGA1ZuAaQdt"
	}
}`


var a FilmsApp

func TestMain(m *testing.M) {
	a = FilmsApp{}
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
	a.DB.Exec("DELETE FROM AcmeDemo.Films")
}

func clearDB() {
	a.DB.Exec("DROP DATABASE IF EXISTS AcmeDemo")
	a.DB.Exec("CREATE DATABASE AcmeDemo")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS AcmeDemo.films (
    FilmsUUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    films JSONB
);`


const indexCreate = `
CREATE INDEX IF NOT EXISTS filmsIdx ON AcmeDemo.films USING GIN (films);`


func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/namespace/pavedroad.io/filmsLIST", nil)
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
		"/api/v1/namespace/pavedroad.io/films/43ae99c9", nil)
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
  nt := NewFilms()
  addFilms(nt)
  badUid := "00000000-d01d-4c09-a4e7-59026d143b89"

	statement := fmt.Sprintf(FilmsURL, badUid)

	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
// TestCreate
// Use sample data from newFilmsJSON) to create
// a new record.
// TODO:
//  need to assert tests for subattributes being present
//
func TestCreateFilms(t *testing.T) {
	clearTable()

	payload := []byte(newFilmsJSON)

	req, _ := http.NewRequest("POST", "/api/v1/namespace/pavedroad.io/films", bytes.NewBuffer(payload))
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

func TestMarshallFilms(t *testing.T) {
	nt := NewFilms()
	_, err := json.Marshal(nt)
	if err != nil {
		t.Errorf("Marshal of Films failed: Got '%v'", err)
	}
}

// addFilms
// Inserts a new user into the database and returns the UUID
// for the record that was created
//
func addFilms(t *films) (string) {

  statement := fmt.Sprintf("INSERT INTO AcmeDemo.films(films) VALUES('%s') RETURNING filmsUUID", newFilmsJSON)
  rows, er1 := a.DB.Query(statement)

	if er1 != nil {
		log.Printf("Insert failed error %s", er1)
		return ""
	}

	defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.FilmsUUID)
    if err != nil {
      return ""
    }
  }

  return t.FilmsUUID
}

// NewFilms
// Create a new instance of Films
// Iterate over the struct setting random values
// 
func NewFilms() (t *films) {
	var n films
  json.Unmarshal([]byte(newFilmsJSON), &n) 
	return &n
}

//test getting a films
func TestGetFilms(t *testing.T) {
	clearTable()
	nt := NewFilms()
	uid := addFilms(nt)
	statement := fmt.Sprintf(FilmsURL, uid)

	req, err := http.NewRequest("GET", statement, nil)
  if err != nil {
		panic(err)
  }

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
// TestUpdateFilms
func TestUpdatefilms(t *testing.T) {
	clearTable()
	nt := NewFilms()
	uid := addFilms(nt)

	statement := fmt.Sprintf(FilmsURL, uid)
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

func TestDeletefilms(t *testing.T) {
	clearTable()
	nt := NewFilms()
	uid := addFilms(nt)

	statement := fmt.Sprintf(FilmsURL, uid)
	req, _ := http.NewRequest("DELETE", statement, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", statement, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

/*
func TestDumpFilms(t *testing.T) {
	nt := NewFilms()

  err := dumpUser(*nt)

	if err != nil {
		t.Errorf("Expected erro to be nill. Got %v", err)
	}
}
*/
