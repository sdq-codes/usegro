package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/usegro/services/crm/internal/apps/crm/dto"
	"github.com/usegro/services/crm/internal/apps/crm/services"
	"github.com/usegro/services/crm/internal/apps/crm/validation"
	"github.com/usegro/services/crm/internal/helper/auth"
	"github.com/usegro/services/crm/internal/interface/response"
	"github.com/usegro/services/crm/pkg/exception"
)

type TagController struct {
	service *services.TagService
}

func NewTagController(service *services.TagService) *TagController {
	return &TagController{service: service}
}

// CreateTag godoc
// @Summary      Create a tag
// @Description  Creates a new tag for a CRM. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Param        X-CRM-ID  header    string            true  "CRM ID"
// @Param        payload   body      dto.TagCreateDTO  true  "Tag Data"
// @Success      201  {object} response.CommonResponse
// @Failure      400  {object} response.CommonResponse
// @Router       /api/v1/crm/tags [post]
func (ctl *TagController) CreateTag(c *fiber.Ctx) error {
	var req dto.TagCreateDTO

	authUser, err := auth.AuthUser(c)
	if err != nil {
		return exception.UnauthorizedError
	}

	crmId := c.Locals("crmID").(string)

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.CreateTagValidation(req); err != nil {
		return err
	}

	_, err = ctl.service.CreateTag(c.Context(), req, crmId, authUser.User.ID)
	if err != nil {
		return err
	}

	tags, err := ctl.service.ListTagsByCRM(c.Context(), crmId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Tags retrieved successfully",
		Data:            tags,
	})
}

// GetTag godoc
// @Summary      Get a tag
// @Description  Retrieves a tag by ID. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Param        X-CRM-ID  header    string  true  "CRM ID"
// @Param        tagId     path      string  true  "Tag ID"
// @Success      200  {object} response.CommonResponse
// @Failure      400  {object} response.CommonResponse
// @Router       /api/v1/crm/tags/{tagId} [get]
func (ctl *TagController) GetTag(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)
	tagID := c.Params("tagId")

	tag, err := ctl.service.FetchTag(c.Context(), crmId, tagID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Tag retrieved successfully",
		Data:            tag,
	})
}

// ListTags godoc
// @Summary      List CRM tags
// @Description  Retrieves all tags for a CRM. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Param        X-CRM-ID  header    string  true  "CRM ID"
// @Success      200  {object} response.CommonResponse
// @Failure      400  {object} response.CommonResponse
// @Router       /api/v1/crm/tags [get]
func (ctl *TagController) ListTags(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)

	tags, err := ctl.service.ListTagsByCRM(c.Context(), crmId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Tags retrieved successfully",
		Data:            tags,
	})
}

// UpdateTagName godoc
// @Summary      Update tag name
// @Description  Updates the name of a tag. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Tags
// @Accept       json
// @Produce      json
// @Param        X-CRM-ID  header    string           true  "CRM ID"
// @Param        tagId     path      string           true  "Tag ID"
// @Param        payload   body      dto.TagUpdateDTO true  "Tag Update Data"
// @Success      200  {object} response.CommonResponse
// @Failure      400  {object} response.CommonResponse
// @Router       /api/v1/crm/tags/{tagId}/name [patch]
func (ctl *TagController) UpdateTagName(c *fiber.Ctx) error {
	var req dto.TagUpdateDTO

	crmID := c.Locals("crmID").(string)
	tagID := c.Params("tagId")

	if _, err := uuid.Parse(tagID); err != nil {
		return exception.InvalidUUIDError
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.UpdateTagNameValidation(req); err != nil {
		return err
	}

	if err := ctl.service.UpdateTagName(c.Context(), crmID, tagID, req.Name); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Tag name updated successfully",
	})
}

// UpdateTagStatus godoc
// @Summary      Toggle tag status
// @Description  Toggles the active/inactive status of a tag. No request body required. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Tags
// @Produce      json
// @Param        X-CRM-ID  header  string  true  "CRM ID"
// @Param        tagId     path    string  true  "Tag ID"
// @Success      200  {object} response.CommonResponse
// @Failure      400  {object} response.CommonResponse
// @Router       /api/v1/crm/tags/{tagId}/status [patch]
func (ctl *TagController) UpdateTagStatus(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)
	tagID := c.Params("tagId")

	if _, err := uuid.Parse(tagID); err != nil {
		return exception.InvalidUUIDError
	}

	err := ctl.service.UpdateTagStatus(c.Context(), crmId, tagID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Tag status updated successfully",
	})
}
