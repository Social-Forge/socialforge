BEGIN;

DROP TRIGGER IF EXISTS update_ai_agents_modtime ON ai_agents;

DROP FUNCTION IF EXISTS update_ai_agents_modtime();

DROP INDEX IF EXISTS idx_ai_agents_tenant_id;
DROP INDEX IF EXISTS idx_ai_agents_provider;
DROP INDEX IF EXISTS idx_ai_agents_name;
DROP INDEX IF EXISTS idx_ai_agents_model;
DROP INDEX IF EXISTS idx_ai_agents_is_active;
DROP INDEX IF EXISTS idx_ai_agents_created_at;
DROP INDEX IF EXISTS idx_ai_agents_updated_at;

DROP TABLE IF EXISTS ai_agents CASCADE;

COMMIT;