package services

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
	formModels "github.com/usegro/services/crm/internal/apps/form/models"
	"github.com/usegro/services/crm/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CrmCustomerService struct {
	crmCustomerRepo repositories.CRMCustomerRepositoryInterface
}

func NewCrmCustomerService(db *gorm.DB, rdb *redis.Client, mongoDB *mongo.Database) *CrmCustomerService {
	return &CrmCustomerService{
		crmCustomerRepo: repositories.NewCRMCustomerRepository(db, rdb, mongoDB),
	}
}

func (c *CrmCustomerService) FetchPublishedCreateCustomerForm(ctx context.Context, crmId string) (*formModels.CompleteForm, error) {
	return c.crmCustomerRepo.FetchPublishedCreateCustomerForm(ctx, crmId)
}

func (c *CrmCustomerService) FetchCrmCustomers(ctx context.Context, crmId string, page, limit int) (*repositories.PaginatedCustomers, error) {
	result, err := c.crmCustomerRepo.FetchCrmCustomers(ctx, crmId, page, limit)
	if err != nil {
		logger.Log.Error("CRM customers could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("CRM customers could not be fetched")
	}
	return result, nil
}

func (c *CrmCustomerService) ArchiveCrmCustomer(ctx context.Context, submissionID, formId, crmId string) error {
	err := c.crmCustomerRepo.ArchiveCrmCustomer(ctx, submissionID, formId, crmId)
	if err != nil {
		logger.Log.Error("CRM customer could not be archived", zap.Error(err))
		return fmt.Errorf("CRM customers could not be deleted")
	}
	return nil
}

func (c *CrmCustomerService) GetCrmCustomer(ctx context.Context, submissionID, formId, crmId string) (*formModels.FormSubmission, error) {
	customer, err := c.crmCustomerRepo.GetCrmCustomer(ctx, submissionID, formId, crmId)
	if err != nil {
		logger.Log.Error("CRM customer could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("CRM customers could not be deleted")
	}
	return customer, nil
}
