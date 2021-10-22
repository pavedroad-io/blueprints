{{define "template-deployment-stag-aws.tpl"}}
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
      containers:
      - image: 400276217548.dkr.ecr.us-west-1.amazonaws.com/io.pavedroad.stagging/{{.Info.Name | ToLower}}
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
        - name: PRLOG_AUTOINIT
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: prlog-auto-init
        - name: PRLOG_CFGTYPE
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: prlog-conf-type
        - name: PRLOG_ENABLEKAFKA
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: prlog-enable-kafka
        - name: PRLOG_FILEFORMAT
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: prlog-file-format
        - name: PRLOG_FILELOCATION
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: prlog-file-location
        - name: PRKAFKA_BROKERS
          valueFrom:
            configMapKeyRef:
              name: {{.Info.Name | ToLower}}-configmap
              key: prlog-kafka-brokers
        name: {{.Info.Name | ToLower}}
        ports:
        - containerPort: {{.Project.Config.HTTPPort}}
        resources: {}
      restartPolicy: Always
status: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
