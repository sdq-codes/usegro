package repositories

import (
	"context"

	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"gorm.io/gorm"
)

type VerificationRepositoryInterface interface {
	Create(ctx context.Context, tx *gorm.DB, verification *models.Verification) error
	Fetch(ctx context.Context, tx *gorm.DB, userID string, verificationType string) (*[]models.Verification, error)
	UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status string) error
	Delete(ctx context.Context, tx *gorm.DB, userId string, verificationType string) error
}

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) VerificationRepositoryInterface {
	return &VerificationRepository{
		db: db,
	}
}

func (r *VerificationRepository) Create(ctx context.Context, tx *gorm.DB, verification *models.Verification) error {
	return tx.WithContext(ctx).Create(verification).Error
}

func (r *VerificationRepository) Fetch(ctx context.Context, tx *gorm.DB, userID string, verificationType string) (*[]models.Verification, error) {
	var v []models.Verification
	err := tx.WithContext(ctx).
		Where("user_id = ? AND type = ?", userID, verificationType).
		Order("created_at DESC").
		Find(&v).Error

	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VerificationRepository) UpdateStatus(ctx context.Context, tx *gorm.DB, id string, status string) error {
	return tx.WithContext(ctx).
		Model(&models.Verification{}).
		Where("id = ?", id).
		Update("status", status).
		Error
}

func (r *VerificationRepository) Delete(ctx context.Context, tx *gorm.DB, userId string, verificationType string) error {
	_ = tx.WithContext(ctx).
		Where("user_id = ? AND type = ?", userId, verificationType).
		Order("created_at DESC").
		Delete(models.VerificationToken{}).Error
	return tx.WithContext(ctx).
		Where("user_id = ? AND type = ?", userId, verificationType).
		Order("created_at DESC").
		Delete(models.Verification{}).Error
}
