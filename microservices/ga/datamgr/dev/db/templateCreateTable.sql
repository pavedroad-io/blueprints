{{define "dev/db/templateCreateTable.sql"}}
CREATE TABLE IF NOT EXISTS {{.Organization}}.{{.Name}} (
    {{.NameExported}}UUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    {{.Name}} JSONB
);

CREATE INDEX IF NOT EXISTS {{.Name}}Idx ON {{.Organization}}.{{.Name}} USING GIN ({{.Name}});
{{end}}
