package router

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/sdq-codes/usegro-api/config"
	_ "github.com/sdq-codes/usegro-api/docs"
	httpInterface "github.com/sdq-codes/usegro-api/internal/interface/http"
	httpError "github.com/sdq-codes/usegro-api/internal/interface/http/error"
	"github.com/sdq-codes/usegro-api/internal/router/middleware"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func NewFiberRouter() *fiber.App {
	r := fiber.New(fiber.Config{
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
		ErrorHandler:          httpError.ErrorHandler,
	})

	// Set up global middleware
	r.Use(cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return origin == config.GetConfig().FrontEnd.Url
		},
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization, X-CRM-ID",
		ExposeHeaders:    "Authorization, X-Request-Id, X-CRM-ID",
		AllowMethods:     "GET, POST, PUT, PATCH, DELETE, OPTIONS",
		AllowCredentials: true,
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
