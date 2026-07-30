CREATE TABLE IF NOT EXISTS conversation_labels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    conversation_id uuid NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    label_id uuid NOT NULL REFERENCES labels(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_labels_conversation_id_label_id ON conversation_labels (conversation_id, label_id);
CREATE INDEX IF NOT EXISTS idx_labels_conversation_id ON conversation_labels (conversation_id);
CREATE INDEX IF NOT EXISTS idx_labels_label_id ON conversation_labels (label_id);
CREATE INDEX IF NOT EXISTS idx_labels_created_at ON conversation_labels (created_at);
CREATE INDEX IF NOT EXISTS idx_labels_updated_at ON conversation_labels (updated_at);


CREATE OR REPLACE FUNCTION update_conversation_labels_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_conversation_labels_modtime
BEFORE UPDATE ON conversation_labels
FOR EACH ROW
EXECUTE FUNCTION update_conversation_labels_modtime();
