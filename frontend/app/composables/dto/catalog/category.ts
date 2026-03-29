export interface Category {
  id: string
  crm_id: string
  name: string
  created_at: string
  updated_at: string
}

export interface AttributeValue {
  id: string
  name: string
  handle: string
}

export interface StandardAttribute {
  id: string
  name: string
  handle: string
  values: AttributeValue[]
}

export interface StandardCategory {
  id: string
  parent_id: string | null
  name: string
  is_leaf: boolean
  attributes?: StandardAttribute[]
  children?: StandardCategory[]
}

export interface StandardCategorySearchResult {
  id: string
  name: string
  path: string
}

export interface CreateCategoryDTO {
  name: string
}

export interface UpdateCategoryDTO {
  name: string
}
