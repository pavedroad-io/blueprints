{{ define "response.go"}}
package main

// genericError
//
// swagger:response genericError
type genericError struct {
	// in: body
	// Error message
	Body struct {
		Error string `json:"error"`
	} `json:"body"`
}

// genericResponse
//
// swagger:response genericResponse
type genericResponse struct {
	// in: body
	Body struct {
		// JSON body
		JSONBody string `json:"json_body"`
	} `json:"body"`
}

// get404Response Not found
//
// swagger:response get404Response
type get404Response struct {
	// The 404 error message
	// in: body

	// in: body
	Body struct {
		// Error message
		// Error message
		Error string `json:"error"`

		// UUID / key that was not found
		KEY string `json:"key"`
	} `json:"body"`
}

// Return list of {{.Name}}s
//
//
// swagger:response {{.NameExported}}List
type listResponse struct {
  // in: body
  UUID  string `json:"key"`
}

// {{.NameExported}} response model
//
// This is used for returning a response with a single {{.Name}} as body
//
// swagger:response {{.NameExported}}Response
type {{.NameExported}}Response struct {
    // in: body
    response string `json:"order"`
}

// metricsResponse
//
// swagger:response metricsResponse
type metricsResponse struct {
	// in: body
	Body struct {
		// Error message
		SchedulerMetrics string `json:"scheduler_metrics"`
		DispatherMetrics string `json:"dispather_metrics"`
	} `json:"body"`
}
{{end}}
