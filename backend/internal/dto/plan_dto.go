package dto

type CreatePlanRequest struct {
	Code     string                 `json:"code" validate:"required,min=1,max=50"`
	Name     string                 `json:"name" validate:"required,min=1,max=100"`
	Price    int                    `json:"price" validate:"min=0"`
	Currency string                 `json:"currency,omitempty"`
	Interval string                 `json:"interval,omitempty" validate:"omitempty,oneof=monthly yearly"`
	Features map[string]interface{} `json:"features,omitempty"`
	IsActive *bool                  `json:"is_active,omitempty"`
	Sort     int                    `json:"sort,omitempty"`
}

type UpdatePlanRequest struct {
	Code     string                 `json:"code" validate:"required,min=1,max=50"`
	Name     string                 `json:"name" validate:"required,min=1,max=100"`
	Price    int                    `json:"price" validate:"min=0"`
	Currency string                 `json:"currency,omitempty"`
	Interval string                 `json:"interval,omitempty" validate:"omitempty,oneof=monthly yearly"`
	Features map[string]interface{} `json:"features,omitempty"`
	IsActive *bool                  `json:"is_active,omitempty"`
	Sort     int                    `json:"sort,omitempty"`
}
