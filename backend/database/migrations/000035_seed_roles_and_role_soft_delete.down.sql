BEGIN;

DELETE FROM roles WHERE name IN ('superadmin', 'tenant_owner', 'supervisor', 'agent');

DROP INDEX IF EXISTS idx_roles_deleted_at;
ALTER TABLE roles DROP COLUMN IF EXISTS deleted_at;

COMMIT;
