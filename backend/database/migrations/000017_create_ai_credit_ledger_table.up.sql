CREATE TABLE IF NOT EXISTS ai_credit_ledgers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  conversation_id UUID REFERENCES conversations(id) ON DELETE SET NULL,
  message_id UUID REFERENCES messages(id) ON DELETE SET NULL,
  delta INTEGER NOT NULL DEFAULT 0,
  balance_after INTEGER NOT NULL DEFAULT 0,
  reason TEXT NOT NULL,
  model VARCHAR(255),
  input_tokens INTEGER,
  output_tokens INTEGER,
  cost_usd DECIMAL(12,6) NOT NULL DEFAULT 0,
  credit DECIMAL(12,6) NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT ai_credit_ledgers_reason_check CHECK (reason IN ('grant', 'topup', 'debit', 'adjustment'))
);

CREATE INDEX IF NOT EXISTS idx_ai_credit_ledgers_tenant_id_created_at ON ai_credit_ledgers(tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_ai_credit_ledgers_conversation_id ON ai_credit_ledgers(conversation_id);
CREATE INDEX IF NOT EXISTS idx_ai_credit_ledgers_message_id ON ai_credit_ledgers(message_id) WHERE message_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_ai_credit_ledgers_updated_at ON ai_credit_ledgers(updated_at);

CREATE OR REPLACE FUNCTION update_ai_credit_ledgers_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_ai_credit_ledgers_modtime
BEFORE UPDATE ON ai_credit_ledgers
FOR EACH ROW
EXECUTE FUNCTION update_ai_credit_ledgers_modtime();
