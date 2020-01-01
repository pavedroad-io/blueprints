{{define "templateResult.go"}}package main

// Result for a given job
type Result interface {
	// Return the orginal job
	Job() Job

	// Return the header/message headers
	Package() map[string]string

	// Return the payload/response data
	Payload() []byte
}{{end}}
