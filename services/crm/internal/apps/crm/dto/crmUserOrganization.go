package dto

import (
	"github.com/google/uuid"
	"github.com/usegro/services/crm/internal/apps/crm/models"
	"time"
)

type CrmUserOrganizationDTI struct {
	FullName     string                  `json:"full_name" validate:"required,max=255,min=2"`
	BusinessName string                  `json:"business_name" validate:"required,min=2"`
	BusinessInfo models.BusinessInfoType `json:"business_info" validate:"required,oneof=starter_package existing_package"`
}

type BusinessNameExistDTI struct {
	BusinessName string `json:"business_name" validate:"required,min=2"`
}

type CrmUserOrganizationUpdateDTI struct {
	FullName     string                  `json:"full_name" validate:"max=255,min=2"`
	BusinessName string                  `json:"business_name" validate:"min=2"`
	BusinessInfo models.BusinessInfoType `json:"business_info" validate:"oneof=starter_package existing_package"`
}

type CrmUserOrganizationSalesChannelTypeDTI struct {
	SalesChannelType []string `json:"sales_channel_type" validate:"required,dive,oneof=online_store social_media_store physical_store existing_website_store unknown_store"`
}

type CrmUserOrganizationStockProductTypeDTI struct {
	ProductType []models.ProductType `json:"product_type" validate:"required,dive,oneof=physical_products digital_products services_products unknown_products"`
}

type CrmUserOrganizationDTO struct {
	ID               uuid.UUID                 `json:"id"`
	FullName         string                    `json:"full_name"`
	BusinessName     string                    `json:"business_name"`
	BusinessInfo     models.BusinessInfoType   `json:"business_info"`
	SalesChannel     []models.SalesChannelType `json:"sales_channel"`
	StockProductType []models.ProductType      `json:"stock_product_type"`
	CreatedAt        time.Time                 `json:"created_at" gorm:"default:now()"`
	UpdatedAt        time.Time                 `json:"updated_at" gorm:"default:now()"`
}
