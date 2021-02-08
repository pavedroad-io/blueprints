{{define "templateHook.go"}}{{.PavedroadInfo}}

// {{.ProjectInfo}}

// {{.Name}}Hooks.go
//   Allows users to add their own business logic
//   This file will not be over written if a blueprint is
//   regenerated.  If a function signature is changed,
//   you can update this file by running: $ go fix

package main

import (
  "net/http"
)


// Will generate the func calls below based on HTTP verbs selected in the definitions file
func (a *{{.NameExported}}App) getExplainPreHook(w http.ResponseWriter, r *http.Request, count, start int) (body []byte, error err)
  var body []byte

  return body, nil
}


{{.EndpointHooks}}

// TODO: Placeholders for generated code
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

// post{{.NameExported}}PreHook
//
func (a *{{.NameExported}}App) post{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request) {
  return
}

// post{{.NameExported}}PostHook
//
func (a *{{.NameExported}}App) post{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request) {
  return
}

// put{{.NameExported}}PreHook
func (a *{{.NameExported}}App) put{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// put{{.NameExported}}PostHook
func (a *{{.NameExported}}App) put{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// patch{{.NameExported}}PreHook
func (a *{{.NameExported}}App) patch{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// patch{{.NameExported}}PostHook
func (a *{{.NameExported}}App) patch{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// delete{{.NameExported}}PreHook
func (a *{{.NameExported}}App) delete{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}
// delete{{.NameExported}}PostHook
func (a *{{.NameExported}}App) delete{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
