<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import CustomFormField from "@/components/forms/CustomForm/CustomFormField.vue"
import { type FormField } from "@/composables/helpers/types/form"

interface Props {
  fields?: FormField[]
  groupedFields?: Record<string, FormField[]>
  title?: string
  isLoading?: boolean
  layout?: 'horizontal' | 'vertical'
  errors?: Record<string, string | string[]>
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Form',
  isLoading: false,
  layout: 'horizontal',
  fields: () => [],
  groupedFields: () => ({}),
  errors: () => ({})
})

// Use v-model for answerMap
const answerMap = defineModel<Record<string, string | string[]>>({
  required: false,
  default: () => ({})
})

const formFields = ref<FormField[]>([])
const fieldsMap = ref<Record<string, FormField>>({})
const customFormRefs = ref({})

// Helper functions defined first
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

const initializeFormData = (fields: FormField[]) => {
  // Only initialize if answerMap is empty
  if (Object.keys(answerMap.value).length === 0) {
    answerMap.value = fields.reduce((acc, field) => {
      acc[field.slug] = getDefaultValue(field)
      return acc
    }, {} as Record<string, string | string[]>)
  }

  fieldsMap.value = fields.reduce((acc, field) => {
    const match = field.SK?.match(/FIELD#(.+)$/)
    if (match) {
      acc[match[1]] = field
    }
    return acc
  }, {} as Record<string, FormField>)
}

// Watch for external field changes (if fields prop is used)
watch(() => props.fields, (newFields) => {
  if (newFields && newFields.length > 0) {
    formFields.value = [...newFields]
    initializeFormData(newFields)
  }
}, { immediate: true, deep: true })

// Computed property for grouped fields
const groupFields = computed(() => {
  // If groupedFields prop is provided, use it directly
  if (props.groupedFields && Object.keys(props.groupedFields).length > 0) {
    return props.groupedFields
  }

  // Otherwise, group from formFields
  const grouped = formFields.value.reduce((acc, field) => {
    if (!acc[field.section]) {
      acc[field.section] = []
    }
    acc[field.section].push(field)
    return acc
  }, {} as Record<string, FormField[]>)

  Object.keys(grouped).forEach((section) => {
    grouped[section].sort((a, b) => a.order - b.order)
  })

  return grouped
})

// Computed layout classes
const containerClasses = computed(() => {
  if (props.layout === 'vertical') {
    return 'grid grid-cols-1 gap-6 mb-6'
  }
  return 'grid grid-cols-1 md:grid-cols-3 gap-6 mb-6'
})

const labelClasses = computed(() => {
  if (props.layout === 'vertical') {
    return 'col-span-1 mb-2'
  }
  return 'col-span-1'
})

const sectionLabelClasses = computed(() => {
  if (props.layout === 'vertical') {
    return 'text-md'
  }
  return 'text-sm'
})

const fieldContainerClasses = computed(() => {
  if (props.layout === 'vertical') {
    return 'col-span-1'
  }
  return 'col-span-2 my-auto'
})

const validateForm = (): boolean => {
  let isValid = true

  if (!customFormRefs.value || Object.keys(customFormRefs.value).length === 0) {
    return true
  }

  for (const key in customFormRefs.value) {
    const formRef = customFormRefs.value[key]
    if (formRef && typeof formRef.validateForm === 'function') {
      const sectionValid = formRef.validateForm()
      if (!sectionValid) {
        isValid = false
      }
    }
  }

  return isValid
}

const resetForm = () => {
  initializeFormData(formFields.value)
  customFormRefs.value = {}
}

const emit = defineEmits<{
  (e: 'update:fields', fields: FormField[]): void
}>()

const handleLabelChanged = (payload: { slug: string; label: string }) => {
  const field = formFields.value.find(f => f.slug === payload.slug)
  if (field) {
    field.label = payload.label
    emit('update:fields', [...formFields.value])
  }
}

// Expose methods for parent component
defineExpose({
  resetForm,
  validateForm,
  getFormData: () => answerMap.value
})
</script>

<template>
  <div class="dynamic-form">
    <div
      v-if="isLoading"
      class="flex items-center justify-center py-8"
    >
      <span class="text-gray-500">Loading form...</span>
    </div>

    <div
      v-else
      class="overflow-y-scroll"
    >
      <div
        v-for="(dynamicFormField, key) in groupFields"
        :key="key"
        :class="containerClasses"
      >
        <div :class="labelClasses">
          <label
            class="font-semibold"
            :class="sectionLabelClasses"
          >
            {{ key }}
          </label>
        </div>
        <div :class="fieldContainerClasses">
          <CustomFormField
            :ref="(el) => { if (el) customFormRefs[key] = el }"
            v-model="answerMap"
            :fields="dynamicFormField"
            :errors="errors"
            :editable-labels="key === 'Extra fields'"
            @label-changed="handleLabelChanged"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dynamic-form {
  width: 100%;
}
</style>
