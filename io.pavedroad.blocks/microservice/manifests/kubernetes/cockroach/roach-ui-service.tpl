{{define "roach-ui-service.tpl"}}
apiVersion: v1
kind: Service
metadata:
  labels:
    pavedraod.service: roach-ui
  name: roach-ui
spec:
  ports:
  - name: "26257"
    port: 26257
    targetPort: 26257
  - name: "6060"
    port: 6060
    targetPort: 8080
  selector:
    pavedroad.service: roach-ui
status:
  loadBalancer: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
