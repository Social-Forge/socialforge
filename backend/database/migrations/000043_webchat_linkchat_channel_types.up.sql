BEGIN;

-- Allow webchat + linkchat as channel types (limits already exist on tenants).
ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_type_check;
ALTER TABLE channels ADD CONSTRAINT channels_type_check
  CHECK (type IN ('whatsapp_waha', 'whatsapp_meta', 'messenger', 'instagram', 'telegram', 'webchat', 'linkchat'));

COMMIT;
