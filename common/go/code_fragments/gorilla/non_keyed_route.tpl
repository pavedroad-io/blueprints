{{define "none_keyed_route.tpl"}}

    uri := "/api/"+{{.APIVersion}} + "/" +
           "namespace"+{{.Namespace}} +"/"+
           {{.EndPointName}}
    a.Router.HandleFunc(uri, a.{{.Method|ToLower}}{{.endPointName}}.Methods("{{.Method}}")
    log.Println("{{.Method}}: ", uri)

{{end}}
