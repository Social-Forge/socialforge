CREATE TABLE IF NOT EXISTS ai_agents (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  name VARCHAR(255) NOT NULL,
  provider VARCHAR(255) NOT NULL DEFAULT 'claude',
  model VARCHAR(255) NOT NULL,
  system_prompt TEXT NOT NULL,
  persona JSONB DEFAULT '{}',
  safety JSONB DEFAULT '{}',
  guardrails JSONB DEFAULT '{}',
  temperature FLOAT8,
  max_tokens INTEGER NOT NULL DEFAULT 1024,
  auto_reply_enabled BOOLEAN NOT NULL DEFAULT TRUE,
  working_hours JSONB DEFAULT '{}',
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT ai_agents_provider_CHECK CHECK (provider IN ('claude', 'openai', 'google'))
);

CREATE INDEX IF NOT EXISTS idx_ai_agents_tenant_id ON ai_agents(tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_agents_provider ON ai_agents(provider);
CREATE INDEX IF NOT EXISTS idx_ai_agents_name ON ai_agents(name);
CREATE INDEX IF NOT EXISTS idx_ai_agents_model ON ai_agents(model);
CREATE INDEX IF NOT EXISTS idx_ai_agents_is_active ON ai_agents(is_active);
CREATE INDEX IF NOT EXISTS idx_ai_agents_created_at ON ai_agents(created_at);
CREATE INDEX IF NOT EXISTS idx_ai_agents_updated_at ON ai_agents(updated_at);


CREATE OR REPLACE FUNCTION update_ai_agents_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_ai_agents_modtime
BEFORE UPDATE ON ai_agents
FOR EACH ROW
EXECUTE FUNCTION update_ai_agents_modtime();
