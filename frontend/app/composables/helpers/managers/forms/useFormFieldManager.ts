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

  const buildCompleteField = (raw: Record<string, unknown>, type: string, index?: number): FormField => {
    const section = (raw.section as string) || "Extra fields"
    const suffix = index !== undefined ? `${fieldCount.value[type as keyof typeof fieldCount.value]}-${index}` : fieldCount.value[type as keyof typeof fieldCount.value]
    return {
      ...raw,
      SK: raw.SK || `FIELD#${type}-${suffix}`,
      slug: raw.slug || `${type}_${suffix}`,
      section,
      order: raw.order || 1,
      fieldTypeName: (raw.fieldTypeName as string) || 'Short Text',
      label: (raw.label as string) || `${type.charAt(0).toUpperCase() + type.slice(1)} ${fieldCount.value[type as keyof typeof fieldCount.value]}`,
      description: (raw.description as string) || '',
      required: (raw.required as boolean) ?? false,
      options: (raw.options as FormField['options']) || [],
      configs: (raw.configs as FormField['configs']) || [],
      validations: (raw.validations as FormField['validations']) || [],
    } as FormField
  }

  const addField = (
    type: string,
    formFields: FormField[],
    answerMap: Record<string, string | string[]>
  ): FormField | FormField[] | null => {
    const result = createFieldFromType(type)
    if (!result) return null

    if (Array.isArray(result)) {
      const fields = result.map((raw, i) => buildCompleteField(raw as Record<string, unknown>, type, i))
      fields.forEach(f => { answerMap[f.slug] = getDefaultValue(f) })
      return fields
    }

    const completeField = buildCompleteField(result as Record<string, unknown>, type)
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
