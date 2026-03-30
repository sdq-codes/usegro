<script setup lang="ts">
defineProps<{
  label?: string
  disabled?: boolean
}>()

const model = defineModel<boolean>({ default: false })
</script>

<template>
  <label class="flex items-center gap-3 select-none w-fit" :class="disabled ? 'cursor-not-allowed opacity-50' : 'cursor-pointer'">
    <button
      type="button"
      role="switch"
      :aria-checked="model"
      :disabled="disabled"
      class="relative w-11 h-6 rounded-full transition-colors duration-200 shrink-0 focus:outline-none focus-visible:ring-2 focus-visible:ring-[#1E212B] focus-visible:ring-offset-2"
      :class="[model ? 'bg-[#1E212B]' : 'bg-[#DBDBDD]', disabled ? 'cursor-not-allowed' : 'cursor-pointer']"
      @click="!disabled && (model = !model)"
    >
      <span
        class="absolute top-0.5 left-0.5 w-5 h-5 bg-white rounded-full shadow-sm transition-transform duration-200"
        :class="model ? 'translate-x-5' : 'translate-x-0'"
      />
    </button>
    <span
      v-if="label || $slots.default"
      class="text-sm text-[#1E212B]"
    >
      <slot>{{ label }}</slot>
    </span>
  </label>
</template>
