<script lang="ts" setup>
import { defineProps, computed } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import {CheckmarkSquare01Icon, AlertSquareIcon, CancelSquareIcon} from '@hugeicons/core-free-icons'

const props = defineProps<{
  placeholder?: string,
  hint?: string,
  disabled?: boolean,
  loading?: boolean,
  color?: 'primary' | 'success' | 'error',
  readonly?: boolean,
}>()

// Base classes for input
const inputClass = computed(() => {
  return [
    'w-full rounded-lg border border-solid px-3 py-2 transition-all duration-200 outline-none bg-[#F6F6F7] focus:bg-[#FFFFFF]',
    'disabled:bg-[#DBDBDD] disabled:text-[#AFB5B8] disabled:border-0 disabled:cursor-not-allowed',
    colorClasses[props.color || 'primary'],
  ].join(' ')
})

const colorClasses: Record<string, string> = {
  primary: 'border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] text-[#4F5435] placeholder-[#939499]',
  success: 'border-[#007458] focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  error: 'border-[#703425] focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  custom: ''
}

const labelClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]',
}

const hintClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]',
}

const labelClass = computed(() => {
  return [
    'text-xs font-medium mb-1',
    labelClasses[props.color || 'primary'],
  ].join(' ')
})

const hintClass = computed(() => {
  return [
    hintClasses[props.color || 'primary'],
  ].join(' ')
})

// Wrapper for positioning trailing icon
const wrapperClass = 'relative flex flex-col gap-1 mt-3'

</script>

<template>
  <div class="w-full">
    <div :class="labelClass">
      <slot />
    </div>


    <!-- Input wrapper with trailing icon -->
    <div :class="wrapperClass">
      <textarea
        type="text"
        rows="5"
        :placeholder="placeholder"
        :disabled="props?.disabled"
        :class="inputClass"
        :readonly="props.readonly"
      />
    </div>

    <!-- Hint with icon -->
    <div
      v-if="props.hint"
      class="flex items-center gap-1 mt-1 text-xs"
      :class="hintClass"
    >
      <HugeiconsIcon
        v-if="props.hint && props.color === 'success'"
        color="#FFFFFF"
        fill="#00916E"
        :icon="CheckmarkSquare01Icon"
      />
      <HugeiconsIcon
        v-else-if="props.hint && props.color === 'error'"
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

      <span>{{ props.hint }}</span>
    </div>
  </div>
</template>
