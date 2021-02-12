{{define "templateCreateTable.tpl"}}
CREATE TABLE IF NOT EXISTS {{.Info.Organization | ToCamel}}.{{.Info.Name}} (
    {{.Info.Name | ToCamel}}UUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    {{.Info.Name}} JSONB
);

CREATE INDEX IF NOT EXISTS {{.Info.Name}}Idx ON {{.Info.Organization | ToCamel}}.{{.Info.Name}} USING GIN ({{.Info.Name}});
{{end}}
