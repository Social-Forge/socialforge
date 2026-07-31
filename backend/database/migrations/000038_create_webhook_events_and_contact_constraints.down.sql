BEGIN;

DROP POLICY IF EXISTS tenant_isolation ON contacts;
ALTER TABLE IF EXISTS contacts DISABLE ROW LEVEL SECURITY;
ALTER TABLE contacts DROP CONSTRAINT IF EXISTS uq_contacts_channel_external;

DROP TABLE IF EXISTS webhook_events;

COMMIT;
