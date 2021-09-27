{{define "microservice-kustomization.tpl"}}
resources:
  - {{.Info.Name}}-deployment.yaml
  - {{.Info.Name}}-service.yaml

commonLabels:
  pavedroad.service: {{.Info.Name | ToLower}}

commonAnnotations:
  pavedroad.roadctl.version: alphav1
  pavedroad.roadctl.web: www.pavedroad.io
  pavedroad.roadctl.support: support@pavedroad.io

configMapGenerator:
- name: {{.Info.Name | ToLower}}-configmap
  literals:
  - database-ip=roach-ui
  - ip=0.0.0.0
  - port=8081
  - prlog-auto-init=true
  - prlog-conf-type=env
  - prlog-enable-kafka=true
  - prlog-file-format=text
  - prlog-file-location=logs/{{.Info.Name | ToLower}}.log"
  - prlog-kafka-brokers=kafka:9092

{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}