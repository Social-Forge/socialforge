package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// ErrNoTenantContext is returned when a tenant-scoped operation is attempted
// without a tenant id present in the context.
var ErrNoTenantContext = errors.New("tenant id not found in context")

// Querier is the minimal database surface shared by *pgxpool.Pool and pgx.Tx.
// scany (pgxscan) and all repositories accept this so a single method body can
// run either directly on the pool or inside a tenant-bound transaction.
type Querier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type ctxKey int

const (
	tenantIDKey ctxKey = iota
	txKey
)

// WithTenantID stores the active tenant id in the context. The TenantGuard
// middleware calls this so the value flows all the way down to repositories,
// where it is applied to Postgres RLS via set_config('app.current_tenant', ...).
func WithTenantID(ctx context.Context, tenantID uuid.UUID) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// TenantIDFromContext returns the active tenant id, if any.
func TenantIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(tenantIDKey).(uuid.UUID)
	if !ok || id == uuid.Nil {
		return uuid.Nil, false
	}
	return id, true
}

// withTx stores an ambient transaction in the context so nested repository
// calls compose within a single tenant-bound transaction.
func withTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// txFromContext returns the ambient transaction, if the caller opened one.
func txFromContext(ctx context.Context) (pgx.Tx, bool) {
	tx, ok := ctx.Value(txKey).(pgx.Tx)
	return tx, ok
}

// The GUC (Grand Unified Configuration setting) name used by every RLS policy.
// Must match the migrations: current_setting('app.current_tenant', true).
const tenantGUC = "app.current_tenant"

// applyTenantGUC binds the given tenant id to the current transaction using
// set_config with is_local=true, so it is scoped to the transaction and reset
// automatically on commit/rollback (safe with a shared connection pool).
func applyTenantGUC(ctx context.Context, tx pgx.Tx, tenantID uuid.UUID) error {
	if _, err := tx.Exec(ctx, "SELECT set_config($1, $2, true)", tenantGUC, tenantID.String()); err != nil {
		return fmt.Errorf("failed to set tenant context: %w", err)
	}
	return nil
}
