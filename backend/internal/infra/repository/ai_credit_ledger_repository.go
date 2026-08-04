package repository

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AICreditLedgerRepository interface {
	BaseRepository
	// Balance returns the tenant's current AI credit balance.
	Balance(ctx context.Context, tenantID uuid.UUID) (int, error)
	// Debit atomically subtracts `amount` credits from the tenant balance and
	// records a `debit` ledger row capturing the token usage. It returns the
	// resulting balance. The balance is allowed to go negative (overage is
	// recorded, not silently dropped); callers gate on Balance beforehand.
	Debit(ctx context.Context, e *entity.AICreditLedger, amount int) (int, error)
}

type aiCreditLedgerRepository struct {
	*baseRepository
}

func NewAICreditLedgerRepository(db *pgxpool.Pool) AICreditLedgerRepository {
	return &aiCreditLedgerRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *aiCreditLedgerRepository) Balance(ctx context.Context, tenantID uuid.UUID) (int, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	var balance int
	err := r.q(subCtx).QueryRow(subCtx, `SELECT ai_credits FROM tenants WHERE id = $1`, tenantID).Scan(&balance)
	if err != nil {
		return 0, fmt.Errorf("failed to read ai credit balance: %w", err)
	}
	return balance, nil
}

func (r *aiCreditLedgerRepository) Debit(ctx context.Context, e *entity.AICreditLedger, amount int) (int, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 10*time.Second)
	defer cancel()

	if amount < 0 {
		amount = 0
	}
	// Atomically decrement the tenant balance and log the debit against the new
	// balance in one statement so the ledger can never drift from the balance.
	query := `
		WITH upd AS (
			UPDATE tenants SET ai_credits = ai_credits - $1, updated_at = now()
			WHERE id = $2 RETURNING ai_credits
		)
		INSERT INTO ai_credit_ledgers
			(tenant_id, conversation_id, message_id, delta, balance_after, reason, model, input_tokens, output_tokens, cost_usd, credit)
		SELECT $2, $3, $4, $5, upd.ai_credits, $6, $7, $8, $9, $10, $11 FROM upd
		RETURNING id, balance_after, created_at, updated_at`

	err := r.q(subCtx).QueryRow(subCtx, query,
		amount,           // $1 subtract
		e.TenantID,       // $2
		e.ConversationID, // $3
		e.MessageID,      // $4
		-amount,          // $5 delta (negative)
		entity.AICreditReasonDebit, // $6
		e.Model,          // $7
		e.InputTokens,    // $8
		e.OutputTokens,   // $9
		e.CostUSD,        // $10
		e.Credit,         // $11
	).Scan(&e.ID, &e.BalanceAfter, &e.CreatedAt, &e.UpdatedAt)
	if err != nil {
		return 0, fmt.Errorf("failed to record ai credit debit: %w", err)
	}
	e.Delta = -amount
	e.Reason = entity.AICreditReasonDebit
	return e.BalanceAfter, nil
}
