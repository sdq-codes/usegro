package router

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	fiberSwagger "github.com/swaggo/fiber-swagger"
	_ "github.com/usegro/services/catalog/docs"
	httpInterface "github.com/usegro/services/catalog/internal/interface/http"
	httpError "github.com/usegro/services/catalog/internal/interface/http/error"
	"github.com/usegro/services/catalog/internal/router/middleware"
)

func NewFiberRouter() *fiber.App {
	r := fiber.New(fiber.Config{
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
		ErrorHandler:          httpError.ErrorHandler,
		BodyLimit:             25 * 1024 * 1024, // 25 MB
	})

	// Set up global middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:  "*", // or restrict to your frontend domain
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, X-CRM-ID",
		ExposeHeaders: "Authorization, X-Request-Id, X-CRM-ID", // <- allow frontend to READ Authorization header
	}))
	r.Use(requestid.New())
	r.Use(recover.New())
	r.Use(idempotency.New())
	// r.Use(cache.New())
	r.Use(middleware.Logger())
	r.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))
	r.Use(middleware.EnhanceSentryEvent())

	r.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Register routes (handlers)
	httpInterface.RegisterRoute(r)

	return r
}
