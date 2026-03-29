package invoiceRoutes

import (
	"github.com/gofiber/fiber/v2"
	invoiceControllers "github.com/usegro/services/billing/int
	invoiceControllers "github.com/usegro/services/billing/internal/apps/invoice/controllers"
	invoiceServices "github.com/usegro/services/billing/internal/apps/invoice/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func InvoiceRouter(v1 fiber.Router, db *gorm.DB, emailCfg invoiceServices.EmailCfg) {
	group := v1.Group("/billing")

	svc := invoiceServices.NewInvoiceService(db, emailCfg)
	ctrl := invoiceControllers.NewInvoiceController(svc)

	invoices := group.Group("/invoices")
	invoices.Post("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), ctrl.CreateInvoice)
	invoices.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), ctrl.ListInvoices)
	invoices.Get("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), ctrl.GetInvoice)
	invoices.Patch("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), ctrl.UpdateInvoice)
	invoices.Delete("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), ctrl.DeleteInvoice)
	invoices.Post("/:id/send", middleware.JwtVerify(), middleware.CRMContextMiddleware(), ctrl.SendInvoice)
}
