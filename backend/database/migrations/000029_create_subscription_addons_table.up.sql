CREATE TABLE IF NOT EXISTS subscription_addons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    type VARCHAR(255) NOT NULL DEFAULT 'channel_slot',
    quantity INTEGER NOT NULL DEFAULT 0,
    meta JSONB DEFAULT '{}', -- { channelType } for channel_slot, { agentSlot } for agent_slot, { credits } for ai_credits
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT chk_subscription_addon_type CHECK (type IN ('channel_slot', 'agent_slot', 'ai_credits'))
);

CREATE INDEX IF NOT EXISTS idx_subscription_addons_tenant_id_type ON subscription_addons (tenant_id, type);
CREATE INDEX IF NOT EXISTS idx_subscription_addons_tenant_id ON subscription_addons (tenant_id);
CREATE INDEX IF NOT EXISTS idx_subscription_addons_created_at ON subscription_addons (created_at);
CREATE INDEX IF NOT EXISTS idx_subscription_addons_updated_at ON subscription_addons (updated_at);



CREATE OR REPLACE FUNCTION update_subscription_addons_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_subscription_addons_modtime
BEFORE UPDATE ON subscription_addons
FOR EACH ROW
EXECUTE FUNCTION update_subscription_addons_modtime();
