{{define "dev/testPutJob.sh"}}#!/bin/bash

host=127.0.0.1
port={{.HTTPPort}}
service="{{.NameExported}}"
namespace="{{.Namespace}}"
uuid=""
newurl="https://geek-jokes.sameerkumar.website/api"

getUUID()
{
  uuid=`curl -H "Content-Type: application/json" -s http://$host:$port/api/v1/namespace/$namespace/$service/jobs"LIST" | jq -r '.[0].id'`
  if [ "$uuid" == "" ]
  then
    echo "UUID lookup failed"
    exit
  fi
}

put()
{
curl -H "Content-Type: application/json" \
     -X PUT \
     -d "$postdata" \
     http://$host:$port/api/v1/namespace/$namespace/$service/jobs/$uuid | jq '.'
}

usage()
{
  echo "usage: testPutJob -k |--k8s"
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
postdata="{\"id\": \"$uuid\", \"url\": \"$newurl\", \"type\": \"httpJob\"}"
put{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
