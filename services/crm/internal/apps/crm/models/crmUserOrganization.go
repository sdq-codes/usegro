package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type BusinessInfoType string
type SalesChannel string
type ProductType string

//type SocialNetwork string

const (
	StarterPackage  BusinessInfoType = "starter_package"
	ExistingPackage BusinessInfoType = "existing_package"

	OnlineStore          SalesChannel = "online_store"
	SocialMediaStore     SalesChannel = "social_media_store"
	PhysicalStore        SalesChannel = "physical_store"
	ExistingWebsiteStore SalesChannel = "existing_website_store"
	UnknownStore         SalesChannel = "unknown_store"

	PhysicalProducts ProductType = "physical_products"
	DigitalProducts  ProductType = "digital_products"
	ServicesProducts ProductType = "services_products"
	UnknownProducts  ProductType = "unknown_products"

	//Instagram SocialNetwork = "instagram"
	//Whatsapp  SocialNetwork = "whatsapp"
	//Tiktok    SocialNetwork = "tiktok"
)

type CrmUserOrganization struct {
	ID     uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:char(36);not null;index: crm_user_organization_user_id" json:"user_id"`

	FullName         string             `json:"full_name" gorm:"type:varchar(255);not null"`
	BusinessName     string             `json:"business_name" gorm:"type:text;not null;uniqueIndex:crm_business_name_unique"`
	BusinessInfo     BusinessInfoType   `json:"business_info" gorm:"type:text;not null"`
	SalesChannel     []SalesChannelType `json:"sales_channel" gorm:"-"`
	StockProductType []ProductType      `json:"stock_product_type" gorm:"-"`

	Active bool `json:"active" gorm:"default:true"`

	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

func (c *CrmUserOrganization) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&CrmUserOrganization{}).Where("id = ?", c.ID).Update("updated_at", time.Time{})
	return nil
}

type SalesChannelType struct {
	ID                    uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CrmUserOrganizationID uuid.UUID `json:"crm_user_organization_id" gorm:"type:varchar(36);not null;index:sales_channel_type_crm_user_organization_id;"`

	SalesChannelType SalesChannel `json:"sales_channel_type" gorm:"type:varchar(255);not null"`

	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

func (s *SalesChannelType) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&CrmUserOrganization{}).Where("id = ?", s.ID).Update("updated_at", time.Time{})
	return nil
}

type StockProductType struct {
	ID                    uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CrmUserOrganizationID uuid.UUID `json:"crm_user_organization_id" gorm:"type:varchar(36);not null;"`

	ProductType string `json:"product_type" gorm:"type:varchar(255);not null"`

	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
}

func (s *StockProductType) AfterUpdate(tx *gorm.DB) (err error) {
	tx.Model(&CrmUserOrganization{}).Where("id = ?", s.ID).Update("updated_at", time.Time{})
	return nil
}

//type ConnectedSocialNetworkType struct {
//	ID                    uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
//	CrmUserOrganizationID uuid.UUID `json:"crm_user_organization_id" gorm:"type:char(36);not null;index:connected_social_network_type1_crm_user_organization_id;"`
//
//	SocialNetwork ProductType `json:"social_network" gorm:"type:varchar(255);not null"`
//
//	CreatedAt time.Time `json:"created_at" gorm:"default:now()"`
//	UpdatedAt time.Time `json:"updated_at" gorm:"default:now()"`
//}
