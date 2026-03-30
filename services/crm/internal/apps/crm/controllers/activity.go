package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/crm/internal/apps/crm/services"
	"github.com/usegro/services/crm/internal/interface/response"
)

type CustomerActivityController struct {
	service *services.CustomerActivityService
}

func NewCustomerActivityController(service *services.CustomerActivityService) *CustomerActivityController {
	return &CustomerActivityController{service: service}
}

// FetchCustomerActivity godoc
// @Summary      Get activity log for a customer
// @Description  Returns all activity events for the given customer (submission). CRM ID is supplied via the X-CRM-ID header.
// @Tags         Customers
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        X-CRM-ID       header  string  true  "CRM ID"
// @Param        submissionId   path    string  true  "Submission ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/crm/customers/{submissionId}/activity [get]
func (c *CustomerActivityController) FetchCustomerActivity(ctx *fiber.Ctx) error {
	submissionID := ctx.Params("submissionId")

	activities, err := c.service.FetchCustomerActivity(ctx.Context(), submissionID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Customer activity retrieved successfully",
		Data:            activities,
	})
}

type logCommentRequest struct {
	Comment string `json:"comment"`
}

// LogComment godoc
// @Summary      Post a comment on a customer
// @Tags         Customers
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        X-CRM-ID       header  string  true  "CRM ID"
// @Param        submissionId   path    string  true  "Submission ID"
// @Param        request        body    logCommentRequest  true  "Comment body"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/crm/customers/{submissionId}/activity [post]
func (c *CustomerActivityController) LogComment(ctx *fiber.Ctx) error {
	submissionID := ctx.Params("submissionId")
	crmID := ctx.Locals("crmID").(string)

	var req logCommentRequest
	if err := ctx.BodyParser(&req); err != nil || req.Comment == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: "comment is required",
		})
	}

	activity, err := c.service.LogComment(ctx.Context(), submissionID, crmID, req.Comment, "")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Comment posted successfully",
		Data:            activity,
	})
}
