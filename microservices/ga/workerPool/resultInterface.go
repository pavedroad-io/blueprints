{{define "resultInterface.go"}}package main

// Result for a given job
type Result interface {
	// Return the original job
	Job() []byte

	// Decode () (Job, error)
	Decode() (Job, error)

	// Return the header/message headers
	MetaData() map[string]string

	// Return the payload/response data
	Payload() []byte
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
