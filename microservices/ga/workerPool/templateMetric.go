{{define "templateMetric.go"}}package main

// Metric returns metrics for a given function
type Metric interface {
	// Return the orginal job
	Get() []byte

	ResetAll()

	Reset(specific string) error

}{{end}}
