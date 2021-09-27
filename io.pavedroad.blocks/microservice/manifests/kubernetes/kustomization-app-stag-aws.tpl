{{ define "kustomization-app-stag-aws.tpl"}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: {{.Project.Kubernetes.Namespace}}

bases:
  - {{.Info.Name}}
  - kafka
  - ingress

resources:
  - namespace.yaml

commonLabels:
  pavedroad.env: stag

commonAnnotations:
  pavedroad.kustomize.base: {{.Info.Name}}/manifests/kubernetes/stag
  pavedroad.kustomize.bases: "{{.Info.Name}}"
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
