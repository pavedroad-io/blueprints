{{define "roach-ui-deployment.tpl"}}
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    pavedroad.service: roach-ui
  name: roach-ui
spec:
  replicas: 1
  strategy:
    type: Recreate
  selector:
    matchLabels:
      pavedroad.service: roach-ui
  template:
    metadata:
      creationTimestamp: null
      labels:
        pavedroad.service: roach-ui
    spec:
      containers:
      - args:
        - start-single-node
        - --insecure
        image: cockroachdb/cockroach
        name: roach-ui
        ports:
        - containerPort: 26257
        - containerPort: 8080
        resources: {}
        volumeMounts:
        - mountPath: /cockroach/cockroach-data
          name: roach-ui-claim0
      restartPolicy: Always
      volumes:
      - name: roach-ui-claim0
        persistentVolumeClaim:
          claimName: roach-ui-claim0
status: {}
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
