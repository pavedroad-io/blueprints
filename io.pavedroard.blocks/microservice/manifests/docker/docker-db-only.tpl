{{define "docker-db-only.tpl"}}
version: '3'

services:
  roach-ui:
    image: cockroachdb/cockroach
    command: start-single-node --insecure
    ports:
     - "26257:26257"
     - "6060:8080"
    volumes:
     - ../volumes/data/db-1:/cockroach/cockroach-data
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
