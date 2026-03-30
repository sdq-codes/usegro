package error

import (
	"github.com/gofiber/fiber/v2"
	"github.com/usegro/services/catalog/internal/interface/response"
	"github.com/usegro/services/catalog/pkg/exception"
)

// Centralized error handler for all routes
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Retrieve neccessary details
	// Status code defaults to 500
	responseCode := fiber.StatusInternalServerError
	responseMessage := err.Error()
	requestID, _ := c.Locals("requestid").(string)

	var cErrs *exception.ExceptionErrors

	// Use response code from ExceptionError
	if e, ok := err.(*exception.ExceptionErrors); ok {
		responseCode = e.HttpStatusCode
		cErrs = e
	} else if e, ok := err.(*fiber.Error); ok {
		responseCode = e.Code
	}

	// Handle 500 error
	return c.Status(responseCode).JSON(
		&response.CommonResponse{
			ResponseCode:    responseCode,
			ResponseMessage: responseMessage,
			Errors:          cErrs,
			RequestID:       requestID,
		},
	)
}
