CREATE TABLE IF NOT EXISTS contacts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  channel_id UUID NOT NULL REFERENCES channels(id) ON DELETE CASCADE,
  external_id TEXT NOT NULL,  -- Provider-side identity: phone (wa), PSID (messenger), IG id, chat id (tg).
  display_name VARCHAR(255),
  avatar_url TEXT,
  is_blocked BOOLEAN DEFAULT FALSE,
  attributes JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_contacts_channel_id_external_id ON contacts(channel_id, external_id);
CREATE INDEX IF NOT EXISTS idx_contacts_tenant_id ON contacts(tenant_id);
CREATE INDEX IF NOT EXISTS idx_contacts_channel_id ON contacts(channel_id);
CREATE INDEX IF NOT EXISTS idx_contacts_external_id ON contacts(external_id) WHERE external_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_contacts_display_name ON contacts(display_name) WHERE display_name IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_contacts_created_at ON contacts(created_at);
CREATE INDEX IF NOT EXISTS idx_contacts_updated_at ON contacts(updated_at);


CREATE OR REPLACE FUNCTION update_contacts_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_contacts_modtime
BEFORE UPDATE ON contacts
FOR EACH ROW
EXECUTE FUNCTION update_contacts_modtime();
