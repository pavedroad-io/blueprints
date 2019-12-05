{{define "dev/db/organizationAdmin.sql"}}
CREATE USER IF NOT EXISTS {{.Organization}}Admin;
{{end}}

