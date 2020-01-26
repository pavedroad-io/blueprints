
//
// Copyright (c) PavedRoad. All rights reserved.
// Licensed under the Apache2. See LICENSE file in the project root for full license information.
//

// Manage database of films

// filmsHooks.go
//   Allows users to add their own business logic
//   This file will not be over written if a template is
//   regenerated.  If a function signature is changed,
//   you can update this file by running: $ go fix

package main

import (
  "net/http"
)

// listFilmsPreHook
//
func (a *FilmsApp) listFilmsPreHook(w http.ResponseWriter, r *http.Request, count, start int) {
  return
}

// listFilmsPostHook
//
func (a *FilmsApp) listFilmsPostHook(w http.ResponseWriter, r *http.Request) {
  return
}

// getFilmsPreHook
//
func (a *FilmsApp) getFilmsPreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// getFilmsPostHook
//
func (a *FilmsApp) getFilmsPostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// createFilmsPreHook
//
func (a *FilmsApp) createFilmsPreHook(w http.ResponseWriter, r *http.Request) {
  return
}

// createFilmsPostHook
//
func (a *FilmsApp) createFilmsPostHook(w http.ResponseWriter, r *http.Request) {
  return
}

// updateFilmsPreHook
func (a *FilmsApp) updateFilmsPreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// updateFilmsPostHook
func (a *FilmsApp) updateFilmsPostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}

// deleteFilmsPreHook
func (a *FilmsApp) deleteFilmsPreHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}
// deleteFilmsPostHook
func (a *FilmsApp) deleteFilmsPostHook(w http.ResponseWriter, r *http.Request, key string) {
  return
}