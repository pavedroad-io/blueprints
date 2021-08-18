{{define "template-service-stag-aws.tpl"}}
apiVersion: v1
kind: Service
metadata:
  name: {{.Info.Name | ToLower}}
spec:
  ports:
  - name: "8081"
    port: 8081
    targetPort: 8081
  selector:
    pavedroad.service: {{.Info.Name | ToLower}}
  type: NodePort
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
