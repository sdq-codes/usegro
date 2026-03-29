package catalogServices

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/apps/catalog/repositories"
	"github.com/usegro/services/catalog/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type StandardCategoryService struct {
	repo repositories.StandardCategoryRepositoryInterface
}

func NewStandardCategoryService(db *gorm.DB) *StandardCategoryService {
	return &StandardCategoryService{
		repo: repositories.NewStandardCategoryRepository(db),
	}
}

func (s *StandardCategoryService) ListRootCategories(ctx context.Context) ([]models.StandardCategory, error) {
	categories, err := s.repo.ListRootCategories(ctx)
	if err != nil {
		logger.Log.Error("failed to list root standard categories", zap.Error(err))
		return nil, fmt.Errorf("could not fetch categories")
	}
	return categories, nil
}

func (s *StandardCategoryService) GetCategoryWithChildren(ctx context.Context, id string) (*models.StandardCategory, error) {
	parsed, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid category id")
	}
	category, err := s.repo.GetCategoryWithChildren(ctx, parsed)
	if err != nil {
		logger.Log.Error("failed to get standard category", zap.Error(err))
		return nil, fmt.Errorf("category not found")
	}
	return category, nil
}

type StandardCategorySearchResult struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

func (s *StandardCategoryService) SearchCategories(ctx context.Context, query string) ([]StandardCategorySearchResult, error) {
	matches, err := s.repo.Search(ctx, query)
	if err != nil {
		logger.Log.Error("failed to search standard categories", zap.Error(err))
		return nil, fmt.Errorf("could not search categories")
	}

	results := make([]StandardCategorySearchResult, 0, len(matches))
	for _, cat := range matches {
		path := cat.Name
		parentID := cat.ParentID
		for parentID != nil {
			parent, err := s.repo.GetByID(ctx, *parentID)
			if err != nil {
				break
			}
			path = parent.Name + " › " + path
			parentID = parent.ParentID
		}
		results = append(results, StandardCategorySearchResult{
			ID:   cat.ID.String(),
			Name: cat.Name,
			Path: path,
		})
	}
	return results, nil
}

func (s *StandardCategoryService) ListByParent(ctx context.Context, parentID string) ([]models.StandardCategory, error) {
	parsed, err := uuid.Parse(parentID)
	if err != nil {
		return nil, fmt.Errorf("invalid parent_id")
	}
	categories, err := s.repo.ListByParent(ctx, parsed)
	if err != nil {
		logger.Log.Error("failed to list categories by parent", zap.Error(err))
		return nil, fmt.Errorf("could not fetch categories")
	}
	return categories, nil
}
