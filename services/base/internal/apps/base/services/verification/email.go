package verification

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"github.com/sdq-codes/usegro-api/internal/apps/base/repositories"
	"github.com/sdq-codes/usegro-api/internal/logger"
	"github.com/sdq-codes/usegro-api/pkg/exception"
	"gorm.io/gorm"
)

const emailVerifyKeyPrefix = "auth:email_verify:"

type EmailService struct {
	db                     *gorm.DB
	rdb                    *redis.Client
	verificationRepository repositories.VerificationRepositoryInterface
}

func NewEmailVerificationService(db *gorm.DB, rdb *redis.Client) *EmailService {
	return &EmailService{
		db:                     db,
		rdb:                    rdb,
		verificationRepository: repositories.NewVerificationRepository(db),
	}
}

func emailVerifyKey(userID string) string {
	h := sha256.Sum256([]byte(userID))
	return emailVerifyKeyPrefix + hex.EncodeToString(h[:])
}

func (s *EmailService) EmailVerification(ctx context.Context, code string, user *models.User) error {
	// Fetch verification status record.
	records, err := s.verificationRepository.Fetch(ctx, s.db, user.ID.String(), "EMAIL")
	if err != nil {
		logger.Log.Error(fmt.Sprintf("EmailVerification: fetch failed for %s: %v", user.Email, err))
		return exception.EmailVerificationError
	}
	if len(*records) == 0 {
		logger.Log.Error(fmt.Sprintf("EmailVerification: no verification record for %s", user.Email))
		return exception.EmailVerificationError
	}
	record := (*records)[0]

	if record.Status == "VERIFIED" {
		return exception.EmailVerificationError
	}

	// Validate and consume the code atomically from Redis.
	stored, err := s.rdb.GetDel(ctx, emailVerifyKey(user.ID.String())).Result()
	if err != nil || stored == "" {
		logger.Log.Info(fmt.Sprintf("EmailVerification: code expired or missing for %s", user.Email))
		return exception.IncorrectEmailVerificationError
	}
	if stored != code {
		logger.Log.Info(fmt.Sprintf("EmailVerification: wrong code for %s", user.Email))
		return exception.IncorrectEmailVerificationError
	}

	// Mark as verified inside a transaction.
	tx := s.db.Begin()
	if err := s.verificationRepository.UpdateStatus(ctx, tx, record.ID.String(), "VERIFIED"); err != nil {
		tx.Rollback()
		logger.Log.Error(fmt.Sprintf("EmailVerification: status update failed for %s: %v", user.Email, err))
		return exception.EmailVerificationError
	}
	tx.Commit()
	return nil
}
