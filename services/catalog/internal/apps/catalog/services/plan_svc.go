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

type PlanService struct {
	repo repositories.PlanRepositoryInterface
}

func NewPlanService(db *gorm.DB) *PlanService {
	return &PlanService{repo: repositories.NewPlanRepository(db)}
}

func (s *PlanService) CreatePlan(ctx context.Context, crmID string, catalogItemID string, d dto.CreatePlanDTO) (*models.Plan, error) {
	plan, err := s.repo.CreatePlan(ctx, crmID, catalogItemID, d)
	if err != nil {
		logger.Log.Error("plan could not be created", zap.Error(err))
		return nil, fmt.Errorf("plan could not be created")
	}
	return plan, nil
}

func (s *PlanService) ListPlans(ctx context.Context, crmID string, catalogItemID string) ([]models.Plan, error) {
	plans, err := s.repo.ListPlans(ctx, crmID, catalogItemID)
	if err != nil {
		logger.Log.Error("plans could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("plans could not be fetched")
	}
	return plans, nil
}

func (s *PlanService) GetPlan(ctx context.Context, crmID string, planID string) (*models.Plan, error) {
	plan, err := s.repo.GetPlan(ctx, crmID, planID)
	if err != nil {
		logger.Log.Error("plan could not be fetched", zap.Error(err))
		return nil, fmt.Errorf("plan not found")
	}
	return plan, nil
}

func (s *PlanService) UpdatePlan(ctx context.Context, crmID string, planID string, d dto.UpdatePlanDTO) (*models.Plan, error) {
	plan, err := s.repo.UpdatePlan(ctx, crmID, planID, d)
	if err != nil {
		logger.Log.Error("plan could not be updated", zap.Error(err))
		return nil, fmt.Errorf("plan could not be updated")
	}
	return plan, nil
}

func (s *PlanService) DeletePlan(ctx context.Context, crmID string, planID string) error {
	if err := s.repo.DeletePlan(ctx, crmID, planID); err != nil {
		logger.Log.Error("plan could not be deleted", zap.Error(err))
		return fmt.Errorf("plan could not be deleted")
	}
	return nil
}
