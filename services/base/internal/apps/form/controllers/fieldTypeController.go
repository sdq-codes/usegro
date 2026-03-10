package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/apps/form/services"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	"net/url"
)

type FieldTypeController struct {
	service services.FieldTypeService
}

func NewFieldTypeController(service services.FieldTypeService) *FieldTypeController {
	return &FieldTypeController{
		service: service,
	}
}

// GetAllFieldTypes godoc
// @Summary      List all field types
// @Description  Returns all field types with their configs and validation rules
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/forms/field/types [get]
func (ctl *FieldTypeController) GetAllFieldTypes(c *fiber.Ctx) error {
	fieldTypes, err := ctl.service.GetAll(c.Context())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Field types retrieved successfully",
		Data:            fieldTypes,
	})
}

// GetFieldTypeByName godoc
// @Summary      Get a field type by name
// @Description  Returns a single field type along with its configs and validation rules
// @Tags         Forms
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        name           path    string  true  "Field type name (URL-encoded)"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/forms/field/types/{name} [get]
func (ctl *FieldTypeController) GetFieldTypeByName(c *fiber.Ctx) error {
	name := c.Params("name")
	if name == "" {
		return exception.InvalidRequestBodyError
	}
	decoded, err := url.QueryUnescape(name)
	if err != nil {
		panic(err)
	}

	fieldType, err := ctl.service.GetByName(c.Context(), decoded)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Field type retrieved successfully",
		Data:            fieldType,
	})
}
