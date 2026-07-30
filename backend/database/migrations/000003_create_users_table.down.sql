BEGIN;

DROP TRIGGER IF EXISTS update_users_modtime ON users;

DROP FUNCTION IF EXISTS update_users_modtime();

DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_status;
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_role_id;
DROP INDEX IF EXISTS idx_users_created_at;
DROP INDEX IF EXISTS idx_users_updated_at;
DROP INDEX IF EXISTS idx_users_deleted_at;

-- ALTER TABLE users DROP IF EXISTS CONSTRAINT users_name_length_check ON users;
-- ALTER TABLE users DROP IF EXISTS CONSTRAINT users_username_length_check ON users;
-- ALTER TABLE users DROP IF EXISTS CONSTRAINT users_email_format_check ON users;

DROP TABLE IF EXISTS users CASCADE;

COMMIT;