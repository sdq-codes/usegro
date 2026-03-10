<template>
  <form>
    <div
      v-for="field in visibleFields"
      :key="field.SK"
      class="mb-4"
    >
      <component
        :is="getComponent(field.fieldTypeName)"
        v-model="model[field.slug]"
        :options="field.options"
        :info="field.description"
        :field="field"
        :required="field.required"
        :type="field.fieldTypeName === 'Date' ? 'date' : field.fieldTypeName === 'Email' ? 'email' : 'text'"
        class="mb-1"
      >
        <template #default>
          <span>{{ field.label }}</span>
          <span
            v-if="field.required"
            class="text-red-500 ml-1"
          >*</span>
        </template>
      </component>
      <p
        v-if="displayErrors[field.slug]"
        class="text-red-500 text-xs mt-1"
      >
        {{ displayErrors[field.slug] }}
      </p>
    </div>
  </form>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import GroBasicInput from '../input/GroBasicInput.vue'
import GroBasicTextArea from '../input/GroBasicTextArea.vue'
import GroBasicRadio from '../select/GroBasicRadio.vue'
import GroTelephoneSelect from '../select/GroTelephoneSelect.vue'
import GroBasicCheckBox from '../select/GroBasicCheckBox.vue'
import GroCountrySelect from '../select/GroCountrySelect.vue'
import GroCountryState from '@/components/forms/select/GroCountryState.vue'
import GroBasicTagSelect from '@/components/forms/select/GroBasicTagSelect.vue'

type Validation = {
  key: string;
  rule: string;
  message: string;
};

type Config = {
  key: string;
  valueType: string;
  fieldId?: string;
  fieldSlug?: string;
  fieldValue?: string;
};

type FormField = {
  SK: string;
  slug: string;
  label: string;
  fieldTypeName: string;
  description?: string;
  required?: boolean;
  options?: Array<{ value: string; label: string }>;
  configs?: Config[];
  validations?: Validation[];
};

interface Props {
  fields: FormField[]
  errors?: Record<string, string>  // Accept errors from parent
}

const props = withDefaults(defineProps<Props>(), {
  errors: () => ({})
})

const model = defineModel<Record<string, string | string[]>>({ required: true })
const internalErrors = ref<Record<string, string>>({})

// Use parent errors if provided, otherwise use internal errors
const displayErrors = computed(() => {
  // If parent errors exist and have values, use them
  if (props.errors && Object.keys(props.errors).length > 0) {
    return props.errors
  }
  // Otherwise use internal errors
  return internalErrors.value
})

// Component mapping
const COMPONENT_MAP = {
  'Short Text': GroBasicInput,
  'City': GroBasicInput,
  'Email': GroBasicInput,
  'Date': GroBasicInput,
  'Long Text': GroBasicTextArea,
  'Radio Button': GroBasicRadio,
  'Phone Number': GroTelephoneSelect,
  'Checkbox': GroBasicCheckBox,
  'Country': GroCountrySelect,
  'State': GroCountryState,
  'Tags': GroBasicTagSelect,
} as const

function getComponent(type: string) {
  return COMPONENT_MAP[type as keyof typeof COMPONENT_MAP] || 'div'
}

function extractShowIf(configs?: Config[]): Config | null {
  if (!configs?.length) return null
  return configs.find(cfg => cfg.key === 'showIf') ?? null
}

function shouldShowField(field: FormField): boolean {
  const showIfConfig = extractShowIf(field.configs)

  if (!showIfConfig) return true

  if (showIfConfig.valueType === 'string' && showIfConfig.fieldSlug) {
    return model.value[showIfConfig.fieldSlug] === showIfConfig.fieldValue
  }

  return false
}

// Computed property for visible fields
const visibleFields = computed(() => {
  return props.fields.filter(shouldShowField)
})

// Clear errors for hidden fields
watch(visibleFields, (newFields) => {
  const visibleSlugs = new Set(newFields.map(f => f.slug))
  Object.keys(internalErrors.value).forEach(slug => {
    if (!visibleSlugs.has(slug)) {
      delete internalErrors.value[slug]
    }
  })
})

// Expose validation function for parent component
defineExpose({
  errors: displayErrors
})
</script>
