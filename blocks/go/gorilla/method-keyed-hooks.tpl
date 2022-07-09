{{define "method-keyed-hooks.tpl"}}

// {{.Method | ToLower}}{{.NameExported}}PreHoobool
func (a *{{.NameExported}}App) {{.Method | ToLower}}{{.NameExported}}PreHook(w http.ResponseWriter, r *http.Request, key string) bool {
  return false
}

// {{.Method | ToLower}}{{.NameExported}}PostHook
func (a *{{.NameExported}}App) {{.Method | ToLower}}{{.NameExported}}PostHook(w http.ResponseWriter, r *http.Request, key string) bool {
  return false
}

{{end}}

