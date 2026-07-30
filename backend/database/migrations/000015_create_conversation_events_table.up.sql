CREATE TABLE IF NOT EXISTS conversation_events (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  actor_id UUID REFERENCES users(id) ON DELETE SET NULL,
  type VARCHAR(255) NOT NULL,
  metadata JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT conversation_events_type_check CHECK (type IN ('assigned', 'unassigned', 'completed', 'reopened', 'archived', 'unarchived', 'labeled', 'note', 'call_rejected', 'auto_response'))
);

CREATE INDEX IF NOT EXISTS idx_conversation_events_conversation_id_created_at ON conversation_events(conversation_id, created_at);
CREATE INDEX IF NOT EXISTS idx_conversation_events_updated_at ON conversation_events(updated_at);


CREATE OR REPLACE FUNCTION update_conversation_events_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_conversation_events_modtime
BEFORE UPDATE ON conversation_events
FOR EACH ROW
EXECUTE FUNCTION update_conversation_events_modtime();
