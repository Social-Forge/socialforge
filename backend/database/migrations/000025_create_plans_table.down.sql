BEGIN;

DROP TRIGGER IF EXISTS update_plans_modtime ON plans;

DROP FUNCTION IF EXISTS update_plans_modtime();

DROP INDEX IF EXISTS idx_plans_code;
DROP INDEX IF EXISTS idx_plans_name;
DROP INDEX IF EXISTS idx_plans_created_at;
DROP INDEX IF EXISTS idx_plans_updated_at;

DROP TABLE IF EXISTS plans CASCADE;

COMMIT;