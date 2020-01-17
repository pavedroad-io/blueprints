{{define "dev/getk8sip.sh"}}
#!/bin/bash

microk8s.config | grep server | sed -r "s/.*\/\/(.*):.*$/\1/"
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
