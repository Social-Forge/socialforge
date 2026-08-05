-- Seed the plan catalog. Idempotent: re-running updates prices/features.
-- `features` numeric entitlements mirror the tenant Max* limits + AI quotas.
INSERT INTO plans (code, name, price, currency, interval, features, is_active, sort) VALUES
('free', 'Free', 0, 'IDR', 'monthly',
  '{"divisions":1,"agents":1,"quick_replies":5,"ai_agents":1,"ai_credits":0,"waha_whatsapp":0,"meta_whatsapp":0,"meta_messenger":1,"instagram":1,"telegram":1,"webchat":1,"linkchat":1}',
  TRUE, 0),
('starter', 'Starter', 149000, 'IDR', 'monthly',
  '{"divisions":5,"agents":5,"quick_replies":100,"ai_agents":2,"ai_credits":50000,"waha_whatsapp":1,"meta_whatsapp":1,"meta_messenger":5,"instagram":5,"telegram":5,"webchat":3,"linkchat":3}',
  TRUE, 1),
('pro', 'Pro', 499000, 'IDR', 'monthly',
  '{"divisions":20,"agents":20,"quick_replies":500,"ai_agents":5,"ai_credits":200000,"waha_whatsapp":5,"meta_whatsapp":5,"meta_messenger":10,"instagram":10,"telegram":10,"webchat":10,"linkchat":10}',
  TRUE, 2),
('enterprise', 'Enterprise', 1999000, 'IDR', 'monthly',
  '{"divisions":100,"agents":100,"quick_replies":1000,"ai_agents":50,"ai_credits":1000000,"waha_whatsapp":10,"meta_whatsapp":10,"meta_messenger":100,"instagram":100,"telegram":100,"webchat":100,"linkchat":100}',
  TRUE, 3)
ON CONFLICT (code) DO UPDATE SET
  name = EXCLUDED.name,
  price = EXCLUDED.price,
  features = EXCLUDED.features,
  is_active = EXCLUDED.is_active,
  sort = EXCLUDED.sort,
  updated_at = now();
