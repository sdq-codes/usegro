package httperror

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/shared/pkg/exception"
	"github.com/usegro/services/shared/pkg/response"
)

// ErrorHandler is a centralized error handler for all Fiber routes.
// It converts ExceptionErrors to a CommonResponse and sets the appropriate HTTP status code.
func ErrorHandler(c *fiber.Ctx, err error) error {
	responseCode := fiber.StatusInternalServerError
	responseMessage := err.Error()
	requestID := c.Locals("requestid").(string)

	var cErrs *exception.ExceptionErrors

	cErrs, ok := err.(*exception.ExceptionErrors)
	if ok {
		responseCode = cErrs.HttpStatusCode
	}

	return c.Status(responseCode).JSON(
		&response.CommonResponse{
			ResponseCode:    responseCode,
			ResponseMessage: responseMessage,
			Errors:          cErrs,
			RequestID:       requestID,
		},
	)
}
