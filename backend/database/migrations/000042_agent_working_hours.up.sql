BEGIN;

-- Per-agent working hours. Auto-assign only routes to agents within a window
-- (agents with no hours configured are treated as always available).
CREATE TABLE IF NOT EXISTS agent_working_hours (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  day_of_week SMALLINT NOT NULL,       -- 0=Sunday .. 6=Saturday (matches EXTRACT(DOW))
  start_time TIME NOT NULL,
  end_time TIME NOT NULL,
  is_active BOOLEAN NOT NULL DEFAULT TRUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT chk_working_hours_dow CHECK (day_of_week BETWEEN 0 AND 6),
  CONSTRAINT chk_working_hours_range CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_agent_working_hours_user ON agent_working_hours(user_id);
CREATE INDEX IF NOT EXISTS idx_agent_working_hours_tenant ON agent_working_hours(tenant_id);

ALTER TABLE IF EXISTS agent_working_hours ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS agent_working_hours FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON agent_working_hours USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;
