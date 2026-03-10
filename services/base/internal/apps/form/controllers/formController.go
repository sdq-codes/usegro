package controllers

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/apps/form/dtos"
	"github.com/sdq-codes/usegro-api/internal/apps/form/services"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	"log"
)

type FormController struct {
	service services.FormService
}

func NewFormController(service services.FormService) *FormController {
	return &FormController{
		service: service,
	}
}

// CreateForm godoc
// @Summary      Create a new form
// @Description  Creates a new form with an initial draft version. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                       true  "Bearer access token"
// @Param        X-CRM-ID       header  string                       true  "CRM ID"
// @Param        request        body    dtos.CreateVersionRequestDTI true  "Form payload"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/forms [post]
func (ctl *FormController) CreateForm(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)

	var req dtos.CreateVersionRequestDTI
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := ctl.service.CreateForm(c.Context(), req, crmId); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Form version created successfully",
	})
}

// FetchForm godoc
// @Summary      Fetch latest published form
// @Description  Returns a form with its latest published version and all fields/logic
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        formId         path    string  true  "Form ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId} [get]
func (ctl *FormController) FetchForm(c *fiber.Ctx) error {
	formID := c.Params("formID")
	if formID == "" {
		return exception.InvalidRequestBodyError
	}

	form, err := ctl.service.FetchForm(c.Context(), formID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(form)
}

// FetchDraftForm godoc
// @Summary      Fetch latest draft form
// @Description  Returns a form with its latest draft version
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        formId         path    string  true  "Form ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/draft [get]
func (ctl *FormController) FetchDraftForm(c *fiber.Ctx) error {
	formID := c.Params("formID")
	if formID == "" {
		return exception.InvalidRequestBodyError
	}

	form, err := ctl.service.FetchDraftForm(c.Context(), formID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(form)
}

// PublishFormVersion godoc
// @Summary      Publish a form version
// @Description  Publishes the latest draft version of a form, making it publicly accessible
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        formId         path    string  true  "Form ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/publish [post]
func (ctl *FormController) PublishFormVersion(c *fiber.Ctx) error {
	formID := c.Params("formId")
	if formID == "" {
		return exception.InvalidRequestBodyError
	}

	err := ctl.service.PublishFormVersion(c.Context(), formID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: fmt.Sprintf("Form %s published successfully", formID),
	})
}

// CreateVersion godoc
// @Summary      Create a new form version
// @Description  Creates a new draft version for an existing form
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                       true  "Bearer access token"
// @Param        formId         path    string                       true  "Form ID"
// @Param        request        body    dtos.CreateVersionRequestDTI true  "Version payload"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/version [post]
func (ctl *FormController) CreateVersion(c *fiber.Ctx) error {
	formID := c.Params("formID")
	if formID == "" {
		return exception.InvalidRequestBodyError
	}

	var req dtos.CreateVersionRequestDTI
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := ctl.service.CreateVersion(c.Context(), req, formID); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Form version created successfully",
	})
}

// CreateFormVersionFields godoc
// @Summary      Add fields to a form
// @Description  Appends one or more fields to an existing draft form version
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                          true  "Bearer access token"
// @Param        formId         path    string                          true  "Form ID"
// @Param        request        body    dtos.CreateVersionFieldInputDTI true  "Field payload"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/fields [post]
func (ctl *FormController) CreateFormVersionFields(c *fiber.Ctx) error {
	formID := c.Params("formId")
	if formID == "" {
		return exception.InvalidRequestBodyError
	}

	var req dtos.CreateVersionFieldInputDTI
	if err := c.BodyParser(&req); err != nil {
		log.Println(err)
		return exception.InvalidRequestBodyError
	}

	if err := ctl.service.CreateFormVersionField(c.Context(), formID, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Form version fields created successfully",
	})
}

// DeleteFormVersionField godoc
// @Summary      Delete a field from a form
// @Description  Permanently removes a field from the given draft form version
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        formId         path    string  true  "Form ID"
// @Param        fieldId        path    string  true  "Field ID"
// @Success      204  "No Content"
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/fields/{fieldId} [delete]
func (ctl *FormController) DeleteFormVersionField(c *fiber.Ctx) error {
	formID := c.Params("formId")
	fieldID := c.Params("fieldId")
	if formID == "" || fieldID == "" {
		return exception.InvalidRequestBodyError
	}

	if err := ctl.service.DeleteFormVersionField(c.Context(), formID, fieldID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// FetchFormVersion godoc
// @Summary      Fetch a specific form version
// @Description  Returns a specific version of a form with its fields and logic
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        formId         path    string  true  "Form ID"
// @Param        versionId      path    string  true  "Version ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/version/{versionId} [get]
func (ctl *FormController) FetchFormVersion(c *fiber.Ctx) error {
	formID := c.Params("formId")
	versionID := c.Params("versionId")
	if formID == "" || versionID == "" {
		return exception.InvalidRequestBodyError
	}

	formVersion, err := ctl.service.FetchFormVersion(c.Context(), formID, versionID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(formVersion)
}

// UpdateFormVersionField godoc
// @Summary      Update a field in a form version
// @Description  Partially updates one or more properties of a specific field within a form version
// @Tags         Forms
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string                 true  "Bearer access token"
// @Param        formId         path    string                 true  "Form ID"
// @Param        versionId      path    string                 true  "Version ID"
// @Param        fieldId        path    string                 true  "Field ID"
// @Param        request        body    map[string]interface{} true  "Fields to update"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/{formId}/version/{versionId}/field/{fieldId} [patch]
func (ctl *FormController) UpdateFormVersionField(c *fiber.Ctx) error {
	formID := c.Params("formId")
	versionID := c.Params("versionId")
	fieldID := c.Params("fieldId")

	if formID == "" || versionID == "" || fieldID == "" {
		return exception.InvalidRequestBodyError
	}

	var req map[string]interface{}
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := ctl.service.UpdateFormField(c.Context(), formID, versionID, fieldID, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Form version field updated successfully",
	})
}
