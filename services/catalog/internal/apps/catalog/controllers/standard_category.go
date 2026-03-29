package catalogControllers

import (
	"github.com/gofiber/fiber/v2"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/interface/response"
)

type StandardCategoryController struct {
	service *catalogServices.StandardCategoryService
}

func NewStandardCategoryController(service *catalogServices.StandardCategoryService) *StandardCategoryController {
	return &StandardCategoryController{service: service}
}

// ListRootCategories godoc
// @Summary List top-level standard categories
// @Tags standard-categories
// @Produce json
// @Success 200 {object} response.CommonResponse
// @Router /catalog/standard-categories [get]
func (cc *StandardCategoryController) ListRootCategories(c *fiber.Ctx) error {
	categories, err := cc.service.ListRootCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Categories retrieved successfully",
		Data:            categories,
	})
}

// GetCategory godoc
// @Summary Get a category with its immediate children
// @Tags standard-categories
// @Produce json
// @Param id path string true "Category UUID"
// @Success 200 {object} response.CommonResponse
// @Router /catalog/standard-categories/{id} [get]
func (cc *StandardCategoryController) GetCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Category ID is required",
		})
	}

	category, err := cc.service.GetCategoryWithChildren(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Category retrieved successfully",
		Data:            category,
	})
}

// SearchCategories godoc
// @Summary Search standard categories by name
// @Tags standard-categories
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {object} response.CommonResponse
// @Router /catalog/standard-categories/search [get]
func (cc *StandardCategoryController) SearchCategories(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Query parameter 'q' is required",
		})
	}

	results, err := cc.service.SearchCategories(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.INTERNAL_SERVER_ERROR,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Categories retrieved successfully",
		Data:            results,
	})
}

// ListChildren godoc
// @Summary List children of a category
// @Tags standard-categories
// @Produce json
// @Param id path string true "Parent Category UUID"
// @Success 200 {object} response.CommonResponse
// @Router /catalog/standard-categories/{id}/children [get]
func (cc *StandardCategoryController) ListChildren(c *fiber.Ctx) error {
	parentID := c.Params("id")
	if parentID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Parent category ID is required",
		})
	}

	categories, err := cc.service.ListByParent(c.Context(), parentID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Children retrieved successfully",
		Data:            categories,
	})
}
