BEGIN;

-- Labels are unique per tenant (spec: unique within 1 tenant).
ALTER TABLE labels
  ADD CONSTRAINT uq_labels_tenant_name UNIQUE (tenant_id, name);

-- A label can be attached to a conversation at most once.
ALTER TABLE conversation_labels
  ADD CONSTRAINT uq_conversation_labels UNIQUE (conversation_id, label_id);

CREATE INDEX IF NOT EXISTS idx_conversation_labels_conversation ON conversation_labels(conversation_id);
CREATE INDEX IF NOT EXISTS idx_conversation_labels_label ON conversation_labels(label_id);

-- conversation_labels is tenant-scoped -> enforce RLS like the other tables.
ALTER TABLE IF EXISTS conversation_labels ENABLE ROW LEVEL SECURITY;
ALTER TABLE IF EXISTS conversation_labels FORCE ROW LEVEL SECURITY;
CREATE POLICY tenant_isolation ON conversation_labels USING (tenant_id = current_setting('app.current_tenant', true)::uuid) WITH CHECK (tenant_id = current_setting('app.current_tenant', true)::uuid);

COMMIT;
