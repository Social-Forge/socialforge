CREATE TABLE IF NOT EXISTS plans (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL UNIQUE, -- free | pro | ... — stable identifier used across the app. --
    name TEXT NOT NULL,
    price INTEGER NOT NULL, -- Price in minor units (IDR has none, so this is whole rupiah). --
    currency TEXT NOT NULL DEFAULT 'IDR',
    interval TEXT NOT NULL DEFAULT 'monthly',
    features JSONB DEFAULT '{}', -- Entitlements: { channels: {type: n}, agents, aiCredits, aiAgents, quickReplies } --
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    sort INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_plans_code ON plans (code);
CREATE INDEX IF NOT EXISTS idx_plans_name ON plans (name);
CREATE INDEX IF NOT EXISTS idx_plans_created_at ON plans (created_at);
CREATE INDEX IF NOT EXISTS idx_plans_updated_at ON plans (updated_at);


CREATE OR REPLACE FUNCTION update_plans_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_plans_modtime
BEFORE UPDATE ON plans
FOR EACH ROW
EXECUTE FUNCTION update_plans_modtime();
