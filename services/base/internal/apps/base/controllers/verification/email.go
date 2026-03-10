package verification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	"github.com/sdq-codes/usegro-api/internal/apps/base/services/verification"
	"github.com/sdq-codes/usegro-api/internal/apps/base/validation"
	"github.com/sdq-codes/usegro-api/internal/helper/auth"
	"github.com/sdq-codes/usegro-api/internal/interface/response"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	"gorm.io/gorm"
)

type Controller struct {
	emailService        *verification.EmailService
	verificationService *verification.Service
}

func NewEmailVerificationController(db *gorm.DB, rdb *redis.Client) *Controller {
	return &Controller{
		emailService:        verification.NewEmailVerificationService(db, rdb),
		verificationService: verification.NewVerificationService(db, rdb),
	}
}

// Verify godoc
// @Summary      Verify email address
// @Description  Verifies a user's email address using the 6-digit code sent to their inbox
// @Tags         Verification
// @Accept       json
// @Produce      json
// @Param        Authorization  header    string                    true  "Bearer access token"
// @Param        payload        body      dto.EmailVerificationDTI  true  "Verification code"
// @Success      201  {object}  response.CommonResponse
// @Failure      400  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/verification/email [post]
func (uc *Controller) Verify(c *fiber.Ctx) error {
	var req dto.EmailVerificationDTI
	authUser, err := auth.AuthUser(c)
	if err != nil {
		return exception.UnauthorizedError
	}

	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	emailVerificationValidation := validation.EmailValidation{}
	if err := emailVerificationValidation.EmailVerificationValidation(req); err != nil {
		return err
	}

	err = uc.emailService.EmailVerification(c.Context(), req.TokenHash, &authUser.User)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "User email successfully verified",
		Data: dto.RegisteredUserDTO{
			ID:    authUser.User.ID,
			Email: authUser.User.Email,
		},
	})
}

// Resend godoc
// @Summary      Resend verification email
// @Description  Sends a new verification code to the authenticated user's email address
// @Tags         Verification
// @Produce      json
// @Param        Authorization  header  string  true  "Bearer access token"
// @Success      200  {object}  response.CommonResponse
// @Failure      401  {object}  response.CommonResponse
// @Router       /api/v1/verification/email/resend [get]
func (uc *Controller) Resend(c *fiber.Ctx) error {
	ctx := c.Context()
	authUser, err := auth.AuthUser(c)
	if err != nil {
		return exception.UnauthorizedError
	}

	if err := uc.verificationService.ResendVerificationToken(ctx, &authUser.User); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Verification email resent successfully",
		Data: dto.RegisteredUserDTO{
			ID:    authUser.User.ID,
			Email: authUser.User.Email,
		},
	})
}
