BEGIN;

DROP TRIGGER IF EXISTS update_invoices_modtime ON invoices;

DROP FUNCTION IF EXISTS update_invoices_modtime();

DROP INDEX IF EXISTS idx_invoices_tenant_id_created_at;
DROP INDEX IF EXISTS idx_invoices_provider_invoice_id;
DROP INDEX IF EXISTS idx_invoices_paid_at;
DROP INDEX IF EXISTS idx_invoices_expires_at;

DROP TABLE IF EXISTS invoices CASCADE;

COMMIT;