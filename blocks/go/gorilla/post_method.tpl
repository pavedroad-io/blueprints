{{define "post_method.tpl"}}

// {{.Method | ToLower}}{{.EndPointName | ToCamel}} swagger:route {{.Method}} /api/v1/namespace/{{.Namespace}}/{{.EndPointName | ToCamel}} {{.EndPointName}} {{.Method}}{{.EndPointName}}
//
// Returns a {{.EndPointName}} object
//
// Responses:
//    default: genericError
//        201: {{.EndPointName | ToCamel}}Response


func (a *{{.EndPointName | ToCamel }}App) {{.Method | ToLower}}{{.EndPointName}}(w http.ResponseWriter, r *http.Request) {
    var response []byte

    // Pre-processing hook
    a.{{.Method | ToLower}}{{.EndPointName | ToCamel}}PreHook(w, r)


    // Post-processing hook
    a.{{.Method | ToLower}}{{.EndPointName | ToCamel}}PostHook(w, r)

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
