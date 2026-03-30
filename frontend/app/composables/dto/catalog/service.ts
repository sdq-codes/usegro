export interface CreateServiceVariationDTO { name: string; price: number; position: number }
export interface CreateServiceLocationDTO { location_type: string; address: string; phone_method: string; phone: string }
export interface CreateAdditionalFieldDTO { label: string; field_type: string; value: string; position: number }

export interface CreatePlanDTO {
  name: string
  plan_type: 'subscription' | 'package'
  price?: number
  price_currency?: string
  billing_cycle?: 'monthly' | 'yearly'
  session_count?: number | null  // null = unlimited
  validity_days?: number | null  // package only
}

export interface Plan {
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

export interface CreateServiceDTO {
  name: string
  description?: string
  price?: number
  price_currency?: string
  discount_percent?: number
  status?: string
  show_in_store?: boolean
  service_type?: string
  tagline?: string
  duration?: string
  buffer_time?: string
  price_type?: string
  payment_mode?: string
  custom_price_label?: string
  booking_mode?: string
  tag_ids?: string[]
  plans?: CreatePlanDTO[]
  variations?: CreateServiceVariationDTO[]
  locations?: CreateServiceLocationDTO[]
  additional_fields?: CreateAdditionalFieldDTO[]
}

export type UpdateServiceDTO = CreateServiceDTO
