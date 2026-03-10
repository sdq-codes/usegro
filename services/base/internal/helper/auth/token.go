package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
)

const refreshKeyPrefix = "refresh:"

// ── Access token ─────────────────────────────────────────────────────────────

func tokenExpiry() time.Duration {
	mins := config.GetConfig().Auth.TokenExpiryMinutes
	if mins <= 0 {
		mins = 30
	}
	return time.Duration(mins) * time.Minute
}

// CreateToken mints a signed access JWT using models.Token as the claims struct.
func CreateToken(user models.User) (string, error) {
	claims := &models.Token{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().Auth.ApiSecret))
}

// AuthUser parses the Bearer token from the request and returns the claims.
func AuthUser(c *fiber.Ctx) (*models.Token, error) {
	reqToken := c.Get("Authorization")
	parts := strings.Split(reqToken, "Bearer ")
	if len(parts) != 2 {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing Authorization header")
	}

	tokenStr := strings.TrimSpace(parts[1])
	tk := &models.Token{}
	_, err := jwt.ParseWithClaims(tokenStr, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Auth.ApiSecret), nil
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return tk, nil
}

// ── Refresh token ─────────────────────────────────────────────────────────────

func refreshExpiry() time.Duration {
	days := config.GetConfig().Auth.RefreshTokenExpiryDays
	if days <= 0 {
		days = 7
	}
	return time.Duration(days) * 24 * time.Hour
}

// hashToken returns the SHA-256 hex digest of a raw token string.
// We store the hash in Redis so that a Redis dump doesn't expose usable tokens.
func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func refreshRedisKey(rawToken string) string {
	return refreshKeyPrefix + hashToken(rawToken)
}

// CreateRefreshToken generates a cryptographically random opaque token,
// stores the user payload in Redis keyed by its hash, and returns the raw token.
func CreateRefreshToken(ctx context.Context, rdb redis.Cmdable, user models.User) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	rawToken := base64.URLEncoding.EncodeToString(b)

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	if err := rdb.Set(ctx, refreshRedisKey(rawToken), userJSON, refreshExpiry()).Err(); err != nil {
		return "", err
	}

	return rawToken, nil
}

// ValidateAndRotateRefreshToken validates the incoming refresh token, deletes it
// (rotation — one-time use), and issues a new refresh token + access token pair.
// Returns (user, newAccessToken, newRefreshToken, error).
func ValidateAndRotateRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) (*models.User, string, string, error) {
	key := refreshRedisKey(rawToken)

	userJSON, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, "", "", fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}
	if err != nil {
		return nil, "", "", err
	}

	var user models.User
	if err := json.Unmarshal(userJSON, &user); err != nil {
		return nil, "", "", err
	}

	// Rotate: delete old token immediately.
	rdb.Del(ctx, key)

	newAccessToken, err := CreateToken(user)
	if err != nil {
		return nil, "", "", err
	}

	newRefreshToken, err := CreateRefreshToken(ctx, rdb, user)
	if err != nil {
		return nil, "", "", err
	}

	return &user, newAccessToken, newRefreshToken, nil
}

// RevokeRefreshToken deletes a refresh token from Redis (logout).
func RevokeRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) error {
	return rdb.Del(ctx, refreshRedisKey(rawToken)).Err()
}
