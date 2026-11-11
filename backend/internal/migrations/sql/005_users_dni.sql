-- Migración 005 (Fase 6): registro de pasajeros con DNI (spec 002).
-- dni es NULL para las cuentas creadas antes de esta fase (admin, conductores
-- y pasajeros de prueba).
ALTER TABLE kuska.users ADD COLUMN IF NOT EXISTS dni VARCHAR(8);

ALTER TABLE kuska.users ADD COLUMN IF NOT EXISTS terms_accepted_at TIMESTAMPTZ;

-- Una cuenta por DNI
CREATE UNIQUE INDEX IF NOT EXISTS uq_users_dni ON kuska.users (dni) WHERE dni IS NOT NULL;

-- Formato del DNI y roles válidos también en la base
ALTER TABLE kuska.users DROP CONSTRAINT IF EXISTS chk_users_dni_format;

ALTER TABLE kuska.users ADD CONSTRAINT chk_users_dni_format CHECK (dni IS NULL OR dni ~ '^[0-9]{8}$');

ALTER TABLE kuska.users DROP CONSTRAINT IF EXISTS chk_users_role;

ALTER TABLE kuska.users ADD CONSTRAINT chk_users_role CHECK (role IN ('passenger', 'driver', 'admin'));
