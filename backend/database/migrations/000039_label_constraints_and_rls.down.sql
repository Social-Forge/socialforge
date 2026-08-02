BEGIN;

DROP POLICY IF EXISTS tenant_isolation ON conversation_labels;
ALTER TABLE IF EXISTS conversation_labels DISABLE ROW LEVEL SECURITY;

ALTER TABLE conversation_labels DROP CONSTRAINT IF EXISTS uq_conversation_labels;
ALTER TABLE labels DROP CONSTRAINT IF EXISTS uq_labels_tenant_name;

COMMIT;
