package models

import (
	"time"

	"github.com/google/uuid"
)

// Plan is a subscription or package available for purchase on a service
type Plan struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CatalogItemID uuid.UUID `json:"catalog_item_id" gorm:"type:uuid;not null;index"`
	CRMID         uuid.UUID `json:"crm_id" gorm:"type:uuid;not null;index"`
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	PlanType      string    `json:"plan_type" gorm:"type:varchar(20);not null"` // subscription | package
	Price         float64   `json:"price" gorm:"type:decimal(10,2);default:0"`
	PriceCurrency string    `json:"price_currency" gorm:"type:varchar(10);default:'USD'"`
	BillingCycle  string    `json:"billing_cycle" gorm:"type:varchar(20)"`           // monthly | yearly (subscription only)
	SessionCount  *int      `json:"session_count"`                                   // nil = unlimited; per-cycle for subscription, total for package
	ValidityDays  *int      `json:"validity_days"`                                   // package: days after purchase until expiry
	Status        string    `json:"status" gorm:"type:varchar(20);default:'active'"` // active | archived
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// CustomerPlan tracks a customer's active subscription or package
type CustomerPlan struct {
	ID            uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	PlanID        uuid.UUID  `json:"plan_id" gorm:"type:uuid;not null;index"`
	CRMID         uuid.UUID  `json:"crm_id" gorm:"type:uuid;not null;index"`
	CustomerID    string     `json:"customer_id" gorm:"type:varchar(255);not null;index"` // FormSubmission _id (MongoDB)
	Status        string     `json:"status" gorm:"type:varchar(20);default:'active'"`     // active | expired | cancelled | paused
	PurchasedAt   time.Time  `json:"purchased_at"`
	ExpiresAt     *time.Time `json:"expires_at"`
	NextBillingAt *time.Time `json:"next_billing_at"`
	SessionsTotal *int       `json:"sessions_total"`
	SessionsUsed  int        `json:"sessions_used" gorm:"default:0"`
	Plan          Plan       `json:"plan" gorm:"foreignKey:PlanID"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// Booking is a scheduled session against a service, optionally tied to a CustomerPlan
type Booking struct {
	ID             uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CatalogItemID  uuid.UUID  `json:"catalog_item_id" gorm:"type:uuid;not null;index"`
	CRMID          uuid.UUID  `json:"crm_id" gorm:"type:uuid;not null;index"`
	CustomerID     string     `json:"customer_id" gorm:"type:varchar(255);not null;index"` // FormSubmission _id (MongoDB)
	CustomerPlanID *uuid.UUID `json:"customer_plan_id" gorm:"type:uuid"`                   // nil for per-session bookings
	Status         string     `json:"status" gorm:"type:varchar(20);default:'pending'"`    // pending | confirmed | cancelled | completed | no-show
	ScheduledAt    time.Time  `json:"scheduled_at"`
	Notes          string     `json:"notes" gorm:"type:text"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
