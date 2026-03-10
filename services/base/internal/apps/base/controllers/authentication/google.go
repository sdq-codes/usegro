package authentication

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/config"
	authenticationService "github.com/sdq-codes/usegro-api/internal/apps/base/services/authentication"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
)

const (
	oauthStateTTL    = 10 * time.Minute
	oauthStatePrefix = "oauth:google:state:"
)

type GoogleController struct {
	service *authenticationService.Service
	rdb     *redis.Client
}

func NewGoogleAuthenticationController(service *authenticationService.Service, rdb *redis.Client) *GoogleController {
	return &GoogleController{service: service, rdb: rdb}
}

// HandleGoogleLogin godoc
// @Summary      Initiate Google OAuth login
// @Description  Generates a CSRF state token, stores it in Redis, and redirects the user to Google's OAuth consent screen
// @Tags         Authentication
// @Produce      json
// @Success      307  "Redirect to Google OAuth consent screen"
// @Failure      500  {object}  response.CommonResponse
// @Router       /api/v1/authentication/google/login [get]
func (gc *GoogleController) HandleGoogleLogin(c *fiber.Ctx) error {
	state, err := generateState()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusInternalServerError,
			ResponseMessage: "Failed to generate OAuth state",
		})
	}

	if err := gc.rdb.Set(c.Context(), oauthStatePrefix+state, "1", oauthStateTTL).Err(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.CommonResponse{
			ResponseCode:    fiber.StatusInternalServerError,
			ResponseMessage: "Failed to store OAuth state",
		})
	}

	url := gc.service.GoogleAuthURL(state)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

// HandleGoogleCallback godoc
// @Summary      Google OAuth callback
// @Description  Validates the CSRF state, exchanges the authorization code for tokens, then redirects to the frontend with access_token and refresh_token as query params. On error, redirects with an error param.
// @Tags         Authentication
// @Produce      json
// @Param        code   query  string  true  "Authorization code returned by Google"
// @Param        state  query  string  true  "CSRF state token (must match the one issued by /login)"
// @Success      307  "Redirect to frontend /oauth/callback?access_token=...&refresh_token=..."
// @Failure      307  "Redirect to frontend /oauth/callback?error=..."
// @Router       /api/v1/authentication/google/callback [get]
func (gc *GoogleController) HandleGoogleCallback(c *fiber.Ctx) error {
	frontendURL := config.GetConfig().FrontEnd.Url
	state := c.Query("state")
	code := c.Query("code")

	if state == "" || code == "" {
		return c.Redirect(
			fmt.Sprintf("%s/oauth/callback?error=invalid_request", frontendURL),
			fiber.StatusTemporaryRedirect,
		)
	}

	// Validate and consume the state token (one-time use).
	key := oauthStatePrefix + state
	result, err := gc.rdb.GetDel(c.Context(), key).Result()
	if err != nil || result == "" {
		return c.Redirect(
			fmt.Sprintf("%s/oauth/callback?error=invalid_state", frontendURL),
			fiber.StatusTemporaryRedirect,
		)
	}

	_, accessToken, refreshToken, err := gc.service.GoogleLogin(c.Context(), code)
	if err != nil {
		return c.Redirect(
			fmt.Sprintf("%s/oauth/callback?error=auth_failed", frontendURL),
			fiber.StatusTemporaryRedirect,
		)
	}

	return c.Redirect(
		fmt.Sprintf("%s/oauth/callback?access_token=%s&refresh_token=%s", frontendURL, accessToken, refreshToken),
		fiber.StatusTemporaryRedirect,
	)
}

func generateState() (string, error) {
	b := make([]byte, 24)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
