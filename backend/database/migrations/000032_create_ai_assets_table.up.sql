CREATE TABLE IF NOT EXISTS ai_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    ai_agent_id UUID NOT NULL REFERENCES ai_agents(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(255) NOT NULL DEFAULT 'image',
    storage_key TEXT NOT NULL, -- MinIO object key — a fresh presigned URL is minted at send time.
    mime_type VARCHAR(255),
    size INTEGER,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_ai_asset_type CHECK (type IN ('image', 'video', 'document'))
);

CREATE INDEX IF NOT EXISTS idx_ai_assets_tenant_id ON ai_assets (tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_assets_ai_agent_id ON ai_assets (ai_agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_assets_created_at ON ai_assets (created_at);
CREATE INDEX IF NOT EXISTS idx_ai_assets_updated_at ON ai_assets (updated_at);

CREATE OR REPLACE FUNCTION update_ai_assets_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_ai_assets_modtime
BEFORE UPDATE ON ai_assets
FOR EACH ROW
EXECUTE FUNCTION update_ai_assets_modtime();

