package routes

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	crm "github.com/usegro/services/crm/internal/apps/crm/controllers"
	crmsocials "github.com/usegro/services/crm/internal/apps/crm/controllers/socials"
	crmService "github.com/usegro/services/crm/internal/apps/crm/services"
	"github.com/usegro/services/crm/internal/router/middleware"
	"gorm.io/gorm"
)

func CrmRouter(v1 fiber.Router, db *gorm.DB, rdb *redis.Client, dynamodbForms *dynamodb.Client) {
	crmUserOrganizationAPIGroup := v1.Group("/crm")

	// Tags — nested under /crm/tags; CRM ID supplied via X-CRM-ID header
	crmTagsAPIGroup := crmUserOrganizationAPIGroup.Group("/tags")
	crmTagController := crm.NewTagController(crmService.NewTagService(dynamodbForms))

	crmTagsAPIGroup.Post("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmTagController.CreateTag)
	crmTagsAPIGroup.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmTagController.ListTags)
	crmTagsAPIGroup.Get("/:tagId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmTagController.GetTag)
	crmTagsAPIGroup.Patch("/:tagId/name", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmTagController.UpdateTagName)
	crmTagsAPIGroup.Patch("/:tagId/status", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmTagController.UpdateTagStatus)

	crmCustomerController := crm.NewCRMCustomerController(crmService.NewCrmCustomerService(db, rdb, dynamodbForms))
	crmUserOrganizationCustomerAPIGroup := crmUserOrganizationAPIGroup.Group("/customers")
	crmUserOrganizationCustomerAPIGroup.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmCustomerController.FetchCrmCustomers)
	crmUserOrganizationCustomerAPIGroup.Delete("/:formId/:submissionId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmCustomerController.ArchiveCrmCustomer)
	crmUserOrganizationCustomerAPIGroup.Get("/:formId/:submissionId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmCustomerController.GetCrmCustomer)

	crmFormsAPIGroup := crmUserOrganizationAPIGroup.Group("/forms")
	crmUserOrganization := crm.NewCRMUserOrganizationController(crmService.NewCRMUserOrganizationService(db, rdb))
	// Static routes first
	crmUserOrganizationAPIGroup.Post("/business-name/exist", crmUserOrganization.CheckBusinessNameExist)
	crmUserOrganizationAPIGroup.Post("/", middleware.JwtVerify(), crmUserOrganization.CreateCRMUserOrganization)
	crmUserOrganizationAPIGroup.Get("/", middleware.JwtVerify(), crmUserOrganization.FetchCRMUserOrganization)
	crmUserOrganizationAPIGroup.Patch("/status", middleware.JwtVerify(), crmUserOrganization.ToggleCRMUserOrganizationStatus)

	// Dynamic routes after
	crmUserOrganizationAPIGroup.Get("/:id", middleware.JwtVerify(), crmUserOrganization.GetCRMUserOrganization)
	crmUserOrganizationAPIGroup.Patch("/:id", middleware.JwtVerify(), crmUserOrganization.UpdateCRMUserOrganization)
	crmUserOrganizationAPIGroup.Post("/:id/sales-channels", middleware.JwtVerify(), crmUserOrganization.CreateSalesChannelType)
	crmUserOrganizationAPIGroup.Patch("/:id/sales-channels", middleware.JwtVerify(), crmUserOrganization.UpdateSalesChannelType)
	crmUserOrganizationAPIGroup.Post("/:id/stock-products", middleware.JwtVerify(), crmUserOrganization.CreateStockProductType)
	crmUserOrganizationAPIGroup.Patch("/:id/stock-products", middleware.JwtVerify(), crmUserOrganization.UpdateStockProductType)

	// CRM Socials
	crmSocialsAPIGroup := crmUserOrganizationAPIGroup.Group("/socials")
	crmSocialsInstagram := crmsocials.NewCRMSocialsInstagramController()
	crmSocialsAPIGroup.Get("/instagram", middleware.JwtVerify(), crmSocialsInstagram.FacebookLogin)

	crmCustomerFormsAPIGroup := crmFormsAPIGroup.Group("/customers")
	crmCustomerFormsAPIGroup.Get("/create", middleware.JwtVerify(), middleware.CRMContextMiddleware(), crmCustomerController.FetchPublishedCreateCustomerForm)
}
