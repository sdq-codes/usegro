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
        v-bind="{
          ...(field.fieldTypeName === 'State' ? { countryValue: model[field.configs?.find(c => c.key === 'countrySlug')?.fieldSlug ?? 'country'] } : {}),
          ...(field.fieldTypeName === 'Radio Button' ? { layout: 'horizontal' } : {}),
        }"
        class="mb-1"
      >
        <template #default>
          <template v-if="editableLabels">
            <div class="flex items-center gap-1 group">
              <input
                v-if="editingSlug === field.slug"
                :ref="(el) => { if (el) focusInput(el as HTMLInputElement) }"
                :value="field.label"
                class="border-b border-blue-500 focus:outline-none bg-transparent text-sm"
                @change="emit('label-changed', { slug: field.slug, label: ($event.target as HTMLInputElement).value })"
                @blur="editingSlug = null"
                @keydown.enter="editingSlug = null"
                @keydown.escape="editingSlug = null"
              />
              <span v-else class="text-sm">{{ field.label }}</span>
              <button
                type="button"
                class="inline-flex items-center justify-center w-6 h-6 rounded transition-all cursor-pointer"
                :class="editingSlug === field.slug
                  ? 'bg-blue-500 text-white shadow'
                  : 'bg-gray-100 text-gray-500 hover:bg-blue-100 hover:text-blue-600'"
                @click.prevent="editingSlug = editingSlug === field.slug ? null : field.slug"
              >
                <HugeiconsIcon :icon="PencilEditIcon" :size="13" :stroke-width="2" />
              </button>
            </div>
          </template>
          <template v-else>
            <span>{{ field.label }}</span>
          </template>
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
import { HugeiconsIcon } from '@hugeicons/vue'
import { PencilEditIcon } from '@hugeicons/core-free-icons'
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
  editableLabels?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  errors: () => ({}),
  editableLabels: false
})

const emit = defineEmits<{
  (e: 'label-changed', payload: { slug: string; label: string }): void
}>()

const model = defineModel<Record<string, string | string[]>>({ required: true })
const internalErrors = ref<Record<string, string>>({})
const editingSlug = ref<string | null>(null)

const focusInput = (el: HTMLInputElement) => {
  el.focus()
  el.select()
}

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
