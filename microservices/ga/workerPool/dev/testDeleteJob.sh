{{define "dev/testDeleteJob.sh"}}#!/bin/bash

host=127.0.0.1
port={{.HTTPPort}}
service="{{.NameExported}}"
namespace="{{.Namespace}}"
uuid=""

getUUID()
{
  uuid=`curl -H "Content-Type: application/json" -s http://$host:$port/api/v1/namespace/$namespace/$service/jobs"LIST" | jq -r '.[0].id'`
  if [ "$uuid" == "" ]
  then
    echo "UUID lookup failed"
    exit
  fi
}

delete()
{
curl -H "Content-Type: application/json" \
     -X DELETE \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/jobs/$uuid | jq '.'
}

usage()
{
  echo "usage: testDeleteJob -k |--k8s"
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
getUUID
delete{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
