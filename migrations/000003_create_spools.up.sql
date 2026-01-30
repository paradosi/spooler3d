CREATE TABLE IF NOT EXISTS spools (
    id               SERIAL PRIMARY KEY,
    uid              UUID NOT NULL DEFAULT gen_random_uuid() UNIQUE,
    manufacturer_id  INTEGER REFERENCES manufacturers(id) ON DELETE SET NULL,
    filament_type_id INTEGER REFERENCES filament_types(id) ON DELETE SET NULL,
    color_name       VARCHAR(100),
    color_hex        VARCHAR(7),
    diameter         NUMERIC(4,2) NOT NULL DEFAULT 1.75,
    spool_weight     NUMERIC(8,2),
    net_weight       NUMERIC(8,2) DEFAULT 1000,
    current_weight   NUMERIC(8,2),
    location         VARCHAR(100),
    purchase_date    DATE,
    purchase_price   NUMERIC(8,2),
    notes            TEXT,
    td_code          VARCHAR(100),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_spools_uid ON spools(uid);
CREATE INDEX idx_spools_manufacturer ON spools(manufacturer_id);
CREATE INDEX idx_spools_filament_type ON spools(filament_type_id);
