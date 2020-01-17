
CREATE TABLE IF NOT EXISTS Mirantis.eventCollector (
    EventCollectorUUID UUID DEFAULT uuid_v4()::UUID PRIMARY KEY,
    eventCollector JSONB
);

CREATE INDEX IF NOT EXISTS eventCollectorIdx ON Mirantis.eventCollector USING GIN (eventCollector);
