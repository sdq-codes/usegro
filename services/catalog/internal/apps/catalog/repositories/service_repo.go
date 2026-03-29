package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/usegro/services/catalog/config"
	"github.com/usegro/services/catalog/internal/apps/catalog/dto"
	"github.com/usegro/services/catalog/internal/apps/catalog/models"
	"gorm.io/gorm"
)

type ServiceRepositoryInterface interface {
	CreateService(ctx context.Context, crmID string, d dto.CreateServiceDTO) (*models.CatalogItem, error)
	ListServices(ctx context.Context, crmID string, search string, status string) ([]models.CatalogItem, error)
	GetService(ctx context.Context, crmID string, itemID string) (*models.CatalogItem, error)
	UpdateService(ctx context.Context, crmID string, itemID string, d dto.UpdateServiceDTO) (*models.CatalogItem, error)
	DeleteService(ctx context.Context, crmID string, itemID string) error
}

type ServiceRepository struct {
	db *gorm.DB
}

func NewServiceRepository(db *gorm.DB) ServiceRepositoryInterface {
	return &ServiceRepository{db: db}
}

func (r *ServiceRepository) preload(tx *gorm.DB) *gorm.DB {
	return tx.
		Preload("ServiceDetail.Variations").
		Preload("ServiceDetail.Locations").
		Preload("Plans").
		Preload("Tags").
		Preload("AdditionalFields").
		Preload("Categories").
		Preload("Media")
}

func (r *ServiceRepository) CreateService(ctx context.Context, crmID string, d dto.CreateServiceDTO) (*models.CatalogItem, error) {
	parsedCRMID, err := uuid.Parse(crmID)
	if err != nil {
		return nil, fmt.Errorf("invalid crm_id: %w", err)
	}

	status := d.Status
	if status == "" {
		status = "active"
	}
	currency := d.PriceCurrency
	if currency == "" {
		currency = "USD"
	}

	item := models.CatalogItem{
		CRMID:         parsedCRMID,
		ItemType:      "service",
		Name:          d.Name,
		Description:   d.Description,
		Price:         d.Price,
		PriceCurrency: currency,
		Status:        status,
		ShowInStore:   d.ShowInStore,
	}
	if err := r.db.WithContext(ctx).Create(&item).Error; err != nil {
		return nil, fmt.Errorf("failed to create catalog item: %w", err)
	}

	if len(d.TagIDs) > 0 {
		var tags []models.CatalogTag
		for _, tid := range d.TagIDs {
			if pid, err := uuid.Parse(tid); err == nil {
				tags = append(tags, models.CatalogTag{ID: pid})
			}
		}
		r.db.WithContext(ctx).Model(&item).Association("Tags").Replace(tags)
	}

	serviceType := d.ServiceType
	if serviceType == "" {
		serviceType = "appointment"
	}
	duration := d.Duration
	if duration == "" {
		duration = "1 Hour"
	}
	bufferTime := d.BufferTime
	if bufferTime == "" {
		bufferTime = "No buffer"
	}
	priceType := d.PriceType
	if priceType == "" {
		priceType = "fixed"
	}
	paymentMode := d.PaymentMode
	if paymentMode == "" {
		paymentMode = "per-session"
	}
	bookingMode := d.BookingMode
	if bookingMode == "" {
		bookingMode = "auto"
	}

	detail := models.ServiceDetail{
		ItemID:           item.ID,
		ServiceType:      serviceType,
		Tagline:          d.Tagline,
		Duration:         duration,
		BufferTime:       bufferTime,
		PriceType:        priceType,
		PaymentMode:      paymentMode,
		CustomPriceLabel: d.CustomPriceLabel,
		BookingMode:      bookingMode,
	}
	if err := r.db.WithContext(ctx).Create(&detail).Error; err != nil {
		return nil, fmt.Errorf("failed to create service detail: %w", err)
	}

	for _, vDTO := range d.Variations {
		v := models.ServiceVariation{
			ServiceDetailID: detail.ID,
			Name:            vDTO.Name,
			Price:           vDTO.Price,
			Position:        vDTO.Position,
		}
		if err := r.db.WithContext(ctx).Create(&v).Error; err != nil {
			return nil, fmt.Errorf("failed to create variation: %w", err)
		}
	}

	for _, lDTO := range d.Locations {
		l := models.ServiceLocation{
			ServiceDetailID: detail.ID,
			LocationType:    lDTO.LocationType,
			Address:         lDTO.Address,
			PhoneMethod:     lDTO.PhoneMethod,
			Phone:           lDTO.Phone,
		}
		if err := r.db.WithContext(ctx).Create(&l).Error; err != nil {
			return nil, fmt.Errorf("failed to create location: %w", err)
		}
	}

	for i := range d.Plans {
		pDTO := d.Plans[i]
		currency := pDTO.PriceCurrency
		if currency == "" {
			currency = "USD"
		}
		p := models.Plan{
			CatalogItemID: item.ID,
			CRMID:         parsedCRMID,
			Name:          pDTO.Name,
			PlanType:      pDTO.PlanType,
			Price:         pDTO.Price,
			PriceCurrency: currency,
			BillingCycle:  pDTO.BillingCycle,
			SessionCount:  pDTO.SessionCount,
			ValidityDays:  pDTO.ValidityDays,
			Status:        "active",
		}
		if err := r.db.WithContext(ctx).Create(&p).Error; err != nil {
			return nil, fmt.Errorf("failed to create plan: %w", err)
		}
	}

	for _, fDTO := range d.AdditionalFields {
		f := models.CatalogAdditionalField{
			ItemID:    item.ID,
			Label:     fDTO.Label,
			FieldType: fDTO.FieldType,
			Value:     fDTO.Value,
			Position:  fDTO.Position,
		}
		if err := r.db.WithContext(ctx).Create(&f).Error; err != nil {
			return nil, fmt.Errorf("failed to create additional field: %w", err)
		}
	}

	r2cfg := config.GetConfig().R2
	for i, key := range d.MediaKeys {
		media := models.CatalogItemMedia{
			ItemID:       item.ID,
			URL:          fmt.Sprintf("%s/%s", r2cfg.PublicURL, key),
			Key:          key,
			Position:     i,
			DisplayImage: key == d.DisplayImageKey,
		}
		if err := r.db.WithContext(ctx).Create(&media).Error; err != nil {
			return nil, fmt.Errorf("failed to create media record: %w", err)
		}
	}
	if d.DisplayImageKey == "" && len(d.MediaKeys) > 0 {
		if err := r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
			Where("item_id = ? AND position = 0", item.ID).
			Update("display_image", true).Error; err != nil {
			return nil, fmt.Errorf("failed to set default display image: %w", err)
		}
	}

	var result models.CatalogItem
	if err := r.preload(r.db.WithContext(ctx)).First(&result, "id = ?", item.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload: %w", err)
	}
	return &result, nil
}

func (r *ServiceRepository) ListServices(ctx context.Context, crmID string, search string, status string) ([]models.CatalogItem, error) {
	tx := r.db.WithContext(ctx).Where("crm_id = ? AND item_type = 'service'", crmID)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	if search != "" {
		tx = tx.Where("name ILIKE ?", "%"+search+"%")
	}
	var items []models.CatalogItem
	if err := r.preload(tx).Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	return items, nil
}

func (r *ServiceRepository) GetService(ctx context.Context, crmID string, itemID string) (*models.CatalogItem, error) {
	var item models.CatalogItem
	err := r.preload(r.db.WithContext(ctx)).
		Where("id = ? AND crm_id = ? AND item_type = 'service'", itemID, crmID).
		First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("service not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}
	return &item, nil
}

func (r *ServiceRepository) UpdateService(ctx context.Context, crmID string, itemID string, d dto.UpdateServiceDTO) (*models.CatalogItem, error) {
	var item models.CatalogItem
	err := r.db.WithContext(ctx).Where("id = ? AND crm_id = ? AND item_type = 'service'", itemID, crmID).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("service not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find service: %w", err)
	}

	status := d.Status
	if status == "" {
		status = "active"
	}
	currency := d.PriceCurrency
	if currency == "" {
		currency = "USD"
	}

	r.db.WithContext(ctx).Model(&item).Updates(map[string]interface{}{
		"name":           d.Name,
		"description":    d.Description,
		"price":          d.Price,
		"price_currency": currency,
		"status":         status,
		"show_in_store":  d.ShowInStore,
	})

	var tags []models.CatalogTag
	for _, tid := range d.TagIDs {
		if pid, err := uuid.Parse(tid); err == nil {
			tags = append(tags, models.CatalogTag{ID: pid})
		}
	}
	r.db.WithContext(ctx).Model(&item).Association("Tags").Replace(tags)

	// Find and replace service detail
	var detail models.ServiceDetail
	r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&detail)
	r.db.WithContext(ctx).Where("service_detail_id = ?", detail.ID).Delete(&models.ServiceVariation{})
	r.db.WithContext(ctx).Where("service_detail_id = ?", detail.ID).Delete(&models.ServiceLocation{})
	r.db.WithContext(ctx).Where("item_id = ?", item.ID).Delete(&models.ServiceDetail{})

	serviceType := d.ServiceType
	if serviceType == "" {
		serviceType = "appointment"
	}
	newDetail := models.ServiceDetail{
		ItemID:           item.ID,
		ServiceType:      serviceType,
		Tagline:          d.Tagline,
		Duration:         d.Duration,
		BufferTime:       d.BufferTime,
		PriceType:        d.PriceType,
		PaymentMode:      d.PaymentMode,
		CustomPriceLabel: d.CustomPriceLabel,
		BookingMode:      d.BookingMode,
	}
	r.db.WithContext(ctx).Create(&newDetail)

	for _, vDTO := range d.Variations {
		r.db.WithContext(ctx).Create(&models.ServiceVariation{
			ServiceDetailID: newDetail.ID,
			Name:            vDTO.Name,
			Price:           vDTO.Price,
			Position:        vDTO.Position,
		})
	}
	for _, lDTO := range d.Locations {
		r.db.WithContext(ctx).Create(&models.ServiceLocation{
			ServiceDetailID: newDetail.ID,
			LocationType:    lDTO.LocationType,
			Address:         lDTO.Address,
			PhoneMethod:     lDTO.PhoneMethod,
			Phone:           lDTO.Phone,
		})
	}

	// Replace plans
	r.db.WithContext(ctx).Where("catalog_item_id = ?", item.ID).Delete(&models.Plan{})
	for i := range d.Plans {
		pDTO := d.Plans[i]
		currency := pDTO.PriceCurrency
		if currency == "" {
			currency = "USD"
		}
		r.db.WithContext(ctx).Create(&models.Plan{
			CatalogItemID: item.ID,
			CRMID:         item.CRMID,
			Name:          pDTO.Name,
			PlanType:      pDTO.PlanType,
			Price:         pDTO.Price,
			PriceCurrency: currency,
			BillingCycle:  pDTO.BillingCycle,
			SessionCount:  pDTO.SessionCount,
			ValidityDays:  pDTO.ValidityDays,
			Status:        "active",
		})
	}

	r.db.WithContext(ctx).Where("item_id = ?", item.ID).Delete(&models.CatalogAdditionalField{})
	for _, fDTO := range d.AdditionalFields {
		r.db.WithContext(ctx).Create(&models.CatalogAdditionalField{
			ItemID:    item.ID,
			Label:     fDTO.Label,
			FieldType: fDTO.FieldType,
			Value:     fDTO.Value,
			Position:  fDTO.Position,
		})
	}

	// Reset display flags then set specified display image
	r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).Where("item_id = ?", item.ID).Update("display_image", false)
	if d.DisplayImageID != "" {
		r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
			Where("id = ? AND item_id = ?", d.DisplayImageID, item.ID).
			Update("display_image", true)
	}

	// Append new media
	if len(d.MediaKeys) > 0 {
		var maxPos int
		r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
			Where("item_id = ?", item.ID).
			Select("COALESCE(MAX(position), -1)").
			Scan(&maxPos)
		r2cfg := config.GetConfig().R2
		for i, key := range d.MediaKeys {
			media := models.CatalogItemMedia{
				ItemID:       item.ID,
				URL:          fmt.Sprintf("%s/%s", r2cfg.PublicURL, key),
				Key:          key,
				Position:     maxPos + 1 + i,
				DisplayImage: key == d.DisplayImageKey,
			}
			r.db.WithContext(ctx).Create(&media)
		}
	}

	var result models.CatalogItem
	if err := r.preload(r.db.WithContext(ctx)).First(&result, "id = ?", item.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload: %w", err)
	}
	return &result, nil
}

func (r *ServiceRepository) DeleteService(ctx context.Context, crmID string, itemID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND crm_id = ? AND item_type = 'service'", itemID, crmID).Delete(&models.CatalogItem{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete service: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("service not found")
	}
	return nil
}
