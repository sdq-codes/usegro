package catalogControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/interface/response"
)

type ServiceController struct {
	service *catalogServices.CatalogServiceService
}

func NewServiceController(service *catalogServices.CatalogServiceService) *ServiceController {
	return &ServiceController{service: service}
}

func (sc *ServiceController) CreateService(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	var body dto.CreateServiceDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}
	if body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Service name is required",
		})
	}
	item, err := sc.service.CreateService(c.Context(), crmID, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Service created successfully",
		Data:            item,
	})
}

func (sc *ServiceController) ListServices(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	items, err := sc.service.ListServices(c.Context(), crmID, c.Query("search"), c.Query("status"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Services retrieved successfully",
		Data:            items,
	})
}

func (sc *ServiceController) GetService(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	item, err := sc.service.GetService(c.Context(), crmID, c.Params("serviceId"))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Service retrieved successfully",
		Data:            item,
	})
}

func (sc *ServiceController) UpdateService(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	var body dto.UpdateServiceDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}
	item, err := sc.service.UpdateService(c.Context(), crmID, c.Params("serviceId"), body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Service updated successfully",
		Data:            item,
	})
}

func (sc *ServiceController) DeleteService(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	if err := sc.service.DeleteService(c.Context(), crmID, c.Params("serviceId")); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Service deleted successfully",
	})
}
