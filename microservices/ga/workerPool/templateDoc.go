{{define "templateDoc.go"}}
// Package classification {{.Named}} API.
//
// Micro service for managing a pool of workers
//
// A scheduler go routine writes jobs to be performed to
// a dispatcher.  The dispatcher manages and forwards jobs
// to a number N number of workers using a buffered channel.
//
// Workers read the jobs, perform the tasks, and log the 
// results. The log code, logs to one or more configured
// destinations.  This can include local file system, stdout,
// or a Kafka topic.
//
// Jobs, Scheduler, are both defined as interfaces enabling
// them to be customized to specific tasks.
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
