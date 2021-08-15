{{ define "kustomization-debug.tpl"}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

namespace: {{.Project.Kubernetes.Namespace}}-debug

patchesStrategicMerge:
  - debug.yaml

resources:
  - ../dev

{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
