package entity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

const (
	PlanFree       = "free"
	PlanStarter    = "starter"
	PlanPro        = "pro"
	PlanEnterprise = "enterprise"
)

const (
	StatusActive    = "active"
	StatusSuspended = "suspended"
	StatusCancelled = "canceled"
	StatusExpired   = "expired"
)

type Tenant struct {
	ID                 uuid.UUID            `json:"id" db:"id"`
	Name               string               `json:"name" db:"name" validate:"required,max=255"`
	Slug               string               `json:"slug" db:"slug" validate:"required,max=100"`
	MaxDivisions       int                  `json:"max_divisions" db:"max_divisions"`
	MaxAgents          int                  `json:"max_agents" db:"max_agents"`
	MaxQuickReplies    int                  `json:"max_quick_replies" db:"max_quick_replies"`
	MaxWahaWhatsApp    int                  `json:"max_waha_whatsapp" db:"max_waha_whatsapp"`
	MaxMetaWhatsApp    int                  `json:"max_meta_whatsapp" db:"max_meta_whatsapp"`
	MaxMetaMessenger   int                  `json:"max_meta_messenger" db:"max_meta_messenger"`
	MaxInstagram       int                  `json:"max_instagram" db:"max_instagram"`
	MaxTelegram        int                  `json:"max_telegram" db:"max_telegram"`
	MaxWebChat         int                  `json:"max_webchat" db:"max_webchat"`
	MaxLinkChat        int                  `json:"max_linkchat" db:"max_linkchat"`
	AiCredits          int                  `json:"ai_credits" db:"ai_credits"`
	SubscriptionPlan   string               `json:"subscription_plan" db:"subscription_plan" validate:"required,oneof=free starter pro enterprise"`
	SubscriptionStatus string               `json:"subscription_status" db:"subscription_status" validate:"required,oneof=active suspended canceled expired"`
	TrialEndsAt        NullTime             `json:"trial_ends_at,omitempty" db:"trial_ends_at"`
	IsActive           bool                 `json:"is_active" db:"is_active"`
	Settings           *TenantSettingConfig `json:"settings,omitempty" db:"settings"`
	CreatedAt          time.Time            `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at" db:"updated_at"`
	DeletedAt          NullTime             `json:"deleted_at,omitempty" db:"deleted_at"`
	// Transient profile fields backed by the settings JSONB (see UpdateLogo),
	// not dedicated columns. db:"-" so they are ignored by SELECT * scans.
	Description NullString `json:"description,omitempty" db:"-"`
	LogoURL     NullString `json:"logo_url,omitempty" db:"-"`
}

type TenantSettingConfig map[string]interface{}

func (ts TenantSettingConfig) Value() (driver.Value, error) {
	if ts == nil {
		return nil, nil
	}
	return json.Marshal(ts)
}

func (ts *TenantSettingConfig) Scan(value interface{}) error {
	if value == nil {
		*ts = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, ts)
}

func (Tenant) TableName() string {
	return "tenants"
}

func (t *Tenant) IsDeleted() bool {
	return t.DeletedAt.Valid
}

func (t *Tenant) IsSubscriptionActive() bool {
	return t.SubscriptionStatus == StatusActive && t.IsActive
}

func (t *Tenant) IsSubscriptionExpired() bool {
	return t.SubscriptionStatus == StatusExpired && t.IsActive
}

func (t *Tenant) IsOnTrial() bool {
	if !t.TrialEndsAt.Valid {
		return false
	}
	return time.Now().Before(t.TrialEndsAt.Time)
}

func (t *Tenant) IsTrialExpired() bool {
	if !t.TrialEndsAt.Valid {
		return false
	}
	return time.Now().After(t.TrialEndsAt.Time)
}

func (t *Tenant) CanCreateDivision(currentCount int) bool {
	return currentCount < t.MaxDivisions
}

func (t *Tenant) CanAddAgent(currentCount int) bool {
	return currentCount < t.MaxAgents
}

func (t *Tenant) CanAddChannel(channelType string, currentCount int) bool {
	switch channelType {
	case "waha_whatsapp":
		return currentCount < t.MaxWahaWhatsApp
	case "meta_whatsapp":
		return currentCount < t.MaxMetaWhatsApp
	case "meta_messenger":
		return currentCount < t.MaxMetaMessenger
	case "instagram":
		return currentCount < t.MaxInstagram
	case "telegram":
		return currentCount < t.MaxTelegram
	case "webchat":
		return currentCount < t.MaxWebChat
	case "linkchat":
		return currentCount < t.MaxLinkChat
	default:
		return false
	}
}

func (t *Tenant) GetPlanLimits() map[string]int {
	return map[string]int{
		"divisions":      t.MaxDivisions,
		"agents":         t.MaxAgents,
		"quick_replies":  t.MaxQuickReplies,
		"waha_whatsapp":  t.MaxWahaWhatsApp,
		"meta_whatsapp":  t.MaxMetaWhatsApp,
		"meta_messenger": t.MaxMetaMessenger,
		"instagram":      t.MaxInstagram,
		"telegram":       t.MaxTelegram,
		"webchat":        t.MaxWebChat,
		"linkchat":       t.MaxLinkChat,
	}
}

func (t *Tenant) UpgradePlan(plan string) {
	t.SubscriptionPlan = plan

	// Set limits based on plan
	switch plan {
	case PlanStarter:
		t.MaxDivisions = 5
		t.MaxAgents = 5
		t.MaxQuickReplies = 100
		t.MaxWahaWhatsApp = 1
		t.MaxMetaWhatsApp = 1
		t.MaxMetaMessenger = 5
		t.MaxInstagram = 5
		t.MaxTelegram = 5
	case PlanPro:
		t.MaxDivisions = 20
		t.MaxAgents = 20
		t.MaxQuickReplies = 500
		t.MaxWahaWhatsApp = 5
		t.MaxMetaWhatsApp = 5
		t.MaxMetaMessenger = 10
		t.MaxInstagram = 10
		t.MaxTelegram = 10
	case PlanEnterprise:
		t.MaxDivisions = 100
		t.MaxAgents = 100
		t.MaxQuickReplies = 1000
		t.MaxWahaWhatsApp = 10
		t.MaxMetaWhatsApp = 10
		t.MaxMetaMessenger = 100
		t.MaxInstagram = 100
		t.MaxTelegram = 100
	default: // free
		t.MaxDivisions = 1
		t.MaxAgents = 1
		t.MaxQuickReplies = 5
		t.MaxWahaWhatsApp = 0
		t.MaxMetaWhatsApp = 0
		t.MaxMetaMessenger = 1
		t.MaxInstagram = 1
		t.MaxTelegram = 1
	}
}
