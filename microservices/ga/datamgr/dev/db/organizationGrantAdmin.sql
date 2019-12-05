{{define "dev/db/organizationGrantAdmin.sql"}}
GRANT ALL ON DATABASE {{.Organization}} TO {{.Organization}}Admin;
{{end}}
