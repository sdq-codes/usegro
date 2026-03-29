package invoiceServices

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	invoiceDTO "github.com/usegro/services/billing/
	invoiceDTO "github.com/usegro/services/billing/internal/apps/invoice/dto"
	invoiceModels "github.com/usegro/services/billing/internal/apps/invoice/models"
	invoiceRepositories "github.com/usegro/services/billing/internal/apps/invoice/repositories"
	"github.com/usegro/services/billing/internal/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type EmailCfg struct {
	Env       string
	Region    string
	FromEmail string
}

type InvoiceServiceInterface interface {
	CreateInvoice(crmID string, req invoiceDTO.CreateInvoiceRequest) (*invoiceModels.Invoice, error)
	GetInvoice(id uuid.UUID, crmID string) (*invoiceModels.Invoice, error)
	ListInvoices(crmID string, page, limit int, filter invoiceRepositories.InvoiceFilter) ([]invoiceModels.Invoice, int64, error)
	UpdateInvoice(id uuid.UUID, crmID string, req invoiceDTO.UpdateInvoiceRequest) (*invoiceModels.Invoice, error)
	DeleteInvoice(id uuid.UUID, crmID string) error
	SendInvoice(id uuid.UUID, crmID string) (*invoiceModels.Invoice, error)
}

type invoiceService struct {
	repo     invoiceRepositories.InvoiceRepositoryInterface
	emailCfg EmailCfg
}

func NewInvoiceService(db *gorm.DB, emailCfg EmailCfg) InvoiceServiceInterface {
	return &invoiceService{
		repo:     invoiceRepositories.NewInvoiceRepository(db),
		emailCfg: emailCfg,
	}
}

// ── Create ────────────────────────────────────────────────────────────────────

func (s *invoiceService) CreateInvoice(crmID string, req invoiceDTO.CreateInvoiceRequest) (*invoiceModels.Invoice, error) {
	lineItems, subtotal := buildLineItems(req.LineItems)
	taxAmount := subtotal * req.TaxRate / 100
	total := subtotal + taxAmount

	invoice := &invoiceModels.Invoice{
		CRMID:              crmID,
		InvoiceNumber:      generateInvoiceNumber(),
		Status:             invoiceModels.InvoiceStatusDraft,
		CustomerIDs:        req.CustomerIDs,
		CustomerNames:      req.CustomerNames,
		CustomerEmails:     req.CustomerEmails,
		LineItems:          lineItems,
		TaxRate:            req.TaxRate,
		Subtotal:           subtotal,
		TaxAmount:          taxAmount,
		Total:              total,
		Subject:            req.Subject,
		TermsAndConditions: req.TermsAndConditions,
		ReferenceNumber:    req.ReferenceNumber,
		Memo:               req.Memo,
		DueDate:            parseDate(req.DueDate),
		RecurringStartDate: parseDate(req.RecurringStartDate),
	}

	if err := s.repo.Create(invoice); err != nil {
		return nil, err
	}
	return invoice, nil
}

// ── Read ──────────────────────────────────────────────────────────────────────

func (s *invoiceService) GetInvoice(id uuid.UUID, crmID string) (*invoiceModels.Invoice, error) {
	inv, err := s.repo.FindByID(id, crmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, err
	}
	return inv, nil
}

func (s *invoiceService) ListInvoices(crmID string, page, limit int, filter invoiceRepositories.InvoiceFilter) ([]invoiceModels.Invoice, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return s.repo.FindAll(crmID, page, limit, filter)
}

// ── Update ────────────────────────────────────────────────────────────────────

func (s *invoiceService) UpdateInvoice(id uuid.UUID, crmID string, req invoiceDTO.UpdateInvoiceRequest) (*invoiceModels.Invoice, error) {
	inv, err := s.repo.FindByID(id, crmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, err
	}

	if inv.Status != invoiceModels.InvoiceStatusDraft {
		return nil, fmt.Errorf("only draft invoices can be edited")
	}

	if req.CustomerIDs != nil {
		inv.CustomerIDs = req.CustomerIDs
	}
	if req.CustomerNames != nil {
		inv.CustomerNames = req.CustomerNames
	}
	if req.CustomerEmails != nil {
		inv.CustomerEmails = req.CustomerEmails
	}
	if req.Subject != nil {
		inv.Subject = *req.Subject
	}
	if req.TermsAndConditions != nil {
		inv.TermsAndConditions = *req.TermsAndConditions
	}
	if req.ReferenceNumber != nil {
		inv.ReferenceNumber = *req.ReferenceNumber
	}
	if req.Memo != nil {
		inv.Memo = *req.Memo
	}
	if req.DueDate != nil {
		inv.DueDate = parseDate(req.DueDate)
	}
	if req.RecurringStartDate != nil {
		inv.RecurringStartDate = parseDate(req.RecurringStartDate)
	}

	if req.LineItems != nil {
		newItems, subtotal := buildLineItems(req.LineItems)
		taxRate := inv.TaxRate
		if req.TaxRate != nil {
			taxRate = *req.TaxRate
			inv.TaxRate = taxRate
		}
		taxAmount := subtotal * taxRate / 100
		inv.Subtotal = subtotal
		inv.TaxAmount = taxAmount
		inv.Total = subtotal + taxAmount

		for i := range newItems {
			newItems[i].InvoiceID = inv.ID
		}
		if err := s.repo.ReplaceLineItems(inv.ID, newItems); err != nil {
			return nil, err
		}
		inv.LineItems = newItems
	} else if req.TaxRate != nil {
		inv.TaxRate = *req.TaxRate
		inv.TaxAmount = inv.Subtotal * inv.TaxRate / 100
		inv.Total = inv.Subtotal + inv.TaxAmount
	}

	if err := s.repo.Update(inv); err != nil {
		return nil, err
	}
	return inv, nil
}

// ── Delete ────────────────────────────────────────────────────────────────────

func (s *invoiceService) DeleteInvoice(id uuid.UUID, crmID string) error {
	_, err := s.repo.FindByID(id, crmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("invoice not found")
		}
		return err
	}
	return s.repo.Delete(id, crmID)
}

// ── Send ──────────────────────────────────────────────────────────────────────

func (s *invoiceService) SendInvoice(id uuid.UUID, crmID string) (*invoiceModels.Invoice, error) {
	inv, err := s.repo.FindByID(id, crmID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("invoice not found")
		}
		return nil, err
	}

	if len(inv.CustomerEmails) == 0 {
		return nil, fmt.Errorf("no customer emails on this invoice")
	}

	html := buildInvoiceEmailHTML(inv)

	for _, toEmail := range inv.CustomerEmails {
		mail := &sharedemail.MailRequest{
			From:    s.emailCfg.FromEmail,
			To:      []string{toEmail},
			Subject: fmt.Sprintf("Invoice from Gro — #%s", inv.InvoiceNumber),
			Body:    html,
		}

		cfg := sharedemail.EmailConfig{
			Env:       s.emailCfg.Env,
			Region:    s.emailCfg.Region,
			FromEmail: s.emailCfg.FromEmail,
		}

		// Send inline (Body is pre-built, skip template parsing)
		mail.SendBody(cfg)
		logger.Log.Info("Invoice email sent", zap.String("to", toEmail), zap.String("invoice", inv.InvoiceNumber))
	}

	now := time.Now()
	inv.SentAt = &now
	inv.Status = invoiceModels.InvoiceStatusSent
	if err := s.repo.Update(inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func buildLineItems(reqs []invoiceDTO.LineItemRequest) ([]invoiceModels.InvoiceLineItem, float64) {
	items := make([]invoiceModels.InvoiceLineItem, len(reqs))
	var subtotal float64
	for i, r := range reqs {
		amount := r.Qty * r.Rate
		subtotal += amount
		items[i] = invoiceModels.InvoiceLineItem{
			Name:         r.Name,
			Type:         r.Type,
			CatalogID:    r.CatalogID,
			Qty:          r.Qty,
			Rate:         r.Rate,
			Amount:       amount,
			BillingType:  invoiceModels.BillingType(r.BillingType),
			BillingCycle: r.BillingCycle,
			ImageURL:     r.ImageURL,
		}
	}
	return items, subtotal
}

func generateInvoiceNumber() string {
	now := time.Now()
	return fmt.Sprintf("INV-%d%02d%02d-%04d", now.Year(), now.Month(), now.Day(), now.UnixMilli()%10000)
}

func parseDate(s *string) *time.Time {
	if s == nil || *s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", *s)
	if err != nil {
		return nil
	}
	return &t
}

func fmtCurrency(v float64) string {
	return fmt.Sprintf("$%.2f", v)
}

func buildInvoiceEmailHTML(inv *invoiceModels.Invoice) string {
	billedTo := strings.Join(inv.CustomerNames, ", ")
	if billedTo == "" {
		billedTo = "Customer"
	}

	genDate := inv.CreatedAt.Format("Jan 2, 2006")

	var dueLine string
	if inv.DueDate != nil {
		dueLine = fmt.Sprintf(`<p style="font-size:12px;color:#9CA3AF;margin:4px 0 0;">Due %s</p>`, inv.DueDate.Format("Jan 2, 2006"))
	}

	// Recurring summary
	recurringGroups := map[string]float64{}
	for _, li := range inv.LineItems {
		if li.BillingType == invoiceModels.BillingTypeRecurring {
			cycle := li.BillingCycle
			if cycle == "" {
				cycle = "monthly"
			}
			recurringGroups[cycle] += li.Amount
		}
	}

	hasOneTime := false
	var oneTimeTotal float64
	for _, li := range inv.LineItems {
		if li.BillingType == invoiceModels.BillingTypeOneTime {
			hasOneTime = true
			oneTimeTotal += li.Amount
		}
	}
	oneTimeTax := oneTimeTotal * inv.TaxRate / 100
	oneTimeTotal += oneTimeTax

	var amountLabel, amountValue string
	if hasOneTime {
		amountLabel = "Amount due"
		amountValue = fmtCurrency(oneTimeTotal)
	} else if len(recurringGroups) > 0 {
		amountLabel = "Per cycle"
		for _, v := range recurringGroups {
			amountValue = fmtCurrency(v)
			break
		}
	} else {
		amountLabel = "Amount due"
		amountValue = fmtCurrency(inv.Total)
	}

	// Line items rows
	var itemRows strings.Builder
	for _, li := range inv.LineItems {
		cycleTag := ""
		if li.BillingType == invoiceModels.BillingTypeRecurring {
			cycleTag = fmt.Sprintf(` <span style="font-size:10px;color:#2176AE;">(%s)</span>`, li.BillingCycle)
		}
		itemRows.WriteString(fmt.Sprintf(`
		<tr>
			<td style="padding:8px 0;border-bottom:1px solid #F3F4F6;font-size:13px;color:#1E212B;">%s%s</td>
			<td style="padding:8px 0;border-bottom:1px solid #F3F4F6;font-size:13px;color:#6B7280;text-align:center;">%g</td>
			<td style="padding:8px 0;border-bottom:1px solid #F3F4F6;font-size:13px;color:#6B7280;text-align:right;">%s</td>
			<td style="padding:8px 0;border-bottom:1px solid #F3F4F6;font-size:13px;font-weight:600;color:#1E212B;text-align:right;">%s</td>
		</tr>`, li.Name, cycleTag, li.Qty, fmtCurrency(li.Rate), fmtCurrency(li.Amount)))
	}

	var noteSection string
	if inv.Subject != "" {
		noteSection = fmt.Sprintf(`
		<div style="background:#EFF6FF;border-radius:8px;padding:12px 16px;margin-bottom:20px;">
			<p style="font-size:10px;font-weight:700;color:#2176AE;text-transform:uppercase;letter-spacing:0.05em;margin:0 0 4px;">Note</p>
			<p style="font-size:12px;color:#4B5563;margin:0;">%s</p>
		</div>`, inv.Subject)
	}

	ctaLabel := "View &amp; Pay Invoice"
	if len(recurringGroups) > 0 && !hasOneTime {
		ctaLabel = "View Invoice &amp; Subscribe"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1"></head>
<body style="margin:0;padding:0;background:#F6F8FA;font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;">
  <div style="max-width:560px;margin:32px auto;background:#F6F8FA;padding:0 16px 32px;">
    <div style="background:#ffffff;border-radius:16px;overflow:hidden;box-shadow:0 4px 24px rgba(0,0,0,0.06);">

      <!-- Header -->
      <div style="background:#E87117;padding:28px 32px;text-align:center;">
        <p style="font-size:20px;font-weight:800;color:#ffffff;margin:0;">Gro</p>
      </div>

      <!-- Body -->
      <div style="padding:28px 32px;">
        <p style="font-size:14px;color:#374151;margin:0 0 4px;">Hi %s,</p>
        <p style="font-size:14px;color:#6B7280;line-height:1.6;margin:0 0 24px;">
          You have a new invoice waiting for you. Please review the details below.
        </p>

        <!-- Amount card -->
        <div style="background:#FAFAFA;border:1px solid #EDEDEE;border-radius:12px;padding:20px;text-align:center;margin-bottom:24px;">
          <p style="font-size:10px;color:#9CA3AF;text-transform:uppercase;letter-spacing:0.08em;margin:0 0 6px;">%s</p>
          <p style="font-size:28px;font-weight:800;color:#E87117;margin:0;">%s</p>
          %s
        </div>

        <!-- Invoice details -->
        <table style="width:100%%;border-collapse:collapse;margin-bottom:24px;">
          <tr>
            <td style="font-size:12px;color:#9CA3AF;padding:4px 0;">Invoice #</td>
            <td style="font-size:12px;font-weight:600;color:#1E212B;text-align:right;padding:4px 0;">%s</td>
          </tr>
          <tr>
            <td style="font-size:12px;color:#9CA3AF;padding:4px 0;">Date issued</td>
            <td style="font-size:12px;font-weight:600;color:#1E212B;text-align:right;padding:4px 0;">%s</td>
          </tr>
          <tr>
            <td style="font-size:12px;color:#9CA3AF;padding:4px 0;">Billed to</td>
            <td style="font-size:12px;font-weight:600;color:#1E212B;text-align:right;padding:4px 0;">%s</td>
          </tr>
        </table>

        <!-- Line items -->
        <table style="width:100%%;border-collapse:collapse;margin-bottom:24px;">
          <thead>
            <tr>
              <th style="font-size:10px;font-weight:700;color:#9CA3AF;text-transform:uppercase;letter-spacing:0.05em;padding:0 0 8px;text-align:left;">Item</th>
              <th style="font-size:10px;font-weight:700;color:#9CA3AF;text-transform:uppercase;letter-spacing:0.05em;padding:0 0 8px;text-align:center;">Qty</th>
              <th style="font-size:10px;font-weight:700;color:#9CA3AF;text-transform:uppercase;letter-spacing:0.05em;padding:0 0 8px;text-align:right;">Rate</th>
              <th style="font-size:10px;font-weight:700;color:#9CA3AF;text-transform:uppercase;letter-spacing:0.05em;padding:0 0 8px;text-align:right;">Amount</th>
            </tr>
          </thead>
          <tbody>%s</tbody>
        </table>

        <!-- Totals -->
        <table style="width:100%%;border-collapse:collapse;margin-bottom:24px;">
          <tr>
            <td style="font-size:12px;color:#9CA3AF;padding:3px 0;text-align:right;">Subtotal</td>
            <td style="font-size:12px;font-weight:600;color:#1E212B;text-align:right;padding:3px 0;width:100px;">%s</td>
          </tr>
          <tr>
            <td style="font-size:12px;color:#9CA3AF;padding:3px 0;text-align:right;">Tax (%.0f%%)</td>
            <td style="font-size:12px;font-weight:600;color:#1E212B;text-align:right;padding:3px 0;">%s</td>
          </tr>
          <tr style="border-top:1px solid #F3F4F6;">
            <td style="font-size:14px;font-weight:700;color:#1E212B;padding:8px 0 0;text-align:right;">Total</td>
            <td style="font-size:16px;font-weight:800;color:#E87117;padding:8px 0 0;text-align:right;">%s</td>
          </tr>
        </table>

        %s

        <!-- CTA -->
        <a href="#" style="display:block;background:#E87117;color:#ffffff;text-decoration:none;text-align:center;font-size:14px;font-weight:700;padding:14px;border-radius:12px;">
          %s
        </a>

        <p style="font-size:11px;color:#9CA3AF;text-align:center;margin:16px 0 0;line-height:1.6;">
          Sent via Gro. Questions? Reply to this email.
        </p>
      </div>

      <!-- Footer -->
      <div style="background:#F9FAFB;border-top:1px solid #F3F4F6;padding:12px 32px;text-align:center;">
        <p style="font-size:10px;color:#D1D5DB;margin:0;">Powered by Gro</p>
      </div>
    </div>
  </div>
</body>
</html>`,
		billedTo,
		amountLabel, amountValue, dueLine,
		inv.InvoiceNumber, genDate, billedTo,
		itemRows.String(),
		fmtCurrency(inv.Subtotal), inv.TaxRate, fmtCurrency(inv.TaxAmount), fmtCurrency(inv.Total),
		noteSection,
		ctaLabel,
	)
}
