package services

import (
	"context"

	"github.com/usegro/services/crm/internal/apps/crm/models"
	"github.com/usegro/services/crm/internal/apps/crm/repositories"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerActivityService struct {
	activityRepo repositories.CustomerActivityRepositoryInterface
}

func NewCustomerActivityService(mongoDB *mongo.Database) *CustomerActivityService {
	return &CustomerActivityService{
		activityRepo: repositories.NewCustomerActivityRepository(mongoDB),
	}
}

func (s *CustomerActivityService) FetchCustomerActivity(ctx context.Context, customerID string) ([]models.CustomerActivity, error) {
	return s.activityRepo.FetchCustomerActivity(ctx, customerID)
}

func (s *CustomerActivityService) LogComment(ctx context.Context, customerID, crmID, comment, performedBy string) (*models.CustomerActivity, error) {
	activity := models.CustomerActivity{
		ActivityType: models.ActivityTypeComment,
		Description:  comment,
		CrmID:        crmID,
		CustomerID:   customerID,
		PerformedBy:  performedBy,
	}
	if err := s.activityRepo.LogActivity(ctx, activity); err != nil {
		return nil, err
	}
	return &activity, nil
}
