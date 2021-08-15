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
    a.{{.Method | ToLower }}{{.NameExported}}PreHook(w, r, key)


    // Post-processing hook
    a.{{.Method | ToLower }}{{.NameExported}}PostHook(w, r, key)

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
