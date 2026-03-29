package catalogControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/interface/response"
)

type ProductController struct {
	service *catalogServices.ProductService
}

func NewProductController(service *catalogServices.ProductService) *ProductController {
	return &ProductController{service: service}
}

func (pc *ProductController) CreateProduct(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)

	var body dto.CreateProductDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}

	if body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Product name is required",
		})
	}

	product, err := pc.service.CreateProduct(c.Context(), crmID, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Product created successfully",
		Data:            product,
	})
}

func (pc *ProductController) ListProducts(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	search := c.Query("search")
	status := c.Query("status")
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	products, err := pc.service.ListProducts(c.Context(), crmID, search, status, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Products retrieved successfully",
		Data:            products,
	})
}

func (pc *ProductController) GetProduct(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Product ID is required",
		})
	}

	product, err := pc.service.GetProduct(c.Context(), crmID, productID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Product retrieved successfully",
		Data:            product,
	})
}

func (pc *ProductController) UpdateProduct(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Product ID is required",
		})
	}

	var body dto.UpdateProductDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}

	product, err := pc.service.UpdateProduct(c.Context(), crmID, productID, body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Product updated successfully",
		Data:            product,
	})
}

func (pc *ProductController) DeleteProduct(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	productID := c.Params("id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Product ID is required",
		})
	}

	if err := pc.service.DeleteProduct(c.Context(), crmID, productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Product deleted successfully",
	})
}
