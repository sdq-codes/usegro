package dto

type CreateProductOptionDTO struct {
	Name     string   `json:"name"`
	Values   []string `json:"values"`
	Position int      `json:"position"`
}

type CreateProductVariantDTO struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
	ImageKey string  `json:"image_key"`
	ImageURL string  `json:"image_url"`
}

type CreateVariantGroupDTO struct {
	OptionName string `json:"option_name"`
	Value      string `json:"value"`
	ImageKey   string `json:"image_key"`
	ImageURL   string `json:"image_url"`
}

type CreateAdditionalFieldDTO struct {
	Label     string `json:"label"`
	FieldType string `json:"field_type"`
	Value     string `json:"value"`
	Position  int    `json:"position"`
}

type CreateProductDTO struct {
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	PriceCurrency string  `json:"price_currency"`
	CostPerItem   float64 `json:"cost_per_item"`
	Barcode       string  `json:"barcode"`
	ChargeTax     bool    `json:"charge_tax"`
	Status        string  `json:"status"`
	ShowInStore   bool    `json:"show_in_store"`

	// Product-specific
	Brand                         string `json:"brand"`
	Ribbon                        string `json:"ribbon"`
	ItemSubType                   string `json:"item_sub_type"` // physical | digital
	SKU                           string `json:"sku"`
	TrackInventory                bool   `json:"track_inventory"`
	StockStatus                   string `json:"stock_status"`
	Quantity                      int    `json:"quantity"`
	ContinueSellingWhenOutOfStock bool   `json:"continue_selling_when_out_of_stock"`

	TagIDs             []string                   `json:"tag_ids"`
	StandardCategoryID string                     `json:"standard_category_id"`
	Options            []CreateProductOptionDTO   `json:"options"`
	Variants           []CreateProductVariantDTO  `json:"variants"`
	AdditionalFields   []CreateAdditionalFieldDTO `json:"additional_fields"`
	MediaKeys          []string                   `json:"media_keys"`
	DisplayImageKey    string                     `json:"display_image_key"` // key of a new upload to mark as display
	DisplayImageID     string                     `json:"display_image_id"`  // ID of an existing media record to mark as display
	VariantGroups      []CreateVariantGroupDTO    `json:"variant_groups"`
}

type UpdateProductDTO = CreateProductDTO
