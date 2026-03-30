package authentication

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	authenticationService "github.com/sdq-codes/usegro-api/internal/apps/base/services/authentication"
	"github.com/sdq-codes/usegro-api/internal/apps/base/validation"
	"github.com/sdq-codes/usegro-api/internal/helper/auth"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
	"github.com/sdq-codes/usegro-api/pkg/exception"
)

type Controller struct {
	authenticationService authenticationService.Service
}

func NewAuthenticationController(authenticationService authenticationService.Service) *Controller {
	return &Controller{authenticationService: authenticationService}
}

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

	auth.SetRefreshCookie(c, refreshToken)

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "User registered successfully",
		Data: dto.AuthTokensDTO{
			ID:          user.ID,
			Email:       user.Email,
			AccessToken: accessToken,
		},
	})
}

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

	auth.SetRefreshCookie(c, refreshToken)

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Login successful",
		Data: dto.AuthTokensDTO{
			ID:          user.ID,
			Email:       user.Email,
			AccessToken: accessToken,
		},
	})
}

func (uc *Controller) Refresh(c *fiber.Ctx) error {
	rawToken := auth.GetRefreshCookie(c)
	if rawToken == "" {
		return exception.UnauthorizedError
	}

	accessToken, refreshToken, err := uc.authenticationService.RefreshTokens(c.Context(), rawToken)
	if err != nil {
		auth.ClearRefreshCookie(c)
		return err
	}

	auth.SetRefreshCookie(c, refreshToken)

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Token refreshed",
		Data: dto.AuthTokensDTO{
			AccessToken: accessToken,
		},
	})
}

func (uc *Controller) Logout(c *fiber.Ctx) error {
	rawToken := auth.GetRefreshCookie(c)
	if rawToken != "" {
		_ = uc.authenticationService.Logout(c.Context(), rawToken)
	}

	auth.ClearRefreshCookie(c)

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_UPDATED,
		ResponseMessage: "Logged out successfully",
	})
}

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

	auth.SetRefreshCookie(c, refreshToken)

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Login successful",
		Data: dto.AuthTokensDTO{
			ID:          user.ID,
			Email:       user.Email,
			AccessToken: accessToken,
		},
	})
}
