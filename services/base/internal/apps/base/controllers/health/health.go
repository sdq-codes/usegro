package health

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
)

type Controller struct{}

func NewHealthController() *Controller {
	return &Controller{}
}

// Show godoc
// @Summary      Health check
// @Description  Returns 200 OK when the service is running
// @Tags         Health
// @Produce      json
// @Success      200  {object}  response.CommonResponse
// @Router       /health [get]
func (h *Controller) Show(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    0,
		ResponseMessage: "OK",
	})
}
