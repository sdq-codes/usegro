package routes

import (
	"github.com/gofiber/fiber/v2"
	healthController "github.com/sdq-codes/usegro-api/internal/apps/base/controllers/health"
)

func Health(api fiber.Router) {
	healthAPI := api.Group("/health")
	healthHandler := healthController.NewHealthController()
	healthAPI.Get("/", healthHandler.Show)
}
