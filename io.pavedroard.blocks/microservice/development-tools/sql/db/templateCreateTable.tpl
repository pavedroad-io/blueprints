{{define "templateCreateTable.tpl"}}
CREATE TABLE IF NOT EXISTS {{.Info.Organization | ToCamel}}.{{.Info.Name | ToCamel}} (
    {{.Info.Name | ToCamel}}UUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    {{.Info.Name | ToCamel}} JSONB
);

CREATE INDEX IF NOT EXISTS {{.Info.Name | ToCamel}}Idx ON {{.Info.Organization | ToCamel}}.{{.Info.Name | ToCamel}} USING GIN ({{.Info.Name | ToCamel}});
{{end}}
