BEGIN;

ALTER TABLE IF EXISTS labels ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS labels FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS quick_replies ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS quick_replies FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON labels USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON quick_replies USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;