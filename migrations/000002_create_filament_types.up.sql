CREATE TABLE IF NOT EXISTS filament_types (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(100) NOT NULL UNIQUE,
    print_temp_min INTEGER,
    print_temp_max INTEGER,
    bed_temp_min   INTEGER,
    bed_temp_max   INTEGER,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO filament_types (name, print_temp_min, print_temp_max, bed_temp_min, bed_temp_max) VALUES
    ('PLA',  190, 220, 50, 65),
    ('PETG', 220, 250, 70, 85),
    ('ABS',  230, 260, 95, 110),
    ('TPU',  210, 230, 40, 60),
    ('ASA',  240, 260, 95, 110)
ON CONFLICT (name) DO NOTHING;
