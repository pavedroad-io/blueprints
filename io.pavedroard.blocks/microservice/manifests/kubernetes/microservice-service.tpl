{{define "microservice-service.tpl"}}
apiVersion: v1
kind: Service
metadata:
  name: {{.Info.Name}}
spec:
  ports:
  - name: "8081"
    port: 8081
    targetPort: 8081
  selector:
    pavedroad.service: {{.Info.Name}}
  type: NodePort
status:
  loadBalancer: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
