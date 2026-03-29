package catalogControllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	catalogServices "github.com/usegro/services/catalog/internal/apps/catalog/services"
	"github.com/usegro/services/catalog/internal/interface/response"
)

type PlanController struct {
	service *catalogServices.PlanService
}

func NewPlanController(service *catalogServices.PlanService) *PlanController {
	return &PlanController{service: service}
}

func (pc *PlanController) CreatePlan(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	serviceID := c.Params("serviceId")
	var body dto.CreatePlanDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}
	if body.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Plan name is required",
		})
	}
	if body.PlanType != "subscription" && body.PlanType != "package" {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "plan_type must be 'subscription' or 'package'",
		})
	}
	plan, err := pc.service.CreatePlan(c.Context(), crmID, serviceID, body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Plan created successfully",
		Data:            plan,
	})
}

func (pc *PlanController) ListPlans(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	serviceID := c.Params("serviceId")
	plans, err := pc.service.ListPlans(c.Context(), crmID, serviceID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Plans retrieved successfully",
		Data:            plans,
	})
}

func (pc *PlanController) UpdatePlan(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	var body dto.UpdatePlanDTO
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: "Invalid request body",
		})
	}
	plan, err := pc.service.UpdatePlan(c.Context(), crmID, c.Params("planId"), body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    response.BAD_REQUEST,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Plan updated successfully",
		Data:            plan,
	})
}

func (pc *PlanController) DeletePlan(c *fiber.Ctx) error {
	crmID := c.Locals("crmID").(string)
	if err := pc.service.DeletePlan(c.Context(), crmID, c.Params("planId")); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(response.CommonResponse{
			ResponseCode:    response.RESOURCE_NOT_FOUND,
			ResponseMessage: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Plan deleted successfully",
	})
}
