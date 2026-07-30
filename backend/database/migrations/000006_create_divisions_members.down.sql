BEGIN;

DROP TRIGGER IF EXISTS update_division_members_modtime ON division_members;

DROP FUNCTION IF EXISTS update_division_members_modtime();

DROP INDEX IF EXISTS idx_division_members_user_tenant;
DROP INDEX IF EXISTS idx_division_members_division;

DROP TABLE IF EXISTS division_members CASCADE;

COMMIT;