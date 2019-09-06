{{define "templateExecuteAll.sh"}}
#!/bin/bash

CMD=`which cockroach`" sql"
PORT="26257"
IP="127.0.0.1"
USER="root"

CMD=`echo $CMD "--insecure" --host=$IP:$PORT`

echo "$CMD"

#1 Create kevlarAdmin if it doesn not already exists
$CMD < {{.Organization}}Admin.sql

#2 Create kevlarWeb db
$CMD < {{.Organization}}.sql

#3 Create kevlarAdmin all on kevlarWeb db
$CMD < {{.Organization}}GrantAdmin.sql

#4 Create prTokenTable 
$CMD < {{.Name}}CreateTable.sql
{{end}}
