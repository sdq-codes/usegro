<script setup lang="ts">
import { computed } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { AlertCircleIcon } from '@hugeicons/core-free-icons'

interface Props {
  variant?: 'info' | 'warning' | 'error' | 'success'
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'info'
})

const emit = defineEmits<{
  upgradeClick: []
}>()

const colorClasses: Record<string, {
  bg: string
  text: string
  icon: string
  iconFill: string
  link: string
}> = {
  info: {
    bg: 'bg-gray-100',
    text: 'text-gray-700',
    icon: '#6B7280',
    iconFill: '#F3F4F6',
    link: 'text-blue-600 hover:text-blue-700'
  },
  warning: {
    bg: 'bg-yellow-50',
    text: 'text-yellow-900',
    icon: '#D97706',
    iconFill: '#FEF3C7',
    link: 'text-yellow-700 hover:text-yellow-800'
  },
  error: {
    bg: 'bg-red-50',
    text: 'text-red-900',
    icon: '#DC2626',
    iconFill: '#FEE2E2',
    link: 'text-red-700 hover:text-red-800'
  },
  success: {
    bg: 'bg-[#F5F5F1]',
    text: 'text-[#6F7177]',
    icon: '#ffffff',
    iconFill: '#7C8453',
    link: 'text-green-700 hover:text-green-800'
  }
}

const classes = computed(() => colorClasses[props.variant])

</script>

<template>
  <div
    class="border border-gray-200 rounded-xl px-4 py-3"
    :class="classes.bg"
  >
    <div class="flex items-center gap-3">
      <HugeiconsIcon
        :icon="AlertCircleIcon"
        class="w-5 h-5 flex-shrink-0"
        :color="classes.icon"
        :fill="classes.iconFill"
      />
      <p
        class="text-sm flex-1"
        :class="classes.text"
      >
        <slot name="text" />
      </p>
    </div>
  </div>
</template>
