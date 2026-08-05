package payments

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// MidtransGateway implements Gateway against Midtrans Snap.
type MidtransGateway struct {
	serverKey  string
	isProd     bool
	httpClient *http.Client
}

func NewMidtransGateway(serverKey string, isProd bool) *MidtransGateway {
	return &MidtransGateway{
		serverKey:  serverKey,
		isProd:     isProd,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (g *MidtransGateway) Name() string { return "midtrans" }

func (g *MidtransGateway) snapURL() string {
	if g.isProd {
		return "https://app.midtrans.com/snap/v1/transactions"
	}
	return "https://app.sandbox.midtrans.com/snap/v1/transactions"
}

type midtransSnapReq struct {
	TransactionDetails struct {
		OrderID     string `json:"order_id"`
		GrossAmount int    `json:"gross_amount"`
	} `json:"transaction_details"`
	CustomerDetails struct {
		Email string `json:"email,omitempty"`
	} `json:"customer_details"`
	Callbacks struct {
		Finish string `json:"finish,omitempty"`
	} `json:"callbacks"`
}

type midtransSnapResp struct {
	Token         string   `json:"token"`
	RedirectURL   string   `json:"redirect_url"`
	ErrorMessages []string `json:"error_messages"`
}

func (g *MidtransGateway) CreateInvoice(ctx context.Context, p CreateInvoiceParams) (*CreatedInvoice, error) {
	if g.serverKey == "" {
		return nil, fmt.Errorf("midtrans server key not configured")
	}
	var reqBody midtransSnapReq
	reqBody.TransactionDetails.OrderID = p.ExternalID
	reqBody.TransactionDetails.GrossAmount = p.Amount
	reqBody.CustomerDetails.Email = p.PayerEmail
	reqBody.Callbacks.Finish = p.SuccessURL

	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal midtrans request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", g.snapURL(), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("create midtrans request: %w", err)
	}
	// Basic auth: server key as username, empty password.
	auth := base64.StdEncoding.EncodeToString([]byte(g.serverKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("midtrans request failed: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var out midtransSnapResp
	_ = json.Unmarshal(body, &out)
	if resp.StatusCode >= 300 || out.RedirectURL == "" {
		return nil, fmt.Errorf("midtrans create transaction failed (status %d): %s", resp.StatusCode, string(body))
	}
	// order_id (our external id) is what the notification echoes back.
	return &CreatedInvoice{ProviderInvoiceID: p.ExternalID, CheckoutURL: out.RedirectURL}, nil
}

type midtransNotif struct {
	OrderID           string `json:"order_id"`
	StatusCode        string `json:"status_code"`
	GrossAmount       string `json:"gross_amount"`
	SignatureKey      string `json:"signature_key"`
	TransactionStatus string `json:"transaction_status"`
	FraudStatus       string `json:"fraud_status"`
}

// VerifyWebhook validates the Midtrans notification signature:
// sha512(order_id + status_code + gross_amount + server_key).
func (g *MidtransGateway) VerifyWebhook(headers map[string]string, body []byte) (*WebhookResult, error) {
	var n midtransNotif
	if err := json.Unmarshal(body, &n); err != nil {
		return nil, fmt.Errorf("parse midtrans notification: %w", err)
	}
	raw := n.OrderID + n.StatusCode + n.GrossAmount + g.serverKey
	sum := sha512.Sum512([]byte(raw))
	expected := hex.EncodeToString(sum[:])
	if subtle.ConstantTimeCompare([]byte(expected), []byte(n.SignatureKey)) != 1 {
		return nil, fmt.Errorf("invalid midtrans signature")
	}

	paid := (n.TransactionStatus == "settlement") ||
		(n.TransactionStatus == "capture" && n.FraudStatus == "accept")
	expired := n.TransactionStatus == "expire" || n.TransactionStatus == "cancel" || n.TransactionStatus == "deny"

	// Event id for audit: order_id + transaction_status (Midtrans has no event id).
	eventID := n.OrderID + ":" + n.TransactionStatus + ":" + strconv.Itoa(len(body))
	return &WebhookResult{
		ProviderInvoiceID: n.OrderID,
		ExternalID:        n.OrderID,
		EventID:           eventID,
		EventType:         "transaction." + n.TransactionStatus,
		Paid:              paid,
		Expired:           expired,
	}, nil
}
