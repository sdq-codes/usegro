<script setup lang="ts">
defineProps<{
  title: string
  open: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()
</script>

<template>
  <button
    class="w-full flex cursor-pointer items-center justify-between px-8 py-5 hover:bg-[#F6F6F7] transition-colors"
    @click="emit('toggle')"
  >
    <h3 class="text-base font-semibold text-[#1E212B]">
      {{ title }}
    </h3>
    <svg
      class="w-4 h-4 text-[#6F7177] transition-transform shrink-0 border border-[#EDEDEE] rounded-full"
      :class="open ? 'rotate-180' : ''"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      stroke-width="2"
    >
      <path d="M6 9l6 6 6-6" />
    </svg>
  </button>
</template>
