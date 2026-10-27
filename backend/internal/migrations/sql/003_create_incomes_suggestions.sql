CREATE SCHEMA IF NOT EXISTS kuska;

-- 1. Registros de ingreso manual del conductor (con snapshot de bus_code)
CREATE TABLE IF NOT EXISTS kuska.income_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    fleet_id UUID REFERENCES kuska.fleet (id) ON DELETE SET NULL,
    bus_code VARCHAR(100) NOT NULL,
    driver_id UUID REFERENCES kuska.users (id) ON DELETE SET NULL,
    record_date DATE NOT NULL DEFAULT CURRENT_DATE,
    passengers INTEGER NOT NULL CHECK (passengers >= 0),
    amount NUMERIC(10, 2) NOT NULL CHECK (amount >= 0),
    notes TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_income_records_date ON kuska.income_records (record_date);

CREATE INDEX IF NOT EXISTS idx_income_records_driver_date ON kuska.income_records (driver_id, record_date);

-- 2. Buzón de sugerencias de usuarios
CREATE TABLE IF NOT EXISTS kuska.suggestions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID REFERENCES kuska.users (id) ON DELETE SET NULL,
    message TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'read')),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_suggestions_status_created ON kuska.suggestions (status, created_at DESC);