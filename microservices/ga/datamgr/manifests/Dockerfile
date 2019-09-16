{{define "Dockerfile"}}
# User golang based off al Alpine
FROM golang:latest

LABEL "vendor": "{{.Organization}}" \
      "microservice": "{{.Name}}" \
      "description": "{{.ProjectInfo}}" \
      "version": "{{.Version}}" \
      "tempalte": "{{.TplName}}" \
      "definition": "{{.DefFile}}" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV ms {{.Name}}
ENV pavedroad /pavedroad
ENV pavedroadbin $pavedroad/$ms

# make working directory, move to it, and copy in prToken ms
RUN mkdir ${pavedroad}
WORKDIR ${pavedroad}
COPY $ms $pavedroad

EXPOSE 8081
CMD ["/bin/sh", "-c", "$pavedroadbin"]
{{end}}
