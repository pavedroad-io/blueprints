{{define "options_route.tpl"}}

	uri = "/api/" + "{{.APIVersion}}" + "/" +
		"namespace" + "/" + "{{.Namespace}}" + "/" +
		"{{.NameExported}}"
		a.Router.HandleFunc(uri, a.{{.Method | ToLower}}{{.NameExported}}).Methods("OPTIONS")
	log.Println("OPTIONS: ", uri)

	uri = "/api/" + "{{.APIVersion}}" + "/" +
		"namespace" + "/" + "{{.Namespace}}" + "/" +
		"{{.NameExported}}" + "/" + "{key}"
		a.Router.HandleFunc(uri, a.{{.Method | ToLower}}{{.NameExported}}).Methods("OPTIONS")
	log.Println("OPTIONS: ", uri)

{{end}}
