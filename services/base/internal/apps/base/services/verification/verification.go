package verification

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/apps/base/repositories"
	notificationModels "github.com/sdq-codes/usegro-api/internal/apps/notifications/models"
	notification "github.com/sdq-codes/usegro-api/internal/apps/notifications/services"
	"github.com/sdq-codes/usegro-api/internal/helper/random"
	"github.com/sdq-codes/usegro-api/internal/interface/resources/templates/emails"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"gorm.io/gorm"
)

const verifyCodeTTL = 20 * time.Minute

type Service struct {
	db                     *gorm.DB
	rdb                    *redis.Client
	verificationRepository repositories.VerificationRepositoryInterface
	jobRepository          repositories.JobRepositoryInterface
}

func NewVerificationService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{
		db:                     db,
		rdb:                    rdb,
		verificationRepository: repositories.NewVerificationRepository(db),
		jobRepository:          repositories.NewJobRepository(db),
	}
}

func (s *Service) CreateVerification(ctx context.Context, verification *models.Verification) error {
	return s.verificationRepository.Create(ctx, s.db, verification)
}

func (s *Service) GetVerification(ctx context.Context, userId, verificationType string) (*[]models.Verification, error) {
	return s.verificationRepository.Fetch(ctx, s.db, userId, verificationType)
}

func (s *Service) UpdateVerificationStatus(ctx context.Context, id, status string) error {
	return s.verificationRepository.UpdateStatus(ctx, s.db, id, status)
}

// StoreEmailVerifyCode saves the code to Redis, overwriting any previous code.
func (s *Service) StoreEmailVerifyCode(ctx context.Context, userID, code string) error {
	return s.rdb.Set(ctx, emailVerifyKey(userID), code, verifyCodeTTL).Err()
}

// ResendVerificationToken generates a fresh code, stores it in Redis, and re-sends the email.
func (s *Service) ResendVerificationToken(ctx context.Context, user *models.User) error {
	code := random.GenerateRandomCode(6)

	if err := s.StoreEmailVerifyCode(ctx, user.ID.String(), code); err != nil {
		return err
	}

	tx := s.db.Begin()
	if err := notification.QueueEmail(ctx, tx, s.rdb, notificationModels.EmailNotification{
		FromEmail: emails.NO_REPLY_EMAIL,
		ToEmails:  []string{user.Email},
		Template:  emails.EMAIL_VERIFICATION_TEMPLATE,
		Data:      map[string]string{"token": code},
		Subject:   "Your useGro verification code",
	}, s.jobRepository); err != nil {
		tx.Rollback()
		logger.Log.Error(fmt.Sprintf("ResendVerificationToken: email queue failed for %s: %v", user.Email, err))
		return err
	}
	tx.Commit()
	return nil
}
