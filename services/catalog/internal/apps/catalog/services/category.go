package catalogServices

import (
	"context"
	"fmt"

	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/apps/catalog/repositories"
	"github.com/usegro/services/catalog/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CategoryService struct {
	categoryRepo repositories.CategoryRepositoryInterface
}

func NewCategoryService(db *gorm.DB) *CategoryService {
	return &CategoryService{
		categoryRepo: repositories.NewCategoryRepository(db),
	}
}

func (s *CategoryService) CreateCategory(ctx context.Context, crmID string, d dto.CreateCategoryDTO) (*models.Category, error) {
	category, err := s.categoryRepo.CreateCategory(ctx, crmID, d)
	if err != nil {
		logger.Log.Error("category could not be created", zap.Error(err))
		return nil, fmt.Errorf("category could not be created")
	}
	return category, nil
}

func (s *CategoryService) ListCategories(ctx context.Context, crmID string) ([]models.Category, error) {
	categories, err := s.categoryRepo.ListCategories(ctx, crmID)
	if err != nil {
		logger.Log.Error("categories could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("categories could not be fetched")
	}
	return categories, nil
}

func (s *CategoryService) GetCategory(ctx context.Context, crmID string, categoryID string) (*models.Category, error) {
	category, err := s.categoryRepo.GetCategory(ctx, crmID, categoryID)
	if err != nil {
		logger.Log.Error("category could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

func (s *CategoryService) UpdateCategory(ctx context.Context, crmID string, categoryID string, d dto.UpdateCategoryDTO) (*models.Category, error) {
	category, err := s.categoryRepo.UpdateCategory(ctx, crmID, categoryID, d)
	if err != nil {
		logger.Log.Error("category could not be updated", zap.Error(err))
		return nil, fmt.Errorf("category could not be updated")
	}
	return category, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, crmID string, categoryID string) error {
	err := s.categoryRepo.DeleteCategory(ctx, crmID, categoryID)
	if err != nil {
		logger.Log.Error("category could not be deleted", zap.Error(err))
		return fmt.Errorf("category could not be deleted")
	}
	return nil
}
