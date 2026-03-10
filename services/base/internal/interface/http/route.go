package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/database"
	"github.com/sdq-codes/usegro-api/database/dynamo"
	httpMiscellaneous "github.com/sdq-codes/usegro-api/internal/apps/base/controllers/miscellaneous"
	workers "github.com/sdq-codes/usegro-api/internal/apps/base/controllers/queue"
	"github.com/sdq-codes/usegro-api/internal/apps/base/routes"
	formRouter "github.com/sdq-codes/usegro-api/internal/apps/form/router"
)

// ====================================================
// =================== DEFINE ROUTE ===================
// ====================================================

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

	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Health API
	routes.Health(r)
	routes.BaseRouter(v1, db, rdb)
	formRouter.FormsRouter(v1, db, dynamodbForms)

	// Error Case Handler
	miscellaneousHandler := httpMiscellaneous.NewMiscellaneousHTTPHandler()
	r.All("*", miscellaneousHandler.NotFound)

	workers.NewWorkerService(rdb, db, []string{database.AddQueuePrefix("crm_queue_email")}).Start()
}
