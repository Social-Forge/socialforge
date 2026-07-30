BEGIN;

DROP TRIGGER IF EXISTS update_divisions_members_modtime ON divisions_members;

DROP FUNCTION IF EXISTS update_divisions_members_modtime();

DROP INDEX IF EXISTS idx_divisions_members_tenant_id_user_id;
DROP INDEX IF EXISTS idx_divisions_members_created_at;
DROP INDEX IF EXISTS idx_divisions_members_updated_at;

DROP TABLE IF EXISTS divisions_members CASCADE;

COMMIT;