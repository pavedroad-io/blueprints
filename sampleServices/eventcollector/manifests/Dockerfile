
# User golang based off al Alpine
FROM golang:latest

LABEL "io.pavedroad.vendor": "mirantis" \
      "io.pavedroad.microservice": "eventcollector" \
      "io.pavedroad.description": "Allocate jobs to workers in a pool" \
      "io.pavedroad.version": "0.0.1" \
      "io.pavedroad.tempalte": "workerPool" \
      "io.pavedroad.definition": "metrics.yaml" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV ms eventcollector
ENV pavedroad /pavedroad
ENV pavedroadbin $pavedroad/$ms

# make working directory, move to it, and copy in prToken ms
RUN mkdir -p ${pavedroad}/logs
WORKDIR ${pavedroad}
COPY $ms $pavedroad

EXPOSE 8081
CMD ["/bin/sh", "-c", "$pavedroadbin"]
