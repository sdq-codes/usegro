package invoiceModels

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvoiceStatus string
type BillingType string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusSent      InvoiceStatus = "sent"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"

	BillingTypeOneTime   BillingType = "one-time"
	BillingTypeRecurring BillingType = "recurring"
)

type Invoice struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CRMID     string         `json:"crm_id" gorm:"not null;index"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Invoice identity
	InvoiceNumber string        `json:"invoice_number" gorm:"not null"`
	Status        InvoiceStatus `json:"status" gorm:"default:'draft'"`

	// Customers (denormalised from MongoDB CRM — stored as JSON arrays)
	CustomerIDs    []string `json:"customer_ids" gorm:"serializer:json"`
	CustomerNames  []string `json:"customer_names" gorm:"serializer:json"`
	CustomerEmails []string `json:"customer_emails" gorm:"serializer:json"`

	// Line items
	LineItems []InvoiceLineItem `json:"line_items" gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE"`

	// Financials
	TaxRate   float64 `json:"tax_rate" gorm:"default:0"`
	Subtotal  float64 `json:"subtotal" gorm:"default:0"`
	TaxAmount float64 `json:"tax_amount" gorm:"default:0"`
	Total     float64 `json:"total" gorm:"default:0"`

	// Notes & metadata
	Subject            string `json:"subject"`
	TermsAndConditions string `json:"terms_and_conditions"`
	ReferenceNumber    string `json:"reference_number"`
	Memo               string `json:"memo"`

	// Dates
	DueDate            *time.Time `json:"due_date"`
	RecurringStartDate *time.Time `json:"recurring_start_date"`
	SentAt             *time.Time `json:"sent_at"`
}

type InvoiceLineItem struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	InvoiceID uuid.UUID `json:"invoice_id" gorm:"type:uuid;not null;index"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name         string      `json:"name" gorm:"not null"`
	Type         string      `json:"type"` // product | service
	CatalogID    string      `json:"catalog_id"`
	Qty          float64     `json:"qty" gorm:"default:1"`
	Rate         float64     `json:"rate" gorm:"not null"`
	Amount       float64     `json:"amount"` // qty * rate
	BillingType  BillingType `json:"billing_type" gorm:"default:'one-time'"`
	BillingCycle string      `json:"billing_cycle"` // monthly | yearly | weekly
	ImageURL     string      `json:"image_url"`
}
