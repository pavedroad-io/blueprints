{{define "microservice-deployment.tpl"}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Info.Name}}
spec:
  replicas: 3
  strategy: {}
  selector:
    matchLabels:
      pavedroad.service: {{.Info.Name}}
  template:
    metadata:
      creationTimestamp: null
      labels:
        pavedroad.service: {{.Info.Name}}
    spec:
      initContainers:
      - image: busybox:1.28
        name: wait-for-cockroach
        command: ['sh', '-c', 'until nslookup roach-ui; do echo waiting for roach-ui; sleep 2; done;']
      - image: localhost:32000/{{.Info.Organization}}/{{.Info.Name}}initdb:0.0
        env:
        - name: COCKROACH_HOST
          valueFrom:
            configMapKeyRef:
              name: cockroach-configmap
              key: host-ip
        name: {{.Info.Name}}dbinit
      containers:
      - image: localhost:32000/{{.Info.Organization}}/{{.Info.Name}}:0.0
        env:
        - name: HTTP_IP_ADDR
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name}}-configmap
              key: ip
        - name: HTTP_IP_PORT
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name}}-configmap
              key: port
        - name: APP_DB_IP
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name}}-configmap
              key: database-ip
        name: {{.Info.Name}}
        ports:
        - containerPort: 8081
        resources: {}
      restartPolicy: Always
status: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
