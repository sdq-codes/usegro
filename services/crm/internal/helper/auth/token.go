package auth

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/config"
	"github.com/usegro/services/crm/internal/apps/base/models"
	sharedauth "github.com/usegro/services/shared/pkg/auth"
)

// CreateToken mints a signed access JWT for the given user.
func CreateToken(user models.User) (string, error) {
	cfg := config.GetConfig().Auth
	return sharedauth.CreateToken(user, cfg.ApiSecret, cfg.TokenExpiryMinutes)
}

// AuthUser parses the Bearer token from the request and returns the claims.
func AuthUser(c *fiber.Ctx) (*models.Token, error) {
	claims, err := sharedauth.AuthUser[models.User](c, config.GetConfig().Auth.ApiSecret)
	if err != nil {
		return nil, err
	}
	return &models.Token{User: claims.User, RegisteredClaims: claims.RegisteredClaims}, nil
}

// CreateRefreshToken generates an opaque refresh token stored in Redis.
func CreateRefreshToken(ctx context.Context, rdb redis.Cmdable, user models.User) (string, error) {
	return sharedauth.CreateRefreshToken(ctx, rdb, user, config.GetConfig().Auth.RefreshTokenExpiryDays)
}

// ValidateAndRotateRefreshToken validates the refresh token, rotates it, and returns new tokens.
func ValidateAndRotateRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) (*models.User, string, string, error) {
	cfg := config.GetConfig().Auth
	return sharedauth.ValidateAndRotateRefreshToken[models.User](ctx, rdb, rawToken, cfg.ApiSecret, cfg.TokenExpiryMinutes, cfg.RefreshTokenExpiryDays)
}

// RevokeRefreshToken deletes the refresh token from Redis (logout).
func RevokeRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) error {
	return sharedauth.RevokeRefreshToken(ctx, rdb, rawToken)
}
