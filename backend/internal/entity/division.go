package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Division represents a team/group within a tenant
type Division struct {
	ID            uuid.UUID      `json:"id" db:"id"`
	TenantID      uuid.UUID      `json:"tenant_id" db:"tenant_id" validate:"required"`
	Name          string         `json:"name" db:"name" validate:"required,max=255"`
	Slug          string         `json:"slug" db:"slug" validate:"required,max=100"`
	Description   NullString     `json:"description,omitempty" db:"description"`
	RoutingType   string         `json:"routing_type" db:"routing_type" validate:"required,oneof=equal percentage priority"`
	RoutingConfig *RoutingConfig `json:"routing_config,omitempty" db:"routing_config"`
	IsActive      bool           `json:"is_active" db:"is_active"`
	LinkURL       NullString     `json:"link_url,omitempty" db:"link_url"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
	DeletedAt     NullTime       `json:"deleted_at,omitempty" db:"deleted_at"`
}

// Routing types
const (
	RoutingEqual      = "equal"
	RoutingPercentage = "percentage"
	RoutingPriority   = "priority"
)

// RoutingConfig holds the routing configuration
type RoutingConfig map[string]interface{}

// Value implements the driver.Valuer interface for RoutingConfig
func (rc RoutingConfig) Value() (driver.Value, error) {
	if rc == nil {
		return nil, nil
	}
	return json.Marshal(rc)
}

// Scan implements the sql.Scanner interface for RoutingConfig
func (rc *RoutingConfig) Scan(value interface{}) error {
	if value == nil {
		*rc = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, rc)
}

// TableName returns the table name for Division
func (Division) TableName() string {
	return "divisions"
}

func (d *Division) IsDeleted() bool {
	return d.DeletedAt.Valid
}

// GenerateLinkURL generates a public link URL for this division
func (d *Division) GenerateLinkURL(baseURL string) string {
	return baseURL + "/c/" + d.Slug
}
