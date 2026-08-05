package payments

import "github/socialforge/config"

// BuildGateways constructs the enabled payment gateways keyed by provider name.
// A provider is included only when its credentials are configured. Midtrans and
// PayPal are added in Fase 6D.
func BuildGateways(cfg *config.PaymentConfig) map[string]Gateway {
	gws := make(map[string]Gateway)
	if cfg.XenditSecretKey != "" {
		gws["xendit"] = NewXenditGateway(cfg.XenditSecretKey, cfg.XenditWebhookToken)
	}
	if cfg.MidtransServerKey != "" {
		gws["midtrans"] = NewMidtransGateway(cfg.MidtransServerKey, cfg.MidtransIsProd)
	}
	if cfg.PaypalClientId != "" && cfg.PaypalClientSecret != "" {
		gws["paypal"] = NewPaypalGateway(cfg.PaypalClientId, cfg.PaypalClientSecret, cfg.PaypalWebhookId, cfg.PaypalisProd)
	}
	return gws
}
