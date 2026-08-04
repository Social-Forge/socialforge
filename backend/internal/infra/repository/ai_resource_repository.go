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

// ============================ AI Knowledge ============================

type AIKnowledgeRepository interface {
	BaseRepository
	Create(ctx context.Context, k *entity.AIKnowledge) error
	ListByAgent(ctx context.Context, agentID uuid.UUID) ([]*entity.AIKnowledge, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.AIKnowledge, error)
	Update(ctx context.Context, k *entity.AIKnowledge) (*entity.AIKnowledge, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type aiKnowledgeRepository struct{ *baseRepository }

func NewAIKnowledgeRepository(db *pgxpool.Pool) AIKnowledgeRepository {
	return &aiKnowledgeRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *aiKnowledgeRepository) Create(ctx context.Context, k *entity.AIKnowledge) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if k.ID == uuid.Nil {
		k.ID = uuid.New()
	}
	query := `
		INSERT INTO ai_knowledge (id, tenant_id, ai_agent_id, title, content, token_count)
		VALUES ($1,$2,$3,$4,$5,$6)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		k.ID, k.TenantID, k.AIAgentID, k.Title, k.Content, k.TokenCount,
	).Scan(&k.ID, &k.CreatedAt, &k.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create ai knowledge: %w", err)
	}
	return nil
}

func (r *aiKnowledgeRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]*entity.AIKnowledge, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out []*entity.AIKnowledge
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out,
		`SELECT id, tenant_id, ai_agent_id, title, content, token_count, created_at, updated_at
		 FROM ai_knowledge WHERE ai_agent_id = $1 ORDER BY created_at DESC`, agentID); err != nil {
		return nil, fmt.Errorf("failed to list ai knowledge: %w", err)
	}
	return out, nil
}

func (r *aiKnowledgeRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.AIKnowledge, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var k entity.AIKnowledge
	if err := pgxscan.Get(subCtx, r.q(subCtx), &k,
		`SELECT id, tenant_id, ai_agent_id, title, content, token_count, created_at, updated_at
		 FROM ai_knowledge WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai knowledge not found")
		}
		return nil, fmt.Errorf("failed to find ai knowledge: %w", err)
	}
	return &k, nil
}

func (r *aiKnowledgeRepository) Update(ctx context.Context, k *entity.AIKnowledge) (*entity.AIKnowledge, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out entity.AIKnowledge
	err := r.q(subCtx).QueryRow(subCtx,
		`UPDATE ai_knowledge SET title=$1, content=$2, token_count=$3 WHERE id=$4
		 RETURNING id, tenant_id, ai_agent_id, title, content, token_count, created_at, updated_at`,
		k.Title, k.Content, k.TokenCount, k.ID,
	).Scan(&out.ID, &out.TenantID, &out.AIAgentID, &out.Title, &out.Content, &out.TokenCount, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai knowledge not found")
		}
		return nil, fmt.Errorf("failed to update ai knowledge: %w", err)
	}
	return &out, nil
}

func (r *aiKnowledgeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM ai_knowledge WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete ai knowledge: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ai knowledge not found")
	}
	return nil
}

// ============================ AI Playbook ============================

type AIPlaybookRepository interface {
	BaseRepository
	Create(ctx context.Context, p *entity.AIPlaybook) error
	ListByAgent(ctx context.Context, agentID uuid.UUID) ([]*entity.AIPlaybook, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.AIPlaybook, error)
	Update(ctx context.Context, p *entity.AIPlaybook) (*entity.AIPlaybook, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type aiPlaybookRepository struct{ *baseRepository }

func NewAIPlaybookRepository(db *pgxpool.Pool) AIPlaybookRepository {
	return &aiPlaybookRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *aiPlaybookRepository) Create(ctx context.Context, p *entity.AIPlaybook) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	query := `
		INSERT INTO ai_playbooks (id, tenant_id, ai_agent_id, name, keywords, instruction, asset_ids, priority, is_active)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		p.ID, p.TenantID, p.AIAgentID, p.Name, p.Keywords, p.Instruction, p.AssetIDs, p.Priority, p.IsActive,
	).Scan(&p.ID, &p.CreatedAt, &p.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create ai playbook: %w", err)
	}
	return nil
}

func (r *aiPlaybookRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]*entity.AIPlaybook, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out []*entity.AIPlaybook
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out,
		`SELECT * FROM ai_playbooks WHERE ai_agent_id = $1 ORDER BY priority DESC, created_at DESC`, agentID); err != nil {
		return nil, fmt.Errorf("failed to list ai playbooks: %w", err)
	}
	return out, nil
}

func (r *aiPlaybookRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.AIPlaybook, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var p entity.AIPlaybook
	if err := pgxscan.Get(subCtx, r.q(subCtx), &p, `SELECT * FROM ai_playbooks WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai playbook not found")
		}
		return nil, fmt.Errorf("failed to find ai playbook: %w", err)
	}
	return &p, nil
}

func (r *aiPlaybookRepository) Update(ctx context.Context, p *entity.AIPlaybook) (*entity.AIPlaybook, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out entity.AIPlaybook
	err := r.q(subCtx).QueryRow(subCtx,
		`UPDATE ai_playbooks SET name=$1, keywords=$2, instruction=$3, asset_ids=$4, priority=$5, is_active=$6
		 WHERE id=$7
		 RETURNING id, tenant_id, ai_agent_id, name, keywords, instruction, asset_ids, priority, is_active, created_at, updated_at`,
		p.Name, p.Keywords, p.Instruction, p.AssetIDs, p.Priority, p.IsActive, p.ID,
	).Scan(&out.ID, &out.TenantID, &out.AIAgentID, &out.Name, &out.Keywords, &out.Instruction,
		&out.AssetIDs, &out.Priority, &out.IsActive, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai playbook not found")
		}
		return nil, fmt.Errorf("failed to update ai playbook: %w", err)
	}
	return &out, nil
}

func (r *aiPlaybookRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM ai_playbooks WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete ai playbook: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ai playbook not found")
	}
	return nil
}

// ============================ AI Asset ============================

type AIAssetRepository interface {
	BaseRepository
	Create(ctx context.Context, a *entity.AIAsset) error
	ListByAgent(ctx context.Context, agentID uuid.UUID) ([]*entity.AIAsset, error)
	FindByID(ctx context.Context, id uuid.UUID) (*entity.AIAsset, error)
	Update(ctx context.Context, a *entity.AIAsset) (*entity.AIAsset, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type aiAssetRepository struct{ *baseRepository }

func NewAIAssetRepository(db *pgxpool.Pool) AIAssetRepository {
	return &aiAssetRepository{baseRepository: NewBaseRepository(db).(*baseRepository)}
}

func (r *aiAssetRepository) Create(ctx context.Context, a *entity.AIAsset) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	query := `
		INSERT INTO ai_assets (id, tenant_id, ai_agent_id, name, type, storage_key, mime_type, size, description)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING id, created_at, updated_at`
	if err := r.q(subCtx).QueryRow(subCtx, query,
		a.ID, a.TenantID, a.AIAgentID, a.Name, a.Type, a.StorageKey, a.MimeType, a.Size, a.Description,
	).Scan(&a.ID, &a.CreatedAt, &a.UpdatedAt); err != nil {
		return fmt.Errorf("failed to create ai asset: %w", err)
	}
	return nil
}

func (r *aiAssetRepository) ListByAgent(ctx context.Context, agentID uuid.UUID) ([]*entity.AIAsset, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out []*entity.AIAsset
	if err := pgxscan.Select(subCtx, r.q(subCtx), &out,
		`SELECT * FROM ai_assets WHERE ai_agent_id = $1 ORDER BY created_at DESC`, agentID); err != nil {
		return nil, fmt.Errorf("failed to list ai assets: %w", err)
	}
	return out, nil
}

func (r *aiAssetRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.AIAsset, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var a entity.AIAsset
	if err := pgxscan.Get(subCtx, r.q(subCtx), &a, `SELECT * FROM ai_assets WHERE id = $1`, id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai asset not found")
		}
		return nil, fmt.Errorf("failed to find ai asset: %w", err)
	}
	return &a, nil
}

func (r *aiAssetRepository) Update(ctx context.Context, a *entity.AIAsset) (*entity.AIAsset, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	var out entity.AIAsset
	err := r.q(subCtx).QueryRow(subCtx,
		`UPDATE ai_assets SET name=$1, type=$2, storage_key=$3, mime_type=$4, size=$5, description=$6
		 WHERE id=$7
		 RETURNING id, tenant_id, ai_agent_id, name, type, storage_key, mime_type, size, description, created_at, updated_at`,
		a.Name, a.Type, a.StorageKey, a.MimeType, a.Size, a.Description, a.ID,
	).Scan(&out.ID, &out.TenantID, &out.AIAgentID, &out.Name, &out.Type, &out.StorageKey,
		&out.MimeType, &out.Size, &out.Description, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("ai asset not found")
		}
		return nil, fmt.Errorf("failed to update ai asset: %w", err)
	}
	return &out, nil
}

func (r *aiAssetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tag, err := r.q(subCtx).Exec(subCtx, `DELETE FROM ai_assets WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete ai asset: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("ai asset not found")
	}
	return nil
}
