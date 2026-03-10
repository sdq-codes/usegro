package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Verification struct {
	ID         uuid.UUID  `gorm:"type:char(36);primaryKey;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID  `gorm:"type:char(36);not null;" json:"user_id"`
	Type       string     `gorm:"column:type;type:verification_type;not null;" json:"type"`
	Status     string     `gorm:"column:status;type:verification_status;not null;default:'PENDING'" json:"status"`
	VerifiedAt *time.Time `json:"verified_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"default:now()"`
}

type VerificationToken struct {
	ID        uuid.UUID  `gorm:"type:char(36);primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID  `gorm:"type:char(36);not null;" json:"user_id"`
	TokenHash string     `gorm:"type:text;not null;" json:"token_hash"`
	Type      string     `gorm:"type:verification_type;not null;" json:"type"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"default:now()"`
}

func CreateVerificationEnums(db *gorm.DB) error {
	err := db.Exec(`
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'verification_type') THEN
                CREATE TYPE verification_type AS ENUM ('EMAIL', 'PHONE', '2FA', 'PASSWORD_RESET');
            END IF;
        END$$;
    `).Error
	if err != nil {
		return err
	}

	err = db.Exec(`
        DO $$
        BEGIN
            IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'verification_status') THEN
                CREATE TYPE verification_status AS ENUM ('PENDING', 'VERIFIED', 'EXPIRED', 'FAILED');
            END IF;
        END$$;
    `).Error
	if err != nil {
		return err
	}

	return nil
}
