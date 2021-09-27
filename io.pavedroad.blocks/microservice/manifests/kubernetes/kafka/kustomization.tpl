{{define "kustomization.tpl"}}
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
resources:
  - kafka-broker-deployment.yaml
  - kafka-broker-service.yaml
  - zookepper-deployment.yaml
  - zookepper-service.yaml
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
