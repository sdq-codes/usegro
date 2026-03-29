package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID        uuid.UUID     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CRMID     uuid.UUID     `json:"crm_id" gorm:"type:uuid;not null;index:category_crm_id"`
	Name      string        `json:"name" gorm:"type:varchar(255);not null"`
	Items     []CatalogItem `json:"-" gorm:"many2many:catalog_item_categories"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}
