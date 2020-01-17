{{define "httpResult.go"}}
package main

// Result for a given job
type httpResult struct {
	job      *httpJob
	metaData map[string]string
	payload  []byte
}

func (r *httpResult) Job() Job {
	return r.job
}

// TODO: consider changing this to map[string]interface{}
func (r *httpResult) MetaData() map[string]string {
	return r.metaData
}

func (r *httpResult) Payload() []byte {
	return r.payload
}{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
