package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"gorm.io/gorm"
)

type CategoryRepositoryInterface interface {
	CreateCategory(ctx context.Context, crmID string, d dto.CreateCategoryDTO) (*models.Category, error)
	ListCategories(ctx context.Context, crmID string) ([]models.Category, error)
	GetCategory(ctx context.Context, crmID string, categoryID string) (*models.Category, error)
	UpdateCategory(ctx context.Context, crmID string, categoryID string, d dto.UpdateCategoryDTO) (*models.Category, error)
	DeleteCategory(ctx context.Context, crmID string, categoryID string) error
}

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepositoryInterface {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) CreateCategory(ctx context.Context, crmID string, d dto.CreateCategoryDTO) (*models.Category, error) {
	parsedCRMID, err := uuid.Parse(crmID)
	if err != nil {
		return nil, fmt.Errorf("invalid crm_id: %w", err)
	}

	category := models.Category{
		CRMID: parsedCRMID,
		Name:  d.Name,
	}

	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) ListCategories(ctx context.Context, crmID string) ([]models.Category, error) {
	var categories []models.Category
	if err := r.db.WithContext(ctx).
		Where("crm_id = ?", crmID).
		Order("created_at DESC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}
	return categories, nil
}

func (r *CategoryRepository) GetCategory(ctx context.Context, crmID string, categoryID string) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Where("id = ? AND crm_id = ?", categoryID, crmID).
		First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (r *CategoryRepository) UpdateCategory(ctx context.Context, crmID string, categoryID string, d dto.UpdateCategoryDTO) (*models.Category, error) {
	var category models.Category
	err := r.db.WithContext(ctx).
		Where("id = ? AND crm_id = ?", categoryID, crmID).
		First(&category).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find category: %w", err)
	}

	if err := r.db.WithContext(ctx).Model(&category).Update("name", d.Name).Error; err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func (r *CategoryRepository) DeleteCategory(ctx context.Context, crmID string, categoryID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND crm_id = ?", categoryID, crmID).Delete(&models.Category{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}
