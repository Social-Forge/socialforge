package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/centrifugo"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type LabelService struct {
	labelRepo  repository.LabelRepository
	centrifugo *centrifugo.CentrifugoClient
	logger     *zap.Logger
}

func NewLabelService(labelRepo repository.LabelRepository, centrifugoClient *centrifugo.CentrifugoClient, logger *zap.Logger) *LabelService {
	return &LabelService{labelRepo: labelRepo, centrifugo: centrifugoClient, logger: logger}
}

func (s *LabelService) tctx(ctx context.Context, tenantID string) (context.Context, uuid.UUID, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return ctx, uuid.Nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	return repository.WithTenantID(ctx, tid), tid, nil
}

func (s *LabelService) List(ctx context.Context, tenantID string) ([]*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var labels []*entity.Label
	err = s.labelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		labels, err = s.labelRepo.List(txCtx, tid)
		return err
	})
	return labels, err
}

func (s *LabelService) Create(ctx context.Context, tenantID string, req *dto.CreateLabelRequest) (*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	label := &entity.Label{ID: uuid.New(), TenantID: tid, Name: req.Name, Color: req.Color}
	err = s.labelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.labelRepo.Create(txCtx, label)
	})
	if err != nil {
		return nil, err
	}
	return label, nil
}

func (s *LabelService) Update(ctx context.Context, tenantID, id string, req *dto.UpdateLabelRequest) (*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	labelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid label id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var out *entity.Label
	err = s.labelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.labelRepo.FindByID(txCtx, labelID)
		if err != nil {
			return err
		}
		existing.Name = req.Name
		if req.Color != "" {
			existing.Color = req.Color
		}
		out, err = s.labelRepo.Update(txCtx, existing)
		return err
	})
	return out, err
}

func (s *LabelService) Delete(ctx context.Context, tenantID, id string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	labelID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid label id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return err
	}
	return s.labelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.labelRepo.Delete(txCtx, labelID)
	})
}

func (s *LabelService) ListForConversation(ctx context.Context, tenantID, conversationID string) ([]*entity.Label, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	cid, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, fmt.Errorf("invalid conversation id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var labels []*entity.Label
	err = s.labelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		labels, err = s.labelRepo.ListForConversation(txCtx, cid)
		return err
	})
	return labels, err
}

func (s *LabelService) Attach(ctx context.Context, tenantID, conversationID, labelID string) error {
	return s.toggle(ctx, tenantID, conversationID, labelID, true)
}

func (s *LabelService) Detach(ctx context.Context, tenantID, conversationID, labelID string) error {
	return s.toggle(ctx, tenantID, conversationID, labelID, false)
}

func (s *LabelService) toggle(ctx context.Context, tenantID, conversationID, labelID string, attach bool) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	cid, err := uuid.Parse(conversationID)
	if err != nil {
		return fmt.Errorf("invalid conversation id: %w", err)
	}
	lid, err := uuid.Parse(labelID)
	if err != nil {
		return fmt.Errorf("invalid label id: %w", err)
	}
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return err
	}
	err = s.labelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		if attach {
			return s.labelRepo.Attach(txCtx, &entity.ConversationLabel{TenantID: tid, ConversationID: cid, LabelID: lid})
		}
		return s.labelRepo.Detach(txCtx, cid, lid)
	})
	if err != nil {
		return err
	}
	if s.centrifugo != nil {
		event := "label_attached"
		if !attach {
			event = "label_detached"
		}
		_ = s.centrifugo.BroadcastConversationUpdate(subCtx, cid.String(), map[string]interface{}{"type": event, "label_id": lid.String()})
	}
	return nil
}
