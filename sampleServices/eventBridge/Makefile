#-include .env

VERSION := 0.0.1
BUILD := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")
PROJDIR := $(shell pwd)
TARGET := $(PROJECTNAME)

K8SRUNNING := $(shell dev/microk8sStatus.sh)

ASSETS := $(PROJDIR)/assets/images
ARTIFACTS := $(PROJDIR)/artifacts
BUILDS := $(PROJDIR)/builds
DOCS := $(PROJDIR)/docs
LOGS := $(PROJDIR)/logs
VOLS := $(PROJDIR)/volumes
FOSSATEST := .fossa.ymml
PREFLIGHT := .pr_preflight_check

# Go related variables.
GOBASE := $(shell cd ../../;pwd)
GOPATH := $(GOBASE)
export DOCKER_IP = `(dev/getdockerip.sh)`
GOBIN := $(GOBASE)/bin
GOFILES := $(wildcard *.go)
GOLINT := $(shell which golint)
GOARCH := $(shell go env GOARCH)
GOOS := $(shell go env GOOS)
GOCOVERAGE := $(ARTIFACTS)/coverage.out
GOLINTREPORT := $(ARTIFACTS)/lint.out
GOSECREPORT := $(ARTIFACTS)/gosec.out
GOVETREPORT := $(ARTIFACTS)/govet.out
GOTESTREPORT := https://sonarcloud.io/dashboard?id=PavedRoad_eventbridge

GIT_TAG := $(shell git describe)

SHELL := /bin/bash

# Use linker flags to provide version/build settings to the target
LDFLAGS=-ldflags "-X=main.Version=$(VERSION) -X=main.Build=$(BUILD) -X=main.GitTag=$(GIT_TAG)"

# Make is verbose in Linux. Make it silent.
# MAKEFLAGS += --silent

.PHONY: check build compile sonar-scanner

all: mod-setup $(PREFLIGHT) $(FOSSATEST) compile check

## compile: Compile the binary.
compile: $(LOGS) $(ARTIFACTS) $(ASSETS) $(DOCS) $(BUILDS) api-doc mod-graph
	@echo "  Compiling"
	@-$(MAKE) -s build

## clean: Remove binary(s) and execute go clean
clean:
	@echo "  execute go-clean"
	@-rm $(GOBIN)/$(PROJECTNAME)* 2> /dev/null || true
	@-$(MAKE) go-clean

## build: Build the binary for linux / mac x86 and amd
	$(shell (grep -q '"definitions": {' docs/api.json || sed -i -e 's/"swagger": "2.0",/"swagger": "2.0",\n"definitions": { "Streams": { "type": "object", "allOf": [ { "properties": {} } ] } },/' docs/api.json))
build: mod-setup
	@echo "  >  Building binary..."
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS=$(GOOS) GOARCH=$(GOARCH) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(GOFILES)
# make this conditional on build GOARCH
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) GOOS="darwin" GOARCH="amd64" go build $(LDFLAGS) -o $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(GOFILES)
	cp $(GOBIN)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH)
	cp $(BUILDS)/$(PROJECTNAME)-$(GOOS)-$(GOARCH) $(PROJECTNAME)
	cp $(GOBIN)/$(PROJECTNAME)-"darwin"-"amd64" $(BUILDS)/$(PROJECTNAME)-"darwin"-"amd64"

mod-setup:
	@echo "  >  setting up modules..."
	go mod download
	go mod tidy

mod-graph: mod-setup $(ASSETS)
	@echo "  >  Creating dependencies graph png..."
	$(shell (go mod graph | modgraphviz | dot -T png -o $(ASSETS)/$(PROJECTNAME).png))

## deploy: Deploy image to repository and k8s cluster
deploy:
ifeq "$(K8SRUNNING)" "down"
	@echo "  >  Starting k8s for deployment..."
	dev/microk8sStart.sh
endif
	@echo "  >  Starting k8s is up..."
	@dev/kube-config.sh
#       wait for the registry service to be ready	
	@echo "  >  Wait for registry to come up..."
	@sleep 20
	@echo "  >  Build image and deploy..."
	@skaffold run -f manifests/skaffold.yaml

## deploy-debug: Deploy image to the k8s cluster in headless debug mode
deploy-debug:
ifeq "$(K8SRUNNING)" "down"
	@echo "  >  Starting k8s for deployment..."
	dev/microk8sStart.sh
endif
	@echo "  >  Starting k8s is up..."
	@dev/kube-config.sh
#       wait for the registry service to be ready	
	@echo "  >  Wait for registry to come up..."
	@sleep 20
	@echo "  >  Build image and deploy..."
	@echo "  >  Start image headles for remote debugging..."
	@skaffold run -f manifests/skaffold.yaml -p dev-debug

## deploy-down: Delete and cleanup deployment from the k8s cluster
deploy-down:
	@echo "  >  Cleanup deployment..."
	@skaffold delete -f manifests/skaffold.yaml

## dev-mode: Start deployment in debug mode, watch for file changes and perform a live update
dev-mode:
	@echo "  >  Running deployment in development mode..."
	@echo "  >  Use CTRL-C to exit..."
	skaffold dev -f manifests/skaffold.yaml

## tilt-up: Start service using Tilt UI
tilt-up: dbclean
	@echo "  >  Start with Tilt UI..."
	@tilt up

## tilt-down: Stop and cleanup a Tilt service
tilt-down:
	@echo "  >  Cleanup deployment..."
	@tilt down

api-doc: $(DOCS)
	@echo "  >  Generate swagger specification..."
	$(shell (export GOPATH=$(GOPATH);swagger generate spec -m -t $(PROJECTNAME) -o docs/api.json))
	@echo "  >  Generate HTML..."
	pretty-swag -i docs/api.json -o docs/api.html
	@echo "  >  Done"

## install: Install packages or main
install:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go clean

## check: Start services and execute static code analysis and tests
check: lint docker-build sonar-scanner $(ARTIFACTS) $(LOGS) $(ASSETS) $(DOCS)
	@echo "  >  starting kafka..."
	docker-compose -f manifests/kafka.yaml up -d --remove-orphans
	@echo "  >  running to tests..."
	go test -coverprofile=$(GOCOVERAGE) -v ./...
	@echo "  >  stoping kafka..."
	docker-compose -f manifests/kafka.yaml down

sonar-scanner: $(ARTIFACTS)
	sonar-scanner

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
	@echo -n "  >  running lint..."
	@echo $?
	$(GOLINT) $? > $(GOLINTREPORT)
	@echo "  >  running gosec... > $(GOSECREPORT)"
	$(shell (export GOPATH=$(GOPATH);gosec -fmt=sonarqube -tests -out $(GOSECREPORT) -exclude-dir=.blueprints ./...))
	@echo "  >  running go vet... > $(GOVETREPORT)"
	$(shell (export GOPATH=$(GOPATH);go vet ./... 2> $(GOVETREPORT)))

	@echo "  >  running FOSSA license scan."
	$(shell (export GOPATH=$(GOPATH); @FOSSA_API_KEY=$(FOSSA_API_KEY) fossa analyze))


## fmt: Run gofmt on all code
fmt: $(GOFILES)
	@gofmt -l -w $?

## simplify: Run gofmt with simplify option
simplify: $(GOFILES)
	@gofmt -s -l -w $?

## k8s-start: Start local microk8s server and update configurations
k8s-start:
	@echo "  > dev/microk8sStart.sh"
	dev/microk8sStart.sh

## k8s-stop: Stop local k8s cluster and delete skaffold deployments
k8s-stop:
	skaffold delete -f manifests/skaffold.yaml
	dev/microk8sStop.sh

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

## docker-build: Build docker images for use with docker-compose
docker-build:
	docker build -f manifests/Dockerfile -t acme-demo/eventbridge .
	docker build -f manifests/Dockerfile -t acme-demo/eventbridgeinitdb .

## up: Start service using docker-compose
up: dbclean docker-build
	docker-compose -f manifests/docker-compose.yaml up -d --remove-orphans


## down: Stop service using docker-compose
down:
	docker-compose -f manifests/docker-compose.yaml down

$(ASSETS):
	$(shell mkdir -p $(ASSETS))

$(ARTIFACTS):
	$(shell mkdir -p $(ARTIFACTS))

$(BUILDS):
	$(shell mkdir -p $(BUILDS))

$(DOCS):
	$(shell mkdir -p $(DOCS))

$(LOGS):
	$(shell mkdir -p $(LOGS))

# Null target for roadctl
dbclean:
	@echo ""

## Preflight ensure all requirements are met
## Once they are the $(PREFLIGH) file will be created
$(PREFLIGHT):
	dev/preflight.sh


$(FOSSATEST):
	fossa init

