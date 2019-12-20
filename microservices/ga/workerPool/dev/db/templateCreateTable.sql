{{define "dev/db/templateCreateTable.sql"}}
CREATE TABLE IF NOT EXISTS {{.OrgSQLSafe}}.{{.Name}} (
    {{.NameExported}}UUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    {{.Name}} JSONB
);

CREATE INDEX IF NOT EXISTS {{.Name}}Idx ON {{.OrgSQLSafe}}.{{.Name}} USING GIN ({{.Name}});
{{end}}
