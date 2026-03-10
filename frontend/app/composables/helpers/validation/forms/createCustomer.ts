import { z } from 'zod'
import type {FormField} from "@/composables/helpers/types/form";


type FormData = Record<string, string | string[]>

/**
 * Creates a dynamic Zod schema based on form field definitions
 */
export function createFormValidationSchema(fields: FormField[]) {
  // Create a base schema object
  const schemaObject: Record<string, z.ZodTypeAny> = {}

  fields.forEach((field) => {
    let fieldSchema: z.ZodTypeAny

    // Determine base schema type based on field type
    switch (field.fieldTypeName) {
      case 'Checkbox':
      case 'Tags':
        // Array of strings
        fieldSchema = z.array(z.string())
        break

      case 'Email':
        fieldSchema = z.string().email({ message: 'Must be a valid email address' })
        break

      case 'Phone Number':
        fieldSchema = z.string()
        break

      default:
        // Default to string for text inputs, radio buttons, etc.
        fieldSchema = z.string()
    }

    // Apply custom validations from the validations array
    if (field.validations && Array.isArray(field.validations)) {
      field.validations.forEach((validation) => {
        if (fieldSchema instanceof z.ZodString) {
          switch (validation.key) {
            case 'minLength':
              { const minLength = parseInt(validation.rule)
              fieldSchema = (fieldSchema as z.ZodString).min(minLength, {
                message: validation.message
              })
              break }

            case 'maxLength':
              const maxLength = parseInt(validation.rule)
              fieldSchema = (fieldSchema as z.ZodString).max(maxLength, {
                message: validation.message
              })
              break

            case 'regex':
              fieldSchema = (fieldSchema as z.ZodString).regex(
                new RegExp(validation.rule),
                { message: validation.message }
              )
              break
          }
        }
      })
    }
    // Apply required validation
    if (field.required) {
      if (fieldSchema instanceof z.ZodArray) {
        fieldSchema = fieldSchema.min(1, { message: `${field.label || field.slug} is required` })
      } else if (fieldSchema instanceof z.ZodString) {
        fieldSchema = fieldSchema.min(1, { message: `${field.label || field.slug} is required` })
      }
    } else {
      // Make optional if not required
      if (fieldSchema instanceof z.ZodString) {
        fieldSchema = fieldSchema.optional().or(z.literal(''))
      } else{
        fieldSchema = fieldSchema.optional()
      }
    }

    schemaObject[field.slug] = fieldSchema
  })

  return z.object(schemaObject)
}

/**
 * Validates form data with conditional field visibility
 */
export function validateFormWithConditions(
  fields: FormField[],
  formData: FormData
): { success: boolean; errors: Record<string, string> } {
  const errors: Record<string, string> = {}

  // Filter visible fields based on showIf conditions
  const visibleFields = fields.filter((field) => {
    if (!field.configs) return true

    const showIfConfig = field.configs.find((config) => config.key === 'showIf')
    if (!showIfConfig) return true

    // Check if the condition is met
    const dependentFieldValue = formData[showIfConfig.fieldSlug || '']
    return dependentFieldValue === showIfConfig.fieldValue
  })

  // Create schema only for visible fields
  const schema = createFormValidationSchema(visibleFields)

  // Validate using Zod
  const result = schema.safeParse(formData)

  if (!result.success) {
    result.error.issues.forEach((issue) => {
      const fieldPath = issue.path[0] as string
      errors[fieldPath] = issue.message
    })
    return { success: false, errors }
  }

  return { success: true, errors: {} }
}

/**
 * Checks if a field should be visible based on its showIf config
 */
export function isFieldVisible(field: FormField, formData: FormData): boolean {
  if (!field.configs) return true

  const showIfConfig = field.configs.find((config) => config.key === 'showIf')
  if (!showIfConfig) return true

  const dependentFieldValue = formData[showIfConfig.fieldSlug || '']
  return dependentFieldValue === showIfConfig.fieldValue
}

/**
 * Gets all visible fields based on current form data
 */
export function getVisibleFields(fields: FormField[], formData: FormData): FormField[] {
  return fields.filter((field) => isFieldVisible(field, formData))
}

/**
 * Validates a single field
 */
export function validateField(
  field: FormField,
  value: string | string[],
  formData: FormData
): string | null {
  // Check if field should be visible
  if (!isFieldVisible(field, formData)) {
    return null // Don't validate hidden fields
  }

  // Create schema for this field
  const schema = createFormValidationSchema([field])
  const result = schema.safeParse({ [field.slug]: value })

  if (!result.success) {
    const error = result.error.errors[0]
    return error?.message || 'Validation error'
  }

  return null
}

// Example usage
export function exampleUsage() {
  const fields: FormField[] = [
    {
      slug: 'customer_type',
      fieldTypeName: 'Radio Button',
      required: true,
      validations: [],
      configs: null,
      options: [
        { label: 'Individual', value: 'individual' },
        { label: 'Business', value: 'business' }
      ]
    },
    {
      slug: 'first_name',
      fieldTypeName: 'Short Text',
      required: false,
      validations: [
        {
          key: 'minLength',
          rule: '1',
          message: 'First name must be greater than 1'
        }
      ],
      configs: [
        {
          key: 'showIf',
          fieldSlug: 'customer_type',
          fieldValue: 'individual',
          valueType: 'string'
        }
      ]
    },
    {
      slug: 'email',
      fieldTypeName: 'Email',
      required: false,
      validations: [
        {
          key: 'regex',
          rule: '^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$',
          message: 'Must be a valid email Address'
        }
      ],
      configs: [
        {
          key: 'multiples',
          fieldValue: 'true',
          valueType: 'boolean'
        }
      ]
    }
  ]

  const formData: FormData = {
    customer_type: 'individual',
    first_name: 'John',
    email: ['john@example.com']
  }

  const validation = validateFormWithConditions(fields, formData)

  return validation
}
