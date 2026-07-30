BEGIN;

DROP TRIGGER IF EXISTS update_audit_logs_modtime ON audit_logs;

DROP FUNCTION IF EXISTS update_audit_logs_modtime();

DROP INDEX IF EXISTS idx_audit_logs_tenant_id_created_at;
DROP INDEX IF EXISTS idx_audit_logs_actor_id;
DROP INDEX IF EXISTS idx_audit_logs_action;
DROP INDEX IF EXISTS idx_audit_logs_entity_type;
DROP INDEX IF EXISTS idx_audit_logs_entity_id;

DROP TABLE IF EXISTS audit_logs CASCADE;

COMMIT;