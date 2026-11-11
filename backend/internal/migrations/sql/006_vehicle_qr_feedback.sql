-- Migración 006 (Fase 7): QR por vehículo y calificaciones de pasajeros.
-- Tablas separadas de kuska.suggestions, que queda para opiniones generales del servicio.
CREATE TABLE IF NOT EXISTS kuska.vehicle_qr_codes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID NOT NULL REFERENCES kuska.fleet (id) ON DELETE CASCADE,
    token VARCHAR(22) NOT NULL,
    created_by UUID REFERENCES kuska.users (id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMPTZ,
    CONSTRAINT uq_vehicle_qr_token UNIQUE (token),
    CONSTRAINT chk_vehicle_qr_token_format CHECK (token ~ '^[A-Za-z0-9_-]{22}$')
);

-- Un solo QR activo por vehículo
CREATE UNIQUE INDEX IF NOT EXISTS uq_vehicle_qr_active
    ON kuska.vehicle_qr_codes (fleet_id) WHERE revoked_at IS NULL;

CREATE TABLE IF NOT EXISTS kuska.vehicle_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fleet_id UUID REFERENCES kuska.fleet (id) ON DELETE SET NULL,
    qr_code_id UUID REFERENCES kuska.vehicle_qr_codes (id) ON DELETE SET NULL,
    bus_code VARCHAR(100) NOT NULL,      -- snapshot, igual que income_records
    route_name VARCHAR(255) NOT NULL,    -- snapshot
    driver_id UUID REFERENCES kuska.users (id) ON DELETE SET NULL, -- asignado al calificar
    rating SMALLINT NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT NOT NULL DEFAULT '' CHECK (char_length(comment) <= 500),
    device_hash CHAR(64) NOT NULL,
    service_date DATE NOT NULL,          -- fecha de Lima
    bus_live BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(20) NOT NULL DEFAULT 'new' CHECK (status IN ('new', 'read', 'hidden')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_vehicle_feedback_device_day UNIQUE (fleet_id, device_hash, service_date)
);

CREATE INDEX IF NOT EXISTS idx_vehicle_feedback_fleet_created ON kuska.vehicle_feedback (fleet_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_vehicle_feedback_status_created ON kuska.vehicle_feedback (status, created_at DESC);
