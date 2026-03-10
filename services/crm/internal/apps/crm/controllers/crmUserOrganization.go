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

type CRMUserOrganizationController struct {
	service *services.CRMUserOrganizationService
}

func NewCRMUserOrganizationController(service *services.CRMUserOrganizationService) *CRMUserOrganizationController {
	return &CRMUserOrganizationController{
		service: service,
	}
}

// CheckBusinessNameExist godoc
// @Summary      Check if business name exists
// @Description  Returns whether a business name is already taken
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.BusinessNameExistDTI  true  "Business Name"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/crm/business-name/exist [post]
func (ctl *CRMUserOrganizationController) CheckBusinessNameExist(c *fiber.Ctx) error {
	var req dto.BusinessNameExistDTI
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	exists, err := ctl.service.IsBusinessNameExist(c.Context(), req.BusinessName)
	if err != nil {
		return err
	}
	if exists {
		return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_FETCHED,
			ResponseMessage: "Business name already taken",
			Data:            map[string]bool{"exists": true},
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Business name is available",
		Data:            map[string]bool{"exists": false},
	})
}

// CreateCRMUserOrganization godoc
// @Summary      Create a CRM user organization
// @Description  Creates a new CRM user organization record
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                      true  "Bearer access token"
// @Param        payload        body      dto.CrmUserOrganizationDTI  true  "CRM User Organization Data"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/crm [post]
func (ctl *CRMUserOrganizationController) CreateCRMUserOrganization(c *fiber.Ctx) error {
	var req dto.CrmUserOrganizationDTI

	authUser, err := auth.AuthUser(c)
	if err != nil {
		return exception.UnauthorizedError
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.CRMUserOrganizationCreateValidation(req); err != nil {
		return err
	}

	crmUserOrganization, err := ctl.service.CreateCRMUserOrganization(c.Context(), &req, authUser.User.ID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "CRM organization created successfully",
		Data: dto.CrmUserOrganizationDTO{
			ID:           crmUserOrganization.ID,
			FullName:     crmUserOrganization.FullName,
			BusinessName: crmUserOrganization.BusinessName,
			BusinessInfo: crmUserOrganization.BusinessInfo,
			CreatedAt:    crmUserOrganization.CreatedAt,
			UpdatedAt:    crmUserOrganization.UpdatedAt,
		},
	})
}

// UpdateCRMUserOrganization godoc
// @Summary      Update a CRM user organization
// @Description  Updates an existing CRM user organization record by ID
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        id       path      string                    true  "CRM User Organization ID"
// @Param        payload  body      dto.CrmUserOrganizationDTI  true  "CRM User Organization Data"
// @Success 200 {object} response.CommonResponse
// @Failure 400 {object} response.CommonResponse
// @Router       /api/v1/crm/{id} [patch]
func (ctl *CRMUserOrganizationController) UpdateCRMUserOrganization(c *fiber.Ctx) error {
	var req dto.CrmUserOrganizationDTI

	orgIdParam := c.Params("id")
	orgId, err := uuid.Parse(orgIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.CRMUserOrganizationUpdateValidation(req); err != nil {
		return err
	}

	crmUserOrganization, err := ctl.service.UpdateCRMUserOrganization(c.Context(), &req, orgId)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "CRM organization updated successfully",
		Data: dto.CrmUserOrganizationDTO{
			ID:           crmUserOrganization.ID,
			FullName:     crmUserOrganization.FullName,
			BusinessName: crmUserOrganization.BusinessName,
			BusinessInfo: crmUserOrganization.BusinessInfo,
			CreatedAt:    crmUserOrganization.CreatedAt,
			UpdatedAt:    crmUserOrganization.UpdatedAt,
		},
	})
}

// GetCRMUserOrganization godoc
// @Summary      Get a CRM user organization
// @Description  Retrieves a CRM user organization by ID
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "CRM User Organization ID"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/crm/{id} [get]
func (ctl *CRMUserOrganizationController) GetCRMUserOrganization(c *fiber.Ctx) error {
	crmUserOrganizationIdParam := c.Params("id")
	_, err := uuid.Parse(crmUserOrganizationIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	org, err := ctl.service.GetCRMUserOrganization(c.Context(), crmUserOrganizationIdParam)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "CRM organization retrieved successfully",
		Data: dto.CrmUserOrganizationDTO{
			ID:           org.ID,
			FullName:     org.FullName,
			BusinessName: org.BusinessName,
			BusinessInfo: org.BusinessInfo,
			CreatedAt:    org.CreatedAt,
			UpdatedAt:    org.UpdatedAt,
		},
	})
}

// FetchCRMUserOrganization godoc
// @Summary      List all CRM organizations for the authenticated user
// @Description  Retrieves all CRM user organizations belonging to the authenticated user
// @Tags         CRM
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/crm [get]
func (ctl *CRMUserOrganizationController) FetchCRMUserOrganization(c *fiber.Ctx) error {
	authUser, err := auth.AuthUser(c)
	if err != nil {
		return exception.UnauthorizedError
	}

	crmOrganizations, err := ctl.service.FetchCRMUserOrganization(c.Context(), authUser.User.ID.String())
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "CRM organizations retrieved successfully",
		Data:            crmOrganizations,
	})
}

// ToggleCRMUserOrganizationStatus godoc
// @Summary      Toggle CRM user organization status
// @Description  Enables or disables the authenticated user's CRM organization
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      404  {object}  response.CommonResponse
// @Router       /api/v1/crm/status [patch]
func (ctl *CRMUserOrganizationController) ToggleCRMUserOrganizationStatus(c *fiber.Ctx) error {
	crmUserOrganizationIdParam := c.Params("id")
	_, err := uuid.Parse(crmUserOrganizationIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	err = ctl.service.ToggleCRMUserOrganizationStatus(c.Context(), crmUserOrganizationIdParam)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "CRM organization status toggled successfully",
	})
}

// CreateSalesChannelType godoc
// @Summary      Create Sales Channel Types
// @Description  Creates new sales channel types for a CRM user organization
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "CRM User Organization ID"
// @Param        payload  body      []models.SalesChannel true  "Sales Channel Data"
// @Success      201 {object} response.CommonResponse
// @Failure      400 {object} response.CommonResponse
// @Router       /api/v1/crm/{id}/sales-channels [post]
func (ctl *CRMUserOrganizationController) CreateSalesChannelType(c *fiber.Ctx) error {
	crmUserOrganizationIdParam := c.Params("id")
	crmUserOrganizationIdParamUuid, err := uuid.Parse(crmUserOrganizationIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	var req dto.CrmUserOrganizationSalesChannelTypeDTI
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.CRMUserOrganizationSalesChannelTypeValidation(req); err != nil {
		return err
	}

	created, err := ctl.service.CreateSalesChannelType(c.Context(), crmUserOrganizationIdParamUuid, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Sales channel types created successfully",
		Data:            created,
	})
}

// UpdateSalesChannelType godoc
// @Summary      Update Sales Channel Types
// @Description  Updates sales channel types for a CRM user organization
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "CRM User Organization ID"
// @Param        payload  body      []models.SalesChannel true  "Sales Channel Data"
// @Success      200 {object} response.CommonResponse
// @Failure      400 {object} response.CommonResponse
// @Router       /api/v1/crm/{id}/sales-channels [patch]
func (ctl *CRMUserOrganizationController) UpdateSalesChannelType(c *fiber.Ctx) error {
	var req dto.CrmUserOrganizationSalesChannelTypeDTI
	crmUserOrganizationIdParam := c.Params("id")
	crmUserOrganizationIdParamUuid, err := uuid.Parse(crmUserOrganizationIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	if err := validation.CRMUserOrganizationSalesChannelTypeValidation(req); err != nil {
		return err
	}

	updated, err := ctl.service.UpdateSalesChannelType(c.Context(), crmUserOrganizationIdParamUuid, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Sales channel types updated successfully",
		Data:            updated,
	})
}

// CreateStockProductType godoc
// @Summary      Create Stock Product Types
// @Description  Creates new stock product types for a CRM user organization
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "CRM User Organization ID"
// @Param        payload  body      dto.CrmUserOrganizationStockProductTypeDTI true  "Stock Product Type Data"
// @Success      201 {object} response.CommonResponse
// @Failure      400 {object} response.CommonResponse
// @Router       /api/v1/crm/{id}/stock-products [post]
func (ctl *CRMUserOrganizationController) CreateStockProductType(c *fiber.Ctx) error {
	var req dto.CrmUserOrganizationStockProductTypeDTI
	crmUserOrganizationIdParam := c.Params("id")
	crmUserOrganizationIdParamUuid, err := uuid.Parse(crmUserOrganizationIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.CRMUserOrganizationStockProductTypeValidation(req); err != nil {
		return err
	}

	created, err := ctl.service.CreateStockProductType(c.Context(), crmUserOrganizationIdParamUuid, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Stock product types created successfully",
		Data:            created,
	})
}

// UpdateStockProductType godoc
// @Summary      Update Stock Product Types
// @Description  Updates stock product types for a CRM user organization
// @Tags         CRM
// @Accept       json
// @Produce      json
// @Param        id       path      string                        true  "CRM User Organization ID"
// @Param        payload  body      dto.CrmUserOrganizationStockProductTypeDTI true  "Stock Product Type Data"
// @Success      200 {object} response.CommonResponse
// @Failure      400 {object} response.CommonResponse
// @Router       /api/v1/crm/{id}/stock-products [patch]
func (ctl *CRMUserOrganizationController) UpdateStockProductType(c *fiber.Ctx) error {
	var req dto.CrmUserOrganizationStockProductTypeDTI
	crmUserOrganizationIdParam := c.Params("id")
	crmUserOrganizationIdParamUuid, err := uuid.Parse(crmUserOrganizationIdParam)
	if err != nil {
		return exception.InvalidUUIDError
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	if err := validation.CRMUserOrganizationStockProductTypeValidation(req); err != nil {
		return err
	}

	updated, err := ctl.service.UpdateStockProductType(c.Context(), crmUserOrganizationIdParamUuid, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Stock product types updated successfully",
		Data:            updated,
	})
}
