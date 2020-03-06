{{define "README.md"}}
# Working with your new {{.NameExported}} microservice

## Preflighting

*New in version v0.6.0alpha for roadctl*

The preflight.sh script ensures that the following required components
are installed:

- An initialized git repository exists
- The user.name and user.email are set in your git configuration
- A GitHub Personal Access Token is available in the GH_ACCESS_TOKEN env variable
- A .gitignore file with entries
- An initialized templates repository

It doesn't generate a git tag; pick your own. Once preflight.sh
successfully executes, a ".pr_preflight_check" is created.  Until this
file exists, Make will always attempt to run preflight.sh.

If you are running a version prior to v0.6.0alpha, you may run preflight.sh
by hand in the top level of your go micro-service package with:

```bash
dev/preflight.sh
```

### Versioning information
The make file sets three versioning variables; VERSION, BUILD, and GIT_TAG.  These are passed go the go compiler and printed when the -v flag is passed on the command line.  Output is formatted as JSON:

    $ films -v
    {"Version": "1.0.0", "Build": "8755e7f", "GitTag": "v0.0alpha"}

VERSION := 1.0.0
----------------
The version variable is set based on the value you enter in your definitions file.

BUILD := $(shell git rev-parse --short HEAD)
--------------------------------------------
The build is set to the commit ID of your current git HEAD.

GIT_TAG := $(shell git describe)
--------------------------------
GIT_TAG set using the most recent tag if any.  You can add a tag with:

    git tag -a "mytag" -m "message about the tag."
Git push doesn't include tags.  To push tags to the origin use:

    git push origin --tag

## roadctl
The roadctl command is modeled after kubectl.
Use "roadctl help" for a list of top level commands.
This may be different for a specific command or it's sub-commands.

### Top level help
    roadctl help

General help output:

```
roadctl allows you to work with the PavedRoad CNCF low-code environment and the associated CI

  Usage: roadctl [command] [TYPE] [NAME] [flags]

  TYPE specifies a resource type
  NAME is the name of a resource
  flags specify options

Usage:
  roadctl [command]

Available Commands:
  apply       Apply configuration to named resource
  completion  Generate completion scripts on stdout
  config      Manage roadctl global configuration options
  create      Create a new resource
  delete      Delete a resource
  deploy      Deploy a service
  describe    Describe provides detailed information about a resource
  doc         Generate documentation for your service
  edit        Edit the configuration for the specified resource
  events      View events
  explain     Return documentation about a resource
  get         Get an existing object
  help        Help about any command
  init        Initialize roadctl development environment
  logs        Return logs for a resource
  replace     Delete and recreate the named resource
  version     Print the current version

Flags:
      --config string      Config file (default is $HOME/.roadctl.yaml)
      --debug string       Debug level: info(default)|warm|error|critical (default "info")
      --format string      Output format: text(default)|json|yaml (default "text")
  -h, --help               help for roadctl
      --password string    HTTP basic auth password
      --templates string   Set the location of the directory holding roadctl templates
      --token string       OAUTH access token
      --user string        HTTP basic auth user name

Use "roadctl [command] --help" for more information about a command.
```

### Specific command
    roadctl get --help

Specific help output:

```
Return summary information about an existing resource

Usage:
  roadctl get [flags]

Flags:
  -h, --help   help for get

Global Flags:
      --config string      Config file (default is $HOME/.roadctl.yaml)
      --debug string       Debug level: info(default)|warm|error|critical (default "info")
      --format string      Output format: text(default)|json|yaml (default "text")
      --password string    HTTP basic auth password
      --templates string   Set the location of the directory holding roadctl templates
      --token string       OAUTH access token
      --user string        HTTP basic auth user name
```

## Generating your service
The roadctl CLI is used to create new services.
It has two fundamental concepts:

- templates: Contain logic need to generate a service
- definitions: Define your custom logic, integrations, and organizational information

A sample definitions is available to help you get started.

### Initialize template repository
    roadctl init

### List available templates
    roadctl get templates

### Create a copy of the sample definition
    roadctl describe templates datamgr > myservice.yaml

Note: edit myservice.yaml to customize your create below.

### Get definitions of attributes in your myservice.txt
    roadctl explain templates datamgr > myservice.txt

### Create your microservice
    roadctl create templates datamgr -f myservice.yaml
### Build and test
Executing make will compile and test your service.
Optionally, you can do `make compile` followed by `make check`.

    make

## Directories

| Name | Contents |
| --------- | -------- |
| artifacts | Outputs from static code analysis and tests |
| assets | Generate assets such as images |
| builds | Executables for supported platforms, Mac/Linux x86/amd64 |
| dev | Generated helper scripts and sample data |
| dev/db | Generated SQL statements |
| docs | Generated documentation |
| logs | Logs generated by the microservice |
| manifests | Docker and docker-composes manifest |
| manifests/kubernetes | Kubernetes manifests for deploying this microservice |
| vendor | Vendor dependencies |

## SQL
To get an SQL prompt, use:

    /bin/sql.sh

## dev/testXXXXX.sh scripts
The following scripts work with your local docker images using
docker-compose or with the local microk8s cluster.  By default they
use the local docker image.  To use the microk8s cluster, use the -k
command line option/flag.

- dev/testAll.sh
- dev/testPost.sh
- dev/testPut.sh
- dev/testGet.sh
- dev/testGetList.sh

## make
Use **make help** to get a list of options:

    make help

Help output:

```
  Choose a command run in films:

  compile         Compile the binary.
  clean           Remove dep, vendor, binary(s), and execute go clean
  build           Build the binary for linux / mac x86 and amd
  deploy          Deploy image to repository and k8s cluster
  install         Install packages or main
  check           Start services and execute static code analysis and tests
  show-coverage   Show go code coverage in browser
  show-test       Show sonarcloud test report
  show-devkit     Show documentation for Devkit
  fmt             Run gofmt on all code
  simplify        Run gofmt with simplify option
  k8s-start       Start local microk8s server and update configurations
  k8s-stop        Stop local k8s cluster and delete skaffold deployments
  k8s-status      Print the status of the local cluster up or down
  help            Print possible commands
```

## Skaffold CI/CD
Skaffold is integrated into your project.  You can use the following commands:

### development mode
Monitors source code and when it changes builds and pushes a new image

```bash
skaffold dev -f manifests/skaffold.yaml
```
### run
Build and push the image when executed

```bash
skaffold run -f manifests/skaffold.yaml
```

### delete
Deletes all deployed resources

```bash
skaffold delete -f manifests/skaffold.yaml
```

## Linters
Three lint applications are integrated to assist in code reviews.

- Go lint checks for conformance with effective go programming recommendations and Go code review suggestions.
- Gosec tests your code against go recommended security practices
- Govet inspects code for constructs that might break.
- FOSSA license scanner
- SonarCloud scanner

The location of each lint's output is below along with links to the rules they enforce.

### golint
artifacts/lint.out
[Effective Go](https://golang.org/doc/effective_go.html)
[Go code review comments](https://github.com/golang/go/wiki/CodeReviewComments)

### gosec
artifacts/gosec.out
[Rules](https://securego.io/docs/rules/rule-intro.html)

### go vet
artifacts/govet.out
[Go vet rules](https://golang.org/cmd/vet/)

### SonarCloud and FOSSA
docs/service.html
Badges for both with links to details can be found in the generated
service.html in the docs directory.

# SonarCloud
SonarCloud provides free code analysis for open-source projects.  By default,
the following tools are included:

- quality gate
- bugs
- code smells
- coverage
- lines of code
- duplicate lines of code
- security
- technical debt
- vulnerabilities


Support for SonarCloud is pre-integrated in the generated Makefile.

You need to set a valid sonarcloud token before executing make in
your .bashrc file:

export SONARCLOUD_TOKEN=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX

## Getting a token
Go to https://sonarcloud.io.  Then login with your GitHub account.

Next, go to https://sonarcloud.io/account/security/.
In the Generate Tokens dialog, enter a name for your token
and click the "Generate" button.


## sonar-project.properties

Controls the executing of an analysis run.  Documentation is
available [here](https://docs.sonarqube.org/latest/analysis/analysis-parameters/).
The default configuration provides extended support for code coverage and go lint reporting.


## Run by hand using
The sonarcloud.sh is provided for executing an analysis by hand.

```bash
sonarcloud.sh

```
# FOSSA
FOSSA provides free license scanning for open-source projects.   The [fossa-cli](https://github.com/fossas/fossa-cli/) documentation is covers basic usage.  Support for fossa is pre-integrated in the generated Makefile.  You need to set a valid fossa token before executing make in your .bashrc file:

export FOSSA_API_KEY=XXXXXXXXXXXXXXXXXXXXXXXXXX

## Getting a token

Go to https://app.fossa.com.  Then login with your GitHub account.

Next, go to https://app.fossa.com/account/settings/integrations/api_tokens.
Use the "Add another Token" button to create your token.


## Run by hand using

```bash
FOSSA_API_KEY=$(FOSSA_API_KEY) fossa analyze
```

# GitHub token

Go to https://github.com and login.  Then go to, https://github.com/settings/tokens.
Use the "Generate new token" button to create your new token.

```bash
# add line to your .bashrc
export ACCESS_TOKEN=####################
```
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
