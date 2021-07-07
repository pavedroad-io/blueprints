{{define "skaffold.tpl"}}
apiVersion: skaffold/v2beta10
kind: Config
build:
  insecureRegistries:
    - localhost:32000
  tagPolicy:
    sha256: {}
  artifacts:
  - image: localhost:32000/{{.Info.Organization}}/{{.Info.Name | ToLower}}
    context: .
    custom:
      dependencies:
        paths:
          - "**.go"
  - image: localhost:32000/{{.Info.Organization}}/{{.Info.Name | ToLower}}initdb
    context: .
    docker:
      dockerfile: manifests/InitDbDockerFile
deploy:
  kustomize:
    paths:
    - "manifests/kubernetes/dev"
profiles:
  - name: dev-debug
    activation:
    - env: GODEBUG=true
    deploy:
      kustomize:
        paths:
        - "manifests/kubernetes/dev-debug"
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
