{{define "manifests/Dockerfile"}}
# User golang based off al Alpine
FROM golang:latest

LABEL "io.pavedroad.vendor": "{{.Organization}}" \
      "io.pavedroad.microservice": "{{.Name}}" \
      "io.pavedroad.description": "{{.ProjectInfo}}" \
      "io.pavedroad.version": "{{.Version}}" \
      "io.pavedroad.tempalte": "{{.TplName}}" \
      "io.pavedroad.definition": "{{.DefFile}}" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV ms {{.Name}}
ENV pavedroad /pavedroad
ENV pavedroadbin $pavedroad/$ms

# make working directory, move to it, and copy in prToken ms
RUN mkdir -p ${pavedroad}/logs
WORKDIR ${pavedroad}
COPY $ms $pavedroad

EXPOSE 8081
CMD ["/bin/sh", "-c", "$pavedroadbin"]
{{end}}
