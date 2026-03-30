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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/config"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
)

const (
	refreshKeyPrefix    = "refresh:"
	refreshFamilyPrefix = "refresh:family:"
	refreshCookieName   = "refresh_token"
)

// refreshTokenData is the value stored in Redis for each refresh token.
type refreshTokenData struct {
	User             models.User `json:"user"`
	FamilyID         string      `json:"family_id"`
	SessionCreatedAt time.Time   `json:"session_created_at"`
}

// ── Expiry helpers ────────────────────────────────────────────────────────────

func tokenExpiry() time.Duration {
	mins := config.GetConfig().Auth.TokenExpiryMinutes
	if mins <= 0 {
		mins = 15
	}
	return time.Duration(mins) * time.Minute
}

func refreshExpiry() time.Duration {
	days := config.GetConfig().Auth.RefreshTokenExpiryDays
	if days <= 0 {
		days = 2
	}
	return time.Duration(days) * 24 * time.Hour
}

func maxSessionDuration() time.Duration {
	days := config.GetConfig().Auth.MaxSessionDays
	if days <= 0 {
		days = 30
	}
	return time.Duration(days) * 24 * time.Hour
}

// ── Access token ──────────────────────────────────────────────────────────────

// CreateToken mints a signed JWT with slim claims (user ID + email only).
func CreateToken(user models.User) (string, error) {
	claims := &models.TokenClaims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExpiry())),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.GetConfig().Auth.ApiSecret))
}

// AuthUser parses the Bearer token from the request and returns the slim claims.
func AuthUser(c *fiber.Ctx) (*models.TokenClaims, error) {
	reqToken := c.Get("Authorization")
	parts := strings.Split(reqToken, "Bearer ")
	if len(parts) != 2 {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing Authorization header")
	}

	tokenStr := strings.TrimSpace(parts[1])
	tk := &models.TokenClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.GetConfig().Auth.ApiSecret), nil
	})
	if err != nil {
		return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
	}

	return tk, nil
}

// ── Refresh token ─────────────────────────────────────────────────────────────

func hashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// CreateRefreshToken generates an opaque refresh token and stores it in Redis.
// Pass familyID="" and zero sessionCreatedAt to start a new session.
// Pass existing values to continue a session (token rotation).
func CreateRefreshToken(ctx context.Context, rdb redis.Cmdable, user models.User, familyID string, sessionCreatedAt time.Time) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	rawToken := base64.URLEncoding.EncodeToString(b)
	hash := hashToken(rawToken)

	if familyID == "" {
		familyID = uuid.New().String()
		sessionCreatedAt = time.Now()
	}

	data := refreshTokenData{
		User:             user,
		FamilyID:         familyID,
		SessionCreatedAt: sessionCreatedAt,
	}
	dataJSON, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	tokenTTL := refreshExpiry()
	familyTTL := maxSessionDuration()

	pipe := rdb.Pipeline()
	pipe.Set(ctx, refreshKeyPrefix+hash, dataJSON, tokenTTL)
	// Family key always tracks the current valid hash for replay detection.
	// TTL is refreshed on every successful rotation.
	pipe.Set(ctx, refreshFamilyPrefix+familyID, hash, familyTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		return "", err
	}

	return rawToken, nil
}

// ValidateAndRotateRefreshToken validates the token, checks for replay attacks
// and session expiry, then issues a new token pair (rotation).
func ValidateAndRotateRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) (*models.User, string, string, error) {
	hash := hashToken(rawToken)
	key := refreshKeyPrefix + hash

	dataJSON, err := rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, "", "", fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired refresh token")
	}
	if err != nil {
		return nil, "", "", err
	}

	var data refreshTokenData
	if err := json.Unmarshal(dataJSON, &data); err != nil {
		return nil, "", "", err
	}

	// Fix 5: Family check — detect replay attacks.
	// The family key always holds the hash of the current valid token.
	// If it doesn't match, this token was already rotated → replay detected.
	familyKey := refreshFamilyPrefix + data.FamilyID
	currentHash, err := rdb.Get(ctx, familyKey).Result()
	if err != nil || currentHash != hash {
		// Revoke the entire family to protect the real user.
		rdb.Del(ctx, familyKey)
		rdb.Del(ctx, key)
		return nil, "", "", fiber.NewError(fiber.StatusUnauthorized, "Session revoked due to suspicious activity")
	}

	// Fix 2: Enforce maximum session lifetime.
	if time.Since(data.SessionCreatedAt) > maxSessionDuration() {
		rdb.Del(ctx, key)
		rdb.Del(ctx, familyKey)
		return nil, "", "", fiber.NewError(fiber.StatusUnauthorized, "Session expired, please log in again")
	}

	// Consume the old token immediately (single-use).
	rdb.Del(ctx, key)

	newAccessToken, err := CreateToken(data.User)
	if err != nil {
		return nil, "", "", err
	}

	// Rotate: issue new token in the same family, preserving session start time.
	newRefreshToken, err := CreateRefreshToken(ctx, rdb, data.User, data.FamilyID, data.SessionCreatedAt)
	if err != nil {
		return nil, "", "", err
	}

	return &data.User, newAccessToken, newRefreshToken, nil
}

// RevokeRefreshToken deletes a refresh token and its entire family (logout).
func RevokeRefreshToken(ctx context.Context, rdb redis.Cmdable, rawToken string) error {
	hash := hashToken(rawToken)
	key := refreshKeyPrefix + hash

	dataJSON, err := rdb.Get(ctx, key).Bytes()
	if err == nil {
		var data refreshTokenData
		if json.Unmarshal(dataJSON, &data) == nil {
			rdb.Del(ctx, refreshFamilyPrefix+data.FamilyID)
		}
	}
	return rdb.Del(ctx, key).Err()
}

// ── Cookie helpers ────────────────────────────────────────────────────────────

// SetRefreshCookie writes the refresh token as an HttpOnly cookie.
func SetRefreshCookie(c *fiber.Ctx, token string) {
	cfg := config.GetConfig()
	secure := cfg.Env == "production"
	sameSite := "Lax"
	if cfg.Env == "production" {
		sameSite = "Strict"
	}
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(refreshExpiry().Seconds()),
		Secure:   secure,
		HTTPOnly: true,
		SameSite: sameSite,
	})
}

// ClearRefreshCookie expires the refresh token cookie immediately.
func ClearRefreshCookie(c *fiber.Ctx) {
	c.Cookie(&fiber.Cookie{
		Name:     refreshCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
	})
}

// GetRefreshCookie reads the refresh token from the request cookie.
func GetRefreshCookie(c *fiber.Ctx) string {
	return c.Cookies(refreshCookieName)
}
