{{define "dev/make-stagging-repo.sh"}}
#!/bin/bash
aws ecr create-repository --repository-name io.pavedroad.stagging/{{.Name}}initdb --region us-west-1 > ../docs/staggging-repo-{{.Name}}initdb.json
aws ecr create-repository --repository-name io.pavedroad.stagging/{{.Name}} --region us-west-1 > ../docs/staggging-repo-{{.Name}}.json
{{end}}
