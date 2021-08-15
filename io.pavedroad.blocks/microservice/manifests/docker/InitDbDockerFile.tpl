{{define "InitDbDockerFile.tpl"}}
# pavedroad image based on centos with cockroachdb installed
FROM pavedroadio/cockroachdb-client:0.3

LABEL "io.pavedroad.vendor": "{{.Info.Organization}}" \
      "io.pavedroad.init.db": "{{.Info.Name}}db" \
      "io.pavedroad.description": "{{.Project.Description}}" \
      "io.pavedroad.version": "{{.Info.Version}}" \
      "io.pavedroad.tempalte": "{{.Info.ID}}" \
      "io.pavedroad.definition": "{{.DefinitionFile}}" \
      "io.pavedroad.env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV scripts dev/db
ENV cmd {{.Info.Name}}ExecuteAll.sh

# make working directory, move to it, and copy in the microservice
RUN mkdir -p pavedroad
WORKDIR pavedroad
COPY dev/db .

CMD ["/bin/sh", "-c", "./$cmd -k"]
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
