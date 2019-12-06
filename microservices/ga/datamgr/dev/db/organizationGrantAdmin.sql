{{define "dev/db/organizationGrantAdmin.sql"}}
GRANT ALL ON DATABASE {{.OrgSQLSafe}} TO {{.OrgSQLSafe}}Admin;
{{end}}
