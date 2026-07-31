BEGIN;

ALTER TABLE IF EXISTS ai_playbooks ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS ai_playbooks FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS ai_assets ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS ai_assets FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON ai_playbooks USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON ai_assets USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;