BEGIN;

ALTER TABLE channels DROP CONSTRAINT IF EXISTS channels_type_check;
ALTER TABLE channels ADD CONSTRAINT channels_type_check
  CHECK (type IN ('whatsapp_waha', 'whatsapp_meta', 'messenger', 'instagram', 'telegram'));

COMMIT;
