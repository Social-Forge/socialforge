BEGIN;

DROP TRIGGER IF EXISTS update_tenants_modtime ON tenants;

DROP FUNCTION IF EXISTS update_tenants_modtime();

DROP INDEX IF EXISTS idx_tenants_slug CASCADE;
DROP INDEX IF EXISTS idx_tenants_owner_id CASCADE;
DROP INDEX IF EXISTS idx_tenants_is_active CASCADE;
DROP INDEX IF EXISTS idx_tenants_created_at CASCADE;
DROP INDEX IF EXISTS idx_tenants_updated_at CASCADE;
DROP INDEX IF EXISTS idx_tenants_deleted_at CASCADE;
DROP INDEX IF EXISTS idx_tenants_subscription_status CASCADE;

-- ALTER TABLE tenants DROP CONSTRAINT IF EXISTS chk_subscription_plan;
-- ALTER TABLE tenants DROP CONSTRAINT IF EXISTS chk_subscription_status;

DROP TABLE IF EXISTS tenants CASCADE;

COMMIT;