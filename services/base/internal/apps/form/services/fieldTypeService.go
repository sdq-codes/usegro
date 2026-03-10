package services

import (
	"context"
	"errors"
	"github.com/sdq-codes/usegro-api/internal/apps/form/models"
	"github.com/sdq-codes/usegro-api/internal/apps/form/repositories"
	"gorm.io/gorm"
)

type FieldTypeService interface {
	GetByName(ctx context.Context, name string) (*models.FieldType, error)
	GetAll(ctx context.Context) ([]models.FieldType, error)
	ValidateExists(ctx context.Context, name string) error
}

type fieldTypeService struct {
	repo repositories.FieldTypeRepository
	db   *gorm.DB
}

func NewFieldTypeService(db *gorm.DB) FieldTypeService {
	return &fieldTypeService{repo: repositories.NewFieldTypeRepository(), db: db}
}

func (s *fieldTypeService) GetByName(ctx context.Context, name string) (*models.FieldType, error) {
	if name == "" {
		return nil, errors.New("field type name is required")
	}

	return s.repo.GetByName(ctx, s.db, name)
}

func (s *fieldTypeService) GetAll(ctx context.Context) ([]models.FieldType, error) {
	return s.repo.GetAll(ctx, s.db)
}

func (s *fieldTypeService) ValidateExists(ctx context.Context, name string) error {
	_, err := s.GetByName(ctx, name)
	if err != nil {
		return errors.New("field type does not exist")
	}
	return nil
}
