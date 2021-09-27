{{define "doc.go"}}
// Package classification {{.Name}} API.
//
// Micro service for managing a generic service
// 
// Support capabilities include:
// - Kubernetes life cycle events
// - Docker image generation
// - docker-compose with dependent services
// - Metrics API
// - Management API
// - Swagger generated documentation
// - Dependency management with Go modules
// - SonarCloud for public projects
// - FOSSA for public projects
// - Go SCA
//
// Terms Of Service:
//
// there are no TOS at this moment, use at your own risk we take no responsibility
//
//     Schemes: http, https
//     Host: {{.TLD}}
//     BasePath: /api/v1/namespace/{{.Namespace}}/{{.Name}}
//     Version: 1.0.0
//     License: Apache 2
//     Contact: {{.MaintainerName}}<{{.MaintainerEmail}}> {{.MaintainerWeb}}
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
//{{.PavedroadInfo}}
//
// Licensed under the Apache License Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
