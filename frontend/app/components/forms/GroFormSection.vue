<script setup lang="ts">
import GroSectionHeader from './GroSectionHeader.vue'

defineProps<{
  title: string
  open: boolean
  bordered?: boolean
}>()

const emit = defineEmits<{
  (e: 'toggle'): void
}>()
</script>

<template>
  <div :class="bordered !== false ? 'border-b border-[#EDEDEE]' : ''">
    <GroSectionHeader
      :title="title"
      :open="open"
      @toggle="emit('toggle')"
    />
    <div
      v-show="open"
      class="px-8 pb-8"
    >
      <slot />
    </div>
  </div>
</template>
