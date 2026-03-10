import { ref } from 'vue'
import type { FormField } from "@/helpers/types/form"
import {
  CREATE_EMAIL_FIELDS,
  CREATE_PHONE_FIELDS,
  CREATE_ADDRESS_FIELDS,
  CREATE_COMPANY_FIELDS,
  CREATE_POSITION_FIELDS,
  CREATE_BIRTHDATE_FIELDS,
  CREATE_NOTES_FIELDS,
  CREATE_CUSTOM_FIELDS
} from "@/constants/forms/fields/fields"

export function useFormFieldManager() {
  const fieldCount = ref({
    email: 1,
    phone: 0,
    address: 0,
    company: 0,
    position: 0,
    birthdate: 0,
    notes: 0,
    custom: 0,
  })

  const getDefaultValue = (field: FormField): string | string[] => {
    if (field.fieldTypeName === "Checkbox") {
      return []
    } else if (field.fieldTypeName === "Radio Button") {
      return field.options?.[0]?.value ?? ''
    } else if (field.fieldTypeName === "Tags") {
      return []
    }
    return ''
  }

  const createFieldFromType = (type: string) => {
    const fieldCreators = {
      'email': CREATE_EMAIL_FIELDS,
      'phone': CREATE_PHONE_FIELDS,
      'address': CREATE_ADDRESS_FIELDS,
      'company': CREATE_COMPANY_FIELDS,
      'position': CREATE_POSITION_FIELDS,
      'birthdate': CREATE_BIRTHDATE_FIELDS,
      'notes': CREATE_NOTES_FIELDS,
      'custom': CREATE_CUSTOM_FIELDS,
    }

    const creator = fieldCreators[type as keyof typeof fieldCreators]
    if (!creator) {
      console.error(`Unknown field type: ${type}`)
      return null
    }

    fieldCount.value[type as keyof typeof fieldCount.value] += 1
    return creator(fieldCount.value[type as keyof typeof fieldCount.value])
  }

  const addField = (
    type: string,
    formFields: FormField[],
    answerMap: Record<string, string | string[]>
  ): FormField | null => {
    const newField = createFieldFromType(type)
    if (!newField) return null

    const section = newField.section || "Extra fields"
    const extraFields = formFields.filter(f => f.section === section)

    const completeField: FormField = {
      ...newField,
      SK: newField.SK || `FIELD#${type}-${fieldCount.value[type as keyof typeof fieldCount.value]}`,
      slug: newField.slug || `${type}_${fieldCount.value[type as keyof typeof fieldCount.value]}`,
      section: section,
      order: newField.order || extraFields.length + 1,
      fieldTypeName: newField.fieldTypeName || 'Short Text',
      label: newField.label || `${type.charAt(0).toUpperCase() + type.slice(1)} ${fieldCount.value[type as keyof typeof fieldCount.value]}`,
      description: newField.description || '',
      required: newField.required ?? false,
      options: newField.options || [],
      configs: newField.configs || [],
      validations: newField.validations || []
    }

    // Update answerMap with default value
    answerMap[completeField.slug] = getDefaultValue(completeField)

    return completeField
  }

  const resetFieldCounts = () => {
    fieldCount.value = {
      email: 1,
      phone: 0,
      address: 0,
      company: 0,
      position: 0,
      birthdate: 0,
      notes: 0,
      custom: 0,
    }
  }

  return {
    fieldCount,
    getDefaultValue,
    createFieldFromType,
    addField,
    resetFieldCounts
  }
}
