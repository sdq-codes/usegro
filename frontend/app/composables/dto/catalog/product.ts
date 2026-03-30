export interface ProductDetail {
  id: string
  item_id: string
  brand: string
  ribbon: string
  item_sub_type: string
  sku: string
  track_inventory: boolean
  stock_status: string
  quantity: number
  options: ProductOption[]
  variants: ProductVariant[]
}

export interface ProductOption {
  id: string
  item_id: string
  name: string
  position: number
  values: ProductOptionValue[]
}

export interface ProductOptionValue {
  id: string
  option_id: string
  value: string
  position: number
}

export interface ProductVariant {
  id: string
  item_id: string
  name: string
  sku: string
  price: number
  quantity: number
}

export interface ServiceDetail {
  id: string
  item_id: string
  service_type: string
  tagline: string
  duration: string
  buffer_time: string
  price_type: string
  payment_mode: string
  custom_price_label: string
  booking_mode: string
  variations: ServiceVariation[]
  locations: ServiceLocation[]
}

export interface ServiceVariation {
  id: string
  service_detail_id: string
  name: string
  price: number
  position: number
}

export interface ServiceLocation {
  id: string
  service_detail_id: string
  location_type: string
  address: string
  phone_method: string
  phone: string
}

export interface CatalogAdditionalField {
  id: string
  item_id: string
  label: string
  field_type: string
  value: string
  position: number
}

export interface CatalogCategory {
  id: string
  crm_id: string
  name: string
}

export interface CatalogItemMedia {
  id: string
  item_id: string
  url: string
  key: string
  mime_type: string
  size: number
  position: number
  display_image: boolean
  created_at: string
}

export interface CatalogItemPlan {
  id: string
  catalog_item_id: string
  crm_id: string
  name: string
  plan_type: 'subscription' | 'package'
  price: number
  price_currency: string
  billing_cycle?: string
  session_count?: number | null
  validity_days?: number | null
  status: string
  created_at: string
  updated_at: string
}

export interface CatalogItem {
  id: string
  crm_id: string
  item_type: string
  name: string
  description: string
  price: number
  price_currency: string
  discount_percent: number
  cost_per_item: number
  status: string
  show_in_store: boolean
  standard_category_id?: string | null
  standard_category?: { id: string; name: string; parent_id: string | null } | null
  product_detail?: ProductDetail
  service_detail?: ServiceDetail
  plans?: CatalogItemPlan[]
  additional_fields: CatalogAdditionalField[]
  categories: CatalogCategory[]
  tags: Array<{ id: string; name: string }>
  media: CatalogItemMedia[]
  created_at: string
  updated_at: string
}

// Keep legacy Product alias for backward compatibility with existing pages
export type Product = CatalogItem

// Product DTOs
export interface CreateProductOptionDTO { name: string; values: string[]; position: number }
export interface CreateProductVariantDTO { name: string; sku: string; price: number; quantity: number }
export interface CreateAdditionalFieldDTO { label: string; field_type: string; value: string; position: number }

export interface CreateProductDTO {
  name: string
  description?: string
  price?: number
  price_currency?: string
  discount_percent?: number
  cost_per_item?: number
  status?: string
  show_in_store?: boolean
  brand?: string
  ribbon?: string
  item_sub_type?: string
  sku?: string
  track_inventory?: boolean
  stock_status?: string
  quantity?: number
  tag_ids?: string[]
  standard_category_id?: string | null
  options?: CreateProductOptionDTO[]
  variants?: CreateProductVariantDTO[]
  additional_fields?: CreateAdditionalFieldDTO[]
  media_keys?: string[]
  display_image_key?: string
  display_image_id?: string
}

export type UpdateProductDTO = CreateProductDTO
