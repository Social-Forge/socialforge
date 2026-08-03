package dto

type WorkingHourItem struct {
	DayOfWeek int    `json:"day_of_week" validate:"min=0,max=6"`
	StartTime string `json:"start_time" validate:"required"` // HH:MM
	EndTime   string `json:"end_time" validate:"required"`   // HH:MM
	IsActive  bool   `json:"is_active"`
}

type SetWorkingHoursRequest struct {
	Hours []WorkingHourItem `json:"hours" validate:"dive"`
}
