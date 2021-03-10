{{define "kafka-broker-service.tpl"}}
apiVersion: v1
kind: Service
metadata:
  name: kafka
  labels:
    app: kafka
    component: kafka-broker
spec:
  ports:
  - port: 9092
    name: kafka-port
    targetPort: 9092
    protocol: TCP
  selector:
    app: kafka
    component: kafka-broker
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}