package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	userService "github.com/sdq-codes/usegro-api/internal/apps/base/services/user"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
)

type Controller struct {
	userService userService.Service
}

func NewUserController(userService userService.Service) *Controller {
	return &Controller{userService: userService}
}

// GetLoggedInUser godoc
// @Summary      Get logged-in user
// @Description  Returns the profile of the currently authenticated user
// @Tags         User
// @Produce      json
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/user [get]
func (uc *Controller) GetLoggedInUser(c *fiber.Ctx) error {
	user, orgs, err := uc.userService.GetLoggedInUser(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(response.CommonResponse{
			ResponseCode:    response.VALIDATION_FAILED,
			ResponseMessage: "Unauthorized access",
			Data:            nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "User fetched successfully",
		Data: dto.UserDTO{
			ID:            user.ID,
			Email:         user.Email,
			Verifications: user.Verifications,
			Organizations: orgs,
		},
	})
}
