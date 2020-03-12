

// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root
// for full license information.
//
// Apache2

// User project / copyright / usage information
// Manage database of films

package main

import (
	"database/sql"
  "encoding/json"
  "github.com/google/uuid"
	"errors"
	"fmt"
  "time"
	"log"
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

// Return list of filmss
//
// TODO: add method of including sub attributes
//
// swagger:response filmsList
type listResponse struct {
  // in: body
  UUID  string `json:"uuid"`
}

// Generated structures with Swagger docs
// swagger:response metadata
type metadata struct {
// Author
	Author string	`json:"author"`
// Genre
	Genre string	`json:"genre"`
// Rating
	Rating string	`json:"rating"`
}

// swagger:response films
type films struct {
// FilmsUUID into JSONB

	FilmsUUID string `json:"filmsuuid"`
	Metadata metadata	`json:"metadata"`
// Id
	Id string	`json:"id"`
// Title
	Title string	`json:"title"`
// Updated
	Updated time.Time	`json:"updated"`
// Created
	Created time.Time	`json:"created"`
}



// Films response model
//
// This is used for returning a response with a single films as body
//
// swagger:response filmsResponse
type FilmsResponse struct {
	// in: body
	response string `json:"order"`
}

// updateFilms in database
func (t *films) updateFilms(db *sql.DB, key string) error {
	update := `
	UPDATE AcmeDemo.films
    SET films = '%s'
  WHERE filmsUUID = '%s';`

  jb, err := json.Marshal(t)
  if err != nil {
    log.Println("marshal failed")
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

// createFilms in database
func (t *films) createFilms(db *sql.DB) (string, error) {
  jb, err := json.Marshal(t)
  if err != nil {
    panic(err)
  }

  statement := fmt.Sprintf("INSERT INTO AcmeDemo.films(films) VALUES('%s') RETURNING FilmsUUID", jb)
  rows, er1 := db.Query(statement)

  if er1 != nil {
    log.Printf("Insert failed for: %s", t.FilmsUUID)
    log.Printf("SQL Error: %s", er1)
    return "", er1
  }

  defer rows.Close()

  for rows.Next() {
    err := rows.Scan(&t.FilmsUUID)
    if err != nil {
      return "", err
    }
  }

  return t.FilmsUUID, nil

}

// listFilms: return a list of films
//
func (t *films) listFilms(db *sql.DB, start, count int) ([]listResponse, error) {
/*
    qry := `select uuid,
          films ->> 'active' as active,
          films -> 'Metadata' ->> 'name' as name
          from AcmeDemo.films LIMIT %d OFFSET %d;`
*/
    qry := `select FilmsUUID
          from AcmeDemo.films LIMIT %d OFFSET %d;`
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

// getFilms: return a films based on the key
//
func (t *films) getFilms(db *sql.DB, key string, method int) error {
    var statement string

  switch method {
  case UUID:
    _, err := uuid.Parse(key)
    if err != nil {
      m := fmt.Sprintf("400: invalid UUID: %s", key)
      return errors.New(m)
    }
    statement = fmt.Sprintf(`
  SELECT FilmsUUID, films
  FROM AcmeDemo.films
  WHERE FilmsUUID = '%s';`, key)
  }

  row := db.QueryRow(statement)

  // Fill in mapper
  var jb []byte
  var uid string
  switch err := row.Scan(&uid, &jb); err {

  case sql.ErrNoRows:
    m := fmt.Sprintf("404:name %s does not exist", key)
    return errors.New(m)
  case nil:
    err = json.Unmarshal(jb, t)
    if err != nil {
      m := fmt.Sprintf("400:Unmarshal failed %s", key)
      return errors.New(m)
    }
    t.FilmsUUID = uid
    break
  default:
    //Some error to catch
    panic(err)
  }

  return nil
}

// deleteFilms: return a films based on UUID
//
func (t *films) deleteFilms(db *sql.DB, key string) error {
	statement := fmt.Sprintf("DELETE FROM AcmeDemo.films WHERE FilmsUUID = '%s'", key)
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
