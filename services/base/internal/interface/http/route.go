package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/database"
	"github.com/sdq-codes/usegro-api/database/mongodb"
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

	mongoCfg := config.GetConfig().MongoDB
	err = mongodb.InitMongoClient(mongoCfg.URI, mongoCfg.Database)
	if err != nil {
		return
	}
	mongoDB := mongodb.GetMongoDatabase()

	api := r.Group("/api")
	v1 := api.Group("/v1")
	base := v1.Group("/base")

	// Health API
	routes.Health(r)
	routes.BaseRouter(base, db, rdb)
	formRouter.FormsRouter(base, db, mongoDB)

	// Error Case Handler
	miscellaneousHandler := httpMiscellaneous.NewMiscellaneousHTTPHandler()
	r.All("*", miscellaneousHandler.NotFound)

	workers.NewWorkerService(rdb, db, []string{database.AddQueuePrefix("crm_queue_email")}).Start()
}
