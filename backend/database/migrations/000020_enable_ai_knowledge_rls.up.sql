BEGIN;

ALTER TABLE IF EXISTS ai_knowledge ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS ai_knowledge FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON ai_knowledge USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;