{{define "templateExecuteAll.tpl"}}
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

# 1 Create {{.Info.Organization}}Admin if it doesn not already exists
  $CMD < {{.Info.Organization}}Admin.sql

# 2 Create {{.Info.Organization}}Web db
  $CMD < {{.Info.Organization}}.sql

# 3 Create A{{.Info.Organization}}dmin all on kevlarWeb db
  $CMD < {{.Info.Organization}}GrantAdmin.sql

# 4 Create {{.Info.Organization}}Table 
  $CMD < {{.Info.Name}}CreateTable.sql
}

usage()
{
  echo "usage: {{.Info.Organization}}ExecuteAll.sh -k |--k8s"
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
{{end}}
