{{define "list_route.tpl"}}

    uri = "/api/" + "{{.APIVersion}}" + "/" +
           "namespace" + "/" + "{{.Namespace}}" + "/" +
           "{{.NameExported}}LIST"
    a.Router.HandleFunc(uri, a.{{.Method | ToLower}}{{.NameExported}}).Methods("GET")
    log.Println("GET: ", uri)

{{end}}
