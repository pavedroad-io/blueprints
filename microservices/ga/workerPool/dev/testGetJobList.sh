{{define "dev/testGetJobList.sh"}}#!/bin/bash
# /api/v1/namespace/mirantis/eventCollector/jobsLIST

host=127.0.0.1
port={{.HTTPPort}}
service="{{.NameExported}}"
namespace="{{.Namespace}}"

get()
{
curl -H "Content-Type: application/json" \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/jobsLIST | jq '.'
}

usage()
{
  echo "usage: testGet -k |--k8s"
  echo "    -k locates ands posts to local k8s cluster"
  echo "    Otherwise, it will post to $host on port $port"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -k | --k8s ) shift
      host="$(./getk8sip.sh)"
      port="$(./getNodePort.sh $service $namespace)"
      echo $host
      echo $port
      ;;
  -h | --help ) usage
    exit
    ;;
  * ) shift
    ;;
  esac
done

# Get UUID and call get
get{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
