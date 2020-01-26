
#!/bin/bash

SQL="cockroach sql"
PORT="26257"
HOST="127.0.0.1"
USER="root"
CMD=`echo $SQL "--insecure" --host=$HOST:$PORT`

buildCommand()
{
  CMD=`echo $SQL "--insecure" --host=$HOST:$PORT`
}

all()
{

  echo "========================================"
  echo " Initializing tables"
  start=`date`
  echo " Starting at : $start"
  echo " Using: $CMD"
  echo ""

# 1 Create acme-demoAdmin if it doesn not already exists
  $CMD < acme-demoAdmin.sql

# 2 Create acme-demoWeb db
  $CMD < acme-demo.sql

# 3 Create Aacme-demodmin all on kevlarWeb db
  $CMD < acme-demoGrantAdmin.sql

# 4 Create acme-demoTable 
  $CMD < filmsCreateTable.sql
}

usage()
{
  echo "usage: acme-demoExecuteAll.sh -k |--k8s"
  echo "    Created database and users as needed"
  echo "    -k locates and posts to local k8s cluster"
  echo "    it will default to $host on port $port"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -k | --k8s ) shift

      if [ -z "$COCKROACH_HOST" ]
      then
        echo "COCKROACH_HOST must be set"
        exit
      else
        HOST=$COCKROACH_HOST
      fi

      if [ -v "$COCKROACH_PORT" ]
      then
        PORT=$COCKROACH_PORT
      fi
      buildCommand
      ;;
  -h | --help ) usage
    exit
    ;;
  * ) shift
    ;;
  esac
done

# call all
all
