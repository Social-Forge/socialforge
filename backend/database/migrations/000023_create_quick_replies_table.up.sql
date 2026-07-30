CREATE TABLE IF NOT EXISTS quick_replies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    shortcut VARCHAR(255) UNIQUE NOT NULL,
    content_type VARCHAR(255) NOT NULL DEFAULT 'text',
    body TEXT,
    media JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT quick_replies_content_type_check CHECK (content_type IN ('text', 'image', 'video', 'document'))
);

CREATE INDEX IF NOT EXISTS idx_quick_replies_tenant_id_shortcut ON quick_replies (tenant_id, shortcut);
CREATE INDEX IF NOT EXISTS idx_quick_replies_tenant_id ON quick_replies (tenant_id);
CREATE INDEX IF NOT EXISTS idx_quick_replies_created_at ON quick_replies (created_at);
CREATE INDEX IF NOT EXISTS idx_quick_replies_updated_at ON quick_replies (updated_at);


CREATE OR REPLACE FUNCTION update_quick_replies_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_quick_replies_modtime
BEFORE UPDATE ON quick_replies
FOR EACH ROW
EXECUTE FUNCTION update_quick_replies_modtime();

