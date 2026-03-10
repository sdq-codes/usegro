<script lang="ts" setup>
import {defineProps, defineEmits, ref, onMounted, computed, defineModel} from 'vue'
import {AlertSquareIcon, CancelSquareIcon, CheckmarkSquare01Icon} from "@hugeicons/core-free-icons";
import {HugeiconsIcon} from "@hugeicons/vue";

const props = defineProps<{
  length?: number,
  hint?: string,
  disabled?: string,
  readonly?: boolean,
  color?: 'primary' | 'success' | 'error',
  value?: string
}>()

const model = defineModel<string>();

const emit = defineEmits<{
  (e: 'update:value', value: string): void
}>()

const inputs = ref<(HTMLInputElement | null)[]>([])
const values = ref<string[]>(Array(props.length || 4).fill(''))

const colorClasses: Record<string, string> = {
  primary: 'border-[#EDEDEE] hover:border-[#94BDD8] focus:border-[#1E212B] text-[#4F5435] placeholder-[#939499]',
  success: 'border-[#007458] focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
  error: 'border-[#703425] focus:bg-[#FFFFFF] text-[#4F5435] placeholder-[#939499]',
}

const baseClass =
  'flex-1 min-w-0 h-10 w-6 sm:h-12 text-center text-base sm:text-sm md:text-lg rounded-lg border border-solid outline-none bg-[#F6F6F7] focus:bg-[#FFFFFF] transition-all duration-200'

const labelClasses: Record<string, string> = {
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

function onInput(e: Event, index: number) {
  const target = e.target as HTMLInputElement
  const val = target.value

  values.value[index] = val
  emit('update:value', values.value.join(''))
  model.value = values.value.join('')

  if (val && index < (props.length || 4) - 1) {
    inputs.value[index + 1]?.focus()
  }
}

function onKeyDown(e: KeyboardEvent, index: number) {
  if (e.key === 'Backspace' && !values.value[index] && index > 0) {
    inputs.value[index - 1]?.focus()
  }
}

onMounted(() => {
  if (props.value) {
    values.value = props.value.split('').slice(0, props.length || 4)
  }
})
</script>

<template>
  <div>
    <div :class="labelClass">
      <slot />
    </div>
    <div class="flex w-full justify-between gap-1.5 sm:gap-2">
      <input
        v-for="(_, i) in (props.length || 4)"
        :key="i"
        ref="inputs"
        v-model="values[i]"
        type="text"
        maxlength="1"
        :disabled="props?.disabled"
        :readonly="props?.readonly"
        :class="[baseClass, colorClasses[props.color || 'primary'], props.disabled && 'cursor-not-allowed disabled:bg-[#DBDBDD] text-[#AFB5B8]']"
        @input="onInput($event, i)"
        @keydown="onKeyDown($event, i)"
      >
    </div>
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
