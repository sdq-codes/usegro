package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"gorm.io/gorm"
)

type TagRepositoryInterface interface {
	CreateTag(ctx context.Context, crmID string, name string) (*models.CatalogTag, error)
	ListTags(ctx context.Context, crmID string) ([]models.CatalogTag, error)
	DeleteTag(ctx context.Context, crmID string, tagID string) error
}

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepositoryInterface {
	return &TagRepository{db: db}
}

func (r *TagRepository) CreateTag(ctx context.Context, crmID string, name string) (*models.CatalogTag, error) {
	parsedCRMID, err := uuid.Parse(crmID)
	if err != nil {
		return nil, fmt.Errorf("invalid crm_id: %w", err)
	}
	tag := models.CatalogTag{CRMID: parsedCRMID, Name: name}
	if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}
	return &tag, nil
}

func (r *TagRepository) ListTags(ctx context.Context, crmID string) ([]models.CatalogTag, error) {
	var tags []models.CatalogTag
	if err := r.db.WithContext(ctx).Where("crm_id = ?", crmID).Order("name ASC").Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	return tags, nil
}

func (r *TagRepository) DeleteTag(ctx context.Context, crmID string, tagID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND crm_id = ?", tagID, crmID).Delete(&models.CatalogTag{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete tag: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("tag not found")
	}
	return nil
}
