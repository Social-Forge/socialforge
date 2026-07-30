CREATE TABLE IF NOT EXISTS division_members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_tenant_id UUID NOT NULL REFERENCES user_tenants(id) ON DELETE CASCADE,
    division_id UUID NOT NULL REFERENCES divisions(id) ON DELETE CASCADE,
    is_active BOOLEAN DEFAULT TRUE,
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ,
    CONSTRAINT division_members_unique UNIQUE (user_tenant_id, division_id)
);
CREATE INDEX idx_division_members_user_tenant ON division_members(user_tenant_id);
CREATE INDEX idx_division_members_division ON division_members(division_id);

CREATE OR REPLACE FUNCTION update_division_members_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_division_members_modtime
BEFORE UPDATE ON division_members
FOR EACH ROW
EXECUTE FUNCTION update_division_members_modtime();
