{{define "organization.tpl"}}
CREATE DATABASE IF NOT EXISTS {{.Info.Organization | ToCamel}};
{{end}}
