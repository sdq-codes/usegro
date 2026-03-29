package models

import (
	"time"

	"github.com/google/uuid"
)

// ServiceDetail holds service-specific fields (appointments, classes, courses)
type ServiceDetail struct {
	ID               uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ItemID           uuid.UUID `json:"item_id" gorm:"type:uuid;not null;uniqueIndex"`
	ServiceType      string    `json:"service_type" gorm:"type:varchar(50);default:'appointment'"` // appointment | class | course
	Tagline          string    `json:"tagline" gorm:"type:varchar(255)"`
	Duration         string    `json:"duration" gorm:"type:varchar(50);default:'1 Hour'"`
	BufferTime       string    `json:"buffer_time" gorm:"type:varchar(50);default:'No buffer'"`
	PriceType        string    `json:"price_type" gorm:"type:varchar(20);default:'fixed'"`         // fixed | variable | free | custom
	PaymentMode      string    `json:"payment_mode" gorm:"type:varchar(30);default:'per-session'"` // per-session | with-plan | per-session-or-plan
	CustomPriceLabel string    `json:"custom_price_label" gorm:"type:varchar(255)"`
	BookingMode      string    `json:"booking_mode" gorm:"type:varchar(20);default:'auto'"` // auto | manual

	Variations []ServiceVariation `json:"variations" gorm:"foreignKey:ServiceDetailID;constraint:OnDelete:CASCADE"`
	Locations  []ServiceLocation  `json:"locations" gorm:"foreignKey:ServiceDetailID;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ServiceVariation struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ServiceDetailID uuid.UUID `json:"service_detail_id" gorm:"type:uuid;not null;index"`
	Name            string    `json:"name" gorm:"type:varchar(255);not null"`
	Price           float64   `json:"price" gorm:"type:decimal(10,2);default:0"`
	Position        int       `json:"position" gorm:"default:0"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type ServiceLocation struct {
	ID              uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ServiceDetailID uuid.UUID `json:"service_detail_id" gorm:"type:uuid;not null;index"`
	LocationType    string    `json:"location_type" gorm:"type:varchar(30);not null"` // zoom | phone | in-person | google-meet | ms-teams
	Address         string    `json:"address" gorm:"type:varchar(500)"`
	PhoneMethod     string    `json:"phone_method" gorm:"type:varchar(20)"` // require | provide
	Phone           string    `json:"phone" gorm:"type:varchar(50)"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
