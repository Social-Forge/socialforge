BEGIN;

DROP TRIGGER IF EXISTS update_quick_replies_modtime ON quick_replies;

DROP FUNCTION IF EXISTS update_quick_replies_modtime();

DROP INDEX IF EXISTS idx_quick_replies_tenant_id_shortcut;
DROP INDEX IF EXISTS idx_quick_replies_tenant_id;
DROP INDEX IF EXISTS idx_quick_replies_created_at;
DROP INDEX IF EXISTS idx_quick_replies_updated_at;

DROP TABLE IF EXISTS quick_replies CASCADE;

COMMIT;