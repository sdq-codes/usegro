<script lang="ts" setup>
/**
 * CountryTextInput.vue
 *
 * Props:
 *  - options: { value, label, flag? }[]  (flag can be emoji or image URL)
 *  - modelValue: string                  (text input's v-model)
 *  - countryValue: string|number|null    (selected country value)
 *  - placeholder?: string
 *  - hint?: string
 *  - charLimit?: number
 *  - disabled?: boolean
 *  - loading?: boolean
 *  - color?: 'primary'|'success'|'error'
 *  - readonly?: boolean
 *
 * Emits:
 *  - update:modelValue
 *  - update:countryValue
 */

import { ref, watch, computed } from 'vue'
import Multiselect from 'vue-multiselect'
import 'vue-multiselect/dist/vue-multiselect.css'
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkSquare01Icon, AlertSquareIcon, CancelSquareIcon } from '@hugeicons/core-free-icons'
import { STATES } from "@/constants/states"
import { COUNTRY_ISO_TO_STATE_ID } from "@/constants/countryStateMapping"

// All state labels (fallback when no country selected)
const allOptions: string[] = STATES.map(state => state.label)

const props = defineProps<{
  modelValue?: string | null
  countryValue?: string | number | null
  placeholder?: string
  hint?: string
  charLimit?: number
  disabled?: boolean
  loading?: boolean
  color?: 'primary' | 'success' | 'error'
  readonly?: boolean
}>()

const emit = defineEmits(['update:modelValue', 'update:countryValue'])

const country = defineModel<string>()

// Filter states to only those matching the selected country ISO code
const options = computed<string[]>(() => {
  const iso = props.countryValue as string | null | undefined
  if (!iso) return allOptions
  const countryId = COUNTRY_ISO_TO_STATE_ID[iso]
  if (!countryId) return allOptions
  return STATES.filter(s => s.countryId === countryId).map(s => s.label)
})

// Reset state when country changes
watch(() => props.countryValue, () => {
  country.value = undefined
})

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

const selectClass = computed(() => {
  return [
    'w-full rounded-lg duration-200 outline-none appearance-none',
    'disabled:bg-[#DBDBDD] disabled:text-[#AFB5B8] disabled:cursor-not-allowed',
    colorClasses[props.color || 'primary'],
  ].join(' ')
})

// // helpers for Multiselect rendering: option label (flag + label)
// function renderOptionLabel(option: { label: string; flag?: string }) {
//   if (!option) return ''
//   return option.flag ? `${option.flag}  ${option.label}` : option.label
// }
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
    <div :class="['flex items-center gap-3']">
      <!-- small country select -->
      <Multiselect
        v-model="country"
        :options="options"
        placeholder=""
        :show-labels="false"
        :searchable="true"
        :close-on-select="true"
        :clear-on-select="false"
        :preserve-search="true"
        :class="selectClass"
        :disabled="props.disabled"
      />
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
