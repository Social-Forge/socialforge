BEGIN;

DROP TRIGGER IF EXISTS update_roles_modtime ON roles;

DROP FUNCTION IF EXISTS update_roles_modtime();

DROP INDEX IF EXISTS idx_roles_level;
DROP INDEX IF EXISTS idx_roles_created_at;
DROP INDEX IF EXISTS idx_roles_slug;
DROP INDEX IF EXISTS idx_roles_name;

-- ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_name_length_check ON roles;
-- ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_slug_length_check ON roles;
-- ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_level_check ON roles;
-- ALTER TABLE roles DROP CONSTRAINT IF EXISTS roles_name_check ON roles;

DROP TABLE IF EXISTS roles CASCADE;

COMMIT;