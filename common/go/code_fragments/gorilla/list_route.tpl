{{define "list_route.tpl"}}

    uri := "/api/"+{{.APIVersion}} + "/" +
           "namespace" + "/" +{{.Namespace}} + "/" +
           "{{.EndPointName}}" + "/" + "{{.EndPointName}}LIST"
    a.Router.HandleFunc(uri, a.{{.Method | ToLower}}{{.EndPointName}}.Methods("GET")
    log.Println("GET: ", uri)

{{end}}
