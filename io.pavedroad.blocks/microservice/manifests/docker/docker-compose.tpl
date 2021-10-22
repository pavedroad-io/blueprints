{{define "docker-compose.tpl"}}
version: '3'

services:
  {{.Info.Name | ToLower}}:
    image: {{.Info.GitHubOrg}}/{{.Info.Name | ToLower}}
    expose:
     - "{{.Project.Config.HTTPPort}}"
    ports: 
     - {{.Project.Config.HTTPPort}}:{{.Project.Config.HTTPPort}}
    environment:
     - HTTP_IP_ADDR={{.Project.Config.HTTPHost}}
     - HTTP_IP_PORT={{.Project.Config.HTTPPort}}
     - APP_DB_IP=manifests_roach-ui_1
     - PRLOG_AUTOINIT=true
     - PRLOG_CFGTYPE=env
     - PRLOG_ENABLEKAFKA=true
     - PRLOG_FILEFORMAT=text
     - PRLOG_FILELOCATION=logs/{{.Info.Name | ToLower}}.log
     - PRKAFKA_BROKERS=kafka:9092
    depends_on:
     - kafka
  roach-ui:
    image: cockroachdb/cockroach
    command: start-single-node --insecure
    expose:
     - "8080"
     - "26257"
    ports:
     - "26257:26257"
     - "6060:8080"
    volumes:
     - /tmp/volumes/data/db-1:/cockroach/cockroach-data
  zookeeper:
    image: wurstmeister/zookeeper
    ports:
      - "2181:2181"
  kafka:
    image: wurstmeister/kafka
    ports:
      - "9092"
    expose:
      - "9092"
    environment:
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_HOST_NAME: kafka
      KAFKA_ADVERTISED_PORT: 9092
      KAFKA_CREATE_TOPICS: logs:1:1
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
