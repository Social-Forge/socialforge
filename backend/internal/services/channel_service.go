package services

import (
	"context"
	"fmt"
	"github/socialforge/internal/dto"
	"github/socialforge/internal/entity"
	"github/socialforge/internal/infra/channels"
	"github/socialforge/internal/infra/contextpool"
	"github/socialforge/internal/infra/repository"
	"github/socialforge/internal/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ChannelService struct {
	channelRepo  repository.ChannelRepository
	divisionRepo repository.DivisionRepository
	tenantRepo   repository.TenantRepository
	connectors   map[string]channels.Connector
	appURL       string // public base URL for provider webhooks (Telegram, Meta)
	internalURL  string // in-network base URL for WAHA -> backend webhooks
	logger       *zap.Logger
}

func NewChannelService(
	channelRepo repository.ChannelRepository,
	divisionRepo repository.DivisionRepository,
	tenantRepo repository.TenantRepository,
	connectors map[string]channels.Connector,
	appURL, internalURL string,
	logger *zap.Logger,
) *ChannelService {
	if internalURL == "" {
		internalURL = "http://backend:8080"
	}
	return &ChannelService{
		channelRepo:  channelRepo,
		divisionRepo: divisionRepo,
		tenantRepo:   tenantRepo,
		connectors:   connectors,
		appURL:       appURL,
		internalURL:  internalURL,
		logger:       logger,
	}
}

// Connect registers the channel with its provider (webhook + session) and
// updates its status. Returns provider info (e.g. WAHA QR url).
func (s *ChannelService) Connect(ctx context.Context, tenantID, id string) (map[string]interface{}, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 25*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	channelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	var channel *entity.Channel
	err = s.channelRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		channel, err = s.channelRepo.FindByID(txCtx, channelID)
		return err
	})
	if err != nil {
		return nil, err
	}

	connector := s.connectors[channel.Type]
	if connector == nil {
		return nil, fmt.Errorf("channel type %s does not support connect yet", channel.Type)
	}

	provider := providerForChannelType(channel.Type)
	secret := ""
	if channel.WebhookSecret.Valid {
		secret = channel.WebhookSecret.String
	}
	// WAHA calls back over the docker network; Telegram/Meta need a public URL.
	base := s.appURL
	if channel.Type == entity.ChannelTypeWhatsAppWaha {
		base = s.internalURL
	}
	webhookURL := fmt.Sprintf("%s/api/webhooks/%s/%s", strings.TrimRight(base, "/"), provider, channel.ID)

	info, status, connErr := connector.Connect(subCtx, channel, webhookURL, secret)
	// Persist the resulting status regardless of connect outcome.
	_ = s.channelRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		return s.channelRepo.UpdateStatus(txCtx, channelID, status)
	})
	if connErr != nil {
		return nil, connErr
	}
	if info == nil {
		info = map[string]interface{}{}
	}
	info["status"] = status
	info["webhook_url"] = webhookURL
	return info, nil
}

// channelQuota maps a channel type to its per-tenant plan limit.
func channelQuota(t *entity.Tenant, channelType string) (int, bool) {
	switch channelType {
	case entity.ChannelTypeWhatsAppWaha:
		return t.MaxWahaWhatsApp, true
	case entity.ChannelTypeWhatsAppMeta:
		return t.MaxMetaWhatsApp, true
	case entity.ChannelTypeMessenger:
		return t.MaxMetaMessenger, true
	case entity.ChannelTypeInstagram:
		return t.MaxInstagram, true
	case entity.ChannelTypeTelegram:
		return t.MaxTelegram, true
	default:
		return 0, false
	}
}

func (s *ChannelService) Create(ctx context.Context, tenantID string, req *dto.CreateChannelRequest) (*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	divisionID, err := uuid.Parse(req.DivisionID)
	if err != nil {
		return nil, fmt.Errorf("invalid division id: %w", err)
	}

	tenant, err := s.tenantRepo.FindByID(subCtx, tid)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}
	limit, ok := channelQuota(tenant, req.Type)
	if !ok {
		return nil, fmt.Errorf("unsupported channel type: %s", req.Type)
	}

	credentials := entity.CredentialsConfig(req.Credentials)
	if credentials == nil {
		credentials = entity.CredentialsConfig{}
	}
	settings := entity.ChannelSetting{}

	channel := &entity.Channel{
		ID:            uuid.New(),
		TenantID:      tid,
		DivisionID:    divisionID,
		Type:          req.Type,
		Name:          req.Name,
		Status:        entity.ChannelStatusDisconnected,
		WebhookSecret: entity.NewNullString(utils.GenerateSecureToken(24)),
		Credentials:   &credentials,
		Settings:      &settings,
	}
	if req.Type == entity.ChannelTypeWhatsAppWaha {
		engine := req.WahaEngine
		if engine == "" {
			engine = "GOWS"
		}
		channel.WahaEngine = entity.NewNullString(engine)
		channel.WahaSessionName = entity.NewNullString(fmt.Sprintf("sf_%s", strings.ToLower(utils.GenerateSecureToken(10))))
	}

	tctx := repository.WithTenantID(subCtx, tid)
	err = s.channelRepo.RunInTenantTx(tctx, func(txCtx context.Context) error {
		// Division must exist within this tenant.
		if _, err := s.divisionRepo.FindByID(txCtx, divisionID); err != nil {
			return fmt.Errorf("division not found in this tenant: %w", err)
		}
		// Plan quota per channel type.
		count, err := s.channelRepo.CountByType(txCtx, tid, req.Type)
		if err != nil {
			return err
		}
		if count >= limit {
			return fmt.Errorf("channel quota reached (max %d %s) for plan %s", limit, req.Type, tenant.SubscriptionPlan)
		}
		return s.channelRepo.Create(txCtx, channel)
	})
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (s *ChannelService) List(ctx context.Context, tenantID, chType, divisionID string) ([]*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	var divPtr *uuid.UUID
	if divisionID != "" {
		if d, err := uuid.Parse(divisionID); err == nil {
			divPtr = &d
		}
	}

	var channels []*entity.Channel
	err = s.channelRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		channels, err = s.channelRepo.List(txCtx, tid, chType, divPtr)
		return err
	})
	if err != nil {
		return nil, err
	}
	return channels, nil
}

func (s *ChannelService) Get(ctx context.Context, tenantID, id string) (*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	channelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	var channel *entity.Channel
	err = s.channelRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		channel, err = s.channelRepo.FindByID(txCtx, channelID)
		return err
	})
	if err != nil {
		return nil, err
	}
	return channel, nil
}

func (s *ChannelService) Update(ctx context.Context, tenantID, id string, req *dto.UpdateChannelRequest) (*entity.Channel, error) {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return nil, fmt.Errorf("invalid tenant id: %w", err)
	}
	channelID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	var updated *entity.Channel
	err = s.channelRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		existing, err := s.channelRepo.FindByID(txCtx, channelID)
		if err != nil {
			return err
		}
		existing.Name = req.Name
		if req.DivisionID != "" {
			did, err := uuid.Parse(req.DivisionID)
			if err != nil {
				return fmt.Errorf("invalid division id: %w", err)
			}
			if _, err := s.divisionRepo.FindByID(txCtx, did); err != nil {
				return fmt.Errorf("division not found in this tenant: %w", err)
			}
			existing.DivisionID = did
		}
		if req.AIAgentID != nil && *req.AIAgentID != "" {
			aid, err := uuid.Parse(*req.AIAgentID)
			if err != nil {
				return fmt.Errorf("invalid ai_agent id: %w", err)
			}
			existing.AIAgentID = &aid
		}
		updated, err = s.channelRepo.Update(txCtx, existing)
		return err
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

func (s *ChannelService) Delete(ctx context.Context, tenantID, id string) error {
	subCtx, cancel := contextpool.WithTimeoutIfNone(ctx, 15*time.Second)
	defer cancel()

	tid, err := uuid.Parse(tenantID)
	if err != nil {
		return fmt.Errorf("invalid tenant id: %w", err)
	}
	channelID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid channel id: %w", err)
	}
	return s.channelRepo.RunInTenantTx(repository.WithTenantID(subCtx, tid), func(txCtx context.Context) error {
		return s.channelRepo.Delete(txCtx, channelID)
	})
}
