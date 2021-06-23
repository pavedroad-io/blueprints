#!/bin/bash

host=127.0.0.1
port=8081
service="Eventbridge"
namespace="pavedroad"

postdata1="{\"command\": \"stop_scheduler\", \"field\": \"\", \"field_value\": 0}"
postdata2="{\"command\": \"start_scheduler\", \"field\": \"\", \"field_value\": 0}"
postdata3="{\"command\": \"stop_workers\", \"field\": \"\", \"field_value\": 0}"
postdata4="{\"command\": \"start_workers\", \"field\": \"\", \"field_value\": 0}"
postdata5="{\"command\": \"shutdown\", \"field\": \"\", \"field_value\": 0}"
postdata6="{\"command\": \"shutdown_now\", \"field\": \"\", \"field_value\": 0}"
postdata7="{\"command\": \"set\", \"field\": \"graceful_shutdown_seconds\", \"field_value\": 5}"
postdata8="{\"command\": \"set\", \"field\": \"hard_shutdown_seconds\", \"field_value\": 5}"
postdata9="{\"command\": \"set\", \"field\": \"number_of_workers\", \"field_value\": 10}"
postdata10="{\"command\": \"set\", \"field\": \"scheduler_channel_size\", \"field_value\": 10}"
postdata11="{\"command\": \"set\", \"field\": \"result_channel_size\", \"field_value\": 10}"

put()
{
curl -H "Content-Type: application/json" \
     -X PUT \
     -d "$postdata" \
     http://$host:$port/api/v1/namespace/$namespace/$service/management | jq '.'
}

usage()
{
  echo "usage: testPutManagement -k |--k8s"
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

# stop_scheduler
#postdata=$postdata1
#put

# start_scheduler
#postdata=$postdata2
#put

# stop_workers
#postdata=$postdata3
#put

# start_workers
#postdata=$postdata4
#put

# shutdown
#postdata=$postdata5
#put

# shutdown_now
#postdata=$postdata6
#put

# graceful_shutdown_second
#postdata=$postdata7
#put

# hard_shutdown_seconds
#postdata=$postdata8
#put

# number_of_workers
#postdata=$postdata9
#put

# scheduler_channel_size
#postdata=$postdata10
#put

# result_channel_size
#postdata=$postdata11
#put