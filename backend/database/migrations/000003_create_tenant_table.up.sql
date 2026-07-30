CREATE TABLE IF NOT EXISTS tenants (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  slug VARCHAR(255) UNIQUE NOT NULL,
  max_divisions INTEGER NOT NULL DEFAULT 1,
  max_agents INTEGER NOT NULL DEFAULT 1,
  max_quick_replies INTEGER NOT NULL DEFAULT 20,
  max_waha_whatsapp INTEGER NOT NULL DEFAULT 0,
  max_meta_whatsapp INTEGER NOT NULL DEFAULT 0,
  max_meta_messenger INTEGER NOT NULL DEFAULT 1,
  max_instagram INTEGER NOT NULL DEFAULT 1,
  max_telegram INTEGER NOT NULL DEFAULT 1,
  max_webchat INTEGER NOT NULL DEFAULT 0,
  max_linkchat INTEGER NOT NULL DEFAULT 0,
  ai_credits INTEGER NOT NULL DEFAULT 0,
  subscription_plan VARCHAR(255) NOT NULL DEFAULT 'free',
  subscription_status VARCHAR(255) NOT NULL DEFAULT 'active',
  trial_ends_at TIMESTAMPTZ,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  settings JSONB DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,
  CONSTRAINT chk_subscription_plan CHECK (subscription_plan IN ('free', 'starter', 'pro', 'enterprise')),
  CONSTRAINT chk_subscription_status CHECK (subscription_status IN ('active', 'suspended', 'canceled', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_tenants_slug ON tenants(slug);
CREATE INDEX IF NOT EXISTS idx_tenants_is_active ON tenants(is_active);
CREATE INDEX IF NOT EXISTS idx_tenants_created_at ON tenants(created_at);
CREATE INDEX IF NOT EXISTS idx_tenants_updated_at ON tenants(updated_at);
CREATE INDEX IF NOT EXISTS idx_tenants_deleted_at ON tenants(deleted_at);
CREATE INDEX IF NOT EXISTS idx_tenants_subscription_status ON tenants(subscription_status);

CREATE OR REPLACE FUNCTION update_tenants_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_tenants_modtime
BEFORE UPDATE ON tenants
FOR EACH ROW
EXECUTE FUNCTION update_tenants_modtime();

COMMENT ON TABLE tenants IS 'Multi-tenant organizations/companies';
COMMENT ON COLUMN tenants.slug IS 'URL-friendly unique identifier';
COMMENT ON COLUMN tenants.subscription_plan IS 'Subscription tier: free, starter, pro, enterprise';
COMMENT ON COLUMN tenants.subscription_status IS 'Subscription status: active, suspended, canceled, expired';