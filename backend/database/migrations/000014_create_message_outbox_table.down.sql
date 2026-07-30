BEGIN;

DROP TRIGGER IF EXISTS update_message_outboxes_modtime ON message_outboxes;

DROP FUNCTION IF EXISTS update_message_outboxes_modtime();

DROP INDEX IF EXISTS idx_message_outboxes_status_next_retry_at;
DROP INDEX IF EXISTS idx_message_outboxes_created_at;
DROP INDEX IF EXISTS idx_message_outboxes_updated_at;

DROP TABLE IF EXISTS message_outboxes CASCADE;

COMMIT;