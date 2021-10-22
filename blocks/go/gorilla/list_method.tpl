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
    count := 10
    start := 0
    var response []byte

    c := r.URL.Query().Get("count")
    if c != "" {
        count, _ = strconv.Atoi(c)
    }
    s := r.URL.Query().Get("start")
    if s != "" {
        start, _ = strconv.Atoi(s)
    }


    // Pre-processing hook
    if final := a.list{{.NameExported}}PreHook(w, r, count, start); final {
	return
	}


    // Post-processing hook
    if final := a.list{{.NameExported}}PostHook(w, r); final {
	return
	}

    respondWithByte(w, http.StatusOK, response)
}
{{end}}
