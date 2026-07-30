BEGIN;

DROP TRIGGER IF EXISTS update_payment_events_modtime ON payment_events;

DROP FUNCTION IF EXISTS update_payment_events_modtime();

DROP INDEX IF EXISTS idx_payment_events_tenant_id_created_at;
DROP INDEX IF EXISTS idx_payment_events_invoice_id;
DROP INDEX IF EXISTS idx_payment_events_provider;
DROP INDEX IF EXISTS idx_payment_events_event_type;
DROP INDEX IF EXISTS idx_payment_events_external_id;

DROP TABLE IF EXISTS payment_events CASCADE;

COMMIT;