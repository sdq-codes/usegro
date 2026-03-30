package middleware

import "github.com/gofiber/fiber/v2"

func CRMContextMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		crmID := c.Get("X-CRM-ID")
		if crmID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Missing CRM ID")
		}

		c.Locals("crmID", crmID)
		return c.Next()
	}
}
