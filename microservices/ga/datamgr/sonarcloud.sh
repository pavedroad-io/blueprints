{{define "sonarcloud.sh"}}
sonar-scanner \
  -Dsonar.projectKey=PavedRoad_{{.Name}} \
  -Dsonar.organization={{.Organization}} \
  -Dsonar.sources=. \
  -Dsonar.host.url=https://sonarcloud.io \
  -Dsonar.login={{.SonarLogin}}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
