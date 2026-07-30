CREATE TABLE IF NOT EXISTS oauth_providers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  provider_name VARCHAR(255) NOT NULL,
  provider_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_oauth_provider_provider_id UNIQUE (provider_id, user_id);
  CONSTRAINT chk_oauth_provider_provider_name CHECK (provider_name IN ('google', 'facebook', 'github'))
);

CREATE INDEX IF NOT EXISTS idx_oauth_providers_provider_id ON oauth_providers(provider_id);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_provider_name ON oauth_providers(provider_name);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_user_id ON oauth_providers(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_created_at ON oauth_providers(created_at);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_updated_at ON oauth_providers(updated_at);

CREATE OR REPLACE FUNCTION update_oauth_providers_modtime()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER update_oauth_providers_modtime
BEFORE UPDATE ON oauth_providers
FOR EACH ROW
EXECUTE FUNCTION update_oauth_providers_modtime();
