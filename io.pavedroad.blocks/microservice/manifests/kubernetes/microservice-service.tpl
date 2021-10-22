{{define "microservice-service.tpl"}}
apiVersion: v1
kind: Service
metadata:
  name: {{.Info.Name | ToLower}}
spec:
  ports:
  - name: "{{.Project.Config.HTTPPort}}"
    port: {{.Project.Config.HTTPPort}}
    targetPort: {{.Project.Config.HTTPPort}}
  selector:
    pavedroad.service: {{.Info.Name | ToLower}}
  type: NodePort
status:
  loadBalancer: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
