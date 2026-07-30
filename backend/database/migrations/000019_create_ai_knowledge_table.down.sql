BEGIN;

DROP TRIGGER IF EXISTS update_ai_knowledge_modtime ON ai_knowledge;

DROP FUNCTION IF EXISTS update_ai_knowledge_modtime();

DROP INDEX IF EXISTS idx_ai_knowledge_ai_agent_id;
DROP INDEX IF EXISTS idx_ai_knowledge_tenant_id;
DROP INDEX IF EXISTS idx_ai_knowledge_created_at;
DROP INDEX IF EXISTS idx_ai_knowledge_updated_at;

DROP TABLE IF EXISTS ai_knowledge CASCADE;

COMMIT;