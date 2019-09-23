{{define "Makefile"}}
#-include .env

VERSION := {{.Version}}
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
PROJDIR := $(shell pwd)
TARGET := $(PROJECTNAME)

K8SRUNNING := $(shell dev/microk8sStatus.sh)
K8SSTART := $(shell dev/microk8sStart.sh)
K8SSTOP := $(shell dev/microk8sStop.sh)

ASSETS := $(PROJDIR)/assets/images
ARTIFACTS := $(PROJDIR)/artifacts
BUILDS := $(PROJDIR)/builds
DOCS := $(PROJDIR)/docs
LOGS := $(PROJDIR)/logs
MANIFEST := $(PROJDIR)/manifests
K8S = ~"$(PROJDIR)/manifest/kubernetes"

# Go related variables.
GOBASE := $(shell cd ../../;pwd)
GOPATH := $(GOBASE)
export GOPATH = $(GOBASE)
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOLINT := $(shell which golint)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
GOCOVERAGE := $(ARTIFACTS)/coverage.out
GOLINTREPORT := $(ARTIFACTS)/lint.out
GOTESTREPORT := https://sonarcloud.io/dashboard?id={{.Organization}}_{{.Name}}

GIT_TAG := $(shell git describe)

SHELL := /bin/bash

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.GitTag=$(GIT_TAG)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := ./$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := /tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

.PHONY: check build compile sonar-scanner

all: compile check

#.DEFAULT_GOAL: $(TARGET)

## compile: Compile the binary.
compile: $(LOGS) $(ARTIFACTS) $(ASSETS) $(DOCS) $(BUILDS)
	@echo "  Compiling"
	@-$(MAKE) -s build

## clean: Remove dep, vendor, binary(s), and executs go clean
clean:
	@echo "  execute go-clean"
	@-rm $(GOBIN)/$(PROJECTNAME)* 2> /dev/null || true
	@-rm -R vendor Gopkt.* 2> /dev/null || true
	@-$(MAKE) go-clean

## build: Build the binary for linux / mac x86 and amd
build: go-get api-doc $(K8S) build
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(GOFILES)
# make this conditional on build GOARCH
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS="darwin" GOARCH="amd64" go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(GOFILES)
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS="darwin" GOARCH="386" go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-"darwin"-"386" $(GOFILES)
	cp $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH)
	cp $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(PROJECTNAME)
	cp $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(BUILDS)/$(PROJECTNAME)-"darwin"-"amd64"
	cp $(GOBIN)/$(PROJECTNAME)-"darwin"-"386" $(BUILDS)/$(PROJECTNAME)-"darwin"-"386"


## deploy: Deploy image to repository and k8s cluster
deploy:
ifeq "$(K8SRUNNING)" "down"
	@echo "  >  Starting k8s for deployment..."
	dev/microk8sStart.sh
endif
	skaffold run -f manifests/skaffold.yaml

Gopkg.toml:
	@echo "  >  initialize dep support..."
	$(shell (export GOPATH=$(GOPATH);dep init))

go-get: Gopkg.toml get-deps $(ASSETS)
	@echo "  >  Creating dependencies graph png..."
	$(shell (export GOPATH=$(GOPATH);dep status -dot | dot -T png -o $(ASSETS)/$(PROJECTNAME).png))

api-doc: $(DOCS)
	@echo "  >  Generate swagger specification..."
	$(shell (export GOPATH=$(GOPATH);swagger generate spec -m -t $(PROJECTNAME) -o docs/api.json))
	@echo "  >  Generate HTML..."
	pretty-swag -i docs/api.json -o docs/api.html
	@echo "  >  Done"

get-deps:
	@echo "  >  dep ensure..."
	$(shell (GOPATH=$(GOPATH);dep ensure $?))

## install: Install packages or main
install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

# TODO: enable sonar scanner once we get API key
# check: lint sonar-scanner $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)
## check: Start services and execute static code analysis and tests
check: lint $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)
	@echo "  >  Starting cockroach DB..."
	docker-compose -f manifests/docker-db-only.yaml up -d
	@echo "  >  running to tests..."
	go test -coverprofile=$(GOCOVERAGE) -v ./...
	@echo "  >  Stopping cockroach DB..."
	docker-compose -f manifests/docker-db-only.yaml down

sonar-scanner: $(ARTIFACTS)
	sonarcloud.sh

## show-coverage: Show go code coverage in browser
show-coverage:
	go tool cover -html=$(GOCOVERAGE)

## show-test: Show sonarcloud test report
show-test:
	xdg-open $(GOTESTREPORT)

## show-devkit: Show documenation for Devkit
show-devkit:
	xdg-open http://localhost:5000/microk8sDevKit.html


lint: $(GOFILES)
	@echo "  >  running lint..."
	@echo $?
	$(GOLINT) $? > $(GOLINTREPORT)

$(K8S):
	$(shell mkdir -p manifests/kubernetes)
	$(shell cd manifests/kubernetes;kompose convert -f ../docker-compose.yaml)

## fmt: Run gofmt on all code
fmt: $(GOFILES)
	@gofmt -l -w $?

## simplify: Run gofmt with simplify option
simplify: $(GOFILES)
	@gofmt -s -l -w $?

## k8s-start: Start local microk8s server and update configurations
k8s-start:
	$(K8SSTART)

## k8s-stop: Stop local k8s cluster and delete skaffold deployments
k8s-stop:
	skaffold delete -f manifests/skaffold.yaml
	$(K8SSTOP)

## k8s-status: Print the status of the local cluster up or down
k8s-status:
	@echo -n "  >  microk8s is "
	@echo $(K8SRUNNING)

## help: Print possible commands
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

$(ASSETS):
	@echo "  >  Creating assets directory"
	$(shell mkdir -p $(ASSETS))

$(ARTIFACTS):
	@echo "  >  Creating artifacts directory"
	$(shell mkdir -p $(ARTIFACTS))

$(BUILDS):
	@echo "  >  Creating $(BUILDS) directory"
	$(shell mkdir -p $(BUILDS))

$(DOCS):
	@echo "  >  Creating docs directory"
	$(shell mkdir -p $(DOCS))

$(LOGS):
	@echo "  >  Creating logs directory"
	$(shell mkdir -p $(LOGS))
{{end}}
