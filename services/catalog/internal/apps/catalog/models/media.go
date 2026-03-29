package models

import (
	"time"

	"github.com/google/uuid"
)

type CatalogItemMedia struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID       uuid.UUID `json:"item_id" gorm:"type:uuid;not null;index"`
	URL          string    `json:"url" gorm:"type:text;not null"`
	Key          string    `json:"key" gorm:"type:text;not null"`
	MimeType     string    `json:"mime_type" gorm:"type:varchar(100)"`
	Size         int64     `json:"size"`
	Position     int       `json:"position" gorm:"default:0"`
	DisplayImage bool      `json:"display_image" gorm:"default:false"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}
