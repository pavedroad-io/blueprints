{{define "template-service-stag-aws.tpl"}}
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
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
