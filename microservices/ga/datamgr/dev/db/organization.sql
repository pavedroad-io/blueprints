{{define "organization.sql"}}
CREATE DATABASE IF NOT EXISTS {{.Organization}};
{{end}}
