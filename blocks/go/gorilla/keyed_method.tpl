{{define "keyed_method.tpl"}}

// {{.Method | ToLower}}{{.EndPointName}} swagger:route {{.Method}} /api/v1/namespace/{{.Namespace}}/{{.EndPointName}} {{.EndPointName}} {{.Method | ToLower}}{{.EndPointName}}
//
// Returns a {{.EndPointName}} object
//
// Responses:
//    default: genericError
//        200: {{.EndPointName | ToCamel}}Response
//        400: genericError


func (a *{{.EndPointName | ToCamel}}App) {{.Method | ToLower}}{{.EndPointName}}(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["key"]
    var response []byte

    // Pre-processing hook
    a.{{.Method | ToLower }}{{.EndPointName | ToCamel}}PreHook(w, r, key)


    // Post-processing hook
    a.{{.Method | ToLower }}{{.EndPointName | ToCamel}}PostHook(w, r, key)

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
