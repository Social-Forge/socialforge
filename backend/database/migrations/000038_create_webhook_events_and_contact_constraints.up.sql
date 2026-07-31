BEGIN;

-- Webhook-level idempotency: dedup inbound provider events before processing.
CREATE TABLE IF NOT EXISTS webhook_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  provider VARCHAR(50) NOT NULL,
  provider_event_id TEXT NOT NULL,
  event_type VARCHAR(100),
  payload JSONB NOT NULL DEFAULT '{}',
  processed_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_webhook_events_channel_event UNIQUE (channel_id, provider_event_id)
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_channel_id ON webhook_events(channel_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_created_at ON webhook_events(created_at);

-- Contact upsert key: one contact per (channel, external identity).
ALTER TABLE contacts
  ADD CONSTRAINT uq_contacts_channel_external UNIQUE (channel_id, external_id);

-- Contacts hold customer PII and are tenant-scoped -> enforce RLS like the
-- other messaging tables.
ALTER TABLE IF EXISTS contacts ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS contacts FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON contacts USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;
