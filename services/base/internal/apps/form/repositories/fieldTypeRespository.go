package repositories

import (
	"context"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"gorm.io/gorm"
)

// FieldTypeRepository interface
type FieldTypeRepository interface {
	GetByName(ctx context.Context, tx *gorm.DB, name string) (*models.FieldType, error)
	GetAll(ctx context.Context, tx *gorm.DB) ([]models.FieldType, error)
}

type fieldTypeRepository struct{}

func NewFieldTypeRepository() FieldTypeRepository {
	return &fieldTypeRepository{}
}

func (r *fieldTypeRepository) GetByName(ctx context.Context, tx *gorm.DB, name string) (*models.FieldType, error) {
	var fieldType *models.FieldType
	err := tx.WithContext(ctx).Preload("Configs").Preload("Validations").Where("name = ?", name).First(&fieldType).Error
	if err != nil {
		return nil, err
	}
	return fieldType, nil
}

func (r *fieldTypeRepository) GetAll(ctx context.Context, tx *gorm.DB) ([]models.FieldType, error) {
	var fieldTypes []models.FieldType
	err := tx.WithContext(ctx).
		Preload("Configs").
		Preload("Validations").
		Find(&fieldTypes).Error
	if err != nil {
		return nil, err
	}
	return fieldTypes, nil
}
