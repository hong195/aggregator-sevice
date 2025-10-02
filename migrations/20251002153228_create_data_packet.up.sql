CREATE TABLE data_packets
(
    id UUID PRIMARY KEY,
    ts TIMESTAMPTZ NOT NULL,
    max_value BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_data_packets_ts_id ON data_packets (ts ASC, id ASC);

CREATE INDEX IF NOT EXISTS idx_data_packets_max_value ON data_packets (max_value);