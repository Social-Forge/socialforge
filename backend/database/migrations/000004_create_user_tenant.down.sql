BEGIN;

DROP TRIGGER IF EXISTS update_user_tenants_modtime ON user_tenants;

DROP FUNCTION IF EXISTS update_user_tenants_modtime();

DROP INDEX IF EXISTS idx_user_tenants_user_id;
DROP INDEX IF EXISTS idx_user_tenants_tenant_id;
DROP INDEX IF EXISTS idx_user_tenants_role_id;

DROP TABLE IF EXISTS user_tenants CASCADE;

COMMIT;