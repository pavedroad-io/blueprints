{{define "testGetReadiness.sh"}}#!/bin/bash
# curl -v http://127.0.0.1:8081/api/v1/namespace/mirantis/eventCollector/ready

host=127.0.0.1
port=8081
service="{{.Name}}"
namespace="{{.Namespace}}"

get()
{
curl -H "Content-Type: application/json" \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/ready | jq '.'
}

usage()
{
  echo "usage: testGetReadiness -k |--k8s"
  echo "    -k locates ands posts to local k8s cluster"
  echo "    Otherwise, it will post to $host on port $port"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -k | --k8s ) shift
      host="$(./getk8sip.sh)"
      port="$(./getNodePort.sh $service)"
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

# Call get
get{{end}}
