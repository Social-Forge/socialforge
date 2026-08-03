BEGIN;

-- Per-channel automatic first-reply for new customers. Overridden when an AI
-- agent is active on the channel (checked in the ingestion pipeline).
CREATE TABLE IF NOT EXISTS auto_responses (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  is_enabled BOOLEAN NOT NULL DEFAULT FALSE,
  content_type VARCHAR(50) NOT NULL DEFAULT 'text',
  body TEXT,
  media JSONB DEFAULT '[]',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT uq_auto_responses_channel UNIQUE (channel_id),
  CONSTRAINT chk_auto_responses_content_type CHECK (content_type IN ('text', 'image', 'video', 'document'))
);

CREATE INDEX IF NOT EXISTS idx_auto_responses_tenant_id ON auto_responses(tenant_id);

ALTER TABLE IF EXISTS auto_responses ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS auto_responses FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON auto_responses USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

CREATE OR REPLACE TRIGGER update_auto_responses_modtime
BEFORE UPDATE ON auto_responses
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();

COMMIT;
