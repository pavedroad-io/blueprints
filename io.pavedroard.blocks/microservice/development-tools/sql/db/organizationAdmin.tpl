{{define "organizationAdmin.tpl"}}
CREATE USER IF NOT EXISTS {{.Info.Organization | ToCamel}}Admin;
{{end}}

