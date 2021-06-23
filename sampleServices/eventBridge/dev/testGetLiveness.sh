#!/bin/bash
# curl -v http://127.0.0.1:8081/api/v1/namespace/mirantis/eventCollector/liveness

host=127.0.0.1
port=8081
service="Eventbridge"
namespace="pavedroad"

get()
{
curl -H "Content-Type: application/json" \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/liveness | jq '.'
}

usage()
{
  echo "usage: testGetLiveness -k |--k8s"
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

# Call get
get