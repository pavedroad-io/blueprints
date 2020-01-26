
#!/bin/bash

getNodePort()
{
  kubectl get svc $name -o json | jq '.spec.ports[0].nodePort'
}

usage()
{
  echo "getNodePort service-name"
}

#### Main

if [ "$1" == "" ]
then
  usage
  exit
fi

name="$1"

getNodePort
