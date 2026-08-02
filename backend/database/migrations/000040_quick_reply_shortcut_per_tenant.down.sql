BEGIN;

ALTER TABLE quick_replies DROP CONSTRAINT IF EXISTS uq_quick_replies_tenant_shortcut;
ALTER TABLE quick_replies ADD CONSTRAINT quick_replies_shortcut_key UNIQUE (shortcut);

COMMIT;
