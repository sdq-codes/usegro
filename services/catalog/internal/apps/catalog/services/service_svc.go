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

type CatalogServiceService struct {
	repo repositories.ServiceRepositoryInterface
}

func NewCatalogServiceService(db *gorm.DB) *CatalogServiceService {
	return &CatalogServiceService{repo: repositories.NewServiceRepository(db)}
}

func (s *CatalogServiceService) CreateService(ctx context.Context, crmID string, d dto.CreateServiceDTO) (*models.CatalogItem, error) {
	item, err := s.repo.CreateService(ctx, crmID, d)
	if err != nil {
		logger.Log.Error("service could not be created", zap.Error(err))
		return nil, fmt.Errorf("service could not be created")
	}
	return item, nil
}

func (s *CatalogServiceService) ListServices(ctx context.Context, crmID string, search string, status string) ([]models.CatalogItem, error) {
	items, err := s.repo.ListServices(ctx, crmID, search, status)
	if err != nil {
		logger.Log.Error("services could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("services could not be fetched")
	}
	return items, nil
}

func (s *CatalogServiceService) GetService(ctx context.Context, crmID string, itemID string) (*models.CatalogItem, error) {
	item, err := s.repo.GetService(ctx, crmID, itemID)
	if err != nil {
		logger.Log.Error("service could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("service not found")
	}
	return item, nil
}

func (s *CatalogServiceService) UpdateService(ctx context.Context, crmID string, itemID string, d dto.UpdateServiceDTO) (*models.CatalogItem, error) {
	item, err := s.repo.UpdateService(ctx, crmID, itemID, d)
	if err != nil {
		logger.Log.Error("service could not be updated", zap.Error(err))
		return nil, fmt.Errorf("service could not be updated")
	}
	return item, nil
}

func (s *CatalogServiceService) DeleteService(ctx context.Context, crmID string, itemID string) error {
	err := s.repo.DeleteService(ctx, crmID, itemID)
	if err != nil {
		logger.Log.Error("service could not be deleted", zap.Error(err))
		return fmt.Errorf("service could not be deleted")
	}
	return nil
}
