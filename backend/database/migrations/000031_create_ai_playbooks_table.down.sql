BEGIN;

DROP TRIGGER IF EXISTS update_ai_playbooks_modtime ON ai_playbooks;

DROP FUNCTION IF EXISTS update_ai_playbooks_modtime();

DROP INDEX IF EXISTS idx_ai_playbooks_tenant_id;
DROP INDEX IF EXISTS idx_ai_playbooks_ai_agent_id;
DROP INDEX IF EXISTS idx_ai_playbooks_created_at;
DROP INDEX IF EXISTS idx_ai_playbooks_updated_at;

DROP TABLE IF EXISTS ai_playbooks CASCADE;

COMMIT;