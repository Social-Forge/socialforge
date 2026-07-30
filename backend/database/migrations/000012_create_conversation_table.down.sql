BEGIN;

DROP TRIGGER IF EXISTS update_conversations_modtime ON conversations;

DROP FUNCTION IF EXISTS update_conversations_modtime();

DROP INDEX IF EXISTS idx_conversations_channel_id_contact_id;
DROP INDEX IF EXISTS idx_conversations_tenant_id_status;
DROP INDEX IF EXISTS idx_conversations_tenant_id;
DROP INDEX IF EXISTS idx_conversations_channel_id;
DROP INDEX IF EXISTS idx_conversations_contact_id;
DROP INDEX IF EXISTS idx_conversations_assigned_agent_id;
DROP INDEX IF EXISTS idx_conversations_status;
DROP INDEX IF EXISTS idx_conversations_last_message_at;
DROP INDEX IF EXISTS idx_conversations_created_at;
DROP INDEX IF EXISTS idx_conversations_updated_at;

DROP TABLE IF EXISTS conversations CASCADE;

COMMIT;