package repositories

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"github.com/usegro/services/crm/internal/logger"
	"gorm.io/gorm"
)

type CRMUserOrganizationRepositoryInterface interface {
	CreateCRMUserOrganization(ctx context.Context, tx *gorm.DB, organization *models.CrmUserOrganization) (*models.CrmUserOrganization, error)
	GetCRMUserOrganization(ctx context.Context, tx *gorm.DB, organizationId string) (*models.CrmUserOrganization, error)
	FetchCRMUserOrganization(ctx context.Context, tx *gorm.DB, userId string) (*[]models.CrmUserOrganization, error)
	IsBusinessNameExist(ctx context.Context, tx *gorm.DB, businessName string) (bool, error)
	UpdateCRMUserOrganization(ctx context.Context, tx *gorm.DB, fullName, businessName string, businessInfo *models.CrmUserOrganization) (*models.CrmUserOrganization, error)
	UpdateCRMUserOrganizationStatus(ctx context.Context, tx *gorm.DB, organizationId string) error
	CreateCRMUserOrganizationSalesChannelType(ctx context.Context, tx *gorm.DB, crmUserOrganizationId string, salesChannel []models.SalesChannel) (*[]models.SalesChannelType, error)
	UpdateCRMUserOrganizationSalesChannelType(ctx context.Context, tx *gorm.DB, crmUserOrganizationId string, salesChannel []models.SalesChannel) (*[]models.SalesChannelType, error)
	CreateCRMUserOrganizationStockProductType(ctx context.Context, tx *gorm.DB, crmUserOrganizationId string, productType []models.ProductType) (*[]models.StockProductType, error)
	UpdateCRMUserOrganizationStockProductType(ctx context.Context, tx *gorm.DB, crmUserOrganizationId string, productType []models.ProductType) (*[]models.StockProductType, error)
}

type CRMUserOrganizationRepository struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewCRMUserOrganizationRepository(db *gorm.DB, rdb *redis.Client) *CRMUserOrganizationRepository {
	return &CRMUserOrganizationRepository{
		db:  db,
		rdb: rdb,
	}
}

func (crm *CRMUserOrganizationRepository) CreateCRMUserOrganization(ctx context.Context, tx *gorm.DB, organization *models.CrmUserOrganization) error {
	return tx.WithContext(ctx).Create(organization).Error
}

func (crm *CRMUserOrganizationRepository) UpdateCRMUserOrganization(ctx context.Context, tx *gorm.DB, businessInfo *models.CrmUserOrganization) error {
	return tx.WithContext(ctx).Save(businessInfo).Error
}

func (crm *CRMUserOrganizationRepository) GetCRMUserOrganization(ctx context.Context, tx *gorm.DB, organizationId string) (*models.CrmUserOrganization, error) {
	crmUserOrganization := &models.CrmUserOrganization{}
	err := tx.WithContext(ctx).First(crmUserOrganization, "id = ?", organizationId).Error
	return crmUserOrganization, err
}

func (crm *CRMUserOrganizationRepository) FetchCRMUserOrganization(ctx context.Context, tx *gorm.DB, userId string) (*[]models.CrmUserOrganization, error) {
	crmUserOrganization := &[]models.CrmUserOrganization{}
	logger.Log.Info("Fetching CRM User Organization")
	err := tx.WithContext(ctx).Find(crmUserOrganization, "user_id = ?", userId).Error
	logger.Log.Info(fmt.Sprintf("%v", crmUserOrganization))
	return crmUserOrganization, err
}

func (crm *CRMUserOrganizationRepository) IsBusinessNameExist(ctx context.Context, tx *gorm.DB, businessName string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&models.CrmUserOrganization{}).Where("business_name = ?", businessName).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (crm *CRMUserOrganizationRepository) UpdateCRMUserOrganizationStatus(ctx context.Context, tx *gorm.DB, organizationId string) error {
	organization, err := crm.GetCRMUserOrganization(ctx, tx, organizationId)
	if err != nil {
		return err
	}
	organization.Active = !organization.Active

	return tx.WithContext(ctx).Save(&organization).Error
}

func (crm *CRMUserOrganizationRepository) CreateCRMUserOrganizationSalesChannelType(ctx context.Context, tx *gorm.DB, salesChannel []models.SalesChannelType) (*[]models.SalesChannelType, error) {
	err := tx.WithContext(ctx).Create(&salesChannel).Error
	return &salesChannel, err
}

func (crm *CRMUserOrganizationRepository) UpdateCRMUserOrganizationSalesChannelType(ctx context.Context, tx *gorm.DB, crmUserOrganizationId string, salesChannel []models.SalesChannelType) (*[]models.SalesChannelType, error) {
	err := tx.WithContext(ctx).Where("crm_user_organization_id = ?", crmUserOrganizationId).Delete(&models.SalesChannelType{}).Error
	if err != nil {
		return nil, err
	}
	return crm.CreateCRMUserOrganizationSalesChannelType(ctx, tx, salesChannel)
}

func (crm *CRMUserOrganizationRepository) CreateCRMUserOrganizationStockProductType(ctx context.Context, tx *gorm.DB, productType []models.StockProductType) (*[]models.StockProductType, error) {
	err := tx.WithContext(ctx).Create(&productType).Error
	return &productType, err
}

func (crm *CRMUserOrganizationRepository) UpdateCRMUserOrganizationStockProductType(ctx context.Context, tx *gorm.DB, crmUserOrganizationId string, productType []models.StockProductType) (*[]models.StockProductType, error) {
	err := tx.WithContext(ctx).Where("crm_user_organization_id = ?", crmUserOrganizationId).Delete(&models.StockProductType{}).Error
	if err != nil {
		return nil, err
	}
	return crm.CreateCRMUserOrganizationStockProductType(ctx, tx, productType)
}
