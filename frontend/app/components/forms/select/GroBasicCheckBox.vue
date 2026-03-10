<script setup lang="ts">
import { defineProps, defineModel, computed } from 'vue'

const props = defineProps<{
  options: { value: string | number; label: string }[]
  name?: string
  color?: 'primary' | 'success' | 'error'
  hint?: string
  disabled?: boolean
}>()

// array of selected values
const model = defineModel<(string | number)[]>()

const wrapperClass = 'grid grid-cols-1 gap-2'

const labelClasses: Record<string, string> = {
  primary: 'text-[#4B4D55]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]',
}

const hintClasses: Record<string, string> = {
  primary: 'text-[#1E212B]',
  success: 'text-[#00916E]',
  error: 'text-[#AF513A]',
}

const hintClass = computed(() => {
  return [
    'flex items-center gap-1 mt-1 text-xs',
    hintClasses[props.color || 'primary'],
  ].join(' ')
})
</script>

<template>
  <div class="w-full">
    <div :class="wrapperClass">
      <label
        v-for="opt in props.options"
        :key="opt.value"
        class="flex items-center space-x-2 cursor-pointer"
      >
        <!-- Outer box -->
        <div
          class="h-5 w-5 rounded-md border-2 border-[#DBDBDD] flex items-center justify-center"
        >
          <!-- Checkbox input -->
          <input
            v-model="model"
            type="checkbox"
            :name="props.name"
            :value="opt.value"
            class="hidden peer"
            :disabled="props.disabled"
          >
          <!-- Checkmark -->
          <svg
            class="h-3 w-3 text-[#1E212B] hidden peer-checked:block"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M5 13l4 4L19 7"
            />
          </svg>
        </div>

        <!-- Label -->
        <h6
          class="text-xs font-medium"
          :class="labelClasses[props.color || 'primary']"
        >
          {{ opt.label }}
        </h6>
      </label>
    </div>

    <!-- Hint -->
    <div
      v-if="props.hint"
      :class="hintClass"
    >
      <span>{{ props.hint }}</span>
    </div>
  </div>
</template>
