BEGIN;

-- shortcut must be unique PER TENANT, not globally (multi-tenant fix).
ALTER TABLE quick_replies DROP CONSTRAINT IF EXISTS quick_replies_shortcut_key;
ALTER TABLE quick_replies
  ADD CONSTRAINT uq_quick_replies_tenant_shortcut UNIQUE (tenant_id, shortcut);

CREATE INDEX IF NOT EXISTS idx_quick_replies_tenant_id ON quick_replies(tenant_id);

COMMIT;
