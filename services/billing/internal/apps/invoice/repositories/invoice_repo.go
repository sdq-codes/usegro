package invoiceRepositories

import (
	"github.com/google/uuid"
	invoiceModels 
	"gorm.io/gorm"
)

// InvoiceFilter holds all optional query filters for listing invoices.
type InvoiceFilter struct {
	// Status / display-status label (e.g. "draft", "Overdue", "Not due yet", "Due today")
	Status string

	// Customer name substring (searches JSON array column as text)
	CustomerName string

	// Invoice number substring
	InvoiceNumber string

	// Due-date window (inclusive, ISO date strings "2006-01-02")
	DueDateFrom string
	DueDateTo   string

	// Created-at window
	CreatedFrom string
	CreatedTo   string

	// Amount range
	AmountMin *float64
	AmountMax *float64

	// Billing type display label: "One-time" | "Monthly" | "Weekly" | "Yearly" | "Daily" | "Bi-Weekly"
	BillingType string
}

type InvoiceRepositoryInterface interface {
	Create(invoice *invoiceModels.Invoice) error
	FindByID(id uuid.UUID, crmID string) (*invoiceModels.Invoice, error)
	FindAll(crmID string, page, limit int, filter InvoiceFilter) ([]invoiceModels.Invoice, int64, error)
	Update(invoice *invoiceModels.Invoice) error
	Delete(id uuid.UUID, crmID string) error
	ReplaceLineItems(invoiceID uuid.UUID, items []invoiceModels.InvoiceLineItem) error
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepositoryInterface {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *invoiceModels.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) FindByID(id uuid.UUID, crmID string) (*invoiceModels.Invoice, error) {
	var invoice invoiceModels.Invoice
	err := r.db.
		Preload("LineItems").
		Where("id = ? AND crm_id = ?", id, crmID).
		First(&invoice).Error
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) FindAll(crmID string, page, limit int, filter InvoiceFilter) ([]invoiceModels.Invoice, int64, error) {
	var invoices []invoiceModels.Invoice
	var total int64

	q := r.db.Model(&invoiceModels.Invoice{}).Where("invoices.crm_id = ?", crmID)

	// ── Status / display-status ──────────────────────────────────
	switch filter.Status {
	case "Overdue":
		q = q.Where("invoices.status = 'overdue' OR (invoices.status = 'sent' AND invoices.due_date IS NOT NULL AND invoices.due_date < NOW())")
	case "Not due yet":
		q = q.Where("invoices.status = 'sent' AND (invoices.due_date IS NULL OR invoices.due_date > NOW())")
	case "Due today":
		q = q.Where("invoices.status = 'sent' AND DATE(invoices.due_date) = CURRENT_DATE")
	case "Paid":
		q = q.Where("invoices.status = 'paid'")
	case "Draft":
		q = q.Where("invoices.status = 'draft'")
	case "Cancelled":
		q = q.Where("invoices.status = 'cancelled'")
	case "Sent":
		q = q.Where("invoices.status = 'sent'")
	default:
		if filter.Status != "" {
			q = q.Where("invoices.status = ?", filter.Status)
		}
	}

	// ── Customer name (JSON array stored as text) ─────────────────
	if filter.CustomerName != "" {
		q = q.Where("invoices.customer_names::text ILIKE ?", "%"+filter.CustomerName+"%")
	}

	// ── Invoice number ────────────────────────────────────────────
	if filter.InvoiceNumber != "" {
		q = q.Where("invoices.invoice_number ILIKE ?", "%"+filter.InvoiceNumber+"%")
	}

	// ── Due-date range ────────────────────────────────────────────
	if filter.DueDateFrom != "" {
		q = q.Where("invoices.due_date >= ?", filter.DueDateFrom)
	}
	if filter.DueDateTo != "" {
		q = q.Where("invoices.due_date <= ?", filter.DueDateTo+" 23:59:59")
	}

	// ── Created-at range ──────────────────────────────────────────
	if filter.CreatedFrom != "" {
		q = q.Where("invoices.created_at >= ?", filter.CreatedFrom)
	}
	if filter.CreatedTo != "" {
		q = q.Where("invoices.created_at <= ?", filter.CreatedTo+" 23:59:59")
	}

	// ── Amount range ──────────────────────────────────────────────
	if filter.AmountMin != nil {
		q = q.Where("invoices.total >= ?", *filter.AmountMin)
	}
	if filter.AmountMax != nil {
		q = q.Where("invoices.total <= ?", *filter.AmountMax)
	}

	// ── Billing type (via line items subquery) ────────────────────
	if filter.BillingType != "" {
		if filter.BillingType == "One-time" {
			q = q.Where(`EXISTS (
				SELECT 1 FROM invoice_line_items
				WHERE invoice_line_items.invoice_id = invoices.id
				  AND invoice_line_items.billing_type = 'one-time'
				  AND invoice_line_items.deleted_at IS NULL
			)`)
		} else {
			// Recurring with a specific cycle label
			cycleMap := map[string]string{
				"Monthly":   "monthly",
				"Weekly":    "weekly",
				"Yearly":    "yearly",
				"Daily":     "daily",
				"Bi-Weekly": "biweekly",
			}
			if cycle, ok := cycleMap[filter.BillingType]; ok {
				q = q.Where(`EXISTS (
					SELECT 1 FROM invoice_line_items
					WHERE invoice_line_items.invoice_id = invoices.id
					  AND invoice_line_items.billing_type = 'recurring'
					  AND LOWER(invoice_line_items.billing_cycle) = ?
					  AND invoice_line_items.deleted_at IS NULL
				)`, cycle)
			}
		}
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := q.
		Preload("LineItems").
		Order("invoices.created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&invoices).Error

	return invoices, total, err
}

func (r *invoiceRepository) Update(invoice *invoiceModels.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *invoiceRepository) Delete(id uuid.UUID, crmID string) error {
	return r.db.
		Where("id = ? AND crm_id = ?", id, crmID).
		Delete(&invoiceModels.Invoice{}).Error
}

func (r *invoiceRepository) ReplaceLineItems(invoiceID uuid.UUID, items []invoiceModels.InvoiceLineItem) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("invoice_id = ?", invoiceID).Delete(&invoiceModels.InvoiceLineItem{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}
