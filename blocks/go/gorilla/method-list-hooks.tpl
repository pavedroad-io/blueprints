{{define "method-list-hooks.tpl"}}

// {{.Method | ToLower}}{{.NameExported}}PreHook
func (a *{{.NameExported}}App) {{.Method | ToLower}}{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, count, start int) {
  return
}

// {{.Method | ToLower}}{{.NameExported}}PostHook
func (a *{{.NameExported}}App) {{.Method | ToLower}}{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request) {
  return
}

{{end}}

