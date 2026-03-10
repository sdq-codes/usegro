package repositories

import (
	"context"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"gorm.io/gorm"
)

type VerificationTokenRepositoryInterface interface {
	Create(ctx context.Context, tx *gorm.DB, token *models.VerificationToken) error
	Fetch(ctx context.Context, tx *gorm.DB, userID string, tokenType string, token string) (*models.VerificationToken, error)
	FetchToken(ctx context.Context, tx *gorm.DB, tokenType string, token string) (*models.VerificationToken, error)
	Delete(ctx context.Context, tx *gorm.DB, userID string, tokenType string, token string) error
	Clear(ctx context.Context, tx *gorm.DB, userID string, tokenType string) error
}

type VerificationTokenRepository struct {
	db *gorm.DB
}

func NewVerificationTokenRepository(db *gorm.DB) *VerificationTokenRepository {
	return &VerificationTokenRepository{
		db: db,
	}
}

func (v *VerificationTokenRepository) Create(ctx context.Context, tx *gorm.DB, token *models.VerificationToken) error {
	return tx.WithContext(ctx).Create(token).Error
}

func (v *VerificationTokenRepository) Fetch(ctx context.Context, tx *gorm.DB, userID string, tokenType string, token string) (*models.VerificationToken, error) {
	var vt models.VerificationToken
	err := tx.WithContext(ctx).
		Where("token_hash = ? AND user_id = ? AND type = ? AND expires_at > NOW()", token, userID, tokenType).
		First(&vt).Error

	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &vt, nil
}

func (v *VerificationTokenRepository) FetchToken(ctx context.Context, tx *gorm.DB, tokenType string, token string) (*models.VerificationToken, error) {
	var vt models.VerificationToken
	err := tx.WithContext(ctx).
		Where("token_hash = ? AND type = ? AND expires_at > NOW()", token, tokenType).
		First(&vt).Error

	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		}
		return nil, err
	}
	return &vt, nil
}

func (v *VerificationTokenRepository) Delete(ctx context.Context, tx *gorm.DB, userID string, tokenType string, token string) error {
	return tx.WithContext(ctx).
		Where("token_hash = ? AND user_id = ? AND type = ?", token, userID, tokenType).
		Delete(&models.VerificationToken{}).Error
}

func (v *VerificationTokenRepository) Clear(ctx context.Context, tx *gorm.DB, userID string, tokenType string) error {
	return tx.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, tokenType).
		Delete(&models.VerificationToken{}).Error
}
