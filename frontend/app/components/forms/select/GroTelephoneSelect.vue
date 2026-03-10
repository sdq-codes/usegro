<script lang="ts" setup>
/**
 * CountryTextInput.vue
 *
 * Combined country code + telephone input that emits a single formatted value
 * Format: "+[country_code][telephone_number]"
 * Example: "+2348012345678"
 */

import { defineProps, ref, watch, computed, onMounted } from 'vue'
import Multiselect from 'vue-multiselect'
import 'vue-multiselect/dist/vue-multiselect.css'
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkSquare01Icon, AlertSquareIcon, CancelSquareIcon } from '@hugeicons/core-free-icons'
import { countryTelephoneOptions } from "@/constants/telephoneOptions"

const options: { value: string | number; label: string; flag?: string }[] = countryTelephoneOptions

const props = defineProps<{
  placeholder?: string
  hint?: string
  charLimit?: number
  disabled?: boolean
  loading?: boolean
  color?: 'primary' | 'success' | 'error'
  readonly?: boolean
}>()

// Use defineModel for two-way binding
const modelValue = defineModel<string>()

// Internal reactive states
const internalCountry = ref<{
  value: string | number
  label: string
  flag?: string
} | null>(null)
const telephone = ref<string>('')

// Flag to prevent infinite loops during parsing
const isUpdating = ref(false)

// Initialize with first country on mount
onMounted(() => {
  if (options.length > 0) {
    internalCountry.value = options[0]
  }

  // Parse existing modelValue if provided (format: "+234801234567")
  if (modelValue.value) {
    parseModelValue(modelValue.value)
  }
})

// Parse incoming modelValue to split country code and telephone
const parseModelValue = (value: string) => {
  if (!value || !value.startsWith('+')) return

  isUpdating.value = true

  // Find matching country by checking if value starts with country label (code)
  const matchingCountry = options.find(opt =>
    value.startsWith('+' + opt.label)
  )

  if (matchingCountry) {
    internalCountry.value = matchingCountry
    // Extract telephone number (everything after country code)
    telephone.value = value.substring(matchingCountry.label.length + 1) // +1 for the '+'
  }

  isUpdating.value = false
}

// Watch for external modelValue changes
watch(modelValue, (newVal) => {
  if (isUpdating.value) return

  if (newVal) {
    parseModelValue(newVal)
  } else {
    telephone.value = ''
  }
})

// Watch for changes in either field and update modelValue
watch([internalCountry, telephone], () => {
  if (isUpdating.value) return

  if (!internalCountry.value || !telephone.value) {
    modelValue.value = ''
    return
  }

  // Format: +[country_code][telephone_number]
  modelValue.value = `${internalCountry.value.label}${telephone.value}`
}, { deep: true })

// char limit
const charLimit = computed(() => props.charLimit ?? 100)

// color classes (keeps visual style from your input component)
const colorClasses: Record<string, string> = {
  primary: 'border-[#EDEDEE] hover:border-[#94BDD8] focus-within:border-[#1E212B] text-[#4F5435] placeholder-[#939499]',
  success: 'border-[#007458] focus-within:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  error: 'border-[#703425] focus-within:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
}

const labelClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]'
}

const hintClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]'
}

const selectWrapperClass = computed(() => [
  'flex items-center justify-center rounded-lg transition',
  props.disabled ? 'opacity-60 cursor-not-allowed' : 'cursor-pointer',
].join(' '))

const inputClass = computed(() => [
  'flex-1 rounded-lg border-0 px-0 md:px-4 py-3  bg-[#F6F6F7] outline-none text-lg placeholder-[#9CA3AF]',
  props.disabled ? 'cursor-not-allowed opacity-60' : '',
  // we intentionally do border on container, not input, for the combined pill look
].join(' '))

const selectClass = computed(() => {
  return [
    'w-full rounded-lg duration-200 outline-none appearance-none',
    'disabled:bg-[#DBDBDD] disabled:text-[#AFB5B8] disabled:cursor-not-allowed',
    colorClasses[props.color || 'primary'],
  ].join(' ')
})

const containerBorder = computed(() => [
  'rounded-lg',
  colorClasses[props.color || 'primary'],
  props.disabled ? 'bg-[#DBDBDD]' : 'bg-transparent'
].join(' '))

// helpers for Multiselect rendering: option label (flag + label)
function renderOptionLabel(option: { label: string; flag?: string }) {
  if (!option) return ''
  return option.flag ? `${option.flag}  ${option.label}` : option.label
}
</script>

<template>
  <div class="w-full">
    <!-- top row: label left, char counter right -->
    <div class="flex items-start mb-1 justify-between">
      <div :class="['text-xs font-medium', labelClasses[props.color || 'primary']]">
        <slot>Label</slot>
      </div>
    </div>

    <!-- main pill: small country select + input -->
    <div :class="[containerBorder, 'flex items-center gap-3']">
      <!-- small country select -->
      <div :class="selectWrapperClass">
        <Multiselect
          v-model="internalCountry"
          :options="options"
          track-by="value"
          placeholder=""
          label="flag"
          :select-label="''"
          :deselect-label="''"
          :show-labels="false"
          :searchable="true"
          :close-on-select="true"
          :clear-on-select="false"
          :preserve-search="true"
          :class="selectClass"
          :disabled="props.disabled"
          :custom-label="renderOptionLabel"
        >
          <!-- custom single label: show flag + caret -->
          <template #singleLabel="{ option }">
            <div class="flex items-center gap-2 w-12">
              <span
                v-if="option?.flag"
                class="text-xl leading-none"
              >{{ option.flag }}</span>
              <span class="text-sm text-[#111827]">{{ option?.label }}</span>
            </div>
          </template>

          <template #option="{ option }">
            <div class="flex items-center gap-2">
              <span
                v-if="option?.flag"
                class="text-xl leading-none"
              >{{ option.flag }}</span>
              <span class="text-sm">{{ option.label }}</span>
            </div>
          </template>
        </Multiselect>
      </div>

      <input
        v-model="telephone"
        type="number"
        :placeholder="props.placeholder"
        :disabled="props.disabled"
        :readonly="props.readonly"
        :class="inputClass"
        :maxlength="charLimit"
        class="h-[42px]"
      >
    </div>

    <!-- hint row -->
    <div
      v-if="props.hint"
      class="flex items-center gap-2 mt-1 text-sm"
      :class="hintClasses[props.color || 'primary']"
    >
      <HugeiconsIcon
        v-if="props.color === 'success'"
        color="#FFFFFF"
        fill="#00916E"
        :icon="CheckmarkSquare01Icon"
      />
      <HugeiconsIcon
        v-else-if="props.color === 'error'"
        color="#FFFFFF"
        fill="#AF513A"
        :icon="CancelSquareIcon"
      />
      <HugeiconsIcon
        v-else
        :icon="AlertSquareIcon"
        fill="#939499"
        color="#FFFFFF"
      />
      <span class="text-sm text-gray-600">{{ props.hint }}</span>
    </div>
  </div>
</template>

<style scoped>
/* Small multiselect adjustments to make it compact like the screenshot */
.multiselect--compact .multiselect__tags {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0;
  border: 0;
  background: transparent;
}

.multiselect--compact .multiselect__select {
  display: none; /* hide default small caret button if any; we use slot caret */
}

.multiselect--compact .multiselect__single {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0;
  margin: 0;
}

.multiselect--compact .multiselect__input {
  display: none;
}

.multiselect__content {
  border-radius: 8px;
  overflow: hidden;
}

.multiselect__content-wrapper {
  box-shadow: 0px 8px 24px 0px rgba(78, 77, 89, 0.12) !important;
  border: 0 !important;
}

/* Make the small select visually pill-like */
.multiselect--compact .multiselect__single, .multiselect--compact .multiselect__tags {
  width: 100%;
  justify-content: center;
}

/* ensure the main pill height aligns */
.multiselect--compact, .multiselect--compact * {
  height: 100%;
}

/* input placeholder color */
input::placeholder {
  color: #9CA3AF;
}

.multiselect .multiselect--compact{
  border-width: 0 !important;
}

.multiselect__tags:hover{
  border:  0 !important;
}

.multiselect__tags{
  border-radius: 8px !important;
  background-color: #F6F6F7 !important;
}

.multiselect__input {
  background-color: #F6F6F7 !important;
}

.multiselect__tags::placeholder {
  color: #939499 !important;
}

.multiselect__option--highlight {
  background-color: #DBDBDD !important;
  color: #000000 !important;
}

.multiselect__content-wrapper {
  box-shadow: 0px 8px 24px 0px rgba(78, 77, 89, 0.12) !important;
  border: 0 !important;
}

.multiselect__single {
  background-color: #F6F6F7 !important;
}
</style>
