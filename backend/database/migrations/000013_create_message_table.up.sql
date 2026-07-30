CREATE TABLE IF NOT EXISTS messages (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  conversation_id UUID NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
  sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
  direction VARCHAR(255) NOT NULL DEFAULT 'in',
  sender_type VARCHAR(255) NOT NULL DEFAULT 'contact',
  content_type VARCHAR(255) NOT NULL DEFAULT 'text',
  body TEXT,
  media JSONB DEFAULT '{}',
  provider_message_id TEXT UNIQUE,
  status VARCHAR(255) NOT NULL DEFAULT 'pending',
  reply_to_id UUID REFERENCES messages(id) ON DELETE SET NULL,
  is_pinned BOOLEAN NOT NULL DEFAULT FALSE,
  error TEXT,
  edited_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,

  CONSTRAINT messages_direction_check CHECK (direction IN ('in', 'out')),
  CONSTRAINT messages_sender_type_check CHECK (sender_type IN ('contact', 'agent', 'ai', 'system')),
  CONSTRAINT messages_content_type_check CHECK (content_type IN ('text', 'image', 'video', 'audio', 'document', 'location', 'contact', 'sticker', 'button', 'list', 'template')),
  CONSTRAINT messages_status_check CHECK (status IN ('pending', 'sent', 'delivered', 'read', 'failed'))
);


CREATE INDEX IF NOT EXISTS idx_messages_conversation_id_created_at ON messages(conversation_id, created_at);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id) WHERE sender_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_messages_status ON messages(status);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);
CREATE INDEX IF NOT EXISTS idx_messages_updated_at ON messages(updated_at);


CREATE OR REPLACE FUNCTION update_messages_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_messages_modtime
BEFORE UPDATE ON messages
FOR EACH ROW
EXECUTE FUNCTION update_messages_modtime();

