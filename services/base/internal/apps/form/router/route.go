package router

import (
	"github.com/gofiber/fiber/v2"
	formController "github.com/sdq-codes/usegro-api/internal/apps/form/controllers"
	"github.com/sdq-codes/usegro-api/internal/apps/form/repositories"
	formService "github.com/sdq-codes/usegro-api/internal/apps/form/services"
	"github.com/sdq-codes/usegro-api/internal/router/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func FormsRouter(v1 fiber.Router, db *gorm.DB, mongoDB *mongo.Database) {
	formAPIGroup := v1.Group("/forms")
	formFieldAPIGroup := formAPIGroup.Group("/field")

	formFieldTypeAPIGroup := formFieldAPIGroup.Group("/types")
	formFieldTypesController := formController.NewFieldTypeController(formService.NewFieldTypeService(db))
	formFieldTypeAPIGroup.Get("", middleware.JwtVerify(), formFieldTypesController.GetAllFieldTypes)
	formFieldTypeAPIGroup.Get("/:name", middleware.JwtVerify(), formFieldTypesController.GetFieldTypeByName)

	formRepo := repositories.NewFormRepository(mongoDB)
	submissionRepo := repositories.NewFormSubmissionRepository(mongoDB)
	activityRepo := repositories.NewCustomerActivityRepository(mongoDB)

	formHandlerController := formController.NewFormController(formService.NewFormService(formRepo))
	formAPIGroup.Post("", middleware.CRMContextMiddleware(), middleware.JwtVerify(), formHandlerController.CreateForm)
	formAPIGroup.Get("/:formId", middleware.JwtVerify(), formHandlerController.FetchForm)
	formAPIGroup.Get("/:formId/draft", middleware.JwtVerify(), formHandlerController.FetchDraftForm)
	formAPIGroup.Post("/:formId/publish", middleware.JwtVerify(), formHandlerController.PublishFormVersion)
	formAPIGroup.Post("/:formId/version", middleware.JwtVerify(), formHandlerController.CreateVersion)
	formAPIGroup.Post("/:formId/fields", middleware.JwtVerify(), formHandlerController.CreateFormVersionFields)
	formAPIGroup.Delete("/:formId/fields/:fieldId", middleware.JwtVerify(), formHandlerController.DeleteFormVersionField)
	formAPIGroup.Get("/:formId/version/:versionId", middleware.JwtVerify(), formHandlerController.FetchFormVersion)
	formAPIGroup.Patch("/:formId/version/:versionId/field/:fieldId", middleware.JwtVerify(), formHandlerController.UpdateFormVersionField)

	formSubmissionController := formController.NewFormSubmissionController(
		formService.NewFormSubmissionService(formRepo, submissionRepo, activityRepo),
	)
	formAPIGroup.Post("/:formId/version/:versionId/submission", middleware.JwtVerify(), middleware.CRMContextMiddleware(), formSubmissionController.CreateSubmission)
	formAPIGroup.Patch("/:formId/submission/:submissionId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), formSubmissionController.UpdateSubmission)
}
