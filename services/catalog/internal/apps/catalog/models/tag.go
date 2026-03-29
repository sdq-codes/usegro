package models

import (
	"time"

	"github.com/google/uuid"
)

// CatalogTag is a shared tag that can be attached to any CatalogItem (product or service)
type CatalogTag struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CRMID     uuid.UUID `json:"crm_id" gorm:"type:uuid;not null;index"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
