{{define "dev/getdockerip.sh"}}#!/bin/bash
ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+'
{{end}}
