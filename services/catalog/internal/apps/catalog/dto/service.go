package dto

type CreateServiceVariationDTO struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Position int     `json:"position"`
}

type CreateServiceLocationDTO struct {
	LocationType string `json:"location_type"`
	Address      string `json:"address"`
	PhoneMethod  string `json:"phone_method"`
	Phone        string `json:"phone"`
}

type CreateServiceDTO struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	PriceCurrency string  `json:"price_currency"`
	Status        string  `json:"status"`
	ShowInStore   bool    `json:"show_in_store"`

	// Service-specific
	ServiceType      string `json:"service_type"` // appointment | class | course
	Tagline          string `json:"tagline"`
	Duration         string `json:"duration"`
	BufferTime       string `json:"buffer_time"`
	PriceType        string `json:"price_type"` // fixed | variable | free | custom
	PaymentMode      string `json:"payment_mode"`
	CustomPriceLabel string `json:"custom_price_label"`
	BookingMode      string `json:"booking_mode"` // auto | manual

	TagIDs           []string                    `json:"tag_ids"`
	Plans            []CreatePlanDTO             `json:"plans"`
	Variations       []CreateServiceVariationDTO `json:"variations"`
	Locations        []CreateServiceLocationDTO  `json:"locations"`
	AdditionalFields []CreateAdditionalFieldDTO  `json:"additional_fields"`
	MediaKeys        []string                    `json:"media_keys"`
	DisplayImageKey  string                      `json:"display_image_key"`
	DisplayImageID   string                      `json:"display_image_id"`
}

type UpdateServiceDTO = CreateServiceDTO
