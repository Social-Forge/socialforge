BEGIN;

DROP TRIGGER IF EXISTS update_conversation_events_modtime ON conversation_events;

DROP FUNCTION IF EXISTS update_conversation_events_modtime();

DROP INDEX IF EXISTS idx_conversation_events_conversation_id_created_at;
DROP INDEX IF EXISTS idx_conversation_events_updated_at;

DROP TABLE IF EXISTS conversation_events CASCADE;

COMMIT;