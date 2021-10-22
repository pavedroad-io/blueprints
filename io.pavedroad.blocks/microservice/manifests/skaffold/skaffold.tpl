{{define "skaffold.tpl"}}
apiVersion: skaffold/v2beta10
kind: Config
build:
  insecureRegistries:
    - localhost:32000
  tagPolicy:
    sha256: {}
  artifacts:
  - image: localhost:32000/{{.Info.GitHubOrg}}/{{.Info.Name | ToLower}}
    context: .
    docker:
      dockerfile: manifests/Dockerfile
  - image: localhost:32000/{{.Info.GitHubOrg}}/{{.Info.Name | ToLower}}initdb
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
  - name: staging
    build:
      artifacts:
      - image: 400276217548.dkr.ecr.us-west-1.amazonaws.com/io.pavedroad.staging/{{.Info.Name | ToLower}}
        context: .
        docker:
          dockerfile: manifests/Dockerfile
      - image: 400276217548.dkr.ecr.us-west-1.amazonaws.com/io.pavedroad.staging/{{.Info.Name | ToLower}}initdb
        context: .
        docker:
          dockerfile: manifests/InitDbDockerFile
    deploy:
      kustomize:
        paths:
        - "manifests/kubernetes/stag"
(END)

{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
