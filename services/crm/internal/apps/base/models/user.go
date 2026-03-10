package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID       uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email    string    `json:"email" gorm:"unique;not null;index:users_email_idx;"`
	Password string    `gorm:"not null;"`

	Verifications []Verification `json:"verifications" gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}
