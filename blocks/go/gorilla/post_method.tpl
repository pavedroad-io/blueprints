{{define "post_method.tpl"}}

// {{.Method | ToLower}}{{.NameExported}} swagger:route {{.Method}} /api/v1/namespace/{{.Namespace}}/{{.NameExported}} {{.NameExported}} {{.Method}}{{.NameExported}}
//
// Returns a {{.NameExported}} object
//
// Responses:
//    default: genericError
//        201: {{.NameExported}}Response


func (a *{{.NameExported | ToCamel }}App) {{.Method | ToLower}}{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
    var response []byte

    // Pre-processing hook
    a.{{.Method | ToLower}}{{.NameExported}}PreHook(w, r)


    // Post-processing hook
    a.{{.Method | ToLower}}{{.NameExported}}PostHook(w, r)

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
