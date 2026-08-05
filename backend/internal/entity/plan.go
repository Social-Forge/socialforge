package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// JSONMap is a generic JSONB object used by billing rows (plan features, invoice
// purpose, addon meta, webhook payload).
type JSONMap map[string]interface{}

func (m JSONMap) Value() (driver.Value, error) {
	if m == nil {
		return "{}", nil
	}
	return json.Marshal(m)
}

func (m *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*m = nil
		return nil
	}
	var b []byte
	switch v := value.(type) {
	case []byte:
		b = v
	case string:
		b = []byte(v)
	default:
		return nil
	}
	return json.Unmarshal(b, m)
}

// Int returns an integer feature value (features are numeric entitlements).
func (m JSONMap) Int(key string) (int, bool) {
	if m == nil {
		return 0, false
	}
	switch v := m[key].(type) {
	case float64:
		return int(v), true
	case int:
		return v, true
	case int64:
		return int(v), true
	}
	return 0, false
}

// Plan billing intervals.
const (
	PlanIntervalMonthly = "monthly"
	PlanIntervalYearly  = "yearly"
)

// Plan is a catalog entry: price + entitlements (features). Plan.Code maps to the
// tenant subscription_plan enum (free/starter/pro/enterprise).
type Plan struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Code      string    `json:"code" db:"code"`
	Name      string    `json:"name" db:"name"`
	Price     int       `json:"price" db:"price"`
	Currency  string    `json:"currency" db:"currency"`
	Interval  string    `json:"interval" db:"interval"`
	Features  JSONMap   `json:"features" db:"features"`
	IsActive  bool      `json:"is_active" db:"is_active"`
	Sort      int       `json:"sort" db:"sort"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
