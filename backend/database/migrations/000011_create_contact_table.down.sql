BEGIN;

DROP TRIGGER IF EXISTS update_contacts_modtime ON contacts;

DROP FUNCTION IF EXISTS update_contacts_modtime();

DROP INDEX IF EXISTS idx_contacts_channel_id_external_id;
DROP INDEX IF EXISTS idx_contacts_tenant_id;
DROP INDEX IF EXISTS idx_contacts_channel_id;
DROP INDEX IF EXISTS idx_contacts_external_id;
DROP INDEX IF EXISTS idx_contacts_display_name;
DROP INDEX IF EXISTS idx_contacts_created_at;
DROP INDEX IF EXISTS idx_contacts_updated_at;

DROP TABLE IF EXISTS contacts CASCADE;

COMMIT;