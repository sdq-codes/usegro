package models

import (
	"time"

	"github.com/google/uuid"
)

// CatalogItem is the polymorphic base for both Products and Services
type CatalogItem struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CRMID         uuid.UUID `json:"crm_id" gorm:"type:uuid;not null;index:catalog_item_crm_id"`
	ItemType      string    `json:"item_type" gorm:"type:varchar(20);not null;default:'product'"` // product | service
	Name          string    `json:"name" gorm:"type:varchar(255);not null"`
	Description   string    `json:"description" gorm:"type:text"`
	Price         float64   `json:"price" gorm:"type:decimal(10,2);default:0"`
	PriceCurrency string    `json:"price_currency" gorm:"type:varchar(10);default:'USD'"`
	CostPerItem   float64   `json:"cost_per_item" gorm:"type:decimal(10,2);default:0"`
	Barcode       string    `json:"barcode" gorm:"type:varchar(255);default:''"`
	ChargeTax     bool      `json:"charge_tax" gorm:"default:true"`
	Status        string    `json:"status" gorm:"type:varchar(50);default:'active'"` // active | draft | archived
	ShowInStore   bool      `json:"show_in_store" gorm:"default:true"`

	StandardCategoryID *uuid.UUID        `json:"standard_category_id,omitempty" gorm:"type:uuid;index"`
	StandardCategory   *StandardCategory `json:"standard_category,omitempty" gorm:"foreignKey:StandardCategoryID"`

	// Polymorphic children
	ProductDetail    *ProductDetail           `json:"product_detail,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	ServiceDetail    *ServiceDetail           `json:"service_detail,omitempty" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Plans            []Plan                   `json:"plans,omitempty" gorm:"foreignKey:CatalogItemID;constraint:OnDelete:CASCADE"`
	AdditionalFields []CatalogAdditionalField `json:"additional_fields" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Categories       []Category               `json:"categories" gorm:"many2many:catalog_item_categories"`
	Tags             []CatalogTag             `json:"tags" gorm:"many2many:catalog_item_tags"`
	Media            []CatalogItemMedia       `json:"media" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CatalogAdditionalField struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID    uuid.UUID `json:"item_id" gorm:"type:uuid;not null;index"`
	Label     string    `json:"label" gorm:"type:varchar(255);not null"`
	FieldType string    `json:"field_type" gorm:"type:varchar(50);not null"`
	Value     string    `json:"value" gorm:"type:text"`
	Position  int       `json:"position" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
