{{ define "rproxy.tpl"}}
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: eo-{{.Info.Name | ToLower}}-reverse-proxy
  namespace: pavedroad
  annotations:
    kubernetes.io/ingress.class: nginx
spec:
  rules:
  - host: api.pavedroad.io
    http:
      paths:
      - path: /api/v1/namespace/{{.Project.Kubernetes.Namespace}}/{{.Info.Name | ToLower}}
        pathType: Prefix
        backend:
          service:
            name: {{.Info.Name}}
            port:
              number: 8082
      - path: /api/v1/namespace/{{.Project.Kubernetes.Namespace}}/{{.Info.Name | ToLower}}LIST
        pathType: Prefix
        backend:
          service:
            name: {{.Info.Name | ToLower}}
            port:
              number: 8082
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
