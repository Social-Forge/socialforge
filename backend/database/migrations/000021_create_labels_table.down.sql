BEGIN;

DROP TRIGGER IF EXISTS update_labels_modtime ON labels;

DROP FUNCTION IF EXISTS update_labels_modtime();

DROP INDEX IF EXISTS idx_labels_tenant_id_name;
DROP INDEX IF EXISTS idx_labels_created_at;
DROP INDEX IF EXISTS idx_labels_updated_at;

DROP TABLE IF EXISTS labels CASCADE;

COMMIT;