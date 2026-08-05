package payments

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const xenditInvoiceURL = "https://api.xendit.co/v2/invoices"

// XenditGateway implements Gateway against Xendit's hosted Invoice API.
type XenditGateway struct {
	secretKey    string
	webhookToken string
	httpClient   *http.Client
}

func NewXenditGateway(secretKey, webhookToken string) *XenditGateway {
	return &XenditGateway{
		secretKey:    secretKey,
		webhookToken: webhookToken,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *XenditGateway) Name() string { return "xendit" }

type xenditCreateReq struct {
	ExternalID         string   `json:"external_id"`
	Amount             int      `json:"amount"`
	Description        string   `json:"description"`
	Currency           string   `json:"currency,omitempty"`
	PayerEmail         string   `json:"payer_email,omitempty"`
	SuccessRedirectURL string   `json:"success_redirect_url,omitempty"`
	FailureRedirectURL string   `json:"failure_redirect_url,omitempty"`
	InvoiceDuration    int      `json:"invoice_duration,omitempty"`
	Items              []string `json:"-"`
}

type xenditCreateResp struct {
	ID         string `json:"id"`
	InvoiceURL string `json:"invoice_url"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error_code"`
}

func (g *XenditGateway) CreateInvoice(ctx context.Context, p CreateInvoiceParams) (*CreatedInvoice, error) {
	if g.secretKey == "" {
		return nil, fmt.Errorf("xendit secret key not configured")
	}
	reqBody := xenditCreateReq{
		ExternalID:         p.ExternalID,
		Amount:             p.Amount,
		Description:        p.Description,
		Currency:           p.Currency,
		PayerEmail:         p.PayerEmail,
		SuccessRedirectURL: p.SuccessURL,
		FailureRedirectURL: p.FailureURL,
		InvoiceDuration:    86400, // 24h
	}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal xendit request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", xenditInvoiceURL, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create xendit request: %w", err)
	}
	// Xendit uses HTTP Basic auth: secret key as username, empty password.
	auth := base64.StdEncoding.EncodeToString([]byte(g.secretKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("xendit request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var out xenditCreateResp
	_ = json.Unmarshal(body, &out)
	if resp.StatusCode >= 300 || out.InvoiceURL == "" {
		msg := out.Message
		if msg == "" {
			msg = string(body)
		}
		return nil, fmt.Errorf("xendit create invoice failed (status %d): %s", resp.StatusCode, msg)
	}
	return &CreatedInvoice{ProviderInvoiceID: out.ID, CheckoutURL: out.InvoiceURL}, nil
}

type xenditWebhookBody struct {
	ID         string `json:"id"`
	ExternalID string `json:"external_id"`
	Status     string `json:"status"`
}

// VerifyWebhook checks the x-callback-token header against the configured token
// (Xendit's webhook auth), then parses the payload.
func (g *XenditGateway) VerifyWebhook(headers map[string]string, body []byte) (*WebhookResult, error) {
	token := headers["X-Callback-Token"]
	if token == "" {
		token = headers["x-callback-token"]
	}
	if g.webhookToken == "" {
		return nil, fmt.Errorf("xendit webhook token not configured")
	}
	if subtle.ConstantTimeCompare([]byte(token), []byte(g.webhookToken)) != 1 {
		return nil, fmt.Errorf("invalid xendit callback token")
	}
	var b xenditWebhookBody
	if err := json.Unmarshal(body, &b); err != nil {
		return nil, fmt.Errorf("parse xendit webhook: %w", err)
	}
	return &WebhookResult{
		ProviderInvoiceID: b.ID,
		ExternalID:        b.ExternalID,
		EventID:           b.ID,
		EventType:         "invoice." + b.Status,
		Paid:              b.Status == "PAID" || b.Status == "SETTLED",
		Expired:           b.Status == "EXPIRED",
	}, nil
}
