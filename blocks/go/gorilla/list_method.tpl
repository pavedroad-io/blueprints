{{define "list_method.tpl"}}
// list{{.EndPointName | ToLower}} swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.EndPointName | ToCamel}}LIST {{.EndPointName}} list{{.EndPointName}}
//
// Returns a list of {{.EndPointName}}
//
// Responses:
//    default: genericError
//        200: {{.EndPointName | ToCamel}}List
//        400: genericError


func (a *{{.EndPointName | ToCamel }}App)list{{.EndPointName}}(w http.ResponseWriter, r *http.Request) {
    count, _ := strconv.Atoi(r.FormValue("count"))
    start, _ := strconv.Atoi(r.FormValue("start"))
    var response []byte

    if count > 10 || count < 1 {
        count = 10
    }
    if start < 0 {
        start = 0
    }

    // Pre-processing hook
    a.list{{.EndPointName | ToCamel}}PreHook(w, r, count, start)


    // Post-processing hook
    a.list{{.EndPointName | ToCamel}}PostHook(w, r)

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
