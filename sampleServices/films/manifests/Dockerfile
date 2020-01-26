
# User golang based off al Alpine
FROM golang:latest

LABEL "io.pavedroad.vendor": "acme-demo" \
      "io.pavedroad.microservice": "films" \
      "io.pavedroad.description": "Manage database of films" \
      "io.pavedroad.version": "1.0.0" \
      "io.pavedroad.tempalte": "datamgr" \
      "io.pavedroad.definition": "films.yaml" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV ms films
ENV pavedroad /pavedroad
ENV pavedroadbin $pavedroad/$ms

# make working directory, move to it, and copy in prToken ms
RUN mkdir -p ${pavedroad}/logs
WORKDIR ${pavedroad}
COPY $ms $pavedroad

EXPOSE 8081
CMD ["/bin/sh", "-c", "$pavedroadbin"]
