{{ define "kustomization.tpl"}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: {{.Project.Kubernetes.Namespace}}

bases:
  - kafka
  - cockroach
  - {{.Info.Name}}

resources:
  - namespace.yaml

commonLabels:
  pavedroad.env: dev

commonAnnotations:
  pavedroad.kustomize.base: {{.Info.Name}}/manifests/kubernetes/dev
  pavedroad.kustomize.bases: "{{.Info.Name}},db,kafka"
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
