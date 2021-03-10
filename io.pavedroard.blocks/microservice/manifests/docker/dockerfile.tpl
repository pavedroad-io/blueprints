{{define "dockerfile.tpl"}}
# User golang based off al Alpine
FROM golang:latest

LABEL "io.pavedroad.vendor": "{{.Info.Organization}}" \
      "io.pavedroad.microservice": "{{.Info.Name | ToLower}}" \
      "io.pavedroad.description": "{{.Project.Description}}" \
      "io.pavedroad.version": "{{.Info.Version}}" \
      "io.pavedroad.tempalte": "{{.Info.ID}}" \
      "io.pavedroad.definition": "{{.DefinitionFile}}" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing the microservice
# docker images must be lowercase
ENV ms {{.Info.Name | ToLower}}
ENV pavedroad /pavedroad
ENV pavedroadbin $pavedroad/$ms

# make working directory, move to it, and copy in the microservice
RUN mkdir -p ${pavedroad}/logs
WORKDIR ${pavedroad}
COPY $ms $pavedroad

EXPOSE 8081
CMD ["/bin/sh", "-c", "$pavedroadbin"]
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
