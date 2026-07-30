BEGIN;

DROP TRIGGER IF EXISTS update_divisions_modtime ON divisions;

DROP FUNCTION IF EXISTS update_divisions_modtime();

DROP INDEX IF EXISTS idx_divisions_tenant_id_name;
DROP INDEX IF EXISTS idx_divisions_tenant_id;
DROP INDEX IF EXISTS idx_divisions_slug;    
DROP INDEX IF EXISTS idx_divisions_is_active;
DROP INDEX IF EXISTS idx_divisions_link_url;
DROP INDEX IF EXISTS idx_divisions_created_at;
DROP INDEX IF EXISTS idx_divisions_updated_at;
DROP INDEX IF EXISTS idx_divisions_deleted_at;

-- ALTER TABLE divisions DROP CONSTRAINT IF EXISTS chk_division_slug_tenant_id;
-- ALTER TABLE divisions DROP CONSTRAINT IF EXISTS chk_division_routing_type;

DROP TABLE IF EXISTS divisions CASCADE;

COMMIT;