BEGIN;

DROP TRIGGER IF EXISTS update_subscriptions_modtime ON subscriptions;

DROP FUNCTION IF EXISTS update_subscriptions_modtime();

DROP INDEX IF EXISTS idx_subscriptions_tenant_id_status;
DROP INDEX IF EXISTS idx_subscriptions_plan_id;
DROP INDEX IF EXISTS idx_subscriptions_status;
DROP INDEX IF EXISTS idx_subscriptions_created_at;
DROP INDEX IF EXISTS idx_subscriptions_updated_at;

DROP TABLE IF EXISTS subscriptions CASCADE;

COMMIT;