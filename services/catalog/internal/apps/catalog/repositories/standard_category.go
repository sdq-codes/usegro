package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"gorm.io/gorm"
)

type StandardCategoryRepositoryInterface interface {
	ListRootCategories(ctx context.Context) ([]models.StandardCategory, error)
	GetCategoryWithChildren(ctx context.Context, id uuid.UUID) (*models.StandardCategory, error)
	ListByParent(ctx context.Context, parentID uuid.UUID) ([]models.StandardCategory, error)
	UpsertBatch(ctx context.Context, categories []models.StandardCategory) error
	Search(ctx context.Context, query string) ([]models.StandardCategory, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.StandardCategory, error)
}

type StandardCategoryRepository struct {
	db *gorm.DB
}

// isLeafExpr computes is_leaf dynamically: true when no child row references this row as parent.
const isLeafExpr = "NOT EXISTS(SELECT 1 FROM standard_categories child WHERE child.parent_id = standard_categories.id)"

// catColumns selects all columns, replacing the stored is_leaf with the computed value.
const catColumns = "id, parent_id, name, full_name, level, (" + isLeafExpr + ") AS is_leaf"

func NewStandardCategoryRepository(db *gorm.DB) StandardCategoryRepositoryInterface {
	return &StandardCategoryRepository{db: db}
}

func (r *StandardCategoryRepository) ListRootCategories(ctx context.Context) ([]models.StandardCategory, error) {
	var categories []models.StandardCategory
	if err := r.db.WithContext(ctx).
		Select(catColumns).
		Where("parent_id IS NULL").
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list root categories: %w", err)
	}
	return categories, nil
}

func (r *StandardCategoryRepository) GetCategoryWithChildren(ctx context.Context, id uuid.UUID) (*models.StandardCategory, error) {
	var category models.StandardCategory
	err := r.db.WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Select(catColumns).Order("name ASC")
		}).
		Preload("Attributes").
		First(&category, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("standard category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get standard category: %w", err)
	}
	return &category, nil
}

func (r *StandardCategoryRepository) ListByParent(ctx context.Context, parentID uuid.UUID) ([]models.StandardCategory, error) {
	var categories []models.StandardCategory
	if err := r.db.WithContext(ctx).
		Select(catColumns).
		Where("parent_id = ?", parentID).
		Order("name ASC").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to list categories by parent: %w", err)
	}
	return categories, nil
}

func (r *StandardCategoryRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.StandardCategory, error) {
	var category models.StandardCategory
	if err := r.db.WithContext(ctx).First(&category, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return &category, nil
}

func (r *StandardCategoryRepository) Search(ctx context.Context, query string) ([]models.StandardCategory, error) {
	var categories []models.StandardCategory
	if err := r.db.WithContext(ctx).
		Select(catColumns).
		Where("name ILIKE ?", "%"+query+"%").
		Order("name ASC").
		Limit(20).
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to search categories: %w", err)
	}
	return categories, nil
}

func (r *StandardCategoryRepository) UpsertBatch(ctx context.Context, categories []models.StandardCategory) error {
	if len(categories) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).
		Save(categories).Error
}
