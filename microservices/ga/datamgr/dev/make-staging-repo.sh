{{define "dev/make-staging-repo.sh"}}
#!/bin/bash
aws ecr create-repository --repository-name io.pavedroad.staging/{{.Name}}initdb --region us-west-1 > ../docs/staging-repo-{{.Name}}initdb.json
aws ecr create-repository --repository-name io.pavedroad.staging/{{.Name}} --region us-west-1 > ../docs/staging-repo-{{.Name}}.json
{{end}}
