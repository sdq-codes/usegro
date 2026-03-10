export interface FormVersionResponse {
  version: FormVersion
  fields: FormField[]
}

export interface FormVersion {
  PK: string
  SK: string
  formID: string
  title: string
  description: string
  formVersionStatus: string
  publishedAt: string
  createdAt: string
  updatedAt: string
}

export interface FormField {
  PK: string
  SK: string
  formVersionID: string
  fieldTypeID: number
  fieldTypeName: string
  label: string
  section: string
  slug: string
  placeholder: string
  configs: FieldConfig[] | null
  options: FieldOption[] | null
  validations: FieldValidation[]
  order: number
  required: boolean
  logic: unknown | null
  createdAt: string
  updatedAt: string
}

export interface FieldConfig {
  key: string
  description?: string
  valueType: string
  fieldId?: string
  fieldValue?: string
}

export interface FieldOption {
  label: string
  value: string
}

export interface FieldValidation {
  key: string
  rule: string
  message: string
}

export interface SubmissionPayload {
  type: string,
  answers: Record<string, string | string[]>,
  versionSnap: FormField[],
}
