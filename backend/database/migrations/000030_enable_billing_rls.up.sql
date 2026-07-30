BEGIN;

ALTER TABLE IF EXISTS subscriptions ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS subscriptions FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS subscription_addons ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS subscription_addons FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS invoices ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS invoices FORCE ROW LEVEL SECURITY;

ALTER TABLE IF EXISTS payment_events ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS payment_events FORCE ROW LEVEL SECURITY;

CREATE POLICY tenant_isolation ON subscriptions USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON subssubscription_addonscriptions USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON invoices USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE POLICY tenant_isolation ON payment_events USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;