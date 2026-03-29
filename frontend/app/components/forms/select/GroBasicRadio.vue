<script setup lang="ts">
import { defineProps, defineModel, computed } from 'vue'
import { HugeiconsIcon } from '@hugeicons/vue'
import { CheckmarkSquare01Icon, CancelSquareIcon, AlertSquareIcon } from '@hugeicons/core-free-icons'

interface Props {
  options: { value: string | number; label: string, description?: string }[]
  name?: string
  color?: 'primary' | 'success' | 'error'
  hint?: string
  borderClass?: string
  disabled?: boolean
  layout: 'horizontal' | 'vertical'
}

const props = withDefaults(defineProps<Props>(), {
  layout: 'horizontal',
  name: '',
  color: 'primary',
  hint: '',
  borderClass: '',
  disabled: false,
})

const model = defineModel<string | number>()

const wrapperClass = 'grid grid-cols-2 md:grid-cols-1 gap-6 h-full'

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
      <div class="col-span-1">
        <div

          class="grid grid-cols-1 h-full gap-x-6"
          :class="[props.layout === 'horizontal' ? 'md:grid-cols-2' : 'md:grid-cols-1']"
        >
          <label
            v-for="opt in props.options"
            :key="opt.value"
            class="col-span-1 mt-0.5"
            :class="[borderClass]"
          >
            <div class="flex items-center gap-x-2 cursor-pointer">
              <input
                v-model="model"
                type="radio"
                :name="props.name"
                :value="opt.value"
                class="hidden peer"
                :disabled="props.disabled"
              >
              <div
                class="h-5 w-5 rounded-full border-2 border-[#DBDBDD] flex items-center justify-center"
              >
                <span
                  class="h-3 w-3 rounded-full transition-colors duration-150"
                  :class="model === opt.value ? 'bg-[#2176AE]' : ''"
                />
              </div>
              <div class="gap-y-3">
                <h6
                  class="text-sm font-semibold"
                  :class="[ labelClasses[props.color || 'primary'], props.layout === 'horizontal' ? '' : 'mt-5']"
                >
                  {{ opt.label }}
                </h6>
                <h6
                  v-if="opt.description"
                  class="text-xs font-regular"
                  :class="[props.layout === 'horizontal' ? '' : 'mt-2']"
                >
                  {{ opt.description }}
                </h6>
              </div>
            </div>
          </label>
        </div>
      </div>
    </div>

    <!-- Hint with icon -->
    <div
      v-if="props.hint"
      :class="hintClass"
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
      <span>{{ props.hint }}</span>
    </div>
  </div>
</template>
