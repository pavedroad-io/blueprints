{{define "microservice-kustomization.tpl"}}
resources:
  - {{.Info.Name}}-deployment.yaml
  - {{.Info.Name}}-service.yaml

commonLabels:
  pavedroad.service: {{.Info.Name}}

commonAnnotations:
  pavedroad.roadctl.version: alphav1
  pavedroad.roadctl.web: www.pavedroad.io
  pavedroad.roadctl.support: support@pavedroad.io

configMapGenerator:
- name: {{.Info.Name}}-configmap
  literals:
  - database-ip=roach-ui
  - ip=0.0.0.0
  - port=8081
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
