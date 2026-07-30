BEGIN;

ALTER TABLE IF EXISTS ai_agents ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS ai_agents FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS ai_credit_ledgers ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS ai_credit_ledgers FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON ai_agents USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON ai_credit_ledgers USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;