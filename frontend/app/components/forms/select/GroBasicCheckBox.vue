<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  options: { value: string | number; label: string }[]
  name?: string
  color?: 'primary' | 'success' | 'error'
  hint?: string
  disabled?: boolean
}>()

const model = defineModel<(string | number)[]>()

const isChecked = (value: string | number) => model.value?.includes(value) ?? false

const toggle = (value: string | number) => {
  if (props.disabled) return
  const current = model.value ?? []
  if (current.includes(value)) {
    model.value = current.filter(v => v !== value)
  } else {
    model.value = [...current, value]
  }
}

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

const hintClass = computed(() =>
  ['flex items-center gap-1 mt-1 text-xs', hintClasses[props.color || 'primary']].join(' '),
)
</script>

<template>
  <div class="w-full">
    <div class="grid grid-cols-1 gap-2">
      <label
        v-for="opt in props.options"
        :key="opt.value"
        class="flex items-center space-x-2 cursor-pointer select-none"
        :class="{ 'opacity-50 cursor-not-allowed': disabled }"
        @click.prevent="toggle(opt.value)"
      >
        <!-- Checkbox box -->
        <div
          class="h-5 w-5 rounded-lg flex items-center justify-center shrink-0 transition-colors duration-150"
          :class="isChecked(opt.value)
            ? 'bg-[#2176AE] border-2 border-[#2176AE]'
            : 'bg-white border-2 border-[#DBDBDD]'"
        >
          <svg
            v-if="isChecked(opt.value)"
            class="h-3 w-3 text-white"
            fill="none"
            stroke="currentColor"
            stroke-width="3"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
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
    <div v-if="props.hint" :class="hintClass">
      <span>{{ props.hint }}</span>
    </div>
  </div>
</template>
