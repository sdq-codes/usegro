package catalogRoutes

import (
	"github.com/gofiber/fiber/v2"
	catalogControllers "github.com/usegro/services/catalog/internal/apps/catalog/controllers"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/router/middleware"
	"gorm.io/gorm"
)

func CatalogRouter(v1 fiber.Router, db *gorm.DB) {
	group := v1.Group("/catalog")

	productController := catalogControllers.NewProductController(catalogServices.NewProductService(db))
	serviceController := catalogControllers.NewServiceController(catalogServices.NewCatalogServiceService(db))
	planController := catalogControllers.NewPlanController(catalogServices.NewPlanService(db))
	tagController := catalogControllers.NewTagController(catalogServices.NewTagService(db))
	categoryController := catalogControllers.NewCategoryController(catalogServices.NewCategoryService(db))
	standardCategoryController := catalogControllers.NewStandardCategoryController(catalogServices.NewStandardCategoryService(db))
	mediaController := catalogControllers.NewMediaController(db)

	// Products
	products := group.Group("/products")
	products.Post("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), productController.CreateProduct)
	products.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), productController.ListProducts)
	products.Get("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), productController.GetProduct)
	products.Patch("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), productController.UpdateProduct)
	products.Delete("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), productController.DeleteProduct)

	// Services
	services := group.Group("/services")
	services.Post("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), serviceController.CreateService)
	services.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), serviceController.ListServices)
	services.Get("/:serviceId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), serviceController.GetService)
	services.Patch("/:serviceId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), serviceController.UpdateService)
	services.Delete("/:serviceId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), serviceController.DeleteService)

	// Plans (nested under services)
	services.Post("/:serviceId/plans", middleware.JwtVerify(), middleware.CRMContextMiddleware(), planController.CreatePlan)
	services.Get("/:serviceId/plans", middleware.JwtVerify(), middleware.CRMContextMiddleware(), planController.ListPlans)
	services.Patch("/:serviceId/plans/:planId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), planController.UpdatePlan)
	services.Delete("/:serviceId/plans/:planId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), planController.DeletePlan)

	// Tags
	tags := group.Group("/tags")
	tags.Post("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), tagController.CreateTag)
	tags.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), tagController.ListTags)
	tags.Delete("/:tagId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), tagController.DeleteTag)

	// Categories
	categories := group.Group("/categories")
	categories.Post("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), categoryController.CreateCategory)
	categories.Get("", middleware.JwtVerify(), middleware.CRMContextMiddleware(), categoryController.ListCategories)
	categories.Patch("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), categoryController.UpdateCategory)
	categories.Delete("/:id", middleware.JwtVerify(), middleware.CRMContextMiddleware(), categoryController.DeleteCategory)

	// Standard Categories (read-only, no CRM scope)
	stdCategories := group.Group("/standard-categories")
	stdCategories.Get("", standardCategoryController.ListRootCategories)
	stdCategories.Get("/search", standardCategoryController.SearchCategories)
	stdCategories.Get("/:id", standardCategoryController.GetCategory)
	stdCategories.Get("/:id/children", standardCategoryController.ListChildren)

	// Media
	media := group.Group("/media")
	media.Post("/upload", middleware.JwtVerify(), middleware.CRMContextMiddleware(), mediaController.UploadMedia)
	media.Post("/upload-temp", middleware.JwtVerify(), middleware.CRMContextMiddleware(), mediaController.UploadTemp)
	media.Post("/presign", middleware.JwtVerify(), middleware.CRMContextMiddleware(), mediaController.PresignUpload)
	media.Delete("/:mediaId", middleware.JwtVerify(), middleware.CRMContextMiddleware(), mediaController.DeleteMedia)
}
