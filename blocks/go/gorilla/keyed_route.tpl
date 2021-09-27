{{define "keyed_route.tpl"}}

    uri = "/api/" + "{{.APIVersion}}" + "/" +
           "namespace" + "/" + "{{.Namespace}}" + "/" +
           "{{.NameExported}}" + "/" + "{key}"
    a.Router.HandleFunc(uri, a.{{.Method | ToLower}}{{.NameExported}}).Methods("{{.Method}}")
    log.Println("{{.Method}}: ", uri)

{{end}}
