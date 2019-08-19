{{define "templateApp.go"}}
// Pavedroad license / copyright information
{{.{{.orginazation}}-info}}

// User project / copyright / usage information
{{.project-info}}

package main

import (
	"database/sql"
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

// Return list of {{.name}}s
//
// swagger:response tokenList
type listResponse struct {
	// in: body
  mappingList []{{.name-exported}} `json:"{{.name}}s"`
}

// Generated structures with Swagger docs
{{.swagger-generated-structs}}

// An {{.name-exported}} response model
//
// This is used for returning a response with a single mapper as body
//
// swagger:response mapperResponse
type {{.name-exported}}Response struct {
	// in: body
	response string `json:"order"`
}

// update{{.name-exported}}{{.name-exported}} in database
func (t *{{.name-exported}}) update{{.name-exported}}{{.name-exported}}(db *sql.DB) error {
	update := `
	UPDATE {{.orginzation}}.{{.name}}
    SET {{.name}} = '%s'
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

// create{{.name-exported}}{{.name-exported}} in database
func (t *{{.name-exported}}) create{{.name-exported}}{{.name-exported}}(db *sql.DB) error {
  jb, err := json.Marshal(t)
  if err != nil {
    panic(err)
  }

  statement := fmt.Sprintf("INSERT INTO {{.orginazation}}.{{.name}}({{.name}}) VALUES('%s') RETURNING uuid", jb)
  rows, er1 := db.Query(statement)

  if er1 != nil {
    log.Printf("Insert failed for: %s", t.Metadata.Name)
    log.Printf("SQL Error: %s", er1)
    return "", er1
  }

  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.Metadata.UUID)
    if err != nil {
      return "", err
    }
  }

  return t.Metadata.UUID, nil

}

// get{{.name-exported}}{{.name-exported}}s: return a list of tokens
//
func (t *{{.name-exported}}) get{{.name-exported}}{{.name-exported}}s(db *sql.DB, start, count int) ([]{{.name-exported}}, error) {
    qry := `select uuid,
          {{.name}} ->> 'active' as active,
          {{.name}} -> 'Metadata' ->> 'name' as name
          from {{.orginazation}}.{{.name}} LIMIT %d OFFSET %d;`
  statement := fmt.Sprintf(qry, count, start)
  rows, err := db.Query(statement)

  if err != nil {
    return nil, err
  }

  defer rows.Close()

  ul := []userList{}

  for rows.Next() {
    var t userList
    err := rows.Scan(&t.UUID, &t.Active, &t.Name)

    if err != nil {
      log.Printf("SQL rows.Scan failed: %s", err)
      return ul, err
    }

    ul = append(ul, t)
  }

  return ul, nil
}

// get{{.name-exported}}{{.name-exported}}: return a token based on credential
//
func (t *{{.name-exported}}) get{{.name-exported}}{{.name-exported}}(db *sql.DB, key string) error {
    var statement string

  switch method {
  case UUID:
    statement = fmt.Sprintf(`
  SELECT uuid, {{.name}}
  FROM {{.orginazation}}.{{.name}}
  WHERE uuid = '%s';`, key)
  case NAME:
    statement = fmt.Sprintf(`
  SELECT uuid, {{.name}}
  FROM {{.orginazation}}.{{.name}}
  WHERE {{.name}} -> 'Metadata' ->> 'name' = '%s';`, key)
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
    t.Metadata.UUID = uid
    break
  default:
    //Some error to catch
    panic(err)
  }

  return nil
}

// delete{{.name-exported}}{{.name-exported}}: return a token based on UID
//
func (t *{{.name-exported}}) delete{{.name-exported}}{{.name-exported}}(db *sql.DB, cred string) error {
	statement := fmt.Sprintf("DELETE FROM {{.orginzation}}.{{.name}} WHERE uuid = '%s'", uuid)
  result, err := db.Exec(statement)
  c, e := result.RowsAffected()

  if e == nil && c == 0 {
    em := fmt.Sprintf("UUID %s does not exist", uuid)
    log.Println(em)
    log.Println(e)
    return errors.New(em)
  }

  return err
}
{{end}}
