-- 1. Tabla de Rutas (Utilizando PostGIS para el trazado espacial)
CREATE TABLE IF NOT EXISTS routes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    -- geometry(LineString, 4326) es el estándar de PostGIS para guardar una secuencia de coordenadas GPS
    path geometry(LineString, 4326) NOT NULL, 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tabla de Flota (Buses asignados a una ruta)
CREATE TABLE IF NOT EXISTS fleet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bus_id VARCHAR(100) UNIQUE NOT NULL, -- Este debe coincidir con el ID que emite el simulador MQTT
    route_id UUID NOT NULL REFERENCES routes(id) ON DELETE CASCADE,
    status VARCHAR(50) DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Índices espaciales para optimizar consultas de mapas en el futuro
CREATE INDEX IF NOT EXISTS routes_path_idx ON routes USING GIST (path);
CREATE INDEX IF NOT EXISTS fleet_bus_id_idx ON fleet (bus_id);