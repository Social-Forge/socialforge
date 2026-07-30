CREATE TABLE IF NOT EXISTS ai_knowledge (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    ai_agent_id UUID NOT NULL REFERENCES ai_agents(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding JSONB DEFAULT '[]',
    token_count INTEGER DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ai_knowledge_ai_agent_id ON ai_knowledge (ai_agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_knowledge_tenant_id ON ai_knowledge (tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_knowledge_created_at ON ai_knowledge (created_at);
CREATE INDEX IF NOT EXISTS idx_ai_knowledge_updated_at ON ai_knowledge (updated_at);


CREATE OR REPLACE FUNCTION update_ai_knowledge_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_ai_knowledge_modtime
BEFORE UPDATE ON ai_knowledge
FOR EACH ROW
EXECUTE FUNCTION update_ai_knowledge_modtime();
