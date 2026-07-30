BEGIN;

DROP TRIGGER IF EXISTS update_oauth_providers_modtime ON oauth_providers;

DROP FUNCTION IF EXISTS update_oauth_providers_modtime();

DROP INDEX IF EXISTS idx_oauth_providers_provider_id;
DROP INDEX IF EXISTS idx_oauth_providers_provider_name;
DROP INDEX IF EXISTS idx_oauth_providers_user_id;
DROP INDEX IF EXISTS idx_oauth_providers_created_at;
DROP INDEX IF EXISTS idx_oauth_providers_updated_at;

DROP TABLE IF EXISTS oauth_providers CASCADE;

COMMIT;