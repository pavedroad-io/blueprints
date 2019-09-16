{{define "organizationAdmin.sql"}}
CREATE USER IF NOT EXISTS {{.Organization}}Admin;
{{end}}

