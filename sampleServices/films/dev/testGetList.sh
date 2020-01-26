
#!/bin/bash
host=127.0.0.1
port=8081
service="films"

get()
{
  curl -v -H "Content-Type: application/json" -s http://$host:$port/api/v1/namespace/pavedroad.io/$service"LIST" | jq '.'
}

usage()
{
  echo "usage: testGetList.sh -k |--k8s"
  echo "    -k locates ands posts to local k8s cluster"
  echo "    Otherwise, it will post to $host on port $port"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -k | --k8s ) shift
      host="$(./getk8sip.sh)"
      port="$(./getNodePort.sh $service)"
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
