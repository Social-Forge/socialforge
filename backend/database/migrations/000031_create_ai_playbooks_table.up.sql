CREATE TABLE IF NOT EXISTS ai_playbooks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    ai_agent_id UUID NOT NULL REFERENCES ai_agents(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    -- Keywords/phrases in the customer message that trigger this playbook. --
    keywords JSONB NOT NULL DEFAULT '[]', 
    -- What the agent should do / say when triggered (scope, tone, offer, etc). --
    instruction TEXT NOT NULL,
    -- Asset ids to send to the customer when this playbook fires. --
    asset_ids JSONB NOT NULL DEFAULT '[]',
    -- Higher priority wins when multiple playbooks match. --
    priority INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_ai_playbooks_tenant_id ON ai_playbooks (tenant_id);
CREATE INDEX IF NOT EXISTS idx_ai_playbooks_ai_agent_id ON ai_playbooks (ai_agent_id);
CREATE INDEX IF NOT EXISTS idx_ai_playbooks_created_at ON ai_playbooks (created_at);
CREATE INDEX IF NOT EXISTS idx_ai_playbooks_updated_at ON ai_playbooks (updated_at);


CREATE OR REPLACE FUNCTION update_ai_playbooks_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_ai_playbooks_modtime
BEFORE UPDATE ON ai_playbooks
FOR EACH ROW
EXECUTE FUNCTION update_ai_playbooks_modtime();
