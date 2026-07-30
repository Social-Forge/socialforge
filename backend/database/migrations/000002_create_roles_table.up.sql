CREATE TABLE IF NOT EXISTS roles (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) UNIQUE NOT NULL DEFAULT 'tenant_owner',
  slug VARCHAR(255) UNIQUE NOT NULL,
  description TEXT,
  level INTEGER NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT roles_name_length_check 
    CHECK (length(trim(name)) > 0),
  CONSTRAINT roles_slug_length_check 
    CHECK (length(trim(slug)) > 0),
  CONSTRAINT roles_level_check 
    CHECK (level IN (0, 1, 2, 3)),
  CONSTRAINT roles_name_check 
    CHECK (name IN ('superadmin', 'tenant_owner', 'supervisor', 'agent'))
);

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_roles_slug ON roles(slug);
CREATE INDEX IF NOT EXISTS idx_roles_created_at ON roles(created_at);
CREATE INDEX IF NOT EXISTS idx_roles_level ON roles(level);

-- Add comments
COMMENT ON TABLE roles IS 'User roles in the system';
COMMENT ON COLUMN roles.level IS '0=superadmin, 1=tenant_owner, 2=supervisor, 3=agent';

CREATE OR REPLACE FUNCTION update_roles_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_roles_modtime
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION update_roles_modtime();
