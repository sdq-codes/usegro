package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/database"
	"github.com/usegro/services/crm/database/dynamo"
	crmRouter "github.com/usegro/services/crm/internal/apps/crm/routes"
)

func RegisterRoute(r *fiber.App) {
	database.SetUpPostgres()
	db := database.PostgressInstance1

	err := database.InitRedisClient(config.GetConfig().Redis)
	if err != nil {
		return
	}
	rdb := database.SingleRdb

	dynamoCfg := config.GetConfig().DynamodbForms
	err = dynamo.InitDynamoFormClient(dynamoCfg.DynamoEndpoint, dynamoCfg.AwsRegion)
	if err != nil {
		return
	}
	dynamodbForms := dynamo.DynamoClient

	r.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	api := r.Group("/api")
	v1 := api.Group("/v1")

	crmRouter.CrmRouter(v1, db, rdb, dynamodbForms)

	r.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Not found"})
	})

}
