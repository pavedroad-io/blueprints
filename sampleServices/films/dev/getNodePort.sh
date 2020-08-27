
#!/bin/bash

getNodePort()
{
  kubectl get svc $name -n $namespace -o json | jq '.spec.ports[0].nodePort'
}

usage()
{
  echo "getNodePort service-name namespace"
}

#### Main

if [ "$1" == "" ]
then
  usage
  exit
fi

name="$1"

if [ "$2" == "" ]
then
  namespace="default"
else
	namespace="$2"
fi


getNodePort
