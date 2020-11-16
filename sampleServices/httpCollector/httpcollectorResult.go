package main

// Result for a given job
type Result interface {
	// Return the original job
	Job() Job

	// Return the header/message headers
	MetaData() map[string]string

	// Return the payload/response data
	Payload() []byte
}