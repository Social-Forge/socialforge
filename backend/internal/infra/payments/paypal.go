package payments

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// PaypalGateway implements Gateway against PayPal Orders v2.
// NOTE: PayPal does not support IDR; plan amounts must be converted to a
// supported currency (e.g. USD) before checkout — currency conversion is a
// refinement. The adapter passes the currency through as given.
type PaypalGateway struct {
	clientID     string
	clientSecret string
	webhookID    string
	isProd       bool
	httpClient   *http.Client
}

func NewPaypalGateway(clientID, clientSecret, webhookID string, isProd bool) *PaypalGateway {
	return &PaypalGateway{
		clientID:     clientID,
		clientSecret: clientSecret,
		webhookID:    webhookID,
		isProd:       isProd,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *PaypalGateway) Name() string { return "paypal" }

func (g *PaypalGateway) base() string {
	if g.isProd {
		return "https://api-m.paypal.com"
	}
	return "https://api-m.sandbox.paypal.com"
}

func (g *PaypalGateway) token(ctx context.Context) (string, error) {
	form := url.Values{"grant_type": {"client_credentials"}}
	req, err := http.NewRequestWithContext(ctx, "POST", g.base()+"/v1/oauth2/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	auth := base64.StdEncoding.EncodeToString([]byte(g.clientID + ":" + g.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("paypal oauth failed (status %d): %s", resp.StatusCode, string(body))
	}
	var out struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &out); err != nil || out.AccessToken == "" {
		return "", fmt.Errorf("paypal oauth: no access token")
	}
	return out.AccessToken, nil
}

func (g *PaypalGateway) CreateInvoice(ctx context.Context, p CreateInvoiceParams) (*CreatedInvoice, error) {
	if g.clientID == "" || g.clientSecret == "" {
		return nil, fmt.Errorf("paypal credentials not configured")
	}
	tok, err := g.token(ctx)
	if err != nil {
		return nil, err
	}

	// Amount value string. IDR has no minor units; USD uses 2 decimals.
	value := strconv.Itoa(p.Amount)
	if p.Currency == "USD" {
		value = fmt.Sprintf("%.2f", float64(p.Amount)/100.0)
	}
	orderReq := map[string]interface{}{
		"intent": "CAPTURE",
		"purchase_units": []map[string]interface{}{{
			"custom_id":   p.ExternalID, // round-trips back in the webhook resource
			"description": p.Description,
			"amount": map[string]string{
				"currency_code": p.Currency,
				"value":         value,
			},
		}},
		"application_context": map[string]string{
			"return_url": p.SuccessURL,
			"cancel_url": p.FailureURL,
		},
	}
	data, _ := json.Marshal(orderReq)
	req, err := http.NewRequestWithContext(ctx, "POST", g.base()+"/v2/checkout/orders", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal create order failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var out struct {
		ID    string `json:"id"`
		Links []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	}
	_ = json.Unmarshal(body, &out)
	if resp.StatusCode >= 300 || out.ID == "" {
		return nil, fmt.Errorf("paypal create order failed (status %d): %s", resp.StatusCode, string(body))
	}
	approve := ""
	for _, l := range out.Links {
		if l.Rel == "approve" || l.Rel == "payer-action" {
			approve = l.Href
			break
		}
	}
	if approve == "" {
		return nil, fmt.Errorf("paypal: no approve link returned")
	}
	return &CreatedInvoice{ProviderInvoiceID: out.ID, CheckoutURL: approve}, nil
}

type paypalWebhookEvent struct {
	ID        string `json:"id"`
	EventType string `json:"event_type"`
	Resource  struct {
		ID       string `json:"id"`
		Status   string `json:"status"`
		CustomID string `json:"custom_id"`
		Amount   struct {
			Value string `json:"value"`
		} `json:"amount"`
		SupplementaryData struct {
			RelatedIDs struct {
				OrderID string `json:"order_id"`
			} `json:"related_ids"`
		} `json:"supplementary_data"`
	} `json:"resource"`
}

// VerifyWebhook authenticates the callback via PayPal's verify-webhook-signature
// API, then parses the event. custom_id round-trips our invoice id.
func (g *PaypalGateway) VerifyWebhook(headers map[string]string, body []byte) (*WebhookResult, error) {
	if g.webhookID == "" {
		return nil, fmt.Errorf("paypal webhook id not configured")
	}
	tok, err := g.token(context.Background())
	if err != nil {
		return nil, err
	}
	var rawEvent json.RawMessage = body
	verifyReq := map[string]interface{}{
		"transmission_id":   headers["Paypal-Transmission-Id"],
		"transmission_time": headers["Paypal-Transmission-Time"],
		"cert_url":          headers["Paypal-Cert-Url"],
		"auth_algo":         headers["Paypal-Auth-Algo"],
		"transmission_sig":  headers["Paypal-Transmission-Sig"],
		"webhook_id":        g.webhookID,
		"webhook_event":     rawEvent,
	}
	data, _ := json.Marshal(verifyReq)
	req, err := http.NewRequestWithContext(context.Background(), "POST", g.base()+"/v1/notifications/verify-webhook-signature", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	req.Header.Set("Content-Type", "application/json")
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("paypal verify webhook failed: %w", err)
	}
	defer resp.Body.Close()
	vb, _ := io.ReadAll(resp.Body)
	var vr struct {
		VerificationStatus string `json:"verification_status"`
	}
	_ = json.Unmarshal(vb, &vr)
	if vr.VerificationStatus != "SUCCESS" {
		return nil, fmt.Errorf("paypal webhook verification failed: %s", vr.VerificationStatus)
	}

	var ev paypalWebhookEvent
	if err := json.Unmarshal(body, &ev); err != nil {
		return nil, fmt.Errorf("parse paypal event: %w", err)
	}
	paid := ev.EventType == "PAYMENT.CAPTURE.COMPLETED" ||
		(ev.EventType == "CHECKOUT.ORDER.APPROVED" && ev.Resource.Status == "APPROVED")
	expired := ev.EventType == "CHECKOUT.ORDER.VOIDED" || ev.EventType == "PAYMENT.CAPTURE.DENIED"

	orderID := ev.Resource.SupplementaryData.RelatedIDs.OrderID
	if orderID == "" {
		orderID = ev.Resource.ID
	}
	return &WebhookResult{
		ProviderInvoiceID: orderID,
		ExternalID:        ev.Resource.CustomID,
		EventID:           ev.ID,
		EventType:         ev.EventType,
		Paid:              paid,
		Expired:           expired,
	}, nil
}
