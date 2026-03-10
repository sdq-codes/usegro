<script lang="ts" setup>
import { computed } from 'vue'
import Multiselect from 'vue-multiselect'
import 'vue-multiselect/dist/vue-multiselect.css'
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkSquare01Icon, AlertSquareIcon, CancelSquareIcon } from '@hugeicons/core-free-icons'
import { COUNTRIES } from "@/constants/countries"

const options: string[] = COUNTRIES.map(country => country.label)

const props = defineProps<{
  placeholder?: string
  hint?: string
  disabled?: boolean
  loading?: boolean
  color?: 'primary' | 'success' | 'error'
  readonly?: boolean
}>()

const modelValue = defineModel<string>()

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
</script>

<template>
  <div class="w-full">
    <!-- top row: label -->
    <div class="flex items-start mb-1 justify-between">
      <div :class="['text-xs font-medium', labelClasses[props.color || 'primary']]">
        <slot>Country</slot>
      </div>
    </div>

    <!-- country select -->
    <div :class="['flex items-center gap-3']">
      <Multiselect
        v-model="modelValue"
        :options="options"
        :placeholder="props.placeholder || 'Select a country'"
        :show-labels="false"
        :searchable="true"
        :close-on-selsfsect="true"
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
/* Small multiselect adjustments */
.multiselect--compact .multiselect__tags {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0;
  border: 0;
  background: transparent;
}

.multiselect--compact .multiselect__select {
  display: none;
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

.multiselect--compact .multiselect__single, .multiselect--compact .multiselect__tags {
  width: 100%;
  justify-content: center;
}

.multiselect--compact, .multiselect--compact * {
  height: 100%;
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
