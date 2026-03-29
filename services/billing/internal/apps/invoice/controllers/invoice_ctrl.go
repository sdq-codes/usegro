package invoiceControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	invoiceDTO "github.com/usegro/services/billing/internal/apps/invoice/dto"
	invoiceRepositories "github.com/us
	invoiceDTO "github.com/usegro/services/billing/internal/apps/invoice/dto"
	invoiceRepositories "github.com/usegro/services/billing/internal/apps/invoice/repositories"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strconv"
)

type InvoiceController struct {
	svc invoiceServices.InvoiceServiceInterface
}

func NewInvoiceController(svc invoiceServices.InvoiceServiceInterface) *InvoiceController {
	return &InvoiceController{svc: svc}
}

// CreateInvoice POST /billing/invoices
func (ctrl *InvoiceController) CreateInvoice(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	var req invoiceDTO.CreateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}

	if len(req.CustomerIDs) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.CommonResponse{
			ResponseCode:    response.VALIDATION_FAILED,
			ResponseMessage: "At least one customer is required",
		})
	}
	if len(req.LineItems) == 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(response.CommonResponse{
			ResponseCode:    response.VALIDATION_FAILED,
			ResponseMessage: "At least one line item is required",
		})
	}

	invoice, err := ctrl.svc.CreateInvoice(crmID, req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Invoice created",
		Data:            invoiceDTO.ToInvoiceResponse(invoice),
	})
}

// ListInvoices GET /billing/invoices
func (ctrl *InvoiceController) ListInvoices(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	filter := invoiceRepositories.InvoiceFilter{
		Status:        c.Query("status", ""),
		CustomerName:  c.Query("customer_name", ""),
		InvoiceNumber: c.Query("invoice_number", ""),
		DueDateFrom:   c.Query("due_date_from", ""),
		DueDateTo:     c.Query("due_date_to", ""),
		CreatedFrom:   c.Query("created_from", ""),
		CreatedTo:     c.Query("created_to", ""),
		BillingType:   c.Query("billing_type", ""),
	}

	if v := c.Query("amount_min", ""); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.AmountMin = &f
		}
	}
	if v := c.Query("amount_max", ""); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			filter.AmountMax = &f
		}
	}

	invoices, total, err := ctrl.svc.ListInvoices(crmID, page, limit, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	data := make([]invoiceDTO.InvoiceResponse, len(invoices))
	for i, inv := range invoices {
		data[i] = invoiceDTO.ToInvoiceResponse(&inv)
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Invoices fetched",
		Data: invoiceDTO.ListInvoicesResponse{
			Data:  data,
			Total: total,
			Page:  page,
			Limit: limit,
		},
	})
}

// GetInvoice GET /billing/invoices/:id
func (ctrl *InvoiceController) GetInvoice(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid invoice ID",
		})
	}

	invoice, err := ctrl.svc.GetInvoice(id, crmID)
	if err != nil {
		if err.Error() == "invoice not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
				ResponseCode:    response.RESOURCE_NOT_FOUND,
				ResponseMessage: "Invoice not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Invoice fetched",
		Data:            invoiceDTO.ToInvoiceResponse(invoice),
	})
}

// UpdateInvoice PATCH /billing/invoices/:id
func (ctrl *InvoiceController) UpdateInvoice(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid invoice ID",
		})
	}

	var req invoiceDTO.UpdateInvoiceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}

	invoice, err := ctrl.svc.UpdateInvoice(id, crmID, req)
	if err != nil {
		if err.Error() == "invoice not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
				ResponseCode:    response.RESOURCE_NOT_FOUND,
				ResponseMessage: "Invoice not found",
			})
		}
		if err.Error() == "only draft invoices can be edited" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.CommonResponse{
				ResponseCode:    response.VALIDATION_FAILED,
				ResponseMessage: "Only draft invoices can be edited",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Invoice updated",
		Data:            invoiceDTO.ToInvoiceResponse(invoice),
	})
}

// DeleteInvoice DELETE /billing/invoices/:id
func (ctrl *InvoiceController) DeleteInvoice(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid invoice ID",
		})
	}

	if err := ctrl.svc.DeleteInvoice(id, crmID); err != nil {
		if err.Error() == "invoice not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
				ResponseCode:    response.RESOURCE_NOT_FOUND,
				ResponseMessage: "Invoice not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_DELETED,
		ResponseMessage: "Invoice deleted",
	})
}

// SendInvoice POST /billing/invoices/:id/send
func (ctrl *InvoiceController) SendInvoice(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid invoice ID",
		})
	}

	invoice, err := ctrl.svc.SendInvoice(id, crmID)
	if err != nil {
		if err.Error() == "invoice not found" {
			return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
				ResponseCode:    response.RESOURCE_NOT_FOUND,
				ResponseMessage: "Invoice not found",
			})
		}
		if err.Error() == "no customer emails on this invoice" {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(response.CommonResponse{
				ResponseCode:    response.VALIDATION_FAILED,
				ResponseMessage: "No customer emails found on this invoice",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Invoice sent",
		Data:            invoiceDTO.ToInvoiceResponse(invoice),
	})
}
