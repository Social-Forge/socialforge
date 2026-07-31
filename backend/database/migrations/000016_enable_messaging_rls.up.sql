BEGIN;

ALTER TABLE IF EXISTS channels ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS channels FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS conversations ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS conversations FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS messages ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS messages FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS message_outboxes ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS message_outboxes FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS conversation_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS conversation_events FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON channels USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON conversations USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON messages USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON message_outboxes USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON conversation_events USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;