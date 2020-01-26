
#!/bin/bash
host=127.0.0.1
port=8081
service="films"
flag=""

post()
{
  if [ $flag == "k8s" ]
  then
    ./testPost.sh -k
  else
    ./testPost.sh
  fi
}

get()
{
  if [ $flag == "k8s" ]
  then
    ./testGet.sh -k
  else
    ./testGet.sh
  fi
}

getList()
{
  if [ $flag == "k8s" ]
  then
    ./testGetList.sh -k
  else
    ./testGetList.sh
  fi
}

put()
{
  if [ $flag == "k8s" ]
  then
    ./testPut.sh -k
  else
    ./testPut.sh
  fi
}

delete()
{
  if [ $flag == "k8s" ]
  then
    ./testDelete.sh -k
  else
    ./testDelete.sh
  fi
}

usage()
{
  echo "usage: testAll -k |--k8s"
  echo "    -k locates ands posts to local k8s cluster"
  echo "    Otherwise, it will post to $host on port $port"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -k | --k8s ) shift
      host="$(./getk8sip.sh)"
      port="$(./getNodePort.sh $service)"
      flag="k8s"
      echo $host
      echo $port
      echo $flag
      ;;
  -h | --help ) usage
    exit
    ;;
  * ) shift
    ;;
  esac
done

# Get UUID and call get
post
post
get
getList
put
delete
