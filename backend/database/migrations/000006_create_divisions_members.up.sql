CREATE TABLE IF NOT EXISTS divisions_members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_tenant_member_tenant_id_user_id UNIQUE (tenant_id, user_id),
  CONSTRAINT chk_tenant_member_user_id CHECK (user_id <> tenant_id)
);

CREATE INDEX IF NOT EXISTS idx_divisions_members_tenant_id_user_id ON divisions_members(tenant_id, user_id);
CREATE INDEX IF NOT EXISTS idx_divisions_members_created_at ON divisions_members(created_at);
CREATE INDEX IF NOT EXISTS idx_divisions_members_updated_at ON divisions_members(updated_at);

CREATE OR REPLACE FUNCTION update_divisions_members_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_divisions_members_modtime
BEFORE UPDATE ON divisions_members
FOR EACH ROW
EXECUTE FUNCTION update_divisions_members_modtime();
