BEGIN;

-- role_repository treats roles as soft-deletable (WHERE deleted_at IS NULL),
-- but the roles table had no such column. Add it for consistency.
ALTER TABLE roles ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;
CREATE INDEX IF NOT EXISTS idx_roles_deleted_at ON roles(deleted_at);

-- Seed the four static system roles used by the level/name-based RBAC.
-- level: 0=superadmin, 1=tenant_owner, 2=supervisor, 3=agent.
INSERT INTO roles (name, slug, description, level) VALUES
  ('superadmin',   'superadmin',   'Platform super administrator', 0),
  ('tenant_owner', 'tenant-owner', 'Tenant owner (full control)',  1),
  ('supervisor',   'supervisor',   'Division supervisor',          2),
  ('agent',        'agent',        'Customer service agent',       3)
ON CONFLICT (name) DO NOTHING;

COMMIT;
