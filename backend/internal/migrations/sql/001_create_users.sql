CREATE SCHEMA IF NOT EXISTS kuska;

-- Limpieza preventiva de tablas residuales en el esquema public
DROP TABLE IF EXISTS public.fleet CASCADE;

DROP TABLE IF EXISTS public.routes CASCADE;

DROP TABLE IF EXISTS public.users CASCADE;

CREATE TABLE IF NOT EXISTS kuska.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL DEFAULT 'passenger',
    created_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP
    WITH
        TIME ZONE DEFAULT CURRENT_TIMESTAMP
);