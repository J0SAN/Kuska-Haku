-- Migración 004 (Fase 5): Historial de recorridos GPS y restricción de un solo bus activo por conductor.
-- 1. Historial de posiciones GPS
-- La clave primaria compuesta (bus_id, ts) garantiza idempotencia ante reintentos HTTP
-- mediante INSERT ... ON CONFLICT DO NOTHING y cumple el requisito de TimescaleDB
-- para convertirse en hypertable con: SELECT create_hypertable('kuska.gps_positions', 'ts');
CREATE TABLE IF NOT EXISTS kuska.gps_positions (
    bus_id VARCHAR(100) NOT NULL,
    route_id VARCHAR(100) NOT NULL,
    ts TIMESTAMPTZ NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    speed_kmh REAL NOT NULL DEFAULT 0,
    heading REAL NOT NULL DEFAULT 0,
    PRIMARY KEY (bus_id, ts)
);

CREATE INDEX IF NOT EXISTS idx_gps_positions_ts ON kuska.gps_positions (ts DESC);

-- 2. Un conductor solo puede tener un bus activo asignado simultáneamente
CREATE UNIQUE INDEX IF NOT EXISTS uq_fleet_active_driver ON kuska.fleet (driver_id)
WHERE
    driver_id IS NOT NULL
    AND status = 'active';