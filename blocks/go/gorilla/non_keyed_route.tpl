{{define "non_keyed_route.tpl"}}

    uri = "/api/"+ "{{.APIVersion}} " + "/" +
           "namespace" + "/" + "{{.Namespace}}" +"/"+
           {{.NameExported}}
    a.Router.HandleFunc(uri, a.{{.Method|ToLower}}{{.NameExported}}).Methods("{{.Method}}")
    log.Println("{{.Method}}: ", uri)

{{end}}
