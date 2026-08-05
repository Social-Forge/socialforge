package repository

import (
	"context"
	"errors"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ============================ Subscription ============================

type SubscriptionRepository interface {
	BaseRepository
	FindActiveByTenant(ctx context.Context, tenantID uuid.UUID) (*entity.Subscription, error)
	Create(ctx context.Context, s *entity.Subscription) error
	Update(ctx context.Context, s *entity.Subscription) (*entity.Subscription, error)
	// ListExpired returns active subscriptions whose period ended before `at`
	// (cross-tenant; used by the worker sweep — runs pool-level, no RLS ctx).
	ListExpired(ctx context.Context, at time.Time) ([]*entity.Subscription, error)
}

type subscriptionRepository struct{ *baseRepository }

func NewSubscriptionRepository(db *pgxpool.Pool) SubscriptionRepository {
	return &subscriptionRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *subscriptionRepository) FindActiveByTenant(ctx context.Context, tenantID uuid.UUID) (*entity.Subscription, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var s entity.Subscription
	query := `SELECT * FROM subscriptions WHERE tenant_id = $1
		ORDER BY (status IN ('active','trailing')) DESC, created_at DESC LIMIT 1`
	if err := pgxscan.Get(subCtx, r.q(subCtx), &s, query, tenantID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find subscription: %w", err)
	}
	return &s, nil
}

func (r *subscriptionRepository) Create(ctx context.Context, s *entity.Subscription) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	query := `
		INSERT INTO subscriptions (id, tenant_id, plan_id, status, current_period_start, current_period_end, cancel_at_period_end)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		s.ID, s.TenantID, s.PlanID, s.Status, s.CurrentPeriodStart, s.CurrentPeriodEnd, s.CancelAtPeriodEnd,
	).Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}
	return nil
}

func (r *subscriptionRepository) Update(ctx context.Context, s *entity.Subscription) (*entity.Subscription, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out entity.Subscription
	query := `
		UPDATE subscriptions SET plan_id=$1, status=$2, current_period_start=$3, current_period_end=$4, cancel_at_period_end=$5
		WHERE id=$6 RETURNING *`
	if err := pgxscan.Get(subCtx, r.q(subCtx), &out, query,
		s.PlanID, s.Status, s.CurrentPeriodStart, s.CurrentPeriodEnd, s.CancelAtPeriodEnd, s.ID,
	); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}
	return &out, nil
}

func (r *subscriptionRepository) ListExpired(ctx context.Context, at time.Time) ([]*entity.Subscription, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 20*time.Second)
	defer cancel()
	var out []*entity.Subscription
	query := `SELECT * FROM subscriptions
		WHERE status IN ('active','trailing','past_due')
		  AND current_period_end IS NOT NULL AND current_period_end < $1`
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out, query, at); err != nil {
		return nil, fmt.Errorf("failed to list expired subscriptions: %w", err)
	}
	return out, nil
}

// ============================ Invoice ============================

type InvoiceRepository interface {
	BaseRepository
	Create(ctx context.Context, inv *entity.Invoice) error
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error)
	// FindByProviderInvoiceID resolves an invoice from a webhook. Runs pool-level
	// (no tenant ctx) — in prod with a non-superuser role this needs a
	// SECURITY DEFINER function (same caveat as webhook channel resolution).
	FindByProviderInvoiceID(ctx context.Context, provider, providerInvoiceID string) (*entity.Invoice, error)
	SetProviderInfo(ctx context.Context, id uuid.UUID, providerInvoiceID, checkoutURL string) error
	MarkPaid(ctx context.Context, id uuid.UUID, paidAt time.Time) (bool, error)
	SetStatus(ctx context.Context, id uuid.UUID, status string) error
	ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Invoice, error)
}

type invoiceRepository struct{ *baseRepository }

func NewInvoiceRepository(db *pgxpool.Pool) InvoiceRepository {
	return &invoiceRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *invoiceRepository) Create(ctx context.Context, inv *entity.Invoice) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if inv.ID == uuid.Nil {
		inv.ID = uuid.New()
	}
	query := `
		INSERT INTO invoices (id, tenant_id, number, status, amount, currency, description, purpose, provider, provider_invoice_id, checkout_url, expires_at)
		VALUES ($1,$2,nextval('invoice_number_seq'),$3,$4,$5,$6,$7,$8,$9,$10,$11)
		RETURNING id, number, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		inv.ID, inv.TenantID, inv.Status, inv.Amount, inv.Currency, inv.Description, inv.Purpose,
		inv.Provider, inv.ProviderInvoiceID, inv.CheckoutURL, inv.ExpiresAt,
	).Scan(&inv.ID, &inv.Number, &inv.CreatedAt, &inv.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create invoice: %w", err)
	}
	return nil
}

func (r *invoiceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Invoice, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var inv entity.Invoice
	if err := pgxscan.Get(subCtx, r.q(subCtx), &inv, `SELECT * FROM invoices WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to find invoice: %w", err)
	}
	return &inv, nil
}

func (r *invoiceRepository) FindByProviderInvoiceID(ctx context.Context, provider, providerInvoiceID string) (*entity.Invoice, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var inv entity.Invoice
	query := `SELECT * FROM invoices WHERE provider = $1 AND provider_invoice_id = $2 LIMIT 1`
	if err := pgxscan.Get(subCtx, r.q(subCtx), &inv, query, provider, providerInvoiceID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, fmt.Errorf("failed to find invoice by provider id: %w", err)
	}
	return &inv, nil
}

func (r *invoiceRepository) SetProviderInfo(ctx context.Context, id uuid.UUID, providerInvoiceID, checkoutURL string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	_, err := r.q(subCtx).Exec(subCtx,
		`UPDATE invoices SET provider_invoice_id=$1, checkout_url=$2 WHERE id=$3`,
		providerInvoiceID, checkoutURL, id)
	if err != nil {
		return fmt.Errorf("failed to set invoice provider info: %w", err)
	}
	return nil
}

// MarkPaid transitions pending -> paid atomically. Returns true only for the
// transition that actually happened (idempotent webhook: a second callback finds
// it already paid and returns false, so effects apply exactly once).
func (r *invoiceRepository) MarkPaid(ctx context.Context, id uuid.UUID, paidAt time.Time) (bool, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tag, err := r.q(subCtx).Exec(subCtx,
		`UPDATE invoices SET status='paid', paid_at=$1 WHERE id=$2 AND status='pending'`, paidAt, id)
	if err != nil {
		return false, fmt.Errorf("failed to mark invoice paid: %w", err)
	}
	return tag.RowsAffected() > 0, nil
}

func (r *invoiceRepository) SetStatus(ctx context.Context, id uuid.UUID, status string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	_, err := r.q(subCtx).Exec(subCtx, `UPDATE invoices SET status=$1 WHERE id=$2`, status, id)
	if err != nil {
		return fmt.Errorf("failed to set invoice status: %w", err)
	}
	return nil
}

func (r *invoiceRepository) ListByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.Invoice, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out []*entity.Invoice
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out,
		`SELECT * FROM invoices WHERE tenant_id = $1 ORDER BY created_at DESC`, tenantID); err != nil {
		return nil, fmt.Errorf("failed to list invoices: %w", err)
	}
	return out, nil
}

// ============================ Payment Event ============================

type PaymentEventRepository interface {
	BaseRepository
	Create(ctx context.Context, e *entity.PaymentEvent) error
}

type paymentEventRepository struct{ *baseRepository }

func NewPaymentEventRepository(db *pgxpool.Pool) PaymentEventRepository {
	return &paymentEventRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *paymentEventRepository) Create(ctx context.Context, e *entity.PaymentEvent) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	query := `
		INSERT INTO payment_events (id, tenant_id, invoice_id, provider, event_type, external_id, payload)
		VALUES ($1,$2,$3,$4,$5,$6,$7)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		e.ID, e.TenantID, e.InvoiceID, e.Provider, e.EventType, e.ExternalID, e.Payload,
	).Scan(&e.ID, &e.CreatedAt, &e.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create payment event: %w", err)
	}
	return nil
}

// ============================ Subscription Addon ============================

type SubscriptionAddonRepository interface {
	BaseRepository
	Create(ctx context.Context, a *entity.SubscriptionAddon) error
	ListActiveByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.SubscriptionAddon, error)
}

type subscriptionAddonRepository struct{ *baseRepository }

func NewSubscriptionAddonRepository(db *pgxpool.Pool) SubscriptionAddonRepository {
	return &subscriptionAddonRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *subscriptionAddonRepository) Create(ctx context.Context, a *entity.SubscriptionAddon) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	query := `
		INSERT INTO subscription_addons (id, tenant_id, type, quantity, meta, expires_at)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		a.ID, a.TenantID, a.Type, a.Quantity, a.Meta, a.ExpiresAt,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create subscription addon: %w", err)
	}
	return nil
}

func (r *subscriptionAddonRepository) ListActiveByTenant(ctx context.Context, tenantID uuid.UUID) ([]*entity.SubscriptionAddon, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out []*entity.SubscriptionAddon
	query := `SELECT * FROM subscription_addons
		WHERE tenant_id = $1 AND (expires_at IS NULL OR expires_at > now())
		ORDER BY created_at DESC`
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out, query, tenantID); err != nil {
		return nil, fmt.Errorf("failed to list addons: %w", err)
	}
	return out, nil
}
