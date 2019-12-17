{{define "templateHook.go"}}{{.PavedroadInfo}}

// {{.ProjectInfo}}

// {{.Name}}Hooks.go
//   Allows users to add their own business logic
//   This file will not be over written if a template is
//   regenerated.  If a function signature is changed,
//   you can update this file by running: $ go fix

package main

import (
  "net/http"
)

// list{{.NameExported}}PreHook
//
func (a *{{.NameExported}}App) list{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, count, start int) {
  return
}

// list{{.NameExported}}PostHook
//
func (a *{{.NameExported}}App) list{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request) {
  return
}

// get{{.NameExported}}PreHook
//
func (a *{{.NameExported}}App) get{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// get{{.NameExported}}PostHook
//
func (a *{{.NameExported}}App) get{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// create{{.NameExported}}PreHook
//
func (a *{{.NameExported}}App) create{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request) {
  return
}

// create{{.NameExported}}PostHook
//
func (a *{{.NameExported}}App) create{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request) {
  return
}

// update{{.NameExported}}PreHook
func (a *{{.NameExported}}App) update{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// update{{.NameExported}}PostHook
func (a *{{.NameExported}}App) update{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// delete{{.NameExported}}PreHook
func (a *{{.NameExported}}App) delete{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}
// delete{{.NameExported}}PostHook
func (a *{{.NameExported}}App) delete{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}{{end}}
