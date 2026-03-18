
CREATE TABLE shipments (
    reference VARCHAR(100) PRIMARY KEY,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    current_status VARCHAR(50) NOT NULL,
    driver VARCHAR(255) NOT NULL,
    unit VARCHAR(100) NOT NULL,
    amount NUMERIC(12,2) NOT NULL,
    driver_revenue NUMERIC(12,2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE shipment_events (
    id BIGSERIAL PRIMARY KEY,
    shipment_reference VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_shipment_events_shipment
        FOREIGN KEY (shipment_reference)
        REFERENCES shipments(reference)
        ON DELETE CASCADE
);

CREATE INDEX idx_shipment_events_reference_created_at
ON shipment_events(shipment_reference, created_at);