export type FormField = {
  SK: string
  slug: string
  section: string
  order: number
  fieldTypeName: string
  label: string
  description?: string
  required?: boolean
  options?: Array<{ value: string; label: string }>
  configs?: Array<{
    key: string
    valueType: string
    fieldId?: string
    fieldSlug?: string
    fieldValue?: string
  }>
  validations?: Array<{
    key: string
    rule: string
    message: string
  }>
  fieldTypeID?: string
}
