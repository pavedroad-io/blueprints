{{define "keyed_route.tpl"}}

    uri := "/api/"+ "{{.APIVersion}}" + "/" +
           "namespace" + "/" + "{{.Namespace}}" + "/" +
           {{.EndPointName}} + "/" + "{key}"
    a.Router.HandleFunc(uri, a.{{.Method | ToLower}}{{.EndPointName}}.Methods("{{.Method}}"))
    log.Println("{{.Method}}: ", uri)

{{end}}
