{{define "metricInterface.go"}}package main

// Metric returns metrics for a given function
type Metric interface {
	// Return the original job
	Get() []byte

	ResetAll()

	Reset(specific string) error

}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
