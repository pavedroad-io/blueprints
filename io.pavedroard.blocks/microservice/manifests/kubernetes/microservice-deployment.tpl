{{define "microservice-deployment.tpl"}}
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Info.Name | ToLower}}
spec:
  replicas: 3
  strategy: {}
  selector:
    matchLabels:
      pavedroad.service: {{.Info.Name | ToLower}}
  template:
    metadata:
      creationTimestamp: null
      labels:
        pavedroad.service: {{.Info.Name | ToLower}}
    spec:
      initContainers:
      - image: busybox:1.28
        name: wait-for-cockroach
        command: ['sh', '-c', 'until nslookup roach-ui; do echo waiting for roach-ui; sleep 2; done;']
      - image: localhost:32000/{{.Info.Organization}}/{{.Info.Name | ToLower}}initdb:0.0
        env:
        - name: COCKROACH_HOST
          valueFrom:
            configMapKeyRef:
              name: cockroach-configmap
              key: host-ip
        name: {{.Info.Name | ToLower}}dbinit
      containers:
      - image: localhost:32000/{{.Info.Organization}}/{{.Info.Name | ToLower}}:0.0
        env:
        - name: HTTP_IP_ADDR
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: ip
        - name: HTTP_IP_PORT
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: port
        - name: APP_DB_IP
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: database-ip
        - name: PRLOG_AUTOINIT
          value: "true"
        - name: PRLOG_CFGTYPE
          value: "env"
        - name: PRLOG_ENABLEKAFKA
          value: "true"
        - name: PRLOG_FILEFORMAT
          value: "text"
        - name: PRLOG_FILELOCATION
          value: "logs/{{.Info.Name | ToLower}}.log"
        - name: PRKAFKA_BROKERS
          value: "kafka:9092"
        name: {{.Info.Name | ToLower}}
        ports:
        - containerPort: 8081
        resources: {}
      restartPolicy: Always
status: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
