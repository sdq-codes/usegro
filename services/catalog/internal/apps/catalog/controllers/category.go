package catalogControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/interface/response"
)

type CategoryController struct {
	service *catalogServices.CategoryService
}

func NewCategoryController(service *catalogServices.CategoryService) *CategoryController {
	return &CategoryController{service: service}
}

func (cc *CategoryController) CreateCategory(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	var body dto.CreateCategoryDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}

	if body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Category name is required",
		})
	}

	category, err := cc.service.CreateCategory(c.Context(), crmID, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Category created successfully",
		Data:            category,
	})
}

func (cc *CategoryController) ListCategories(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	categories, err := cc.service.ListCategories(c.Context(), crmID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Categories retrieved successfully",
		Data:            categories,
	})
}

func (cc *CategoryController) UpdateCategory(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Category ID is required",
		})
	}

	var body dto.UpdateCategoryDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}

	category, err := cc.service.UpdateCategory(c.Context(), crmID, categoryID, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Category updated successfully",
		Data:            category,
	})
}

func (cc *CategoryController) DeleteCategory(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	categoryID := c.Params("id")

	if categoryID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Category ID is required",
		})
	}

	if err := cc.service.DeleteCategory(c.Context(), crmID, categoryID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Category deleted successfully",
	})
}
