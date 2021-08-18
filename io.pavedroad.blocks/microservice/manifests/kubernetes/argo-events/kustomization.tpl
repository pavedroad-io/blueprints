{{ define "kustomization.tpl"}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: argo-events

resources:
  - namespace.yaml
  - install.yaml
  - install-validating-webhook.yaml

commonLabels:
  pavedroad.env: staging

commonAnnotations:
  pavedroad.kustomize.base: {{.Info.Name}}/manifests/kubernetes/stag
  pavedroad.kustomize.bases: argo-events
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
