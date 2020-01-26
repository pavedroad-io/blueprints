
CREATE TABLE IF NOT EXISTS AcmeDemo.films (
    FilmsUUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    films JSONB
);

CREATE INDEX IF NOT EXISTS filmsIdx ON AcmeDemo.films USING GIN (films);
