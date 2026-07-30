CREATE TABLE IF NOT EXISTS message_outboxes (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  message_id UUID NOT NULL REFERENCES messages(id) ON DELETE CASCADE,
  status VARCHAR(255) NOT NULL DEFAULT 'pending',
  attempts INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 3,
  next_retry_at TIMESTAMPTZ,
  last_error TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT message_outboxes_status_check CHECK (status IN ('pending', 'processing', 'sent', 'failed'))
);


CREATE INDEX IF NOT EXISTS idx_message_outboxes_status_next_retry_at ON message_outboxes(status, next_retry_at);
CREATE INDEX IF NOT EXISTS idx_message_outboxes_created_at ON message_outboxes(created_at);
CREATE INDEX IF NOT EXISTS idx_message_outboxes_updated_at ON message_outboxes(updated_at);


CREATE OR REPLACE FUNCTION update_message_outboxes_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_message_outboxes_modtime
BEFORE UPDATE ON message_outboxes
FOR EACH ROW
EXECUTE FUNCTION update_message_outboxes_modtime();
