package services

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/redis/go-redis/v9"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
	formModels "github.com/usegro/services/crm/internal/apps/form/models"
	"github.com/usegro/services/crm/internal/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CrmCustomerService struct {
	crmCustomerRepo *repositories.CRMCustomerRepository
	dynamo          *dynamodb.Client
}

func NewCrmCustomerService(db *gorm.DB, rdb *redis.Client, dynamo *dynamodb.Client) *CrmCustomerService {
	return &CrmCustomerService{
		crmCustomerRepo: repositories.NewCRMCustomerRepository(db, rdb),
		dynamo:          dynamo,
	}
}

// FetchPublishedCreateCustomerForm fetches the published "create_customer" form for the given crmId.
// If none is found for that crmId, it falls back to the "global" form.
func (c *CrmCustomerService) FetchPublishedCreateCustomerForm(ctx context.Context, crmId string) (*formModels.CompleteForm, error) {
	// Try crmId first
	form, err := c.crmCustomerRepo.FetchPublishedCreateCustomerForm(ctx, c.dynamo, crmId)
	if err == nil && form != nil {
		return form, nil
	}
	return form, nil
}

// FetchCrmCustomers fetches the crm "customers" for the given crmId.
func (c *CrmCustomerService) FetchCrmCustomers(ctx context.Context, crmId string) (*[]formModels.FormSubmission, error) {
	customers, err := c.crmCustomerRepo.FetchCrmCustomers(ctx, c.dynamo, crmId)
	if err != nil {
		logger.Log.Error("CRM customers could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("CRM customers could not be fetched")
	}
	return &customers, nil
}

func (c *CrmCustomerService) ArchiveCrmCustomer(ctx context.Context, submissionID, formId, crmId string) error {
	err := c.crmCustomerRepo.ArchiveCrmCustomer(ctx, c.dynamo, submissionID, formId, crmId)
	if err != nil {
		logger.Log.Error("CRM customers could not be fetched", zap.Error(err))
		return fmt.Errorf("CRM customers could not be deleted")
	}
	return nil
}

func (c *CrmCustomerService) GetCrmCustomer(ctx context.Context, submissionID, formId, crmId string) (*formModels.FormSubmission, error) {
	customer, err := c.crmCustomerRepo.GetCrmCustomer(ctx, c.dynamo, submissionID, formId, crmId)
	if err != nil {
		logger.Log.Error("CRM customers could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("CRM customers could not be deleted")
	}
	return customer, nil
}
