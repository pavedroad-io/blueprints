{{define "organizationGrantAdmin.tpl"}}
GRANT ALL ON DATABASE {{.Info.Organization | ToCamel}} TO {{.Info.Organization | ToCamel}}Admin;
{{end}}
