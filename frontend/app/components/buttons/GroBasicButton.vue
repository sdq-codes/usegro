<script lang="ts" setup>
import { defineProps, computed } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { Loading03Icon } from '@hugeicons/core-free-icons'

// Props
const props = defineProps<{
  color?: 'primary' | 'secondary' | 'tertiary' | 'custom',
  size?: 'xs' | 'sm' | 'md' | 'lg',
  iconFront?: string,   // e.g. "fa fa-user" or "mdi mdi-home"
  iconBack?: string,
  loading?: boolean,
  disabled?: boolean,
  shape?: 'round' | 'square' | 'custom',
  borderRadius?: string, // e.g. "8px"
  backgroundColor?: string, // only applies if color = 'custom'
}>()

// Classes for colors
const colorClasses: Record<string, string> = {
  primary: 'bg-[#1E212B] cursor-pointer text-white hover:bg-[#6F7177] font-semibold active:bg-[#4B4D55] disabled:bg-[#B7B8BB] disabled:cursor-not-allowed',
  secondary: 'bg-[#DBDBDD] cursor-pointer text-[#1E212B] hover:bg-[#DBDBDD] hover:text-[#6F7177] active:bg-[#B7B8BB] active:text-[#4B4D55] disabled:bg-[#EDEDEE] disabled:text-[#B7B8BB] disabled:cursor-not-allowed font-semibold',
  tertiary: 'border border-[#DBDBDD] cursor-pointer font-semibold text-[#1E212B] hover:border-[#DBDBDD] hover:bg-[#DBDBDD] active:bg-[#DBDBDD] disabled:bg-[#F6F6F7] disabled:border-0 disabled:text-[#B7B8BB] disabled:cursor-not-allowed',
  custom: '' // handled inline
}

// Classes for sizes
const sizeClasses: Record<string, string> = {
  xs: 'px-0 py-1 text-xs',
  sm: 'px-1 py-1 text-sm',
  md: 'px-2 py-1 text-base',
  lg: 'px-4 py-1 text-lg'
}

// Shape classes
const shapeClasses: Record<string, string> = {
  round: 'rounded-full',
  square: 'rounded-none',
  custom: 'rounded-lg',
}

// Computed class list
const btnClass = computed(() => {
  return [
    'inline-flex items-center gap-1 transition-colors duration-200 px-4 py-2',
    colorClasses[props.color || 'primary'],
    sizeClasses[props.size || 'md'],
    shapeClasses[props.shape || 'custom'],
    props?.loading ? "cursor-not-allowed" : ""
  ].join(' ')
})

// Inline styles (for custom background & border radius)
const btnStyle = computed(() => {
  const style: Record<string, string> = {}
  if (props.color === 'custom' && props.backgroundColor) {
    style.backgroundColor = props.backgroundColor
    style.color = '#fff'
  }
  if (props.shape === 'custom' && props.borderRadius) {
    style.borderRadius = props.borderRadius
  }
  return style
})
</script>

<template>
  <button
    :class="btnClass"
    class="inline-flex w-full justify-center shadow-[0px_1px_0px_0px_#E3E3E3_inset,1px_0px_0px_0px_#E3E3E3_inset,-1px_0px_0px_0px_#E3E3E3_inset,0px_-1px_0px_0px_#B5B5B5_inset]"
    :disabled="props?.disabled"
    :style="btnStyle"
  >
    <slot
      v-if="!props?.loading"
      name="frontIcon"
    />
    <slot
      v-if="!props?.loading"
      name="default"
    />
    <slot
      v-if="!props?.loading"
      name="backIcon"
    />
    <HugeiconsIcon
      v-if="props?.loading"
      :icon="Loading03Icon"
      :size="16"
      color="currentColor"
      class="animate-spin"
      :stroke-width="3"
    />
  </button>
</template>

<style scoped>
</style>


aws s3api put-bucket-versioning \
--bucket usegro-staging-terraform-state \
--versioning-configuration Status=Enabled

