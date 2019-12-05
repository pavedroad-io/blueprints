{{define "dev/testPost.sh"}}
#!/bin/bash

## set default values
host=127.0.0.1
port=8081
service="{{.Name}}"

post()
{
  curl -H "Content-Type: application/json" \
      -X POST \
      -d @{{.Name}}.json \
      -v http://$host:$port/api/v1/namespace/pavedroad.io/{{.Name}}
}

usage()
{
  echo "usage: testPost -k |--k8s"
  echo "    -k locates and posts to local k8s cluster"
  echo "    it will default to $host on port $port"
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

# call post
post

{{end}}
