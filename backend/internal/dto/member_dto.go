package dto

// CreateMemberRequest is used by a tenant owner to create a supervisor or agent
// account within their tenant (user + user_tenants membership).
type CreateMemberRequest struct {
	FullName string `json:"full_name" validate:"required,min=2,max=255"`
	Email    string `json:"email" validate:"required,email,max=255"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=supervisor agent"`
}

// UpdateMemberRequest updates a member's role and/or active state.
type UpdateMemberRequest struct {
	Role     string `json:"role" validate:"required,oneof=supervisor agent"`
	IsActive *bool  `json:"is_active,omitempty"`
}
