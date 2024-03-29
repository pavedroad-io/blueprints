{{define "dev/testPut.sh"}}
#!/bin/bash
host=127.0.0.1
port={{.HTTPPort}}
service="{{.Name}}"
namespace="{{.Namespace}}"
uuid=""

getUUID()
{
  uuid=`curl -H "Content-Type: application/json" -s http://$host:$port/api/v1/namespace/$namespace/$service"LIST" | jq -r '.[0].uuid'`
  echo "UUID for user test is :  $uuid"
}

put()
{
curl -H "Content-Type: application/json" \
     -X PUT \
     -d @$servicePutData.json \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/$uuid
}

usage()
{
  echo "usage: testPut -k |--k8s"
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

# Get UUID and call put
getUUID
put
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
