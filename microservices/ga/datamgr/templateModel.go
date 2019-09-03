{{define "templateModel.go"}}
{{.PavedroadInfo}}
// {{.OrganizationLicense}}

// User project / copyright / usage information
// {{.ProjectInfo}}

package main

import (
	"database/sql"
  "encoding/json"
	"errors"
	"fmt"
	"log"
)

// A GenericError is the default error message that is generated.
// For certain status codes there are more appropriate error structures.
//
// swagger:response genericError
type GenericError struct {
	// in: body
	Body struct {
		Code    int32 `json:"code"`
		Message error `json:"message"`
	} `json:"body"`
}

// Return a basic message as json in the body
//
// swagger:response statusResponse
type statusResponse struct {
	// in: body
  msg error `json:"message"`
}

// Return list of {{.Name}}s
//
// TODO: add method of including subattributes
//
// swagger:response {{.Name}}List
type listResponse struct {
  UUID  string `json:"uuid"`
}

// Generated structures with Swagger docs
{{.SwaggerGeneratedStructs}}

// An {{.NameExported}} response model
//
// This is used for returning a response with a single mapper as body
//
// swagger:response mapperResponse
type {{.NameExported}}Response struct {
	// in: body
	response string `json:"order"`
}

// update{{.Name}}{{.NameExported}} in database
func (t *{{.Name}}) update{{.NameExported}}(db *sql.DB, key string) error {
	update := `
	UPDATE {{.Organization}}.{{.Name}}
    SET {{.Name}} = '%s'
  WHERE UUID = '%s';`

  jb, err := json.Marshal(t)
  if err != nil {
    log.Println("marshall failed")
    panic(err)
  }

  statement := fmt.Sprintf(update, jb, key)
  _, er1 := db.Query(statement)

  if er1 != nil {
    log.Println("Update failed")
    return er1
  }

  return nil
}

// create{{.Name}}{{.NameExported}} in database
func (t *{{.Name}}) create{{.NameExported}}(db *sql.DB) (string, error) {
  jb, err := json.Marshal(t)
  if err != nil {
    panic(err)
  }

  statement := fmt.Sprintf("INSERT INTO {{.Organization}}.{{.Name}}({{.Name}}) VALUES('%s') RETURNING uuid", jb)
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

// get{{.Name}}{{.NameExported}}: return a list of {{.Name}}
//
func (t *{{.Name}}) list{{.NameExported}}(db *sql.DB, start, count int) ([]listResponse, error) {
/*
    qry := `select uuid,
          {{.Name}} ->> 'active' as active,
          {{.Name}} -> 'Metadata' ->> 'name' as name
          from {{.Organization}}.{{.Name}} LIMIT %d OFFSET %d;`
*/
    qry := `select uuid,
          from {{.Organization}}.{{.Name}} LIMIT %d OFFSET %d;`
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

// get{{.NameExported}}{{.NameExported}}: return a {{.Name}} based on the key
//
func (t *{{.Name}}) get{{.NameExported}}(db *sql.DB, key string, method int) error {
    var statement string

  switch method {
  case UUID:
    statement = fmt.Sprintf(`
  SELECT uuid, {{.Name}}
  FROM {{.Organization}}.{{.Name}}
  WHERE uuid = '%s';`, key)
  /*
  case NAME:
    statement = fmt.Sprintf(`
  SELECT uuid, {{.Name}}
  FROM {{.Organization}}.{{.Name}}
  WHERE {{.Name}} -> 'Metadata' ->> 'name' = '%s';`, key)
  */
  }
  row := db.QueryRow(statement)

  // Fill in mapper
  var jb []byte
  var uid string
  switch err := row.Scan(&uid, &jb); err {

  case sql.ErrNoRows:
    m := fmt.Sprintf("name %s does not exist", key)
    return errors.New(m)
  case nil:
    err = json.Unmarshal(jb, t)
    if err != nil {
      m := fmt.Sprintf("unmarshal failed %s", key)
      return errors.New(m)
    }
    t.{{.NameExported}}UUID = uid
    break
  default:
    //Some error to catch
    panic(err)
  }

  return nil
}

// delete{{.Name}}{{.NameExported}}: return a {{.Name}} based on UID
//
func (t *{{.Name}}) delete{{.NameExported}}(db *sql.DB, key string) error {
	statement := fmt.Sprintf("DELETE FROM {{.Organization}}.{{.Name}} WHERE uuid = '%s'", key)
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
{{end}}
