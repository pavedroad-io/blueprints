{{define "keyed_method.tpl"}}

// {{.Method}}{{.EndPointName | ToCamel}} swagger:route {{.Method}} /api/v1/namespace/{{.Namespace}}/{{.EndPointName | ToCamel}} {{.EndPointName}} {{.Method | ToLower}}{{.EndPointName}}
//
// Returns a {{.EndPointName}} object
//
// Responses:
//    default: genericError
//        200: {{.EndPointName | ToCamel}}List


func (a *{{.EndPointName | ToCamel }}) list{{.EndPointName}}(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    // Pre-processing hook
    list{{.EndPointName}}PreHook(w, r, count, start)


    // Post-processing hook
    list{{.EndPointName}}PostHook(w, r)

    respondWithByte(w, http.StatusOK, jl)
}
{{end}}
