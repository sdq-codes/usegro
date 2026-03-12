package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/services"
	"github.com/sdq-codes/usegro-api/internal/helper/auth"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
	"github.com/sdq-codes/usegro-api/pkg/exception"
)

type FormSubmissionController struct {
	service *services.FormSubmissionService
}

func NewFormSubmissionController(service services.FormSubmissionService) *FormSubmissionController {
	return &FormSubmissionController{
		service: &service,
	}
}

// CreateSubmission godoc
// @Summary      Submit form answers
// @Description  Creates a submission for a published form version. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                      true  "Bearer access token"
// @Param        X-CRM-ID       header  string                      true  "CRM ID"
// @Param        formId         path    string                      true  "Form ID"
// @Param        versionId      path    string                      true  "Form Version ID"
// @Param        request        body    dtos.CreateSubmissionInput  true  "Submission answers"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/version/{versionId}/submission [post]
func (ctl *FormSubmissionController) CreateSubmission(c *fiber.Ctx) error {
	formID := c.Params("formID")
	versionID := c.Params("versionID")
	crmId := c.Locals("crmID").(string)

	var userID string
	if claims, err := auth.AuthUser(c); err == nil {
		userID = claims.User.ID.String()
	}

	var req dtos.CreateSubmissionInput
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := ctl.service.CreateSubmission(c.Context(), formID, versionID, req, crmId, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Form submission created successfully",
	})
}
