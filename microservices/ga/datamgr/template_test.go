{{define "template_test.go"}}
// {{.Name}}_test.go

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
  {{.NameExported}}URL string = "/api/v1/namespace/pavedroad.io/{{.Name}}/%s"
)

var new{{.NameExported}}JSON=`{{.PostJSON}}`


var a {{.NameExported}}App

func TestMain(m *testing.M) {
	a = {{.NameExported}}App{}
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

      if _, err := a.DB.Exec("DELETE FROM {{.OrgSQLSafe}}.{{.NameExported}}"); err != nil {
	       fmt.Println("Table delete :", err)
      }
}

func clearDB() {

	 if _, err := a.DB.Exec("DROP DATABASE IF EXISTS {{.OrgSQLSafe}}"); err != nil {
               fmt.Println("Drop table:", err)
      }      
	if _, err := a.DB.Exec("CREATE DATABASE {{.OrgSQLSafe}}");  err != nil {
               fmt.Println("Create table:", err)
      }

}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS {{.OrgSQLSafe}}.{{.Name}} (
    {{.NameExported}}UUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    {{.Name}} JSONB
);`


const indexCreate = `
CREATE INDEX IF NOT EXISTS {{.Name}}Idx ON {{.OrgSQLSafe}}.{{.Name}} USING GIN ({{.Name}});`


func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/api/v1/namespace/pavedroad.io/{{.Name}}LIST", nil)
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

	req, err := http.NewRequest("GET",
		"/api/v1/namespace/pavedroad.io/{{.Name}}/43ae99c9", nil)
	if err != nil {
		fmt.Println("NewRequest:", err)
        }
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

	var m map[string]string
	err = json.Unmarshal(response.Body.Bytes(), &m)
        if err != nil {
                fmt.Println("Unmarshal issue:", err)
        }

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
  nt := New{{.NameExported}}()
  add{{.NameExported}}(nt)
  badUID := "00000000-d01d-4c09-a4e7-59026d143b89"

	statement := fmt.Sprintf({{.NameExported}}URL, badUID)

	req, _ := http.NewRequest("GET", statement, nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)
}
// TestCreate
// Use sample data from new{{.NameExported}}JSON) to create
// a new record.
// TODO:
//  need to assert tests for subattributes being present
//
func TestCreate{{.NameExported}}(t *testing.T) {
	clearTable()

	payload := []byte(new{{.NameExported}}JSON)

	req, err := http.NewRequest("POST", "/api/v1/namespace/pavedroad.io/{{.Name}}", bytes.NewBuffer(payload)) 
	if err != nil {
		fmt.Println("NewRequest Post:", err)
	}
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	//var md map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
                fmt.Println("Unmarshal issue:", err)
        }

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

func TestMarshall{{.NameExported}}(t *testing.T) {
	nt := New{{.NameExported}}()
	_, err := json.Marshal(nt)
	if err != nil {
		t.Errorf("Marshal of {{.NameExported}} failed: Got '%v'", err)
	}
}

// add{{.NameExported}}
// Inserts a new user into the database and returns the UUID
// for the record that was created
//
func add{{.NameExported}}(t *{{.Name}}) (string) {

  statement := `INSERT INTO {{.OrgSQLSafe}}.{{.Name}}({{.Name}}) VALUES($1) RETURNING {{.Name}}UUID`
  rows, er1 := a.DB.Query(statement,new{{.NameExported}}JSON)

	if er1 != nil {
		log.Printf("Insert failed error %s", er1)
		return ""
	}

	defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.{{.NameExported}}UUID)
    if err != nil {
      return ""
    }
  }

  return t.{{.NameExported}}UUID
}

// New{{.NameExported}}
// Create a new instance of {{.NameExported}}
// Iterate over the struct setting random values
// 
func New{{.NameExported}}() (t *{{.Name}}) {
	var N {{.Name}}
        err := json.Unmarshal([]byte(new{{.NameExported}}JSON), &N)
	if err != nil {
                fmt.Println("Unmarshal issue:", err)
        }
	return &N
}

//test getting a {{.Name}}
func TestGet{{.NameExported}}(t *testing.T) {
	clearTable()
	nt := New{{.NameExported}}()
	uid := add{{.NameExported}}(nt)
	statement := fmt.Sprintf({{.NameExported}}URL, uid)

	req, err := http.NewRequest("GET", statement, nil)
  if err != nil {
		panic(err)
  }

	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}
// TestUpdate{{.NameExported}}
func TestUpdate{{.Name}}(t *testing.T) {
	clearTable()
	nt := New{{.NameExported}}()
	uid := add{{.NameExported}}(nt)

	statement := fmt.Sprintf({{.NameExported}}URL, uid)
	req, err := http.NewRequest("GET", statement, nil)
	if err != nil {
		fmt.Println("NewRequest Get:", err)
	}
	response := executeRequest(req)

	err = json.Unmarshal(response.Body.Bytes(), &nt)
	if err != nil {
                fmt.Println("Unmarshal issue:", err)
        }

	ut := nt

	//Update the new struct
	//ut.Active = "eslaf"

	jb, err := json.Marshal(ut)
	if err != nil {
		panic(err)
	}

	req, err = http.NewRequest("PUT", statement, strings.NewReader(string(jb)))
	if err != nil {
		fmt.Println("NewRequest Put:", err)
	}
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	err = json.Unmarshal(response.Body.Bytes(), &m)
	if err != nil {
                fmt.Println("Unmarshal issue:", err)
        }


//	if m["active"] != "eslaf" {
//		t.Errorf("Expected active to be eslaf. Got %v", m["active"])
//	}
}

func TestDelete{{.Name}}(t *testing.T) {
	clearTable()
	nt := New{{.NameExported}}()
	uid := add{{.NameExported}}(nt)

	statement := fmt.Sprintf({{.NameExported}}URL, uid)
	req, _ := http.NewRequest("DELETE", statement, nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", statement, nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

/*
func TestDump{{.NameExported}}(t *testing.T) {
	nt := New{{.NameExported}}()

  err := dumpUser(*nt)

	if err != nil {
		t.Errorf("Expected erro to be nill. Got %v", err)
	}
}
*/
{{end}}
