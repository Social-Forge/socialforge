package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// AIResourceService manages an AI agent's children: knowledge, playbooks and
// assets. Every operation is scoped to the tenant (RLS) and the parent agent,
// and verifies the parent agent exists in the tenant before touching a child.
type AIResourceService struct {
	agentRepo    repository.AIAgentRepository
	knowledgeRepo repository.AIKnowledgeRepository
	playbookRepo repository.AIPlaybookRepository
	assetRepo    repository.AIAssetRepository
	logger       *zap.Logger
}

func NewAIResourceService(
	agentRepo repository.AIAgentRepository,
	knowledgeRepo repository.AIKnowledgeRepository,
	playbookRepo repository.AIPlaybookRepository,
	assetRepo repository.AIAssetRepository,
	logger *zap.Logger,
) *AIResourceService {
	return &AIResourceService{
		agentRepo:    agentRepo,
		knowledgeRepo: knowledgeRepo,
		playbookRepo: playbookRepo,
		assetRepo:    assetRepo,
		logger:       logger,
	}
}

func (s *AIResourceService) prep(ctx context.Context, tenantID, agentID string) (context.Context, context.CancelFunc, uuid.UUID, uuid.UUID, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		cancel()
		return nil, nil, uuid.Nil, uuid.Nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	aid, err := uuid.Parse(agentID)
	if err != nil {
		cancel()
		return nil, nil, uuid.Nil, uuid.Nil, fmt.Errorf("invalid agent id: %w", err)
	}
	return repository.WithTenantID(subCtx, tid), cancel, tid, aid, nil
}

// ensureAgent confirms the parent agent exists within the tenant (RLS + row).
func (s *AIResourceService) ensureAgent(ctx context.Context, agentID uuid.UUID) error {
	_, err := s.agentRepo.FindByID(ctx, agentID)
	return err
}

func estimateTokens(content string) int {
	// Rough heuristic (~4 chars/token) until a real tokenizer is wired.
	n := len(content) / 4
	if n < 1 {
		n = 1
	}
	return n
}

// ============================ Knowledge ============================

func (s *AIResourceService) ListKnowledge(ctx context.Context, tenantID, agentID string) ([]*entity.AIKnowledge, error) {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	var out []*entity.AIKnowledge
	err = s.knowledgeRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.ensureAgent(txCtx, aid); err != nil {
			return err
		}
		out, err = s.knowledgeRepo.ListByAgent(txCtx, aid)
		return err
	})
	return out, err
}

func (s *AIResourceService) CreateKnowledge(ctx context.Context, tenantID, agentID string, req *dto.CreateAIKnowledgeRequest) (*entity.AIKnowledge, error) {
	tctx, cancel, tid, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	k := &entity.AIKnowledge{
		TenantID:   tid,
		AIAgentID:  aid,
		Title:      req.Title,
		Content:    req.Content,
		TokenCount: estimateTokens(req.Content),
	}
	err = s.knowledgeRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.ensureAgent(txCtx, aid); err != nil {
			return err
		}
		return s.knowledgeRepo.Create(txCtx, k)
	})
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (s *AIResourceService) UpdateKnowledge(ctx context.Context, tenantID, agentID, id string, req *dto.UpdateAIKnowledgeRequest) (*entity.AIKnowledge, error) {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	kid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid knowledge id: %w", err)
	}
	var out *entity.AIKnowledge
	err = s.knowledgeRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.knowledgeRepo.FindByID(txCtx, kid)
		if err != nil {
			return err
		}
		if existing.AIAgentID != aid {
			return fmt.Errorf("ai knowledge not found")
		}
		existing.Title = req.Title
		existing.Content = req.Content
		existing.TokenCount = estimateTokens(req.Content)
		out, err = s.knowledgeRepo.Update(txCtx, existing)
		return err
	})
	return out, err
}

func (s *AIResourceService) DeleteKnowledge(ctx context.Context, tenantID, agentID, id string) error {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return err
	}
	defer cancel()
	kid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid knowledge id: %w", err)
	}
	return s.knowledgeRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.knowledgeRepo.FindByID(txCtx, kid)
		if err != nil {
			return err
		}
		if existing.AIAgentID != aid {
			return fmt.Errorf("ai knowledge not found")
		}
		return s.knowledgeRepo.Delete(txCtx, kid)
	})
}

// ============================ Playbook ============================

func (s *AIResourceService) ListPlaybooks(ctx context.Context, tenantID, agentID string) ([]*entity.AIPlaybook, error) {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	var out []*entity.AIPlaybook
	err = s.playbookRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.ensureAgent(txCtx, aid); err != nil {
			return err
		}
		out, err = s.playbookRepo.ListByAgent(txCtx, aid)
		return err
	})
	return out, err
}

func (s *AIResourceService) CreatePlaybook(ctx context.Context, tenantID, agentID string, req *dto.CreateAIPlaybookRequest) (*entity.AIPlaybook, error) {
	tctx, cancel, tid, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	active := true
	if req.IsActive != nil {
		active = *req.IsActive
	}
	p := &entity.AIPlaybook{
		TenantID:    tid,
		AIAgentID:   aid,
		Name:        req.Name,
		Keywords:    entity.JSONStringSlice(req.Keywords),
		Instruction: req.Instruction,
		AssetIDs:    entity.JSONStringSlice(req.AssetIDs),
		Priority:    req.Priority,
		IsActive:    active,
	}
	err = s.playbookRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.ensureAgent(txCtx, aid); err != nil {
			return err
		}
		return s.playbookRepo.Create(txCtx, p)
	})
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *AIResourceService) UpdatePlaybook(ctx context.Context, tenantID, agentID, id string, req *dto.UpdateAIPlaybookRequest) (*entity.AIPlaybook, error) {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	pid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid playbook id: %w", err)
	}
	var out *entity.AIPlaybook
	err = s.playbookRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.playbookRepo.FindByID(txCtx, pid)
		if err != nil {
			return err
		}
		if existing.AIAgentID != aid {
			return fmt.Errorf("ai playbook not found")
		}
		existing.Name = req.Name
		existing.Keywords = entity.JSONStringSlice(req.Keywords)
		existing.Instruction = req.Instruction
		existing.AssetIDs = entity.JSONStringSlice(req.AssetIDs)
		existing.Priority = req.Priority
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		out, err = s.playbookRepo.Update(txCtx, existing)
		return err
	})
	return out, err
}

func (s *AIResourceService) DeletePlaybook(ctx context.Context, tenantID, agentID, id string) error {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return err
	}
	defer cancel()
	pid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid playbook id: %w", err)
	}
	return s.playbookRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.playbookRepo.FindByID(txCtx, pid)
		if err != nil {
			return err
		}
		if existing.AIAgentID != aid {
			return fmt.Errorf("ai playbook not found")
		}
		return s.playbookRepo.Delete(txCtx, pid)
	})
}

// ============================ Asset ============================

func (s *AIResourceService) ListAssets(ctx context.Context, tenantID, agentID string) ([]*entity.AIAsset, error) {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	var out []*entity.AIAsset
	err = s.assetRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.ensureAgent(txCtx, aid); err != nil {
			return err
		}
		out, err = s.assetRepo.ListByAgent(txCtx, aid)
		return err
	})
	return out, err
}

func (s *AIResourceService) CreateAsset(ctx context.Context, tenantID, agentID string, req *dto.CreateAIAssetRequest) (*entity.AIAsset, error) {
	tctx, cancel, tid, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	a := &entity.AIAsset{
		TenantID:    tid,
		AIAgentID:   aid,
		Name:        req.Name,
		Type:        req.Type,
		StorageKey:  req.StorageKey,
		MimeType:    entity.NewNullString(req.MimeType),
		Size:        entity.NewNullInt32(int32(req.Size)),
		Description: entity.NewNullString(req.Description),
	}
	err = s.assetRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if err := s.ensureAgent(txCtx, aid); err != nil {
			return err
		}
		return s.assetRepo.Create(txCtx, a)
	})
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AIResourceService) UpdateAsset(ctx context.Context, tenantID, agentID, id string, req *dto.UpdateAIAssetRequest) (*entity.AIAsset, error) {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return nil, err
	}
	defer cancel()
	asid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid asset id: %w", err)
	}
	var out *entity.AIAsset
	err = s.assetRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.assetRepo.FindByID(txCtx, asid)
		if err != nil {
			return err
		}
		if existing.AIAgentID != aid {
			return fmt.Errorf("ai asset not found")
		}
		existing.Name = req.Name
		existing.Type = req.Type
		existing.StorageKey = req.StorageKey
		existing.MimeType = entity.NewNullString(req.MimeType)
		existing.Size = entity.NewNullInt32(int32(req.Size))
		existing.Description = entity.NewNullString(req.Description)
		out, err = s.assetRepo.Update(txCtx, existing)
		return err
	})
	return out, err
}

func (s *AIResourceService) DeleteAsset(ctx context.Context, tenantID, agentID, id string) error {
	tctx, cancel, _, aid, err := s.prep(ctx, tenantID, agentID)
	if err != nil {
		return err
	}
	defer cancel()
	asid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid asset id: %w", err)
	}
	return s.assetRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.assetRepo.FindByID(txCtx, asid)
		if err != nil {
			return err
		}
		if existing.AIAgentID != aid {
			return fmt.Errorf("ai asset not found")
		}
		return s.assetRepo.Delete(txCtx, asid)
	})
}
