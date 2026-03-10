package repositories

import (
	"context"
	"github.com/sdq-codes/usegro-api/internal/apps/base/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUser(ctx context.Context, tx *gorm.DB, email string) (*models.User, error)
	CreateUser(ctx context.Context, tx *gorm.DB, user *models.User) (*models.User, error)
	IsUserEmailExist(ctx context.Context, email string) (bool, error)
	UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (u *UserRepository) GetUser(ctx context.Context, tx *gorm.DB, email string) (*models.User, error) {
	user := &models.User{}
	if err := tx.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUserById(ctx context.Context, tx *gorm.DB, id string) (*models.User, error) {
	user := &models.User{}
	if err := tx.WithContext(ctx).Preload("Verifications").Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) GetUsersWithPagination(ctx context.Context, limit int, offset int) ([]models.User, error) {
	panic("implement me")
}

func (u *UserRepository) UpdateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {
	err := tx.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"password": user.Password,
		}).Error

	if err != nil {
		return err
	}
	return nil
}

// Add authentication with transaction and return id
func (u *UserRepository) CreateUser(ctx context.Context, tx *gorm.DB, user *models.User) error {
	return tx.WithContext(ctx).Create(user).Error

}

func (u *UserRepository) IsUserEmailExist(ctx context.Context, tx *gorm.DB, email string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
