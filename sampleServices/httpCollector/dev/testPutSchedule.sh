#!/bin/bash
# /api/v1/namespace/mirantis/eventCollector/scheduler
#{
#  "schedule_type": "Constant interval scheduler",
#  "send_interval_seconds": 10
#}

host=127.0.0.1
port=8081
service="httpcollector"
namespace="mirantis"
putdata="{\"schedule_type\": \"Constant interval scheduler\", \"send_interval_seconds\": 30}"

put()
{
curl -H "Content-Type: application/json" \
     -X PUT \
     -d "$putdata" \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/scheduler | jq '.'
}

usage()
{
  echo "usage: testPutSchedule -k |--k8s"
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
put