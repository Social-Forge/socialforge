BEGIN;

DROP TRIGGER IF EXISTS update_channels_modtime ON channels;

DROP FUNCTION IF EXISTS update_channels_modtime();

DROP INDEX IF EXISTS idx_channels_tenant_id;
DROP INDEX IF EXISTS idx_channels_tenant_id_type;
DROP INDEX IF EXISTS idx_channels_ai_agent_id;
DROP INDEX IF EXISTS idx_channels_division_id;
DROP INDEX IF EXISTS idx_channels_type;
DROP INDEX IF EXISTS idx_channels_status;
DROP INDEX IF EXISTS idx_channels_created_at;
DROP INDEX IF EXISTS idx_channels_updated_at;

DROP TABLE IF EXISTS channels CASCADE;

COMMIT;