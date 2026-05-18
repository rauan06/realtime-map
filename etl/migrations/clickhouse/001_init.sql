-- Single generic event sink table for the ETL.
--
-- payload is the source-specific JSON record (lat/lng/heading/severity/...);
-- the dashboard parses it lazily so we can add new sources without migrating.
CREATE TABLE IF NOT EXISTS etl_events (
    source       LowCardinality(String),
    key          String,
    payload      String,
    received_at  DateTime64(3, 'UTC')
)
ENGINE = MergeTree
PARTITION BY toDate(received_at)
ORDER BY (source, key, received_at)
TTL toDateTime(received_at) + INTERVAL 30 DAY;
