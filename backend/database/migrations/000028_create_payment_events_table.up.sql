CREATE TABLE IF NOT EXISTS payment_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    invoice_id UUID REFERENCES invoices(id) ON DELETE SET NULL,
    provider VARCHAR(255) NOT NULL DEFAULT 'xendit',
    event_type VARCHAR(255) NOT NULL,
    external_id TEXT,
    payload JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_payment_event_provider CHECK (provider IN ('xendit', 'midtrans', 'paypal'))
);


CREATE INDEX IF NOT EXISTS idx_payment_events_tenant_id_created_at ON payment_events (tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_payment_events_invoice_id ON payment_events (invoice_id);
CREATE INDEX IF NOT EXISTS idx_payment_events_provider ON payment_events (provider);
CREATE INDEX IF NOT EXISTS idx_payment_events_event_type ON payment_events (event_type);
CREATE INDEX IF NOT EXISTS idx_payment_events_external_id ON payment_events (external_id);


CREATE OR REPLACE FUNCTION update_payment_events_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_payment_events_modtime
BEFORE UPDATE ON payment_events
FOR EACH ROW
EXECUTE FUNCTION update_payment_events_modtime();
