{{define "kustomization-stag-aws.tpl"}}
resources:
  - {{.Info.Name}}-deployment-stag-aws.yaml
  - {{.Info.Name}}-service-stag-aws.yaml

commonLabels:
  pavedroad.service: {{.Info.Name | ToLower}}
  pavedroad.env: stag

commonAnnotations:
  pavedroad.roadctl.version: alphav1
  pavedroad.roadctl.web: www.pavedroad.io
  pavedroad.roadctl.support: support@pavedroad.io

configMapGenerator:
- name: {{.Info.Name | ToLower}}-configmap
  literals:
  - ip=0.0.0.0
  - port={{.Project.Config.HTTPPort}}
  - prlog-auto-init=true
  - prlog-conf-type=env
  - prlog-enable-kafka=false
  - prlog-file-format=text
  - prlog-file-location=logs/{{.Info.Name | ToLower}}.log"
  - prlog-kafka-brokers=kafka:9092

{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
