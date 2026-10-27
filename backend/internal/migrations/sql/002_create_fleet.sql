CREATE SCHEMA IF NOT EXISTS kuska;

-- 1. Tabla de Rutas (PostGIS LineString SRID 4326 - WGS84)
CREATE TABLE IF NOT EXISTS kuska.routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    name VARCHAR(255) NOT NULL,
    path geometry (LineString, 4326) NOT NULL,
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabla de Paraderos (Stops) ordenados por ruta para cálculo de ETA (Fase 4)
CREATE TABLE IF NOT EXISTS kuska.stops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    route_id UUID NOT NULL REFERENCES kuska.routes (id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    location geometry (Point, 4326) NOT NULL,
    stop_order INTEGER NOT NULL CHECK (stop_order >= 1),
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        UNIQUE (route_id, stop_order)
);

-- 3. Tabla de Flota (Buses vinculados a ruta y conductor)
CREATE TABLE IF NOT EXISTS kuska.fleet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    bus_id VARCHAR(100) UNIQUE NOT NULL, -- Coincide con el ID que emite el simulador/conductor MQTT
    route_id UUID NOT NULL REFERENCES kuska.routes (id) ON DELETE CASCADE,
    driver_id UUID REFERENCES kuska.users (id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE kuska.fleet
ADD COLUMN IF NOT EXISTS driver_id UUID REFERENCES kuska.users (id) ON DELETE SET NULL;

-- Índices espaciales GIST y de búsqueda relacional
CREATE INDEX IF NOT EXISTS idx_routes_path ON kuska.routes USING GIST (path);

CREATE INDEX IF NOT EXISTS idx_stops_location ON kuska.stops USING GIST (location);

CREATE INDEX IF NOT EXISTS idx_stops_route_id ON kuska.stops (route_id);

CREATE INDEX IF NOT EXISTS idx_fleet_bus_id ON kuska.fleet (bus_id);

CREATE INDEX IF NOT EXISTS idx_fleet_route_id ON kuska.fleet (route_id);

CREATE INDEX IF NOT EXISTS idx_fleet_driver_id ON kuska.fleet (driver_id);