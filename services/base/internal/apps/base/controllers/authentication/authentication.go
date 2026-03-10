package authentication

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	authenticationService "github.com/sdq-codes/usegro-api/internal/apps/base/services/authentication"
	"github.com/sdq-codes/usegro-api/internal/apps/base/validation"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
	"github.com/sdq-codes/usegro-api/pkg/exception"
)

type Controller struct {
	authenticationService authenticationService.Service
}

func NewAuthenticationController(authenticationService authenticationService.Service) *Controller {
	return &Controller{authenticationService: authenticationService}
}

// Exist godoc
// @Summary      Check if email exists
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.UserExistDTI  true  "Email"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/email/exist [post]
func (uc *Controller) Exist(c *fiber.Ctx) error {
	var req dto.UserExistDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.UserExistValidation(req); err != nil {
		return err
	}

	userExist, err := uc.authenticationService.IsUserEmailExist(c.Context(), req)
	if err != nil {
		return err
	}

	msg := "Email is available"
	if userExist {
		msg = "Email already exists"
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: msg,
	})
}

// Register godoc
// @Summary      Register a new user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterUserDTI  true  "Registration data"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/register [post]
func (uc *Controller) Register(c *fiber.Ctx) error {
	var req dto.RegisterUserDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.CredentialsValidation(req); err != nil {
		return err
	}

	user, accessToken, refreshToken, err := uc.authenticationService.RegisterUser(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "User registered successfully",
		Data: dto.AuthTokensDTO{
			ID:           user.ID,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// Login godoc
// @Summary      Log in a user
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RegisterUserDTI  true  "Credentials"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/login [post]
func (uc *Controller) Login(c *fiber.Ctx) error {
	var req dto.RegisterUserDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.CredentialsValidation(req); err != nil {
		return err
	}

	user, accessToken, refreshToken, err := uc.authenticationService.LoginUser(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Login successful",
		Data: dto.AuthTokensDTO{
			ID:           user.ID,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchanges a valid refresh token for a new access + refresh token pair. The old refresh token is invalidated.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RefreshTokenDTI  true  "Refresh token"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/authentication/refresh [post]
func (uc *Controller) Refresh(c *fiber.Ctx) error {
	var req dto.RefreshTokenDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}
	if req.RefreshToken == "" {
		return exception.InvalidRequestBodyError
	}

	accessToken, refreshToken, err := uc.authenticationService.RefreshTokens(c.Context(), req.RefreshToken)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Token refreshed",
		Data: dto.AuthTokensDTO{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

// Logout godoc
// @Summary      Logout
// @Description  Revokes the refresh token so it cannot be used again.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RefreshTokenDTI  true  "Refresh token"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/logout [post]
func (uc *Controller) Logout(c *fiber.Ctx) error {
	var req dto.RefreshTokenDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}
	if req.RefreshToken == "" {
		return exception.InvalidRequestBodyError
	}

	if err := uc.authenticationService.Logout(c.Context(), req.RefreshToken); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Logged out successfully",
	})
}

// ForgotPassword godoc
// @Summary      Request password reset
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.ForgotPasswordDTI  true  "Email"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/forgot-password [post]
func (uc *Controller) ForgotPassword(c *fiber.Ctx) error {
	var req dto.ForgotPasswordDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.ForgotPasswordEmailValidation(req); err != nil {
		return err
	}

	if err := uc.authenticationService.SendPasswordResetEmail(c.Context(), req.Email); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "Password reset link has been sent.",
	})
}

// ResetPassword godoc
// @Summary      Reset user password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.ResetPasswordDTI  true  "Reset data"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/reset-password [post]
func (uc *Controller) ResetPassword(c *fiber.Ctx) error {
	var req dto.ResetPasswordDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.ResetPasswordEmailValidation(req); err != nil {
		return err
	}

	if err := uc.authenticationService.ResetUserPassword(c.Context(), req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Password has been reset successfully",
	})
}

// RequestEmailCode godoc
// @Summary      Request a login code
// @Description  Sends a 6-character one-time code to the given email. Always returns 200 to avoid leaking whether an email is registered. Code expires in 10 minutes.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.RequestEmailCodeDTI  true  "Email address"
// @Success      200  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Router       /api/v1/authentication/email-code/request [post]
func (uc *Controller) RequestEmailCode(c *fiber.Ctx) error {
	var req dto.RequestEmailCodeDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.RequestEmailCodeValidation(req); err != nil {
		return err
	}

	if err := uc.authenticationService.RequestLoginCode(c.Context(), req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "If that email is registered, a login code has been sent",
	})
}

// VerifyEmailCode godoc
// @Summary      Login with email code
// @Description  Validates the one-time code sent to the user's email and returns an access + refresh token pair.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        payload  body      dto.VerifyEmailCodeDTI  true  "Email and code"
// @Success      200  {object}  response.CommonResponse{data=dto.AuthTokensDTO}
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/authentication/email-code/verify [post]
func (uc *Controller) VerifyEmailCode(c *fiber.Ctx) error {
	var req dto.VerifyEmailCodeDTI

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	authenticationValidation := validation.AuthenticationValidation{}
	if err := authenticationValidation.VerifyEmailCodeValidation(req); err != nil {
		return err
	}

	user, accessToken, refreshToken, err := uc.authenticationService.VerifyLoginCode(c.Context(), req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Login successful",
		Data: dto.AuthTokensDTO{
			ID:           user.ID,
			Email:        user.Email,
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}
