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

type ProductService struct {
	repo repositories.ProductRepositoryInterface
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{repo: repositories.NewProductRepository(db)}
}

func (s *ProductService) CreateProduct(ctx context.Context, crmID string, d dto.CreateProductDTO) (*models.CatalogItem, error) {
	item, err := s.repo.CreateProduct(ctx, crmID, d)
	if err != nil {
		logger.Log.Error("product could not be created", zap.Error(err))
		return nil, fmt.Errorf("product could not be created")
	}
	return item, nil
}

func (s *ProductService) ListProducts(ctx context.Context, crmID string, search string, status string, page, limit int) (*repositories.PaginatedProducts, error) {
	result, err := s.repo.ListProducts(ctx, crmID, search, status, page, limit)
	if err != nil {
		logger.Log.Error("products could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("products could not be fetched")
	}
	return result, nil
}

func (s *ProductService) GetProduct(ctx context.Context, crmID string, itemID string) (*models.CatalogItem, error) {
	item, err := s.repo.GetProduct(ctx, crmID, itemID)
	if err != nil {
		logger.Log.Error("product could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("product not found")
	}
	return item, nil
}

func (s *ProductService) UpdateProduct(ctx context.Context, crmID string, itemID string, d dto.UpdateProductDTO) (*models.CatalogItem, error) {
	item, err := s.repo.UpdateProduct(ctx, crmID, itemID, d)
	if err != nil {
		logger.Log.Error("product could not be updated", zap.Error(err))
		return nil, fmt.Errorf("product could not be updated")
	}
	return item, nil
}

func (s *ProductService) DeleteProduct(ctx context.Context, crmID string, itemID string) error {
	err := s.repo.DeleteProduct(ctx, crmID, itemID)
	if err != nil {
		logger.Log.Error("product could not be deleted", zap.Error(err))
		return fmt.Errorf("product could not be deleted")
	}
	return nil
}
