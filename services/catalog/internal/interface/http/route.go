package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/catalog/config"
	"github.com/usegro/services/catalog/database"
	catalogRouter "github.com/usegro/services/catalog/internal/apps/catalog/routes"
)

func RegisterRoute(r *fiber.App) {
	database.SetUpPostgres()
	db := database.PostgressInstance1

	err := database.InitRedisClient(config.GetConfig().Redis)
	if err != nil {
		return
	}

	r.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	catalogRouter.CatalogRouter(v1, db)

	r.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Not found"})
	})
}
