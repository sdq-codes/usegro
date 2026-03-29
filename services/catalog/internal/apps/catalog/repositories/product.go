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

type PaginatedProducts struct {
	Data       []models.CatalogItem `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

type ProductRepositoryInterface interface {
	CreateProduct(ctx context.Context, crmID string, d dto.CreateProductDTO) (*models.CatalogItem, error)
	ListProducts(ctx context.Context, crmID string, search string, status string, page, limit int) (*PaginatedProducts, error)
	GetProduct(ctx context.Context, crmID string, productID string) (*models.CatalogItem, error)
	UpdateProduct(ctx context.Context, crmID string, productID string, d dto.UpdateProductDTO) (*models.CatalogItem, error)
	DeleteProduct(ctx context.Context, crmID string, productID string) error
}

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepositoryInterface {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) preload(tx *gorm.DB) *gorm.DB {
	return tx.
		Preload("ProductDetail.Options.Values").
		Preload("ProductDetail.Variants").
		Preload("ProductDetail.VariantGroups").
		Preload("AdditionalFields").
		Preload("Categories").
		Preload("Tags").
		Preload("StandardCategory").
		Preload("Media")
}

func (r *ProductRepository) CreateProduct(ctx context.Context, crmID string, d dto.CreateProductDTO) (*models.CatalogItem, error) {
	parsedCRMID, err := uuid.Parse(crmID)
	if err != nil {
		return nil, fmt.Errorf("invalid crm_id: %w", err)
	}

	status := d.Status
	if status == "" {
		status = "active"
	}
	priceCurrency := d.PriceCurrency
	if priceCurrency == "" {
		priceCurrency = "USD"
	}

	item := models.CatalogItem{
		CRMID:         parsedCRMID,
		ItemType:      "product",
		Name:          d.Name,
		Description:   d.Description,
		Price:         d.Price,
		PriceCurrency: priceCurrency,
		CostPerItem:   d.CostPerItem,
		Barcode:       d.Barcode,
		ChargeTax:     d.ChargeTax,
		Status:        status,
		ShowInStore:   d.ShowInStore,
	}
	if d.StandardCategoryID != "" {
		if parsed, err := uuid.Parse(d.StandardCategoryID); err == nil {
			item.StandardCategoryID = &parsed
		}
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

	stockStatus := d.StockStatus
	if stockStatus == "" {
		stockStatus = "in_stock"
	}
	itemSubType := d.ItemSubType
	if itemSubType == "" {
		itemSubType = "physical"
	}

	detail := models.ProductDetail{
		ItemID:                        item.ID,
		Brand:                         d.Brand,
		Ribbon:                        d.Ribbon,
		ItemSubType:                   itemSubType,
		SKU:                           d.SKU,
		TrackInventory:                d.TrackInventory,
		StockStatus:                   stockStatus,
		Quantity:                      d.Quantity,
		ContinueSellingWhenOutOfStock: d.ContinueSellingWhenOutOfStock,
	}
	if err := r.db.WithContext(ctx).Create(&detail).Error; err != nil {
		return nil, fmt.Errorf("failed to create product detail: %w", err)
	}

	for _, optDTO := range d.Options {
		option := models.ProductOption{
			ItemID:   detail.ID,
			Name:     optDTO.Name,
			Position: optDTO.Position,
		}
		if err := r.db.WithContext(ctx).Create(&option).Error; err != nil {
			return nil, fmt.Errorf("failed to create product option: %w", err)
		}
		for pos, val := range optDTO.Values {
			optVal := models.ProductOptionValue{
				OptionID: option.ID,
				Value:    val,
				Position: pos,
			}
			if err := r.db.WithContext(ctx).Create(&optVal).Error; err != nil {
				return nil, fmt.Errorf("failed to create product option value: %w", err)
			}
		}
	}

	for _, varDTO := range d.Variants {
		variant := models.ProductVariant{
			ItemID:   detail.ID,
			Name:     varDTO.Name,
			Price:    varDTO.Price,
			Quantity: varDTO.Quantity,
			ImageKey: varDTO.ImageKey,
			ImageURL: varDTO.ImageURL,
		}
		if err := r.db.WithContext(ctx).Create(&variant).Error; err != nil {
			return nil, fmt.Errorf("failed to create product variant: %w", err)
		}
	}

	for _, grpDTO := range d.VariantGroups {
		group := models.ProductVariantGroup{
			ItemID:     detail.ID,
			OptionName: grpDTO.OptionName,
			Value:      grpDTO.Value,
			ImageKey:   grpDTO.ImageKey,
			ImageURL:   grpDTO.ImageURL,
		}
		if err := r.db.WithContext(ctx).Create(&group).Error; err != nil {
			return nil, fmt.Errorf("failed to create variant group: %w", err)
		}
	}

	for _, fieldDTO := range d.AdditionalFields {
		field := models.CatalogAdditionalField{
			ItemID:    item.ID,
			Label:     fieldDTO.Label,
			FieldType: fieldDTO.FieldType,
			Value:     fieldDTO.Value,
			Position:  fieldDTO.Position,
		}
		if err := r.db.WithContext(ctx).Create(&field).Error; err != nil {
			return nil, fmt.Errorf("failed to create additional field: %w", err)
		}
	}

	// Create media records from presigned upload keys
	for i, key := range d.MediaKeys {
		r2cfg := config.GetConfig().R2
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
	// If no explicit display key was given, mark the first media as display
	if d.DisplayImageKey == "" && len(d.MediaKeys) > 0 {
		if err := r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
			Where("item_id = ? AND position = 0", item.ID).
			Update("display_image", true).Error; err != nil {
			return nil, fmt.Errorf("failed to set default display image: %w", err)
		}
	}

	var result models.CatalogItem
	if err := r.preload(r.db.WithContext(ctx)).First(&result, "id = ? AND item_type = 'product'", item.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload product: %w", err)
	}

	return &result, nil
}

func (r *ProductRepository) ListProducts(ctx context.Context, crmID string, search string, status string, page, limit int) (*PaginatedProducts, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	tx := r.db.WithContext(ctx).Where("crm_id = ? AND item_type = 'product'", crmID)
	if status != "" {
		tx = tx.Where("status = ?", status)
	}
	if search != "" {
		tx = tx.Where("name ILIKE ?", "%"+search+"%")
	}

	var total int64
	if err := tx.Model(&models.CatalogItem{}).Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	if totalPages == 0 {
		totalPages = 1
	}

	var items []models.CatalogItem
	if err := r.preload(tx).Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&items).Error; err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}

	return &PaginatedProducts{
		Data:       items,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

func (r *ProductRepository) GetProduct(ctx context.Context, crmID string, productID string) (*models.CatalogItem, error) {
	var item models.CatalogItem
	err := r.preload(r.db.WithContext(ctx)).
		Where("id = ? AND crm_id = ? AND item_type = 'product'", productID, crmID).
		First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}
	return &item, nil
}

func (r *ProductRepository) UpdateProduct(ctx context.Context, crmID string, productID string, d dto.UpdateProductDTO) (*models.CatalogItem, error) {
	var item models.CatalogItem
	err := r.db.WithContext(ctx).Where("id = ? AND crm_id = ? AND item_type = 'product'", productID, crmID).First(&item).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find product: %w", err)
	}

	status := d.Status
	if status == "" {
		status = "active"
	}
	priceCurrency := d.PriceCurrency
	if priceCurrency == "" {
		priceCurrency = "USD"
	}

	var stdCatID *uuid.UUID
	if d.StandardCategoryID != "" {
		if parsed, err := uuid.Parse(d.StandardCategoryID); err == nil {
			stdCatID = &parsed
		}
	}

	baseUpdates := map[string]interface{}{
		"name":                 d.Name,
		"description":          d.Description,
		"price":                d.Price,
		"price_currency":       priceCurrency,
		"cost_per_item":        d.CostPerItem,
		"barcode":              d.Barcode,
		"charge_tax":           d.ChargeTax,
		"status":               status,
		"show_in_store":        d.ShowInStore,
		"standard_category_id": stdCatID,
	}
	if err := r.db.WithContext(ctx).Model(&item).Updates(baseUpdates).Error; err != nil {
		return nil, fmt.Errorf("failed to update catalog item: %w", err)
	}

	var tags []models.CatalogTag
	for _, tid := range d.TagIDs {
		if pid, err := uuid.Parse(tid); err == nil {
			tags = append(tags, models.CatalogTag{ID: pid})
		}
	}
	r.db.WithContext(ctx).Model(&item).Association("Tags").Replace(tags)

	stockStatus := d.StockStatus
	if stockStatus == "" {
		stockStatus = "in_stock"
	}
	itemSubType := d.ItemSubType
	if itemSubType == "" {
		itemSubType = "physical"
	}

	// Upsert product detail
	detailUpdates := map[string]interface{}{
		"brand":                              d.Brand,
		"ribbon":                             d.Ribbon,
		"item_sub_type":                      itemSubType,
		"sku":                                d.SKU,
		"track_inventory":                    d.TrackInventory,
		"stock_status":                       stockStatus,
		"quantity":                           d.Quantity,
		"continue_selling_when_out_of_stock": d.ContinueSellingWhenOutOfStock,
	}
	if err := r.db.WithContext(ctx).Model(&models.ProductDetail{}).
		Where("item_id = ?", item.ID).
		Updates(detailUpdates).Error; err != nil {
		return nil, fmt.Errorf("failed to update product detail: %w", err)
	}

	// Fetch the product detail to get its ID for child associations
	var detail models.ProductDetail
	if err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).First(&detail).Error; err != nil {
		return nil, fmt.Errorf("failed to find product detail: %w", err)
	}

	// Replace options
	if err := r.db.WithContext(ctx).Where("item_id = ?", detail.ID).Delete(&models.ProductOption{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete old options: %w", err)
	}
	for _, optDTO := range d.Options {
		option := models.ProductOption{
			ItemID:   detail.ID,
			Name:     optDTO.Name,
			Position: optDTO.Position,
		}
		if err := r.db.WithContext(ctx).Create(&option).Error; err != nil {
			return nil, fmt.Errorf("failed to create product option: %w", err)
		}
		for pos, val := range optDTO.Values {
			optVal := models.ProductOptionValue{
				OptionID: option.ID,
				Value:    val,
				Position: pos,
			}
			if err := r.db.WithContext(ctx).Create(&optVal).Error; err != nil {
				return nil, fmt.Errorf("failed to create product option value: %w", err)
			}
		}
	}

	// Replace variants
	if err := r.db.WithContext(ctx).Where("item_id = ?", detail.ID).Delete(&models.ProductVariant{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete old variants: %w", err)
	}
	for _, varDTO := range d.Variants {
		variant := models.ProductVariant{
			ItemID:   detail.ID,
			Name:     varDTO.Name,
			Price:    varDTO.Price,
			Quantity: varDTO.Quantity,
			ImageKey: varDTO.ImageKey,
			ImageURL: varDTO.ImageURL,
		}
		if err := r.db.WithContext(ctx).Create(&variant).Error; err != nil {
			return nil, fmt.Errorf("failed to create product variant: %w", err)
		}
	}

	// Replace variant groups
	if err := r.db.WithContext(ctx).Where("item_id = ?", detail.ID).Delete(&models.ProductVariantGroup{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete old variant groups: %w", err)
	}
	for _, grpDTO := range d.VariantGroups {
		group := models.ProductVariantGroup{
			ItemID:     detail.ID,
			OptionName: grpDTO.OptionName,
			Value:      grpDTO.Value,
			ImageKey:   grpDTO.ImageKey,
			ImageURL:   grpDTO.ImageURL,
		}
		if err := r.db.WithContext(ctx).Create(&group).Error; err != nil {
			return nil, fmt.Errorf("failed to create variant group: %w", err)
		}
	}

	// Replace additional fields
	if err := r.db.WithContext(ctx).Where("item_id = ?", item.ID).Delete(&models.CatalogAdditionalField{}).Error; err != nil {
		return nil, fmt.Errorf("failed to delete old additional fields: %w", err)
	}
	for _, fieldDTO := range d.AdditionalFields {
		field := models.CatalogAdditionalField{
			ItemID:    item.ID,
			Label:     fieldDTO.Label,
			FieldType: fieldDTO.FieldType,
			Value:     fieldDTO.Value,
			Position:  fieldDTO.Position,
		}
		if err := r.db.WithContext(ctx).Create(&field).Error; err != nil {
			return nil, fmt.Errorf("failed to create additional field: %w", err)
		}
	}

	// Reset all display_image flags for this item
	if err := r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
		Where("item_id = ?", item.ID).
		Update("display_image", false).Error; err != nil {
		return nil, fmt.Errorf("failed to reset display image flags: %w", err)
	}
	// Set display_image on existing media if specified
	if d.DisplayImageID != "" {
		if err := r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
			Where("id = ? AND item_id = ?", d.DisplayImageID, item.ID).
			Update("display_image", true).Error; err != nil {
			return nil, fmt.Errorf("failed to set display image: %w", err)
		}
	}

	// Append new media (do not delete existing)
	if len(d.MediaKeys) > 0 {
		var maxPos int
		r.db.WithContext(ctx).Model(&models.CatalogItemMedia{}).
			Where("item_id = ?", item.ID).
			Select("COALESCE(MAX(position), -1)").
			Scan(&maxPos)
		r2cfg := config.GetConfig().R2
		for i, key := range d.MediaKeys {
			isDisplay := key == d.DisplayImageKey
			media := models.CatalogItemMedia{
				ItemID:       item.ID,
				URL:          fmt.Sprintf("%s/%s", r2cfg.PublicURL, key),
				Key:          key,
				Position:     maxPos + 1 + i,
				DisplayImage: isDisplay,
			}
			if err := r.db.WithContext(ctx).Create(&media).Error; err != nil {
				return nil, fmt.Errorf("failed to append media record: %w", err)
			}
		}
	}

	var result models.CatalogItem
	if err := r.preload(r.db.WithContext(ctx)).First(&result, "id = ? AND item_type = 'product'", item.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload product: %w", err)
	}

	return &result, nil
}

func (r *ProductRepository) DeleteProduct(ctx context.Context, crmID string, productID string) error {
	result := r.db.WithContext(ctx).Where("id = ? AND crm_id = ? AND item_type = 'product'", productID, crmID).Delete(&models.CatalogItem{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete product: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
