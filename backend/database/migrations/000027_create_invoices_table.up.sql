CREATE TABLE IF NOT EXISTS invoices (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    number INTEGER NOT NULL UNIQUE,
    status VARCHAR(255) NOT NULL DEFAULT 'pending',
    amount INTEGER NOT NULL,
    currency VARCHAR(255) NOT NULL DEFAULT 'IDR',
    description TEXT NOT NULL,
    purpose JSONB DEFAULT '{}',
    provider VARCHAR(255) NOT NULL DEFAULT 'xendit',
    provider_invoice_id TEXT,
    checkout_url TEXT,
    paid_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_invoice_status CHECK (status IN ('pending', 'paid', 'expired', 'failed')),
    CONSTRAINT chk_invoice_provider CHECK (provider IN ('xendit', 'midtrans', 'paypal'))
);

CREATE INDEX IF NOT EXISTS idx_invoices_tenant_id_created_at ON invoices (tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_invoices_provider_invoice_id ON invoices (provider_invoice_id);
CREATE INDEX IF NOT EXISTS idx_invoices_paid_at ON invoices (paid_at);
CREATE INDEX IF NOT EXISTS idx_invoices_expires_at ON invoices (expires_at);

CREATE OR REPLACE FUNCTION update_invoices_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_invoices_modtime
BEFORE UPDATE ON invoices
FOR EACH ROW
EXECUTE FUNCTION update_invoices_modtime();
