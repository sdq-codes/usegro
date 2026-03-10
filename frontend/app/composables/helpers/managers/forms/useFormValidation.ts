import { ref, watch, type Ref } from 'vue'
import type { FormField } from '@/composables/helpers/types/form'

type Config = {
  key: string
  valueType: string
  fieldId?: string
  fieldSlug?: string
  fieldValue?: string
  description?: string
}

export function useFormValidation(
  answerMap: Ref<Record<string, unknown>>,
  formFields: Ref<FormField[]>
) {
  const errors = ref<Record<string, string>>({})

  // Helper: Extract showIf config
  function extractShowIf(configs?: Config[]): Config | null {
    if (!configs?.length) return null
    return configs.find(cfg => cfg.key === 'showIf') ?? null
  }

  // Helper: Check if field should be visible
  function shouldShowField(field: FormField): boolean {
    const showIfConfig = extractShowIf(field.configs)
    if (!showIfConfig) return true

    if (showIfConfig.valueType === 'string' && showIfConfig.fieldSlug) {
      return answerMap.value[showIfConfig.fieldSlug] === showIfConfig.fieldValue
    }
    return false
  }

  // Helper: Check if field is required
  function isFieldRequired(field: FormField): boolean {
    if (field.required) return true

    const requiredConfig = field.configs?.find(cfg => cfg.key === 'required')
    if (!requiredConfig) return false

    if (requiredConfig.valueType === 'boolean' && requiredConfig.fieldSlug) {
      const dependentFieldValue = answerMap.value[requiredConfig.fieldSlug]
      return dependentFieldValue === requiredConfig.fieldValue || !!dependentFieldValue
    }
    return false
  }

  // Helper: Check if value is empty
  function isEmpty(value): boolean {
    if (value === null || value === undefined || value === '') return true

    // For arrays (Checkbox, Tags)
    if (Array.isArray(value)) return value.length === 0

    // For objects (Phone, Country, State)
    if (typeof value === 'object') {
      return Object.keys(value).length === 0 ||
        !Object.values(value).some(v => v !== null && v !== undefined && v !== '')
    }

    return false
  }

  // Helper: Validate email format
  function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    return emailRegex.test(email)
  }

  // Helper: Validate against regex pattern
  function matchesPattern(value: string, pattern: string): boolean {
    try {
      const regex = new RegExp(pattern)
      return regex.test(value)
    } catch (e) {
      console.error('Invalid regex pattern:', pattern, e)
      return true // Don't fail validation on invalid regex
    }
  }

  // Main validation function
  function validateField(field: FormField): string | null {
    const value = answerMap.value[field.slug]

    // Check if field is required
    if (isFieldRequired(field) && isEmpty(value)) {
      return `${field.label || 'This field'} is required`
    }

    // Skip further validation if field is empty and not required
    if (isEmpty(value)) return null

    // Email validation
    if (field.fieldTypeName === 'Email' && typeof value === 'string') {
      if (!isValidEmail(value)) {
        return 'Please enter a valid email address'
      }
    }

    // String validations (only for string types)
    if (typeof value === 'string' && field.validations?.length) {
      for (const validation of field.validations) {
        switch (validation.key) {
          case 'minLength': {
            const minLength = parseInt(validation.rule)
            if (!isNaN(minLength) && value.length < minLength) {
              return validation.message
            }
            break
          }

          case 'maxLength': {
            const maxLength = parseInt(validation.rule)
            if (!isNaN(maxLength) && value.length > maxLength) {
              return validation.message
            }
            break
          }

          case 'regex':
          case 'pattern': {
            if (!matchesPattern(value, validation.rule)) {
              return validation.message
            }
            break
          }
        }
      }
    }

    // Array validations (for Checkbox, Tags)
    if (Array.isArray(value) && field.validations?.length) {
      for (const validation of field.validations) {
        switch (validation.key) {
          case 'min': {
            const min = parseInt(validation.rule)
            if (!isNaN(min) && value.length < min) {
              return validation.message
            }
            break
          }

          case 'max': {
            const max = parseInt(validation.rule)
            if (!isNaN(max) && value.length > max) {
              return validation.message
            }
            break
          }
        }
      }
    }

    return null
  }

  // Validate all visible fields and update errors
  function validateAllFields() {
    const newErrors: Record<string, string> = {}

    formFields.value.forEach(field => {
      // Only validate visible fields
      if (!shouldShowField(field)) return

      const error = validateField(field)
      if (error) {
        newErrors[field.slug] = error
      }
    })

    errors.value = newErrors
  }

  // Watch answerMap for changes and validate
  watch(
    () => answerMap.value,
    () => {
      validateAllFields()
    },
    { deep: true }
  )

  // Initial validation
  validateAllFields()

  return {
    errors
  }
}
