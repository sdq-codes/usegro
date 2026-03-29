package catalogControllers

import (
	"github.com/gofiber/fiber/v2"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/interface/response"
)

type TagController struct {
	service *catalogServices.TagService
}

func NewTagController(service *catalogServices.TagService) *TagController {
	return &TagController{service: service}
}

func (tc *TagController) CreateTag(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	var body struct {
		Name string `json:"name"`
	}
	if err := c.BodyParser(&body); err != nil || body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Tag name is required",
		})
	}
	tag, err := tc.service.CreateTag(c.Context(), crmID, body.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Tag created successfully",
		Data:            tag,
	})
}

func (tc *TagController) ListTags(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	tags, err := tc.service.ListTags(c.Context(), crmID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Tags retrieved successfully",
		Data:            tags,
	})
}

func (tc *TagController) DeleteTag(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	if err := tc.service.DeleteTag(c.Context(), crmID, c.Params("tagId")); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Tag deleted successfully",
	})
}
