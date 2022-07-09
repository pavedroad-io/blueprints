{{define "method-hooks.tpl"}}

// {{.Method | ToLower}}{{.NameExported}}PreHook
func (a *{{.NameExported}}App) {{.Method | ToLower}}{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request) bool {
  return false
}

// {{.Method | ToLower}}{{.NameExported}}PostHook
func (a *{{.NameExported}}App) {{.Method | ToLower}}{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request) bool {
  return false
}

{{end}}

