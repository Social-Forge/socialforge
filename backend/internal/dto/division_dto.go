package dto

import (
	"github/socialforge/internal/entity"

	"github.com/google/uuid"
)

type CreateDivisionRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
	RoutingType string `json:"routing_type,omitempty" validate:"omitempty,oneof=equal percentage priority"`
}

type UpdateDivisionRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Description string `json:"description,omitempty" validate:"omitempty,max=1000"`
	RoutingType string `json:"routing_type,omitempty" validate:"omitempty,oneof=equal percentage priority"`
	IsActive    *bool  `json:"is_active,omitempty"`
}

// AddDivisionMemberRequest assigns an existing tenant member (identified by its
// user_tenant membership id) to a division.
type AddDivisionMemberRequest struct {
	UserTenantID uuid.UUID `json:"user_tenant_id" validate:"required"`
}

type DivisionResponse struct {
	*entity.Division
	MemberCount int `json:"member_count,omitempty"`
}
