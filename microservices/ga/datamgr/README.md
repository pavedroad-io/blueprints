{{define "README.md"}}
# Working with your new {{.NameExported}} microservice

## git
The build system requires git source code management.  If the directory you choose to generate your service is not under git control, do the following after executing your template.
    git init
    vi .gitignore
    # add .templates/*
    # save your .gitignore
    git add *
    git commit

To utlize the sonar-scanner test functionality and its continuous integrations
features, you will also need a remote repository on github. See sonarcloud.io for information on establishing your client accout, and follow their 'Analayze from your CI' information to install the required scanner. 

Note that you should modify the YAML,sonarcloud.sh and sonar-project.propertiesfiles in your project directory to reflect your new sonar-scanner requirements outlined in their installation instructions. Minor adjustments related to sonar-scanner might also be required in the Makefile (GOTESTREPORT and PATH). 

### Versioning information
The make file sets three versioning variables; VERSION, BUILD, and GIT_TAG.  These are passed go the go compiler and printed when the -v flag is passed on the command line.  Output is formated as JSON.

    films -v
    {"Version": "1.0.0", "Build": "8755e7f", "GitTag": "v0.0alpha"}

VERSION := 1.0.0
----------------
The version variable is set based on the value you enter in your defintions file

BUILD := $(shell git rev-parse --short HEAD)
--------------------------------------------
The build is set to the commit id of your current HEAD

GIT_TAG := $(shell git describe)
--------------------------------
GIT_TAG set using the most recent tag if any.  You can add a tag with:
    git tag -a "mytag" -m "message about the tag."
Git push doesn't include tags.  To push tags to the origin use:
    git push origin --tag

## roadctl
The roadctl is modeled after kubectl. Use roadctl help for a list of top
level commands.  This may be different for a specific commnad or it's 
sub-commands.

### Top level help
roadctl help
roadctl allows you to work with the PavedRoad CNCF low-code environment and the associated CI/CD pipeline

  Usage: roadctl [command] [TYPE] [NAME] [flags]

  TYPE specifies a resource type
  NAME is the name of a resource
  flags specify options

    Usage:
      roadctl [command]
    
    Available Commands:
      apply       Apply configuration to named resource
      config      manage roadctl global configuration options
      create      create a new resource
      delete      delete a resource
      deploy      deploy a service
      describe    describe provides detailed information about a resource
      doc         Generate documentation for your service
      edit        edit the configuration for the specified resource
      events      View events
      explain     return documentation about a resource
      get         get an existing object
      help        Help about any command
      logs        return logs for a resource
      replace     Delete and recreate the named resource
      version     Print the current version
    
    Flags:
          --config string   Config file (default is $HOME/.roadctl.yaml)
          --debug string    Debug level: info(default)|warm|error|critical (default "d")
      --format string   Output format: text(default)|json|yaml (default "f")
  -h, --help            help for roadctl

Use "roadctl [command] --help" for more information about a command.

### Specific command
roadctl get --help
Return summary information about an existing resource

Usage:
  roadctl get [flags]

Flags:
  -h, --help          help for get
  -i, --init          Initialize template repository
  -r, --repo string   Change default repository for templates (default "https://github.pavedroad-io/templates")

Global Flags:
      --config string   Config file (default is $HOME/.roadctl.yaml)
      --debug string    Debug level: info(default)|warm|error|critical (default "d")
      --format string   Output format: text(default)|json|yaml (default "f")

## Generating your service
The roadctl CLI is used to create new services.  It has two 
fundamental concepts:
- templates: Contain logic need to generate a service
- definitions: Define your custom logic, integrations, 
               and oragnizational information

A sample definitions is available to help you get started

### Initialize template repository
    roadctl get templates --init

### List available templates
    roadctl get templates

### Create a copy of the sample definition
    roadctl describe templates datamgr > myservice.yaml

### Get defitions of attributes in your service.yaml
    roadctl explain templates datamgr > myservice.yaml

### Create your microservice
    roadctl create templates --template datamgr --definition myservice.yaml
### Build and test
Executing make will compilte and test your service.  Optionally, you
can do `make compile` followed by `make check`
    make

## Directories

| Name | Contents |
| --------- | -------- |
| artifacts | Outputs from static code analysis and tests |
| assets | Generate assets such as images |
| builds | Executables for supported platforms, Mac/Linux x86/amd64 |
| dev | Generated helper scripts and sample data |
| dev/db | Generated SQL statments |
| docs | Generated documentation |
| logs | Logs generated by the microservice |
| manifests | Docker and docker-composes manifest |
| mainfests/kubernetes | Kubernetes manifests for deploying this microservice |
| vendor | Vendor dependencies |

## SQL
To get an SQL prompt, use:
	bin/sql.sh

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
Use **make help** to get a list of options
make help

    Choose a command run in films:
    
    compile         Compile the binary.
    clean           Remove dep, vendor, binary(s), and executs go clean
    build           Build the binary for linux / mac x86 and amd
    deploy          Deploy image to repository and k8s cluster
    install         Install packages or main
    check           Start services and execute static code analysis and tests
    show-coverage   Show go code coverage in browser
    show-test       Show sonarcloud test report
    show-devkit     Show documenation for Devkit
    fmt             Run gofmt on all code
    simplify        Run gofmt with simplify option
    k8s-start       Start local microk8s server and update configurations
    k8s-stop        Stop local k8s cluster and delete skaffold deployments
    k8s-status      Print the status of the local cluster up or down
    help            Print possible commands

## skaffold CI/CD
Skaffold is integrated into your project.  You can use the following commands:

### development mode
Monitors source code and when it changes builds and pushs a new image

```bash
skaffold dev -f manifests/skaffold.yaml
```
### run
Builds and push image once when executed

```bash
skaffold run -f manifests/skaffold.yaml
```
  
### delete
Deletes all deployed resources

```bash
skaffold delete -f manifests/skaffold.yaml
```
  
## Linter(s)
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
- vunlnerabilites


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

Controlls the executing of an analysis run.  Documentation is
avaiable [here](https://docs.sonarqube.org/latest/analysis/analysis-parameters/).
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
export GITHUB_ACCESS_TOKEN=####################
```
{{end}}
