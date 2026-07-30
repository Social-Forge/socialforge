CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    color VARCHAR(255) NOT NULL DEFAULT '#64748b',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_labels_tenant_id_name ON labels (tenant_id, name);
CREATE INDEX IF NOT EXISTS idx_labels_created_at ON labels (created_at);
CREATE INDEX IF NOT EXISTS idx_labels_updated_at ON labels (updated_at);


CREATE OR REPLACE FUNCTION update_labels_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_labels_modtime
BEFORE UPDATE ON labels
FOR EACH ROW
EXECUTE FUNCTION update_labels_modtime();
