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

type AIAgentService struct {
	repo   repository.AIAgentRepository
	logger *zap.Logger
}

func NewAIAgentService(repo repository.AIAgentRepository, logger *zap.Logger) *AIAgentService {
	return &AIAgentService{repo: repo, logger: logger}
}

func (s *AIAgentService) tctx(ctx context.Context, tenantID string) (context.Context, uuid.UUID, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return ctx, uuid.Nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	return repository.WithTenantID(ctx, tid), tid, nil
}

// defaultModel picks a sensible model when the request omits one.
func defaultModel(provider, model string) string {
	if model != "" {
		return model
	}
	switch provider {
	case "openai":
		return "gpt-4o-mini"
	case "google":
		return "gemini-1.5-flash"
	default:
		return "claude-3-5-sonnet-20241022"
	}
}

func (s *AIAgentService) List(ctx context.Context, tenantID string) ([]*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var out []*entity.AIAgent
	err = s.repo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		out, err = s.repo.List(txCtx, tid)
		return err
	})
	return out, err
}

func (s *AIAgentService) Get(ctx context.Context, tenantID, id string) (*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid agent id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var out *entity.AIAgent
	err = s.repo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		out, err = s.repo.FindByID(txCtx, aid)
		return err
	})
	return out, err
}

func (s *AIAgentService) Create(ctx context.Context, tenantID string, req *dto.CreateAIAgentRequest) (*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	temp := req.Temperature
	if temp == 0 {
		temp = 0.7
	}
	maxTok := req.MaxTokens
	if maxTok == 0 {
		maxTok = 1024
	}
	autoReply := true
	if req.AutoReply != nil {
		autoReply = *req.AutoReply
	}
	agent := &entity.AIAgent{
		ID:               uuid.New(),
		TenantID:         tid,
		Name:             req.Name,
		Provider:         req.Provider,
		Model:            defaultModel(req.Provider, req.Model),
		SystemPrompt:     req.SystemPrompt,
		Persona:          persona(req.Persona),
		Safety:           safety(req.Safety),
		Guardrails:       guardrails(req.Guardrails),
		Temperature:      temp,
		MaxTokens:        maxTok,
		AutoReplyEnabled: autoReply,
		IsActive:         true,
	}
	err = s.repo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.repo.Create(txCtx, agent)
	})
	if err != nil {
		return nil, err
	}
	return agent, nil
}

func (s *AIAgentService) Update(ctx context.Context, tenantID, id string, req *dto.UpdateAIAgentRequest) (*entity.AIAgent, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	aid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid agent id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var out *entity.AIAgent
	err = s.repo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.repo.FindByID(txCtx, aid)
		if err != nil {
			return err
		}
		existing.Name = req.Name
		existing.Provider = req.Provider
		existing.Model = defaultModel(req.Provider, req.Model)
		existing.SystemPrompt = req.SystemPrompt
		if req.Persona != nil {
			existing.Persona = persona(req.Persona)
		}
		if req.Safety != nil {
			existing.Safety = safety(req.Safety)
		}
		if req.Guardrails != nil {
			existing.Guardrails = guardrails(req.Guardrails)
		}
		if req.Temperature > 0 {
			existing.Temperature = req.Temperature
		}
		if req.MaxTokens > 0 {
			existing.MaxTokens = req.MaxTokens
		}
		if req.AutoReply != nil {
			existing.AutoReplyEnabled = *req.AutoReply
		}
		if req.IsActive != nil {
			existing.IsActive = *req.IsActive
		}
		out, err = s.repo.Update(txCtx, existing)
		return err
	})
	return out, err
}

func (s *AIAgentService) Delete(ctx context.Context, tenantID, id string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	aid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid agent id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return err
	}
	return s.repo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.repo.Delete(txCtx, aid)
	})
}

func persona(m map[string]interface{}) *entity.AiPersonaConfig {
	if m == nil {
		return nil
	}
	p := entity.AiPersonaConfig(m)
	return &p
}
func safety(m map[string]interface{}) *entity.AiSafetyConfig {
	if m == nil {
		return nil
	}
	sf := entity.AiSafetyConfig(m)
	return &sf
}
func guardrails(m map[string]interface{}) *entity.AiSafetyGuardrailsConfig {
	if m == nil {
		return nil
	}
	g := entity.AiSafetyGuardrailsConfig(m)
	return &g
}
