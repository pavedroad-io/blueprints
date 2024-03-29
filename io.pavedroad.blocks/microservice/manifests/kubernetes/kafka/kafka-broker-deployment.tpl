{{define "kafka-broker-deployment.tpl"}}
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: kafka
    component: kafka-broker
  name: kafka-broker
spec:
  replicas: 1
  selector:
    matchLabels:
      app: kafka
      component: kafka-broker
  template:
    metadata:
      labels:
        app: kafka
        component: kafka-broker
    spec:
      containers:
      - name: kafka
        image: wurstmeister/kafka:latest
        ports:
        - containerPort: 9092
        env:
          - name: KAFKA_PORT
            value: "9092"
          - name: KAFKA_ADVERTISED_PORT
            value: "9092"
          - name: KAFKA_ADVERTISED_HOST_NAME
            valueFrom:
              fieldRef:
                fieldPath: status.podIP
          - name: KAFKA_ZOOKEEPER_CONNECT
            value: zookeeper:2181
{{end}}{{/* vim: set filetype=gotexttmpl: */ -}}
