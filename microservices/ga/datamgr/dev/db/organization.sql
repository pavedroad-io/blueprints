{{define "dev/db/organization.sql"}}
CREATE DATABASE IF NOT EXISTS {{.OrgSQLSafe}};
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}
