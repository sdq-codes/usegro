package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"unique;not null;index:users_email_idx;"`
	Password     *string   `gorm:"default:null"`
	GoogleID     *string   `json:"google_id" gorm:"default:null;index"`
	AuthProvider string    `json:"auth_provider" gorm:"not null;default:'password'"`

	Verifications []Verification `json:"verifications" gorm:"foreignKey:user_id;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

//func (u *User) AfterUpdate(tx *gorm.DB) (err error) {
//	tx.Model(&User{}).Where("id = ?", u.ID).Update("updated_at", time.Time{})
//	return nil
//}
