CREATE TABLE IF NOT EXISTS channels (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  division_id UUID NOT NULL REFERENCES divisions(id) ON DELETE CASCADE,
  ai_agent_id UUID REFERENCES ai_agents(id) ON DELETE SET NULL,
  type VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  status VARCHAR(255) NOT NULL DEFAULT 'disconnected',
  external_id TEXT,
  waha_engine VARCHAR(255),
  waha_session_name TEXT,
  webhook_secret TEXT,
  credentials JSONB DEFAULT '{}',
  settings JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT channels_name_length_check 
    CHECK (length(trim(name)) > 0),
  CONSTRAINT channels_type_check 
    CHECK (type IN ('whatsapp_waha', 'whatsapp_meta', 'messenger', 'instagram', 'telegram')),
  CONSTRAINT channels_status_check 
    CHECK (status IN ('disconnected', 'connected', 'connecting', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_channels_tenant_id ON channels(tenant_id);
CREATE INDEX IF NOT EXISTS idx_channels_ai_agent_id ON channels(ai_agent_id) WHERE ai_agent_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_channels_tenant_id_type ON channels(tenant_id, type);
CREATE INDEX IF NOT EXISTS idx_channels_division_id ON channels(division_id);
CREATE INDEX IF NOT EXISTS idx_channels_type ON channels(type);
CREATE INDEX IF NOT EXISTS idx_channels_status ON channels(status);
CREATE INDEX IF NOT EXISTS idx_channels_created_at ON channels(created_at);
CREATE INDEX IF NOT EXISTS idx_channels_updated_at ON channels(updated_at);


CREATE OR REPLACE FUNCTION update_channels_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_channels_modtime
BEFORE UPDATE ON channels
FOR EACH ROW
EXECUTE FUNCTION update_channels_modtime();
