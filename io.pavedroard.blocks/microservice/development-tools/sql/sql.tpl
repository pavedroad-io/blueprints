{{define "sql.tpl"}}
#!/bin/bash

#. ./dbcmd.sh

CMD=`which cockroach`" sql"
PORT="26257"
IP="127.0.0.1"
USER="root"

CMD=`echo $CMD "--insecure" --host=$IP:$PORT`

echo "$CMD"
$CMD
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
