{{define "list_method.tpl"}}
// list{{.NameExported}} swagger:route GET /api/v1/namespace/{{.Namespace}}/{{.NameExported}}LIST {{.NameExported}} list{{.NameExported}}
//
// Returns a list of {{.NameExported}}
//
// Responses:
//    default: genericError
//        200: {{.NameExported}}List
//        400: genericError


func (a *{{.NameExported | ToCamel }}App)list{{.NameExported}}(w http.ResponseWriter, r *http.Request) {
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
    a.list{{.NameExported}}PreHook(w, r, count, start)


    // Post-processing hook
    a.list{{.NameExported}}PostHook(w, r)

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
