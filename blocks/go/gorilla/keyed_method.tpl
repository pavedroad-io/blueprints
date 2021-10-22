{{define "keyed_method.tpl"}}

// {{.Method | ToLower}}{{.NameExported}} swagger:route {{.Method}} /api/v1/namespace/{{.Namespace}}/{{.NameExported}} {{.NameExported}} {{.Method | ToLower}}{{.NameExported}}
//
// Returns a {{.NameExported}} object
//
// Responses:
//    default: genericError
//        200: {{.NameExported | ToCamel}}Response
//        400: genericError


func (a *{{.NameExported | ToCamel}}App) {{.Method | ToLower}}{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]
	var response []byte

	// Pre-processing hook
	if final := a.{{.Method | ToLower }}{{.NameExported}}PreHook(w, r, key); final {
		return
	}



	// Post-processing hook
	if final := a.{{.Method | ToLower }}{{.NameExported}}PostHook(w, r, key); final {
		return
	}

	respondWithByte(w, http.StatusOK, response)
}
{{end}}
