BEGIN;

DROP TRIGGER IF EXISTS update_ai_assets_modtime ON ai_assets;

DROP FUNCTION IF EXISTS update_ai_assets_modtime();

DROP INDEX IF EXISTS idx_ai_assets_tenant_id;
DROP INDEX IF EXISTS idx_ai_assets_ai_agent_id;
DROP INDEX IF EXISTS idx_ai_assets_created_at;
DROP INDEX IF EXISTS idx_ai_assets_updated_at;

DROP TABLE IF EXISTS ai_assets CASCADE;

COMMIT;