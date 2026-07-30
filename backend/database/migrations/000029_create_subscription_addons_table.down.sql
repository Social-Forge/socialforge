BEGIN;

DROP TRIGGER IF EXISTS update_subscription_addons_modtime ON subscription_addons;

DROP FUNCTION IF EXISTS update_subscription_addons_modtime();

DROP INDEX IF EXISTS idx_subscription_addons_tenant_id_type;
DROP INDEX IF EXISTS idx_subscription_addons_tenant_id;
DROP INDEX IF EXISTS idx_subscription_addons_created_at;
DROP INDEX IF EXISTS idx_subscription_addons_updated_at;

DROP TABLE IF EXISTS subscription_addons CASCADE;

COMMIT;