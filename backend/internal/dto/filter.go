package dto

import "time"

type FilterQueryRequest struct {
	Search         string                 `json:"search,omitempty" validate:"omitempty"`
	SortBy         string                 `json:"sort_by,omitempty" validate:"omitempty"`
	OrderBy        string                 `json:"order_by,omitempty" validate:"omitempty"`
	Page           int                    `json:"page,omitempty" validate:"required,min=1"`
	Limit          int                    `json:"limit,omitempty" validate:"required,min=1,max=100"`
	Status         string                 `json:"status,omitempty" validate:"omitempty"`
	IncludeDeleted bool                   `json:"include_deleted,omitempty" validate:"omitempty"`
	IsActive       bool                   `json:"is_active,omitempty" validate:"omitempty"`
	IsVerified     bool                   `json:"is_verified,omitempty" validate:"omitempty"`
	TenantID       string                 `json:"tenant_id,omitempty" validate:"omitempty"`
	UserID         string                 `json:"user_id,omitempty" validate:"omitempty"`
	DivisionID     string                 `json:"division_id,omitempty" validate:"omitempty"`
	StartDate      time.Time              `json:"start_date,omitempty" validate:"omitempty"`
	EndDate        time.Time              `json:"end_date,omitempty" validate:"omitempty"`
	Extra          map[string]interface{} `json:"extra,omitempty" validate:"dive,key,required,value"`
}
