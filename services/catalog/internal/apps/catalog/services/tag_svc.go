package catalogServices

import (
	"context"
	"fmt"

	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"github.com/usegro/services/catalog/internal/apps/catalog/repositories"
	"github.com/usegro/services/catalog/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TagService struct {
	repo repositories.TagRepositoryInterface
}

func NewTagService(db *gorm.DB) *TagService {
	return &TagService{repo: repositories.NewTagRepository(db)}
}

func (s *TagService) CreateTag(ctx context.Context, crmID string, name string) (*models.CatalogTag, error) {
	tag, err := s.repo.CreateTag(ctx, crmID, name)
	if err != nil {
		logger.Log.Error("tag could not be created", zap.Error(err))
		return nil, fmt.Errorf("tag could not be created")
	}
	return tag, nil
}

func (s *TagService) ListTags(ctx context.Context, crmID string) ([]models.CatalogTag, error) {
	tags, err := s.repo.ListTags(ctx, crmID)
	if err != nil {
		logger.Log.Error("tags could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("tags could not be fetched")
	}
	return tags, nil
}

func (s *TagService) DeleteTag(ctx context.Context, crmID string, tagID string) error {
	if err := s.repo.DeleteTag(ctx, crmID, tagID); err != nil {
		logger.Log.Error("tag could not be deleted", zap.Error(err))
		return fmt.Errorf("tag could not be deleted")
	}
	return nil
}
