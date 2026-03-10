package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/crm/internal/apps/crm/services"
	"github.com/usegro/services/crm/internal/interface/response"
)

type CRMCustomerController struct {
	service *services.CrmCustomerService
}

func NewCRMCustomerController(service *services.CrmCustomerService) *CRMCustomerController {
	return &CRMCustomerController{
		service: service,
	}
}

// FetchPublishedCreateCustomerForm godoc
// @Summary      Get the published create-customer form
// @Description  Fetches the published "create_customer" form for the CRM. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Customers
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        X-CRM-ID       header  string  true  "CRM ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/crm/forms/customers/create [get]
func (ccc *CRMCustomerController) FetchPublishedCreateCustomerForm(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)

	form, err := ccc.service.FetchPublishedCreateCustomerForm(c.Context(), crmId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Create Customer form retrieved successfully",
		Data:            form,
	})
}

// FetchCrmCustomers godoc
// @Summary      List all CRM customers
// @Description  Returns all customers belonging to the CRM. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Customers
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        X-CRM-ID       header  string  true  "CRM ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/crm/customers [get]
func (ccc *CRMCustomerController) FetchCrmCustomers(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)

	customers, err := ccc.service.FetchCrmCustomers(c.Context(), crmId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "CRM customers retrieved successfully",
		Data:            customers,
	})
}

// ArchiveCrmCustomer godoc
// @Summary      Archive a CRM customer
// @Description  Archives a customer submission. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Customers
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        X-CRM-ID       header  string  true  "CRM ID"
// @Param        formId         path    string  true  "Form ID"
// @Param        submissionId   path    string  true  "Submission ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/crm/customers/{formId}/{submissionId} [delete]
func (ccc *CRMCustomerController) ArchiveCrmCustomer(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)
	submissionId := c.Params("submissionId")
	formId := c.Params("formId")

	if submissionId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: "Submission ID is required",
		})
	}

	err := ccc.service.ArchiveCrmCustomer(c.Context(), submissionId, formId, crmId)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "CRM customer archived successfully",
	})
}

// GetCrmCustomer godoc
// @Summary      Get a CRM customer
// @Description  Fetches a single customer by form and submission ID. CRM ID is supplied via the X-CRM-ID header.
// @Tags         Customers
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Param        X-CRM-ID       header  string  true  "CRM ID"
// @Param        formId         path    string  true  "Form ID"
// @Param        submissionId   path    string  true  "Submission ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/crm/customers/{formId}/{submissionId} [get]
func (ccc *CRMCustomerController) GetCrmCustomer(c *fiber.Ctx) error {
	crmId := c.Locals("crmID").(string)
	submissionId := c.Params("submissionId")
	formId := c.Params("formId")

	if submissionId == "" || formId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: "Form ID and Submission ID are required",
		})
	}

	customer, err := ccc.service.GetCrmCustomer(c.Context(), submissionId, formId, crmId)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "CRM customer fetched successfully",
		Data:            customer,
	})
}
