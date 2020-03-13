{{define "templateModel.go"}}
{{.PavedroadInfo}}
// {{.OrganizationLicense}}

// User project / copyright / usage information
// {{.ProjectInfo}}

package main

import (
	"database/sql"
	"encoding/json"
	"github.com/google/uuid"
	"errors"
	"fmt"
	"time"
	log "github.com/pavedroad-io/core/go/logger"
)

// A GenericError is the default error message that is generated.
// For certain status codes there are more appropriate error structures.
//
// swagger:response genericError
type GenericError struct {
  // The error message
	// in: body
	Body struct {
    // Code: integer code for error message
		Code    int32 `json:"code"`
    // Message: Error message called with Method()
		Message error `json:"message"`
	} `json:"body"`
}

// Return list of {{.Name}}s
//
// TODO: add method of including sub attributes
//
// swagger:response {{.Name}}List
type listResponse struct {
  // in: body
  UUID  string `json:"uuid"`
}

// Generated structures with Swagger docs
{{.SwaggerGeneratedStructs}}

// {{.NameExported}} response model
//
// This is used for returning a response with a single {{.Name}} as body
//
// swagger:response {{.Name}}Response
type {{.NameExported}}Response struct {
	// in: body
	response string `json:"order"`
}

// update{{.NameExported}} in database
func (t *{{.Name}}) update{{.NameExported}}(db *sql.DB, key string) error {
	update := `
	UPDATE {{.OrgSQLSafe}}.{{.Name}}
    SET {{.Name}} = '%s'
  WHERE {{.Name}}UUID = '%s';`

  jb, err := json.Marshal(t)
  if err != nil {
    log.Printf("update{{.Name}} json.Marshal failed; Got (%v)\n", err.Error())
    return(err)
  }

  statement := fmt.Sprintf(update, jb, key)
  _, er1 := db.Query(statement)

  if er1 != nil {
    log.Println("Update failed")
    return er1
  }

  return nil
}

// create{{.NameExported}} in database
func (t *{{.Name}}) create{{.NameExported}}(db *sql.DB) (string, error) {
  jb, err := json.Marshal(t)
  if err != nil {
	msg := fmt.Sprintf("create{{.Name}} json.Marshal failed; Got (%v)\n", err.Error())
    log.Printf(msg)
    return msg, err
  }

  statement := fmt.Sprintf("INSERT INTO {{.OrgSQLSafe}}.{{.Name}}({{.Name}}) VALUES('%s') RETURNING {{.NameExported}}UUID", jb)
  rows, er1 := db.Query(statement)

  if er1 != nil {
    log.Printf("Insert failed for: %s", t.{{.NameExported}}UUID)
    log.Printf("SQL Error: %s", er1)
    return "", er1
  }

  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.{{.NameExported}}UUID)
    if err != nil {
      return "", err
    }
  }

  return t.{{.NameExported}}UUID, nil

}

// list{{.NameExported}}: return a list of {{.Name}}
//
func (t *{{.Name}}) list{{.NameExported}}(db *sql.DB, start, count int) ([]listResponse, error) {
/*
    qry := `select uuid,
          {{.Name}} ->> 'active' as active,
          {{.Name}} -> 'Metadata' ->> 'name' as name
          from {{.OrgSQLSafe}}.{{.Name}} LIMIT %d OFFSET %d;`
*/
    qry := `select {{.NameExported}}UUID
          from {{.OrgSQLSafe}}.{{.Name}} LIMIT %d OFFSET %d;`
  statement := fmt.Sprintf(qry, count, start)
  rows, err := db.Query(statement)

  if err != nil {
    return nil, err
  }

  defer rows.Close()

  ul := []listResponse{}

  for rows.Next() {
    var t listResponse
    err := rows.Scan(&t.UUID)

    if err != nil {
      log.Printf("SQL rows.Scan failed: %s", err)
      return ul, err
    }

    ul = append(ul, t)
  }

  return ul, nil
}

// get{{.NameExported}}: return a {{.Name}} based on the key
//
func (t *{{.Name}}) get{{.NameExported}}(db *sql.DB, key string, method int) error {
    var statement string

  switch method {
  case UUID:
    _, err := uuid.Parse(key)
    if err != nil {
      m := fmt.Sprintf("400: invalid UUID: %s", key)
      return errors.New(m)
    }
    statement = fmt.Sprintf(`
  SELECT {{.NameExported}}UUID, {{.Name}}
  FROM {{.OrgSQLSafe}}.{{.Name}}
  WHERE {{.NameExported}}UUID = '%s';`, key)
  }

  row := db.QueryRow(statement)

  // Fill in mapper
  var jb []byte
  var uid string
  switch err := row.Scan(&uid, &jb); err {

  case sql.ErrNoRows:
    m := fmt.Sprintf("404:Name %s does not exist", key)
    return errors.New(m)
  case nil:
    err = json.Unmarshal(jb, t)
    if err != nil {
      m := fmt.Sprintf("400:Unmarshal failed %s", key)
      return errors.New(m)
    }
    t.{{.NameExported}}UUID = uid
    break
  default:
	  log.Printf("500:get{{.Name}} Select failed; Got (%v)\n", err.Error())
    return(err)
  }

  return nil
}

// delete{{.NameExported}}: return a {{.Name}} based on UUID
//
func (t *{{.Name}}) delete{{.NameExported}}(db *sql.DB, key string) error {
	statement := fmt.Sprintf("DELETE FROM {{.OrgSQLSafe}}.{{.Name}} WHERE {{.NameExported}}UUID = '%s'", key)
  result, err := db.Exec(statement)
  c, e := result.RowsAffected()

  if e == nil && c == 0 {
    em := fmt.Sprintf("UUID %s does not exist", key)
    log.Println(em)
    log.Println(e)
    return errors.New(em)
  }

  return err
}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
