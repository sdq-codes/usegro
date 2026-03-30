package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/billing/config"
	"github.com/usegro/services/billing/database"
	invoiceRoutes "github.com/usegro/services/billing/internal/apps/invoice/routes"
	invoiceServices "github.com/usegro/services/billing/internal/apps/invoice/services"
)

func RegisterRoute(r *fiber.App) {
	database.SetUpPostgres()
	db := database.PostgressInstance1

	err := database.InitRedisClient(config.GetConfig().Redis)
	if err != nil {
		return
	}

	cfg := config.GetConfig()
	emailCfg := invoiceServices.EmailCfg{
		Env:       cfg.Env,
		Region:    cfg.Email.Region,
		FromEmail: cfg.Email.FromEmail,
	}

	r.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	invoiceRoutes.InvoiceRouter(v1, db, emailCfg)

	r.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Not found"})
	})
}
