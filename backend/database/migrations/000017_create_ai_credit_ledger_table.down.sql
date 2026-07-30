BEGIN;

DROP TRIGGER IF EXISTS update_ai_credit_ledgers_modtime ON ai_credit_ledgers;

DROP FUNCTION IF EXISTS update_ai_credit_ledgers_modtime();

DROP INDEX IF EXISTS idx_ai_credit_ledgers_tenant_id_created_at;
DROP INDEX IF EXISTS idx_ai_credit_ledgers_conversation_id;
DROP INDEX IF EXISTS idx_ai_credit_ledgers_message_id;
DROP INDEX IF EXISTS idx_ai_credit_ledgers_updated_at;

DROP TABLE IF EXISTS ai_credit_ledgers CASCADE;

COMMIT;