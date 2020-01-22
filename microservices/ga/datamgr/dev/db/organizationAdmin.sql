{{define "dev/db/organizationAdmin.sql"}}
CREATE USER IF NOT EXISTS {{.OrgSQLSafe}}Admin;
{{/* vim: set filetype=gotexttmpl: */ -}}{{end}}

