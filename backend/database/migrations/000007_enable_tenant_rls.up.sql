BEGIN;

ALTER TABLE IF EXISTS divisions ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS divisions FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON divisions USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;