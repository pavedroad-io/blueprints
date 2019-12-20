{{define "dev/db/organizationAdmin.sql"}}
CREATE USER IF NOT EXISTS {{.OrgSQLSafe}}Admin;
{{end}}

