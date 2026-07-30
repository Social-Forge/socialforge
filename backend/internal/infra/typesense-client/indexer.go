package typesenseclient

import (
	"context"
	"fmt"
	"github/socialforge/internal/entity"

	"github.com/google/uuid"
	"github.com/typesense/typesense-go/v4/typesense/api"
	"go.uber.org/zap"
)

const (
	CollectionUsers     = "users"
	CollectionTenants   = "tenants"
	CollectionDivisions = "divisions"
)

type SearchIndexService struct {
	client *TypesenseClient
	logger *zap.Logger
}

func NewSearchIndexService(client *TypesenseClient, logger *zap.Logger) *SearchIndexService {
	return &SearchIndexService{client: client, logger: logger}
}

func (s *SearchIndexService) Enabled() bool {
	return s != nil && s.client != nil && s.client.IsUp()
}

func (s *SearchIndexService) Bootstrap(ctx context.Context) error {
	if !s.Enabled() {
		return nil
	}

	schemas := []*api.CollectionSchema{
		{
			Name: CollectionUsers,
			Fields: []api.Field{
				{Name: "id", Type: "string"},
				{Name: "email", Type: "string"},
				{Name: "full_name", Type: "string"},
				{Name: "phone", Type: "string"},
				{Name: "avatar_url", Type: "string"},
				{Name: "status", Type: "string", Facet: boolPtr(true)},
				{Name: "is_active", Type: "bool", Facet: boolPtr(true)},
				{Name: "is_verified", Type: "bool", Facet: boolPtr(true)},
				{Name: "tenant_id", Type: "string", Facet: boolPtr(true)},
				{Name: "role_id", Type: "string", Facet: boolPtr(true)},
				{Name: "tenant_names", Type: "string[]", Facet: boolPtr(true)},
				{Name: "division_names", Type: "string[]", Facet: boolPtr(true)},
				{Name: "email_verified_at", Type: "int64", Facet: boolPtr(true)},
				{Name: "last_login_at", Type: "int64", Facet: boolPtr(true)},
				{Name: "created_at", Type: "int64", Sort: boolPtr(true)},
				{Name: "updated_at", Type: "int64", Sort: boolPtr(true)},
			},
			DefaultSortingField: stringPtr("created_at"),
		},
		{
			Name: CollectionTenants,
			Fields: []api.Field{
				{Name: "id", Type: "string"},
				{Name: "name", Type: "string"},
				{Name: "slug", Type: "string"},
				{Name: "max_divisions", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_agents", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_quick_replies", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_waha_whatsapp", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_meta_whatsapp", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_meta_messenger", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_instagram", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_telegram", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_webchat", Type: "int64", Facet: boolPtr(true)},
				{Name: "max_linkchat", Type: "int64", Facet: boolPtr(true)},
				{Name: "ai_credits", Type: "int64", Sort: boolPtr(true)},
				{Name: "subscription_plan", Type: "string", Facet: boolPtr(true)},
				{Name: "subscription_status", Type: "string", Facet: boolPtr(true)},
				{Name: "is_active", Type: "bool", Facet: boolPtr(true)},
				{Name: "created_at", Type: "int64", Sort: boolPtr(true)},
				{Name: "updated_at", Type: "int64", Sort: boolPtr(true)},
			},
			DefaultSortingField: stringPtr("created_at"),
		},
		{
			Name: CollectionDivisions,
			Fields: []api.Field{
				{Name: "id", Type: "string"},
				{Name: "tenant_id", Type: "string", Facet: boolPtr(true)},
				{Name: "name", Type: "string"},
				{Name: "slug", Type: "string"},
				{Name: "description", Type: "string"},
				{Name: "routing_type", Type: "string", Facet: boolPtr(true)},
				{Name: "is_active", Type: "bool", Facet: boolPtr(true)},
				{Name: "link_url", Type: "string"},
				{Name: "created_at", Type: "int64", Sort: boolPtr(true)},
				{Name: "updated_at", Type: "int64", Sort: boolPtr(true)},
			},
			DefaultSortingField: stringPtr("created_at"),
		},
	}

	for _, schema := range schemas {
		if err := s.ensureCollection(ctx, schema); err != nil {
			return err
		}
	}

	return nil
}

func (s *SearchIndexService) UpsertUser(ctx context.Context, user *entity.User) error {
	if !s.Enabled() || user == nil {
		return nil
	}

	var tenantNames []string
	var divisionNames []string
	var tenantID string
	var roleID string

	if len(user.UserTenants) > 0 {
		for _, ut := range user.UserTenants {
			if ut.IsActive && ut.Tenant != nil {
				tenantNames = append(tenantNames, ut.Tenant.Name)
				if tenantID == "" {
					tenantID = ut.TenantID.String()
				}
			}

			if ut.Role != nil && roleID == "" {
				roleID = ut.RoleID.String()
			}

			if len(ut.DivisionMembers) > 0 {
				for _, dm := range ut.DivisionMembers {
					if dm.Division != nil {
						divisionNames = append(divisionNames, dm.Division.Name)
					}
				}
			}
		}
	}

	if tenantID == "" && len(user.UserTenants) > 0 {
		tenantID = user.UserTenants[0].TenantID.String()
		if user.UserTenants[0].Tenant != nil {
			tenantNames = append(tenantNames, user.UserTenants[0].Tenant.Name)
		}
	}

	doc := userSearchDocument{
		ID:            user.ID.String(),
		TenantID:      tenantID,
		RoleID:        roleID,
		Email:         user.Email,
		FullName:      user.FullName,
		Phone:         user.Phone.String,
		AvatarURL:     user.AvatarURL.String,
		Status:        user.Status,
		IsActive:      user.Status == entity.StatusActive,
		IsVerified:    user.EmailVerifiedAt.Valid,
		TenantNames:   tenantNames,
		DivisionNames: divisionNames,
		CreatedAt:     user.CreatedAt.Unix(),
		UpdatedAt:     user.UpdatedAt.Unix(),
	}

	return s.upsert(ctx, CollectionUsers, doc)
}
func (s *SearchIndexService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return s.deleteByID(ctx, CollectionUsers, id.String())
}

func (s *SearchIndexService) UpsertTenant(ctx context.Context, tenant *entity.Tenant) error {
	if !s.Enabled() || tenant == nil {
		return nil
	}

	doc := tenantSearchDocument{
		ID:                 tenant.ID.String(),
		Name:               tenant.Name,
		Slug:               tenant.Slug,
		SubscriptionPlan:   tenant.SubscriptionPlan,
		SubscriptionStatus: tenant.SubscriptionStatus,
		IsActive:           tenant.IsActive,
		AiCredits:          int64(tenant.AiCredits),
		CreatedAt:          tenant.CreatedAt.Unix(),
		UpdatedAt:          tenant.UpdatedAt.Unix(),
	}

	return s.upsert(ctx, CollectionTenants, doc)
}

func (s *SearchIndexService) DeleteTenant(ctx context.Context, id uuid.UUID) error {
	return s.deleteByID(ctx, CollectionTenants, id.String())
}

func (s *SearchIndexService) UpsertDivision(ctx context.Context, division *entity.Division) error {
	if !s.Enabled() || division == nil {
		return nil
	}

	doc := divisionSearchDocument{
		ID:          division.ID.String(),
		TenantID:    division.TenantID.String(),
		Name:        division.Name,
		Slug:        division.Slug,
		RoutingType: division.RoutingType,
		IsActive:    division.IsActive,
		LinkURL:     division.LinkURL.String,
		CreatedAt:   division.CreatedAt.Unix(),
		UpdatedAt:   division.UpdatedAt.Unix(),
	}

	return s.upsert(ctx, CollectionDivisions, doc)
}

func (s *SearchIndexService) DeleteDivision(ctx context.Context, id uuid.UUID) error {
	return s.deleteByID(ctx, CollectionDivisions, id.String())
}

func (s *SearchIndexService) ensureCollection(ctx context.Context, schema *api.CollectionSchema) error {
	_, err := s.client.Client().Collection(schema.Name).Retrieve(ctx)
	if err == nil {
		return nil
	}

	_, createErr := s.client.Client().Collections().Create(ctx, schema)
	if createErr != nil {
		return fmt.Errorf("failed to create typesense collection %s: %w", schema.Name, createErr)
	}

	s.logger.Info("Typesense collection ready", zap.String("collection", schema.Name))
	return nil
}

func (s *SearchIndexService) upsert(ctx context.Context, collection string, doc any) error {
	_, err := s.client.Client().Collection(collection).Documents().Upsert(ctx, doc, &api.DocumentIndexParameters{})
	if err != nil {
		return fmt.Errorf("failed to upsert %s document: %w", collection, err)
	}
	return nil
}

func (s *SearchIndexService) deleteByID(ctx context.Context, collection, id string) error {
	filter := fmt.Sprintf("id:=%s", id)
	ignore := true
	_, err := s.client.Client().Collection(collection).Documents().Delete(ctx, &api.DeleteDocumentsParams{
		FilterBy:       &filter,
		IgnoreNotFound: &ignore,
	})
	if err != nil {
		return fmt.Errorf("failed to delete %s document: %w", collection, err)
	}
	return nil
}

type userSearchDocument struct {
	ID            string   `json:"id"`
	TenantID      string   `json:"tenant_id"`
	RoleID        string   `json:"role_id"`
	Email         string   `json:"email"`
	FullName      string   `json:"full_name"`
	Phone         string   `json:"phone,omitempty"`
	AvatarURL     string   `json:"avatar_url,omitempty"`
	Status        string   `json:"status"`
	IsActive      bool     `json:"is_active"`
	IsVerified    bool     `json:"is_verified"`
	TenantNames   []string `json:"tenant_names,omitempty"`
	DivisionNames []string `json:"division_names,omitempty"`
	CreatedAt     int64    `json:"created_at"`
	UpdatedAt     int64    `json:"updated_at"`
}
type tenantSearchDocument struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Slug               string `json:"slug"`
	SubscriptionPlan   string `json:"subscription_plan"`
	SubscriptionStatus string `json:"subscription_status"`
	IsActive           bool   `json:"is_active"`
	AiCredits          int64  `json:"ai_credits"`
	CreatedAt          int64  `json:"created_at"`
	UpdatedAt          int64  `json:"updated_at"`
}

type divisionSearchDocument struct {
	ID          string `json:"id"`
	TenantID    string `json:"tenant_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	RoutingType string `json:"routing_type"`
	IsActive    bool   `json:"is_active"`
	LinkURL     string `json:"link_url,omitempty"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

func stringPtr(v string) *string { return &v }

func boolPtr(v bool) *bool { return &v }
