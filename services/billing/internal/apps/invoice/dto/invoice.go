package invoiceDTO

import (
	"time"

	"github.com/google/uuid"
	
)

// ── Requests ─────────────────────────────────────────────────────────────────

type LineItemRequest struct {
	Name         string  `json:"name" validate:"required"`
	Type         string  `json:"type" validate:"required,oneof=product service"`
	CatalogID    string  `json:"catalog_id"`
	Qty          float64 `json:"qty" validate:"required,gt=0"`
	Rate         float64 `json:"rate" validate:"gte=0"`
	BillingType  string  `json:"billing_type" validate:"required,oneof=one-time recurring"`
	BillingCycle string  `json:"billing_cycle"` // required when billing_type=recurring
	ImageURL     string  `json:"image_url"`
}

type CreateInvoiceRequest struct {
	CustomerIDs        []string          `json:"customer_ids" validate:"required,min=1"`
	CustomerNames      []string          `json:"customer_names"`
	CustomerEmails     []string          `json:"customer_emails"`
	LineItems          []LineItemRequest `json:"line_items" validate:"required,min=1,dive"`
	TaxRate            float64           `json:"tax_rate" validate:"gte=0"`
	Subject            string            `json:"subject"`
	TermsAndConditions string            `json:"terms_and_conditions"`
	ReferenceNumber    string            `json:"reference_number"`
	Memo               string            `json:"memo"`
	DueDate            *string           `json:"due_date"`             // "2006-01-02"
	RecurringStartDate *string           `json:"recurring_start_date"` // "2006-01-02"
}

type UpdateInvoiceRequest struct {
	CustomerIDs        []string          `json:"customer_ids"`
	CustomerNames      []string          `json:"customer_names"`
	CustomerEmails     []string          `json:"customer_emails"`
	LineItems          []LineItemRequest `json:"line_items"`
	TaxRate            *float64          `json:"tax_rate"`
	Subject            *string           `json:"subject"`
	TermsAndConditions *string           `json:"terms_and_conditions"`
	ReferenceNumber    *string           `json:"reference_number"`
	Memo               *string           `json:"memo"`
	DueDate            *string           `json:"due_date"`
	RecurringStartDate *string           `json:"recurring_start_date"`
}

// ── Responses ────────────────────────────────────────────────────────────────

type LineItemResponse struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	CatalogID    string    `json:"catalog_id"`
	Qty          float64   `json:"qty"`
	Rate         float64   `json:"rate"`
	Amount       float64   `json:"amount"`
	BillingType  string    `json:"billing_type"`
	BillingCycle string    `json:"billing_cycle"`
	ImageURL     string    `json:"image_url"`
}

type InvoiceResponse struct {
	ID                 uuid.UUID          `json:"id"`
	CRMID              string             `json:"crm_id"`
	InvoiceNumber      string             `json:"invoice_number"`
	Status             string             `json:"status"`
	CustomerIDs        []string           `json:"customer_ids"`
	CustomerNames      []string           `json:"customer_names"`
	CustomerEmails     []string           `json:"customer_emails"`
	LineItems          []LineItemResponse `json:"line_items"`
	TaxRate            float64            `json:"tax_rate"`
	Subtotal           float64            `json:"subtotal"`
	TaxAmount          float64            `json:"tax_amount"`
	Total              float64            `json:"total"`
	Subject            string             `json:"subject"`
	TermsAndConditions string             `json:"terms_and_conditions"`
	ReferenceNumber    string             `json:"reference_number"`
	Memo               string             `json:"memo"`
	DueDate            *time.Time         `json:"due_date"`
	RecurringStartDate *time.Time         `json:"recurring_start_date"`
	SentAt             *time.Time         `json:"sent_at"`
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}

type ListInvoicesResponse struct {
	Data  []InvoiceResponse `json:"data"`
	Total int64             `json:"total"`
	Page  int               `json:"page"`
	Limit int               `json:"limit"`
}

// ── Mappers ──────────────────────────────────────────────────────────────────

func ToInvoiceResponse(inv *invoiceModels.Invoice) InvoiceResponse {
	items := make([]LineItemResponse, len(inv.LineItems))
	for i, li := range inv.LineItems {
		items[i] = LineItemResponse{
			ID:           li.ID,
			Name:         li.Name,
			Type:         li.Type,
			CatalogID:    li.CatalogID,
			Qty:          li.Qty,
			Rate:         li.Rate,
			Amount:       li.Amount,
			BillingType:  string(li.BillingType),
			BillingCycle: li.BillingCycle,
			ImageURL:     li.ImageURL,
		}
	}

	return InvoiceResponse{
		ID:                 inv.ID,
		CRMID:              inv.CRMID,
		InvoiceNumber:      inv.InvoiceNumber,
		Status:             string(inv.Status),
		CustomerIDs:        orEmpty(inv.CustomerIDs),
		CustomerNames:      orEmpty(inv.CustomerNames),
		CustomerEmails:     orEmpty(inv.CustomerEmails),
		LineItems:          items,
		TaxRate:            inv.TaxRate,
		Subtotal:           inv.Subtotal,
		TaxAmount:          inv.TaxAmount,
		Total:              inv.Total,
		Subject:            inv.Subject,
		TermsAndConditions: inv.TermsAndConditions,
		ReferenceNumber:    inv.ReferenceNumber,
		Memo:               inv.Memo,
		DueDate:            inv.DueDate,
		RecurringStartDate: inv.RecurringStartDate,
		SentAt:             inv.SentAt,
		CreatedAt:          inv.CreatedAt,
		UpdatedAt:          inv.UpdatedAt,
	}
}

func orEmpty(s []string) []string {
	if s == nil {
		return []string{}
	}
	return s
}
