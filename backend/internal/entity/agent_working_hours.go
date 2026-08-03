package entity

import (
	"time"

	"github.com/google/uuid"
)

// AgentWorkingHours is one working window for an agent on a given weekday.
// day_of_week: 0=Sunday .. 6=Saturday (matches Postgres EXTRACT(DOW)).
type AgentWorkingHours struct {
	ID        uuid.UUID `json:"id" db:"id"`
	TenantID  uuid.UUID `json:"tenant_id" db:"tenant_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	DayOfWeek int       `json:"day_of_week" db:"day_of_week" validate:"min=0,max=6"`
	StartTime string    `json:"start_time" db:"start_time" validate:"required"`
	EndTime   string    `json:"end_time" db:"end_time" validate:"required"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
