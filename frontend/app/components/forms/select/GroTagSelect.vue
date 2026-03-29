<script setup lang="ts">
interface Option {
  value: string | number
  label: string
}

defineProps<{
  options: Option[]
}>()

const model = defineModel<(string | number)[]>({ default: () => [] })

const toggle = (value: string | number) => {
  if (model.value.includes(value)) {
    model.value = model.value.filter(v => v !== value)
  } else {
    model.value = [...model.value, value]
  }
}

const isSelected = (value: string | number) => model.value.includes(value)
</script>

<template>
  <div class="flex flex-wrap gap-2">
    <button
      v-for="opt in options"
      :key="opt.value"
      type="button"
      class="px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors duration-150 cursor-pointer"
      :class="isSelected(opt.value)
        ? 'bg-[#1E212B] text-white'
        : 'bg-[#F6F6F7] text-[#4B4D55] hover:bg-[#EDEDEE]'"
      @click="toggle(opt.value)"
    >
      {{ opt.label }}
    </button>
    <slot name="action" />
  </div>
</template>
