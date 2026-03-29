package verification

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/dto"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
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

func (uc *Controller) Verify(c *fiber.Ctx) error {
	var req dto.EmailVerificationDTI
	claims, err := auth.AuthUser(c)
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

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return exception.UnauthorizedError
	}
	user := &models.User{ID: userID, Email: claims.Email}

	if err = uc.emailService.EmailVerification(c.Context(), req.TokenHash, user); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_CREATED,
		ResponseMessage: "User email successfully verified",
		Data: dto.RegisteredUserDTO{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (uc *Controller) Resend(c *fiber.Ctx) error {
	ctx := c.Context()
	claims, err := auth.AuthUser(c)
	if err != nil {
		return exception.UnauthorizedError
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return exception.UnauthorizedError
	}
	user := &models.User{ID: userID, Email: claims.Email}

	if err := uc.verificationService.ResendVerificationToken(ctx, user); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.CommonResponse{
		ResponseCode:    response.RESOURCE_FETCHED,
		ResponseMessage: "Verification email resent successfully",
		Data: dto.RegisteredUserDTO{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
