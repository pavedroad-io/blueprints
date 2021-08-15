{{define "hook.go"}}{{.PavedroadInfo}}

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
func (a *{{.NameExported}}App) getExplainPreHook(w http.ResponseWriter, r *http.Request) (body []byte, err error) {

  return body, nil
}


{{.EndpointHooks}}

// getLivenessPreHook
//
func getLivenessPreHook(w http.ResponseWriter, r *http.Request) {
	return
}

// getReadinessPreHook
//
func getReadinessPreHook(w http.ResponseWriter, r *http.Request) {
	return
}

// getMetricsPreHook
//
func (a *{{.NameExported}}App) getMetricsPreHook(w http.ResponseWriter, r *http.Request) {
	return
}

// getMetricsPostHook
//
func (a *{{.NameExported}}App) getMetricsPostHook(w http.ResponseWriter, r *http.Request) {
	return
}

// getManagementPreHook
//
func (a *{{.NameExported}}App) getManagementPreHook(w http.ResponseWriter, r *http.Request) {
  return
}

// getManagementPostHook
//
func (a *{{.NameExported}}App) getManagementPostHook(w http.ResponseWriter, r *http.Request) {
  return
}

// putManagementPreHook
//
func (a *{{.NameExported}}App) putManagementPreHook(w http.ResponseWriter, r *http.Request) {
  return
}

// putManagementPostHook
//
func (a *{{.NameExported}}App) putManagementPostHook(w http.ResponseWriter, r *http.Request) {
  return
}

{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
