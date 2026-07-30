BEGIN;

DROP TRIGGER IF EXISTS update_messages_modtime ON messages;

DROP FUNCTION IF EXISTS update_messages_modtime();

DROP INDEX IF EXISTS idx_messages_conversation_id_created_at;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP INDEX IF EXISTS idx_messages_status;
DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_updated_at;

DROP TABLE IF EXISTS messages CASCADE;

COMMIT;