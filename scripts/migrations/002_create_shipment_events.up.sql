CREATE TABLE IF NOT EXISTS shipment_events (
    id             TEXT        NOT NULL PRIMARY KEY,
    shipment_id    TEXT        NOT NULL REFERENCES shipments(id),
    status         TEXT        NOT NULL,
    note           TEXT        NOT NULL DEFAULT '',
    occurred_at    TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_shipment_events_shipment_id ON shipment_events (shipment_id);

