BEGIN;

DROP TRIGGER IF EXISTS update_conversation_labels_modtime ON conversation_labels;

DROP FUNCTION IF EXISTS update_conversation_labels_modtime();

DROP INDEX IF EXISTS idx_labels_conversation_id_label_id;
DROP INDEX IF EXISTS idx_labels_conversation_id;
DROP INDEX IF EXISTS idx_labels_label_id;
DROP INDEX IF EXISTS idx_labels_created_at;
DROP INDEX IF EXISTS idx_labels_updated_at;

DROP TABLE IF EXISTS conversation_labels CASCADE;

COMMIT;