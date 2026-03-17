CREATE TABLE IF NOT EXISTS shipments (
    id               TEXT        NOT NULL PRIMARY KEY,
    reference_number TEXT        NOT NULL,
    origin           TEXT        NOT NULL,
    destination      TEXT        NOT NULL,
    current_status   TEXT        NOT NULL,
    driver_name      TEXT        NOT NULL DEFAULT '',
    unit_number      TEXT        NOT NULL DEFAULT '',
    shipment_amount  NUMERIC     NOT NULL DEFAULT 0,
    driver_revenue   NUMERIC     NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ NOT NULL,
    updated_at       TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_shipments_reference_number ON shipments (reference_number);

