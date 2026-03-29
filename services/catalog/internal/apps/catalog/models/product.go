package models

import (
	"time"

	"github.com/google/uuid"
)

// ProductDetail holds product-specific fields (physical/digital items)
type ProductDetail struct {
	ID                            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID                        uuid.UUID `json:"item_id" gorm:"type:uuid;not null;uniqueIndex"`
	Brand                         string    `json:"brand" gorm:"type:varchar(255)"`
	Ribbon                        string    `json:"ribbon" gorm:"type:varchar(50);default:''"`
	ItemSubType                   string    `json:"item_sub_type" gorm:"type:varchar(20);default:'physical'"` // physical | digital
	SKU                           string    `json:"sku" gorm:"type:varchar(255);default:''"`
	TrackInventory                bool      `json:"track_inventory" gorm:"default:false"`
	StockStatus                   string    `json:"stock_status" gorm:"type:varchar(50);default:'in_stock'"`
	Quantity                      int       `json:"quantity" gorm:"default:0"`
	ContinueSellingWhenOutOfStock bool      `json:"continue_selling_when_out_of_stock" gorm:"default:false"`

	Options       []ProductOption       `json:"options" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	Variants      []ProductVariant      `json:"variants" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	VariantGroups []ProductVariantGroup `json:"variant_groups" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductOption struct {
	ID        uuid.UUID            `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID    uuid.UUID            `json:"item_id" gorm:"type:uuid;not null;index"`
	Name      string               `json:"name" gorm:"type:varchar(255);not null"`
	Position  int                  `json:"position" gorm:"default:0"`
	Values    []ProductOptionValue `json:"values" gorm:"foreignKey:OptionID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
}

type ProductOptionValue struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	OptionID  uuid.UUID `json:"option_id" gorm:"type:uuid;not null;index"`
	Value     string    `json:"value" gorm:"type:varchar(255);not null"`
	Position  int       `json:"position" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariant struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID    uuid.UUID `json:"item_id" gorm:"type:uuid;not null;index"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Price     float64   `json:"price" gorm:"type:decimal(10,2);default:0"`
	Quantity  int       `json:"quantity" gorm:"default:0"`
	ImageKey  string    `json:"image_key" gorm:"type:text;default:''"`
	ImageURL  string    `json:"image_url" gorm:"type:text;default:''"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProductVariantGroup struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID     uuid.UUID `json:"item_id" gorm:"type:uuid;not null;index"`
	OptionName string    `json:"option_name" gorm:"type:varchar(255);not null"`
	Value      string    `json:"value" gorm:"type:varchar(255);not null"`
	ImageKey   string    `json:"image_key" gorm:"type:text;default:''"`
	ImageURL   string    `json:"image_url" gorm:"type:text;default:''"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
