CREATE TABLE IF NOT EXISTS weight_history (
    id          SERIAL PRIMARY KEY,
    spool_id    INTEGER NOT NULL REFERENCES spools(id) ON DELETE CASCADE,
    weight      NUMERIC(8,2) NOT NULL,
    measured_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_weight_history_spool ON weight_history(spool_id);
CREATE INDEX idx_weight_history_measured ON weight_history(measured_at);
