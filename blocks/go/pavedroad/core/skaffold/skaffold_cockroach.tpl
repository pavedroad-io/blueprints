{{define "skaffold_cockroach.tpl"}}
  apiVersion: skaffold/v2beta10
  kind: Config
  build:
    insecureRegistries:
      - localhost:32000
    tagPolicy:
      sha256: {}
    artifacts:
    - image: localhost:32000/{{.Organization}}/{{.Name}}
      context: .
      custom:
        dependencies:
          paths:
            - "**.go"
    - image: localhost:32000/{{.Organization}}/{{.Name}}initdb
      context: .
      docker:
        dockerfile: manifests/InitDbDockerFile
  deploy:
    kustomize:
      paths:
_     - "manifests/kubernetes/dev"
  profiles:
    - name: dev-debug
      activation:
      - env: GODEBUG=true
      deploy:
        kustomize:
          paths:
          - "manifests/kubernetes/dev-debug"
{{end}}
