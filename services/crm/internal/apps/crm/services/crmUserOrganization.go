package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/internal/apps/crm/dto"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
	"github.com/usegro/services/crm/pkg/exception"
	"gorm.io/gorm"
)

type CRMUserOrganizationService struct {
	db   *gorm.DB
	rdb  *redis.Client
	repo *repositories.CRMUserOrganizationRepository
}

func NewCRMUserOrganizationService(db *gorm.DB, rdb *redis.Client) *CRMUserOrganizationService {
	return &CRMUserOrganizationService{
		db:   db,
		rdb:  rdb,
		repo: repositories.NewCRMUserOrganizationRepository(db, rdb),
	}
}

func (s *CRMUserOrganizationService) CreateCRMUserOrganization(ctx context.Context, dti *dto.CrmUserOrganizationDTI, userId uuid.UUID) (*models.CrmUserOrganization, error) {
	exists, err := s.repo.IsBusinessNameExist(ctx, s.db, dti.BusinessName)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, exception.BusinessNameAlreadyTakenError
	}

	tx := s.db.Begin()
	CRMUserOrganization := &models.CrmUserOrganization{
		UserID:       userId,
		FullName:     dti.FullName,
		BusinessName: dti.BusinessName,
		BusinessInfo: dti.BusinessInfo,
	}
	if err := s.repo.CreateCRMUserOrganization(ctx, tx, CRMUserOrganization); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return CRMUserOrganization, nil
}

func (s *CRMUserOrganizationService) UpdateCRMUserOrganization(ctx context.Context, dti *dto.CrmUserOrganizationDTI, crmUserOrganizationId uuid.UUID) (*models.CrmUserOrganization, error) {
	tx := s.db.Begin()
	crmUserOrganization, err := s.GetCRMUserOrganization(ctx, crmUserOrganizationId.String())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if dti.FullName != "" {
		crmUserOrganization.FullName = dti.FullName
	}
	if dti.BusinessName != "" {
		crmUserOrganization.BusinessName = dti.BusinessName
	}
	if dti.BusinessInfo != "" {
		crmUserOrganization.BusinessInfo = dti.BusinessInfo
	}
	if err := s.repo.UpdateCRMUserOrganization(ctx, tx, crmUserOrganization); err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return crmUserOrganization, nil
}

func (s *CRMUserOrganizationService) GetCRMUserOrganization(ctx context.Context, crmUserOrganizationId string) (*models.CrmUserOrganization, error) {
	return s.repo.GetCRMUserOrganization(ctx, s.db, crmUserOrganizationId)
}

func (s *CRMUserOrganizationService) FetchCRMUserOrganization(ctx context.Context, userId string) (*[]models.CrmUserOrganization, error) {
	return s.repo.FetchCRMUserOrganization(ctx, s.db, userId)
}

func (s *CRMUserOrganizationService) IsBusinessNameExist(ctx context.Context, businessName string) (bool, error) {
	return s.repo.IsBusinessNameExist(ctx, s.db, businessName)
}

func (s *CRMUserOrganizationService) ToggleCRMUserOrganizationStatus(ctx context.Context, crmUserOrganizationId string) error {
	return s.repo.UpdateCRMUserOrganizationStatus(ctx, s.db, crmUserOrganizationId)
}

func (s *CRMUserOrganizationService) CreateSalesChannelType(ctx context.Context, crmUserOrganizationId uuid.UUID, salesChannels dto.CrmUserOrganizationSalesChannelTypeDTI) (*[]models.SalesChannelType, error) {
	var salesChannelTypes []models.SalesChannelType

	for _, channel := range salesChannels.SalesChannelType {
		salesChannelTypes = append(salesChannelTypes, models.SalesChannelType{
			CrmUserOrganizationID: crmUserOrganizationId,
			SalesChannelType:      models.SalesChannel(channel),
		})
	}
	return s.repo.CreateCRMUserOrganizationSalesChannelType(ctx, s.db, salesChannelTypes)
}

func (s *CRMUserOrganizationService) UpdateSalesChannelType(ctx context.Context, crmUserOrganizationId uuid.UUID, salesChannels dto.CrmUserOrganizationSalesChannelTypeDTI) (*[]models.SalesChannelType, error) {
	var salesChannelTypes []models.SalesChannelType

	for _, channel := range salesChannels.SalesChannelType {
		salesChannelTypes = append(salesChannelTypes, models.SalesChannelType{
			CrmUserOrganizationID: crmUserOrganizationId,
			SalesChannelType:      models.SalesChannel(channel),
		})
	}
	return s.repo.UpdateCRMUserOrganizationSalesChannelType(ctx, s.db, crmUserOrganizationId.String(), salesChannelTypes)
}

func (s *CRMUserOrganizationService) CreateStockProductType(ctx context.Context, crmUserOrganizationId uuid.UUID, productTypes dto.CrmUserOrganizationStockProductTypeDTI) (*[]models.StockProductType, error) {
	var stockProductTypes []models.StockProductType

	for _, product := range productTypes.ProductType {
		stockProductTypes = append(stockProductTypes, models.StockProductType{
			CrmUserOrganizationID: crmUserOrganizationId,
			ProductType:           string(product),
		})
	}
	return s.repo.CreateCRMUserOrganizationStockProductType(ctx, s.db, stockProductTypes)
}

func (s *CRMUserOrganizationService) UpdateStockProductType(ctx context.Context, crmUserOrganizationId uuid.UUID, productTypes dto.CrmUserOrganizationStockProductTypeDTI) (*[]models.StockProductType, error) {
	var stockProductTypes []models.StockProductType

	for _, product := range productTypes.ProductType {
		stockProductTypes = append(stockProductTypes, models.StockProductType{
			CrmUserOrganizationID: crmUserOrganizationId,
			ProductType:           string(product),
		})
	}
	return s.repo.CreateCRMUserOrganizationStockProductType(ctx, s.db, stockProductTypes)
}
