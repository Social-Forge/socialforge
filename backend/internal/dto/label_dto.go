package dto

type CreateLabelRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Color string `json:"color,omitempty" validate:"omitempty,hexcolor|max=32"`
}

type UpdateLabelRequest struct {
	Name  string `json:"name" validate:"required,min=1,max=255"`
	Color string `json:"color,omitempty" validate:"omitempty,hexcolor|max=32"`
}

type AttachLabelRequest struct {
	LabelID string `json:"label_id" validate:"required,uuid4"`
}
