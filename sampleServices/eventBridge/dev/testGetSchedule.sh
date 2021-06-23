#!/bin/bash
# /api/v1/namespace/mirantis/eventCollector/scheduler/{key}

host=127.0.0.1
port=8081
service="Eventbridge"
namespace="pavedroad"
uuid=""

get()
{
  echo $uuid
curl -H "Content-Type: application/json" \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/scheduler | jq '.'
}

usage()
{
  echo "usage: testGetSchedule -k |--k8s"
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
get