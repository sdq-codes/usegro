package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
)

func JwtVerify() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		parts := strings.SplitN(strings.TrimSpace(header), " ", 2)

		if len(parts) < 2 || parts[1] == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusUnauthorized,
				ResponseMessage: "User Unauthenticated",
			})
		}

		tk := &models.TokenClaims{}
		_, err := jwt.ParseWithClaims(parts[1], tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GetConfig().Auth.ApiSecret), nil
		})
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(response.CommonResponse{
				ResponseCode:    fiber.StatusUnauthorized,
				ResponseMessage: "User Unauthenticated",
			})
		}

		c.Locals("user", tk)
		return c.Next()
	}
}
