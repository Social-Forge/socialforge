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

type QuickReplyService struct {
	quickReplyRepo repository.QuickReplyRepository
	tenantRepo     repository.TenantRepository
	logger         *zap.Logger
}

func NewQuickReplyService(quickReplyRepo repository.QuickReplyRepository, tenantRepo repository.TenantRepository, logger *zap.Logger) *QuickReplyService {
	return &QuickReplyService{quickReplyRepo: quickReplyRepo, tenantRepo: tenantRepo, logger: logger}
}

func (s *QuickReplyService) tctx(ctx context.Context, tenantID string) (context.Context, uuid.UUID, error) {
	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return ctx, uuid.Nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	return repository.WithTenantID(ctx, tid), tid, nil
}

// validateContent enforces the per-type media rules: text needs a body; media
// types need media (hybrid = 1 media + body caption; media-only = up to 5).
func validateContent(contentType, body string, media []map[string]interface{}) error {
	if contentType == entity.QuickReplyTypeText {
		if body == "" {
			return fmt.Errorf("text quick reply requires a body")
		}
		return nil
	}
	if len(media) == 0 {
		return fmt.Errorf("%s quick reply requires at least one media file", contentType)
	}
	if body != "" && len(media) > 1 {
		return fmt.Errorf("hybrid (text + media) quick reply allows only 1 media file")
	}
	if len(media) > 5 {
		return fmt.Errorf("a quick reply allows at most 5 media files")
	}
	return nil
}

func (s *QuickReplyService) List(ctx context.Context, tenantID, search string) ([]*entity.QuickReply, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var out []*entity.QuickReply
	err = s.quickReplyRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		out, err = s.quickReplyRepo.List(txCtx, tid, search)
		return err
	})
	return out, err
}

func (s *QuickReplyService) Create(ctx context.Context, tenantID string, req *dto.CreateQuickReplyRequest) (*entity.QuickReply, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if err := validateContent(req.ContentType, req.Body, req.Media); err != nil {
		return nil, err
	}
	tctx, tid, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	tenant, err := s.tenantRepo.FindByID(subCtx, tid)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	qr := &entity.QuickReply{
		ID:          uuid.New(),
		TenantID:    tid,
		Shortcut:    req.Shortcut,
		ContentType: req.ContentType,
		Body:        entity.NewNullString(req.Body),
		Media:       entity.QuickReplyMedia(req.Media),
	}
	err = s.quickReplyRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		count, err := s.quickReplyRepo.Count(txCtx, tid)
		if err != nil {
			return err
		}
		if count >= tenant.MaxQuickReplies {
			return fmt.Errorf("quick reply quota reached (max %d) for plan %s", tenant.MaxQuickReplies, tenant.SubscriptionPlan)
		}
		return s.quickReplyRepo.Create(txCtx, qr)
	})
	if err != nil {
		return nil, err
	}
	return qr, nil
}

func (s *QuickReplyService) Update(ctx context.Context, tenantID, id string, req *dto.UpdateQuickReplyRequest) (*entity.QuickReply, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	if err := validateContent(req.ContentType, req.Body, req.Media); err != nil {
		return nil, err
	}
	qid, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid quick reply id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return nil, err
	}
	var out *entity.QuickReply
	err = s.quickReplyRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		existing, err := s.quickReplyRepo.FindByID(txCtx, qid)
		if err != nil {
			return err
		}
		existing.Shortcut = req.Shortcut
		existing.ContentType = req.ContentType
		existing.Body = entity.NewNullString(req.Body)
		existing.Media = entity.QuickReplyMedia(req.Media)
		out, err = s.quickReplyRepo.Update(txCtx, existing)
		return err
	})
	return out, err
}

func (s *QuickReplyService) Delete(ctx context.Context, tenantID, id string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()
	qid, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid quick reply id: %w", err)
	}
	tctx, _, err := s.tctx(subCtx, tenantID)
	if err != nil {
		return err
	}
	return s.quickReplyRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		return s.quickReplyRepo.Delete(txCtx, qid)
	})
}
