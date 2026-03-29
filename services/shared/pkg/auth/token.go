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
)

const refreshKeyPrefix = "refresh:"

// Claims is the JWT claims struct wrapping any user type U.
type Claims[U any] struct {
	User U `json:"user"`
	jwt.RegisteredClaims
}

// CreateToken mints a signed access JWT for the given user.
// expiryMins defaults to 30 if <= 0.
func CreateToken[U any](user U, secret string, expiryMins int) (string, error) {
	if expiryMins <= 0 {
		expiryMins = 30
	}
	claims := &Claims[U]{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryMins) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseToken parses and validates a raw JWT string, returning the typed claims.
func ParseToken[U any](tokenStr string, secret string) (*Claims[U], error) {
	tk := &Claims[U]{}
	_, err := jwt.ParseWithClaims(tokenStr, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}
	return tk, nil
}

// AuthUser extracts and validates the Bearer token from the Fiber request context.
func AuthUser[U any](c *fiber.Ctx, secret string) (*Claims[U], error) {
	reqToken := c.Get("Authorization")
	parts := strings.Split(reqToken, "Bearer ")
	if len(parts) != 2 {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing Authorization header")
	}
	return ParseToken[U](strings.TrimSpace(parts[1]), secret)
}

// hashToken returns the SHA-256 hex digest of a raw token string.
func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

func refreshRedisKey(rawToken string) string {
	return refreshKeyPrefix + hashToken(rawToken)
}

// CreateRefreshToken generates a cryptographically random opaque token,
// stores the user payload in Redis keyed by its hash, and returns the raw token.
// expiryDays defaults to 7 if <= 0.
func CreateRefreshToken[U any](ctx context.Context, rdb redis.Cmdable, user U, expiryDays int) (string, error) {
	if expiryDays <= 0 {
		expiryDays = 7
	}
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	rawToken := base64.URLEncoding.EncodeToString(b)

	userJSON, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	expiry := time.Duration(expiryDays) * 24 * time.Hour
	if err := rdb.Set(ctx, refreshRedisKey(rawToken), userJSON, expiry).Err(); err != nil {
		return "", err
	}
	return rawToken, nil
}

// ValidateAndRotateRefreshToken validates the refresh token, rotates it (one-time use),
// and returns the user, new access token, and new refresh token.
func ValidateAndRotateRefreshToken[U any](ctx context.Context, rdb redis.Cmdable, rawToken, secret string, expiryMins, expiryDays int) (*U, string, string, error) {
	key := refreshRedisKey(rawToken)

	userJSON, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, "", "", fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}
	if err != nil {
		return nil, "", "", err
	}

	var user U
	if err := json.Unmarshal(userJSON, &user); err != nil {
		return nil, "", "", err
	}

	rdb.Del(ctx, key)

	newAccessToken, err := CreateToken(user, secret, expiryMins)
	if err != nil {
		return nil, "", "", err
	}

	newRefreshToken, err := CreateRefreshToken(ctx, rdb, user, expiryDays)
	if err != nil {
		return nil, "", "", err
	}

	return &user, newAccessToken, newRefreshToken, nil
}

// RevokeRefreshToken deletes a refresh token from Redis (logout).
func RevokeRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) error {
	return rdb.Del(ctx, refreshRedisKey(rawToken)).Err()
}
