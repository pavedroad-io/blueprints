{{define "make-stagging-repo-app.tpl"}}
#!/bin/bash
aws ecr create-repository --repository-name io.pavedroad.stagging/{{.Info.Name | ToLower}} --region us-west-1 > ../docs/staggging-repo-{{.Info.Name | ToLower}}.json
{{end}}
