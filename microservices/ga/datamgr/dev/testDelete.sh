{{define "dev/testDelete.sh"}}
#!/bin/bash
host=127.0.0.1
port=8082
service="{{.Name}}"
namespace="{{.Namespace}}"
uuid=""

getUUID()
{
  uuid=`curl -H "Content-Type: application/json" -s http://$host:$port/api/v1/namespace/$namespace/$service"LIST" | jq -r '.[0].uuid'`

  if [ $uuid == "" ]
  then
    echo "UUID lookup failed"
    exit
  fi
}
delete()
{
curl -X DELETE \
     -H "Content-Type: application/json" \
     -v http://$host:$port/api/v1/namespace/$namespace/$service/$uuid
}

usage()
{
  echo "usage: testDelete -k |--k8s"
  echo "    -k locates and issues delete to local k8s cluster"
  echo "    Otherwise, it will use $host on port $port"
}

# Main
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
delete
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
