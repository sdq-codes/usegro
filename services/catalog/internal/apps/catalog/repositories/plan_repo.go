package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"gorm.io/gorm"
)

type PlanRepositoryInterface interface {
	CreatePlan(ctx context.Context, crmID string, catalogItemID string, d dto.CreatePlanDTO) (*models.Plan, error)
	ListPlans(ctx context.Context, crmID string, catalogItemID string) ([]models.Plan, error)
	GetPlan(ctx context.Context, crmID string, planID string) (*models.Plan, error)
	UpdatePlan(ctx context.Context, crmID string, planID string, d dto.UpdatePlanDTO) (*models.Plan, error)
	DeletePlan(ctx context.Context, crmID string, planID string) error
}

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) PlanRepositoryInterface {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) CreatePlan(ctx context.Context, crmID string, catalogItemID string, d dto.CreatePlanDTO) (*models.Plan, error) {
	parsedCRMID, err := uuid.Parse(crmID)
	if err != nil {
		return nil, fmt.Errorf("invalid crm_id: %w", err)
	}
	parsedItemID, err := uuid.Parse(catalogItemID)
	if err != nil {
		return nil, fmt.Errorf("invalid catalog_item_id: %w", err)
	}

	currency := d.PriceCurrency
	if currency == "" {
		currency = "USD"
	}

	plan := models.Plan{
		CatalogItemID: parsedItemID,
		CRMID:         parsedCRMID,
		Name:          d.Name,
		PlanType:      d.PlanType,
		Price:         d.Price,
		PriceCurrency: currency,
		BillingCycle:  d.BillingCycle,
		SessionCount:  d.SessionCount,
		ValidityDays:  d.ValidityDays,
		Status:        "active",
	}
	if err := r.db.WithContext(ctx).Create(&plan).Error; err != nil {
		return nil, fmt.Errorf("failed to create plan: %w", err)
	}
	return &plan, nil
}

func (r *PlanRepository) ListPlans(ctx context.Context, crmID string, catalogItemID string) ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.WithContext(ctx).
		Where("crm_id = ? AND catalog_item_id = ? AND status != 'archived'", crmID, catalogItemID).
		Order("created_at ASC").
		Find(&plans).Error
	if err != nil {
		return nil, fmt.Errorf("failed to list plans: %w", err)
	}
	return plans, nil
}

func (r *PlanRepository) GetPlan(ctx context.Context, crmID string, planID string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.WithContext(ctx).Where("id = ? AND crm_id = ?", planID, crmID).First(&plan).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}
	return &plan, nil
}

func (r *PlanRepository) UpdatePlan(ctx context.Context, crmID string, planID string, d dto.UpdatePlanDTO) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.WithContext(ctx).Where("id = ? AND crm_id = ?", planID, crmID).First(&plan).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("plan not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find plan: %w", err)
	}

	currency := d.PriceCurrency
	if currency == "" {
		currency = plan.PriceCurrency
	}

	if err := r.db.WithContext(ctx).Model(&plan).Updates(map[string]interface{}{
		"name":           d.Name,
		"plan_type":      d.PlanType,
		"price":          d.Price,
		"price_currency": currency,
		"billing_cycle":  d.BillingCycle,
		"session_count":  d.SessionCount,
		"validity_days":  d.ValidityDays,
	}).Error; err != nil {
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}
	return &plan, nil
}

func (r *PlanRepository) DeletePlan(ctx context.Context, crmID string, planID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND crm_id = ?", planID, crmID).Delete(&models.Plan{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete plan: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("plan not found")
	}
	return nil
}
